package notification

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	notificationv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notification/v1"
)

const peer = "notification"

type Notifier struct {
	conn      *grpc.ClientConn
	stub      notificationv1.NotificationServiceClient
	timeout   time.Duration
	verifyURL string
}

var _ notary.StatementNotifier = (*Notifier)(nil)

func Dial(settings Settings, options ...grpc.DialOption) (*Notifier, error) {
	const op = "notification.Dial"

	if settings.Timeout <= 0 {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "a call to %s needs a positive timeout", peer)
	}

	if settings.VerifyURL == "" {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "a statement needs the address of the verifier")
	}

	creds, err := mesh.ClientCredentials(settings.Identity, peer)
	if err != nil {
		return nil, err
	}

	options = append([]grpc.DialOption{grpc.WithTransportCredentials(creds)}, options...)

	conn, err := grpc.NewClient(settings.Endpoint, options...)
	if err != nil {
		return nil, notary.Errf(notary.KindInvalidArgument, op, err, "notification endpoint %q", settings.Endpoint)
	}

	return &Notifier{
		conn:      conn,
		stub:      notificationv1.NewNotificationServiceClient(conn),
		timeout:   settings.Timeout,
		verifyURL: settings.VerifyURL,
	}, nil
}

func (n *Notifier) Close() error {
	if err := n.conn.Close(); err != nil {
		return notary.Errf(notary.KindUnknown, "notification.Notifier.Close", err, "")
	}

	return nil
}

func (n *Notifier) NotifyStatement(ctx context.Context, holder notary.Holder, epoch uint64) (string, error) {
	const op = "notification.Notifier.NotifyStatement"

	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	response, err := n.stub.Notify(ctx, &notificationv1.NotifyRequest{
		RecipientPrincipalId: holder.String(),
		Template:             notificationv1.Template_TEMPLATE_LIABILITY_STATEMENT,
		Params: map[string]string{
			"epoch":      strconv.FormatUint(epoch, 10),
			"verify_url": n.verifyURL,
		},
		IdempotencyKey: &commonv1.IdempotencyKey{Key: IdempotencyKey(holder, epoch)},
	})
	if err != nil {
		return "", classify(op, err)
	}

	if response.GetAccepted() {
		if response.GetNotificationId() == "" {
			return "", notary.Errf(
				notary.KindNotifierRefused,
				op,
				nil,
				"notification accepted a statement without naming it",
			)
		}

		return response.GetNotificationId(), nil
	}

	return "", refusal(op, response.GetError())
}

func IdempotencyKey(holder notary.Holder, epoch uint64) string {
	return fmt.Sprintf("liability-statement:%d:%s", epoch, holder)
}

func classify(op string, err error) error {
	switch status.Code(err) {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Aborted, codes.Internal,
		codes.Unknown, codes.Canceled:
		return notary.Errf(notary.KindNotifierUnavailable, op, err, "notification did not answer")
	default:
		return notary.Errf(notary.KindNotifierRefused, op, err, "notification refused the call")
	}
}

func refusal(op string, detail *commonv1.ErrorDetail) error {
	code := detail.GetErrorCode()

	switch code {
	case "RECIPIENT_UNKNOWN", "NO_CHANNEL_AVAILABLE":
		return notary.Errf(notary.KindRecipientUnreachable, op, nil, "%s: %s", code, detail.GetMessage())
	case "DIRECTORY_UNAVAILABLE", "SENDER_UNAVAILABLE", "STORE_UNAVAILABLE":
		return notary.Errf(notary.KindNotifierUnavailable, op, nil, "%s: %s", code, detail.GetMessage())
	case "":
		return notary.Errf(notary.KindNotifierRefused, op, nil, "notification refused a statement without a reason")
	default:
		return notary.Errf(notary.KindNotifierRefused, op, nil, "%s: %s", code, detail.GetMessage())
	}
}
