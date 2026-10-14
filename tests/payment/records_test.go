package payment_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/payment"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func recordsRig(t *testing.T, caller string) (*notarytest.LiabilityServer, *payment.Records) {
	t.Helper()

	pki := notarytest.NewPKI(t)
	server := &notarytest.LiabilityServer{}
	address := notarytest.ServeLiabilities(t, pki.Issue(t, "kinetix.local", "payment"), server)
	records, err := payment.DialRecords(payment.Settings{
		Endpoint: address,
		Identity: mesh.Identity{PKIDir: pki.Issue(t, "kinetix.local", caller), TrustDomains: []string{"kinetix.local"}},
		Timeout:  5 * time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = records.Close() })

	return server, records
}

func TestTheWatchReadsBackWhatPaymentRecordedOfASnapshot(t *testing.T) {
	server, records := recordsRig(t, "notary-watch")
	id := notarytest.SnapshotID(t, 11)
	_, kept := notarytest.SnapshotStream(
		t,
		id,
		&timestamppb.Timestamp{Seconds: 1760000000, Nanos: 123_456_000},
		"IDR",
		[]*paymentv1.LiabilityAccount{owedTo(customer, paymentv1.LiabilityKind_LIABILITY_KIND_CUSTOMER_BALANCE, 10)},
	)
	server.Keep(kept)

	record, err := records.Record(context.Background(), id)

	require.NoError(t, err)
	require.Equal(t, id, record.ID)
	require.Equal(t, uint64(1760000000), record.AsOf.Seconds())
	require.Equal(t, uint32(123_456_000), record.Nanos)
	require.Equal(t, "IDR", record.Currency)
	require.Equal(t, uint64(1), record.Count)
	require.Equal(t, uint64(10), record.Total)
	require.Equal(t, kept.GetDigest(), record.Digest[:])
	require.Equal(t, []string{string(id)}, server.Asked)
}

func TestASnapshotPaymentNeverTookIsUnknownNotUnavailable(t *testing.T) {
	_, records := recordsRig(t, "notary-watch")

	_, err := records.Record(context.Background(), notarytest.SnapshotID(t, 12))

	require.Equal(t, notary.KindSnapshotUnknown, notary.KindOf(err))
}

func TestARecordThatDoesNotAnswerTheQuestionIsRefused(t *testing.T) {
	id := notarytest.SnapshotID(t, 13)
	_, kept := notarytest.SnapshotStream(t, id, &timestamppb.Timestamp{Seconds: 1760000000}, "IDR", nil)

	for name, mutate := range map[string]func(*paymentv1.GetLiabilitySnapshotRecordResponse){
		"another snapshot": func(r *paymentv1.GetLiabilitySnapshotRecordResponse) {
			r.SnapshotId = string(notarytest.SnapshotID(t, 14))
		},
		"no time":         func(r *paymentv1.GetLiabilitySnapshotRecordResponse) { r.AsOf = nil },
		"negative count":  func(r *paymentv1.GetLiabilitySnapshotRecordResponse) { r.AccountCount = -1 },
		"short digest":    func(r *paymentv1.GetLiabilitySnapshotRecordResponse) { r.Digest = r.Digest[:31] },
		"total elsewhere": func(r *paymentv1.GetLiabilitySnapshotRecordResponse) { r.Total.Currency = "USD" },
		"negative total":  func(r *paymentv1.GetLiabilitySnapshotRecordResponse) { r.Total.AmountMinor = -1 },
	} {
		t.Run(name, func(t *testing.T) {
			server, records := recordsRig(t, "notary-watch")
			broken, ok := proto.Clone(kept).(*paymentv1.GetLiabilitySnapshotRecordResponse)
			require.True(t, ok)
			mutate(broken)
			server.Records = map[string]*paymentv1.GetLiabilitySnapshotRecordResponse{string(id): broken}

			_, err := records.Record(context.Background(), id)

			require.Equal(t, notary.KindSourceRefused, notary.KindOf(err))
		})
	}
}

func TestPaymentsStatusDecidesWhetherTheWatchRetries(t *testing.T) {
	for code, kind := range map[codes.Code]notary.Kind{
		codes.PermissionDenied: notary.KindSourceRefused,
		codes.InvalidArgument:  notary.KindSourceRefused,
		codes.Unavailable:      notary.KindSourceUnavailable,
		codes.Internal:         notary.KindSourceUnavailable,
	} {
		t.Run(code.String(), func(t *testing.T) {
			server, records := recordsRig(t, "notary-watch")
			server.RecordErr = status.Error(code, "scripted")

			_, err := records.Record(context.Background(), notarytest.SnapshotID(t, 15))

			require.Equal(t, kind, notary.KindOf(err))
		})
	}
}
