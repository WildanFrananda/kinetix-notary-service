package notification_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notification"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	notificationv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notification/v1"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const verifyURL = "https://docs.example.test/id/verify/proof-of-liabilities/"

func dial(t *testing.T, server *notarytest.NotificationServer, as string) *notification.Notifier {
	t.Helper()

	pki := notarytest.NewPKI(t)
	address := notarytest.ServeMesh(t, pki.Issue(t, "kinetix.local", as), func(s *grpc.Server) {
		notificationv1.RegisterNotificationServiceServer(s, server)
	})

	notifier, err := notification.Dial(notification.Settings{
		Endpoint: address,
		Identity: mesh.Identity{
			PKIDir:       pki.Issue(t, "kinetix.local", "notary"),
			TrustDomains: []string{"kinetix.local"},
		},
		Timeout:   5 * time.Second,
		VerifyURL: verifyURL,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = notifier.Close() })

	return notifier
}

func refusing(code string) func(*notificationv1.NotifyRequest) (*notificationv1.NotifyResponse, error) {
	return func(*notificationv1.NotifyRequest) (*notificationv1.NotifyResponse, error) {
		return &notificationv1.NotifyResponse{Error: &commonv1.ErrorDetail{ErrorCode: code, Message: "refused"}}, nil
	}
}

func TestAStatementAsksNotificationForTheLiabilityTemplateWithoutAnAmount(t *testing.T) {
	server := &notarytest.NotificationServer{}
	notifier := dial(t, server, "notification")
	holder := notarytest.Holder(t, "4d1b6f0e-8a8e-4b1f-9a3e-0f0e5d1c2b3a")

	id, err := notifier.NotifyStatement(context.Background(), holder, 184)
	require.NoError(t, err)
	require.Equal(t, "NTF-1", id)

	requests := server.Requests()
	require.Len(t, requests, 1)
	require.Equal(t, holder.String(), requests[0].GetRecipientPrincipalId())
	require.Equal(t, notificationv1.Template_TEMPLATE_LIABILITY_STATEMENT, requests[0].GetTemplate())
	require.Equal(t, map[string]string{"epoch": "184", "verify_url": verifyURL}, requests[0].GetParams())
	require.Equal(t, "liability-statement:184:"+holder.String(), requests[0].GetIdempotencyKey().GetKey())
	require.Empty(t, requests[0].GetChannels())
	require.Equal(t, []string{"notary"}, server.Callers())
}

func TestEveryRefusalFromNotificationIsReadForWhatItMeans(t *testing.T) {
	cases := map[string]notary.Kind{
		"RECIPIENT_UNKNOWN":      notary.KindRecipientUnreachable,
		"NO_CHANNEL_AVAILABLE":   notary.KindRecipientUnreachable,
		"DIRECTORY_UNAVAILABLE":  notary.KindNotifierUnavailable,
		"SENDER_UNAVAILABLE":     notary.KindNotifierUnavailable,
		"STORE_UNAVAILABLE":      notary.KindNotifierUnavailable,
		"TEMPLATE_NOT_GRANTED":   notary.KindNotifierRefused,
		"TEMPLATE_PARAM_MISSING": notary.KindNotifierRefused,
		"IDEMPOTENCY_KEY_REUSED": notary.KindNotifierRefused,
		"NO_TEMPLATE":            notary.KindNotifierRefused,
		"":                       notary.KindNotifierRefused,
	}

	for code, want := range cases {
		t.Run(code, func(t *testing.T) {
			notifier := dial(t, &notarytest.NotificationServer{Answer: refusing(code)}, "notification")

			_, err := notifier.NotifyStatement(context.Background(), notarytest.Holder(t, "principal-a"), 1)
			require.Equal(t, want, notary.KindOf(err), "%v", err)
		})
	}
}

func TestAnAcceptanceThatNamesNoNotificationIsNotADelivery(t *testing.T) {
	server := &notarytest.NotificationServer{
		Answer: func(*notificationv1.NotifyRequest) (*notificationv1.NotifyResponse, error) {
			return &notificationv1.NotifyResponse{Accepted: true}, nil
		},
	}

	_, err := dial(t, server, "notification").NotifyStatement(context.Background(), notarytest.Holder(t, "p"), 1)
	require.Equal(t, notary.KindNotifierRefused, notary.KindOf(err))
}

func TestAStatusFromTheCallIsTransientOnlyWhenNotificationDidNotAnswer(t *testing.T) {
	cases := map[codes.Code]notary.Kind{
		codes.Unavailable:      notary.KindNotifierUnavailable,
		codes.DeadlineExceeded: notary.KindNotifierUnavailable,
		codes.Internal:         notary.KindNotifierUnavailable,
		codes.PermissionDenied: notary.KindNotifierRefused,
		codes.Unauthenticated:  notary.KindNotifierRefused,
		codes.InvalidArgument:  notary.KindNotifierRefused,
	}

	for code, want := range cases {
		t.Run(code.String(), func(t *testing.T) {
			server := &notarytest.NotificationServer{
				Answer: func(*notificationv1.NotifyRequest) (*notificationv1.NotifyResponse, error) {
					return nil, status.Error(code, "scripted")
				},
			}

			_, err := dial(t, server, "notification").NotifyStatement(
				context.Background(),
				notarytest.Holder(t, "principal-a"),
				1,
			)
			require.Equal(t, want, notary.KindOf(err), "%v", err)
		})
	}
}

func TestAServerThatIsNotNotificationIsNeverAsked(t *testing.T) {
	server := &notarytest.NotificationServer{}
	notifier := dial(t, server, "payment")

	_, err := notifier.NotifyStatement(context.Background(), notarytest.Holder(t, "principal-a"), 1)
	require.Equal(t, notary.KindNotifierUnavailable, notary.KindOf(err))
	require.Empty(t, server.Requests())
}

func TestANotifierNeedsATimeoutAndAVerifier(t *testing.T) {
	_, err := notification.Dial(notification.Settings{Endpoint: "127.0.0.1:1", VerifyURL: verifyURL})
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	_, err = notification.Dial(notification.Settings{Endpoint: "127.0.0.1:1", Timeout: time.Second})
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}
