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
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const (
	snapshot = "5d0c7a52-8e1b-4f6a-9c33-0b7e2d41a9f8"
	customer = "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7"
	merchant = "3c9a77b1-58de-4a01-8f2e-6d4b19c0a8f3"
	purpose  = "UNASSIGNED_DRIVER_SHIPPING_FEE"
)

type rig struct {
	pki    *notarytest.PKI
	server *notarytest.LiabilityServer
	source *payment.Source
}

func newRig(t *testing.T, serverDomain, serverName string, parts ...*paymentv1.StreamLiabilitySnapshotResponse) *rig {
	t.Helper()
	pki := notarytest.NewPKI(t)
	server := &notarytest.LiabilityServer{Parts: parts}
	address := notarytest.ServeLiabilities(t, pki.Issue(t, serverDomain, serverName), server)
	source, err := payment.Dial(payment.Settings{
		Endpoint: address,
		Identity: mesh.Identity{
			PKIDir:       pki.Issue(t, "kinetix.local", "notary"),
			TrustDomains: []string{"kinetix.local"},
		},
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = source.Close() })

	return &rig{pki: pki, server: server, source: source}
}

func header(seconds int64, nanos int32, currency string) *paymentv1.StreamLiabilitySnapshotResponse {
	return &paymentv1.StreamLiabilitySnapshotResponse{Part: &paymentv1.StreamLiabilitySnapshotResponse_Header{
		Header: &paymentv1.LiabilitySnapshotHeader{
			AsOf:       &timestamppb.Timestamp{Seconds: seconds, Nanos: nanos},
			Currency:   currency,
			SnapshotId: snapshot,
		},
	}}
}

func batch(accounts ...*paymentv1.LiabilityAccount) *paymentv1.StreamLiabilitySnapshotResponse {
	return &paymentv1.StreamLiabilitySnapshotResponse{Part: &paymentv1.StreamLiabilitySnapshotResponse_Accounts{
		Accounts: &paymentv1.LiabilityAccountBatch{Accounts: accounts},
	}}
}

func trailer(count, total int64, currency string) *paymentv1.StreamLiabilitySnapshotResponse {
	return &paymentv1.StreamLiabilitySnapshotResponse{Part: &paymentv1.StreamLiabilitySnapshotResponse_Trailer{
		Trailer: &paymentv1.LiabilitySnapshotTrailer{AccountCount: count, Total: money(total, currency)},
	}}
}

func money(amount int64, currency string) *commonv1.Money {
	return &commonv1.Money{AmountMinor: amount, Currency: currency}
}

func owedTo(principal string, kind paymentv1.LiabilityKind, amount int64) *paymentv1.LiabilityAccount {
	return &paymentv1.LiabilityAccount{
		Kind:    kind,
		Holder:  &paymentv1.LiabilityAccount_PrincipalId{PrincipalId: principal},
		Balance: money(amount, "IDR"),
	}
}

func heldFor(purpose string, kind paymentv1.LiabilityKind, amount int64) *paymentv1.LiabilityAccount {
	return &paymentv1.LiabilityAccount{
		Kind:    kind,
		Holder:  &paymentv1.LiabilityAccount_SuspensePurpose{SuspensePurpose: purpose},
		Balance: money(amount, "IDR"),
	}
}

func TestACompleteStreamBecomesTheSnapshotItDescribes(t *testing.T) {
	parts, kept := notarytest.SnapshotStream(
		t,
		notarytest.SnapshotID(t, 7),
		&timestamppb.Timestamp{Seconds: 1760000000, Nanos: 987_000_000},
		"IDR",
		[]*paymentv1.LiabilityAccount{
			owedTo(customer, paymentv1.LiabilityKind_LIABILITY_KIND_CUSTOMER_BALANCE, 35_000_000),
			owedTo(merchant, paymentv1.LiabilityKind_LIABILITY_KIND_MERCHANT_AVAILABLE, 0),
			owedTo(merchant, paymentv1.LiabilityKind_LIABILITY_KIND_MERCHANT_PENDING_ESCROW, 13_000_000),
			heldFor(purpose, paymentv1.LiabilityKind_LIABILITY_KIND_SUSPENSE_PENDING_ESCROW, 2_000_000),
		},
	)
	r := newRig(t, "kinetix.local", "payment", parts...)

	snapshot, err := r.source.Snapshot(context.Background(), "IDR")

	require.NoError(t, err)
	require.Equal(t, uint64(1760000000), snapshot.AsOf().Seconds())
	require.Equal(t, uint64(50_000_000), snapshot.Total())
	require.Equal(t, "IDR", snapshot.Currency())
	liabilities := snapshot.Liabilities()
	require.Len(t, liabilities, 4)
	require.Equal(
		t,
		notary.Liability{
			Holder: notarytest.Holder(t, customer),
			Kind:   notary.AccountCustomerBalance,
			Amount: 35_000_000,
		},
		liabilities[0],
	)
	require.Equal(
		t,
		notary.Liability{
			Holder: notarytest.Holder(t, purpose),
			Kind:   notary.AccountSuspensePendingEscrow,
			Amount: 2_000_000,
		},
		liabilities[3],
	)
	require.Equal(t, []string{"notary"}, r.server.Callers)
	require.Equal(t, []string{"IDR"}, r.server.Requests)

	record := snapshot.Record()
	require.Equal(t, notarytest.SnapshotID(t, 7), record.ID)
	require.Equal(t, uint32(987_000_000), record.Nanos)
	require.Equal(t, uint64(4), record.Count)
	require.Equal(t, uint64(50_000_000), record.Total)
	require.Equal(t, kept.GetDigest(), record.Digest[:])
}

func TestAnEmptyLedgerIsAHeaderAndATrailer(t *testing.T) {
	parts, _ := notarytest.SnapshotStream(
		t,
		notarytest.SnapshotID(t, 8),
		&timestamppb.Timestamp{Seconds: 1760000000},
		"IDR",
		nil,
	)
	r := newRig(t, "kinetix.local", "payment", parts...)

	snapshot, err := r.source.Snapshot(context.Background(), "IDR")

	require.NoError(t, err)
	require.Empty(t, snapshot.Liabilities())
	require.Zero(t, snapshot.Total())
}

func TestEveryBrokenStreamIsRefusedNeverAnchored(t *testing.T) {
	good := owedTo(customer, paymentv1.LiabilityKind_LIABILITY_KIND_CUSTOMER_BALANCE, 10)

	cases := map[string][]*paymentv1.StreamLiabilitySnapshotResponse{
		"no trailer":              {header(1760000000, 0, "IDR"), batch(good)},
		"no header":               {batch(good), trailer(1, 10, "IDR")},
		"trailer first":           {trailer(0, 0, "IDR")},
		"two headers":             {header(1760000000, 0, "IDR"), header(1760000001, 0, "IDR"), trailer(0, 0, "IDR")},
		"after the trailer":       {header(1760000000, 0, "IDR"), trailer(0, 0, "IDR"), batch(good)},
		"short count":             {header(1760000000, 0, "IDR"), batch(good, good), trailer(1, 20, "IDR")},
		"long count":              {header(1760000000, 0, "IDR"), batch(good), trailer(2, 10, "IDR")},
		"wrong total":             {header(1760000000, 0, "IDR"), batch(good), trailer(1, 11, "IDR")},
		"total in another money":  {header(1760000000, 0, "IDR"), batch(good), trailer(1, 10, "USD")},
		"header in another money": {header(1760000000, 0, "USD"), trailer(0, 0, "USD")},
		"no time":                 {header(0, 0, "IDR"), trailer(0, 0, "IDR")},
		"same account twice":      {header(1760000000, 0, "IDR"), batch(good, good), trailer(2, 20, "IDR")},
		"negative balance": {
			header(1760000000, 0, "IDR"),
			batch(owedTo(customer, 1, -1)),
			trailer(1, -1, "IDR"),
		},
		"balance in another money": {
			header(1760000000, 0, "IDR"),
			batch(&paymentv1.LiabilityAccount{
				Kind:    1,
				Holder:  &paymentv1.LiabilityAccount_PrincipalId{PrincipalId: customer},
				Balance: money(10, "USD"),
			}),
			trailer(1, 10, "IDR"),
		},
		"unspecified kind":     {header(1760000000, 0, "IDR"), batch(owedTo(customer, 0, 10)), trailer(1, 10, "IDR")},
		"unknown kind":         {header(1760000000, 0, "IDR"), batch(owedTo(customer, 8, 10)), trailer(1, 10, "IDR")},
		"suspense as a person": {header(1760000000, 0, "IDR"), batch(owedTo(purpose, 6, 10)), trailer(1, 10, "IDR")},
		"person as suspense":   {header(1760000000, 0, "IDR"), batch(heldFor(customer, 1, 10)), trailer(1, 10, "IDR")},
		"nobody": {
			header(1760000000, 0, "IDR"),
			batch(&paymentv1.LiabilityAccount{Kind: 1, Balance: money(10, "IDR")}),
			trailer(1, 10, "IDR"),
		},
		"empty part": {header(1760000000, 0, "IDR"), {}, trailer(0, 0, "IDR")},
	}
	for name, parts := range cases {
		t.Run(name, func(t *testing.T) {
			r := newRig(t, "kinetix.local", "payment", parts...)

			_, err := r.source.Snapshot(context.Background(), "IDR")

			require.Error(t, err)
			require.Equal(t, notary.KindSourceRefused, notary.KindOf(err), err.Error())
		})
	}
}

func TestPaymentsOwnStatusDecidesRetryOrRefusal(t *testing.T) {
	cases := map[codes.Code]notary.Kind{
		codes.FailedPrecondition: notary.KindSourceRefused,
		codes.InvalidArgument:    notary.KindSourceRefused,
		codes.PermissionDenied:   notary.KindSourceRefused,
		codes.Unavailable:        notary.KindSourceUnavailable,
		codes.DeadlineExceeded:   notary.KindSourceUnavailable,
		codes.Internal:           notary.KindSourceUnavailable,
	}
	for code, kind := range cases {
		t.Run(code.String(), func(t *testing.T) {
			r := newRig(t, "kinetix.local", "payment", header(1760000000, 0, "IDR"))
			r.server.Err = status.Error(code, "scripted")

			_, err := r.source.Snapshot(context.Background(), "IDR")

			require.Equal(t, kind, notary.KindOf(err))
		})
	}
}

func TestOnlyPaymentMayAnswer(t *testing.T) {
	for name, server := range map[string][2]string{
		"another service":      {"kinetix.local", "catalog"},
		"another trust domain": {"prod.kinetix", "payment"},
	} {
		t.Run(name, func(t *testing.T) {
			r := newRig(t, server[0], server[1], header(1760000000, 0, "IDR"), trailer(0, 0, "IDR"))

			_, err := r.source.Snapshot(context.Background(), "IDR")

			require.Error(t, err)
			require.Empty(t, r.server.Callers)
		})
	}
}

func TestAStreamWithoutItsRecordIsRefused(t *testing.T) {
	good := owedTo(customer, paymentv1.LiabilityKind_LIABILITY_KIND_CUSTOMER_BALANCE, 10)
	parts, _ := notarytest.SnapshotStream(
		t,
		notarytest.SnapshotID(t, 9),
		&timestamppb.Timestamp{Seconds: 1760000000, Nanos: 5},
		"IDR",
		[]*paymentv1.LiabilityAccount{good},
	)

	unnamed := clone(t, parts[0])
	unnamed.GetHeader().SnapshotId = ""
	r := newRig(t, "kinetix.local", "payment", unnamed, parts[1], parts[2])
	_, err := r.source.Snapshot(context.Background(), "IDR")
	require.Equal(t, notary.KindSourceRefused, notary.KindOf(err))
	require.ErrorContains(t, err, "predates the snapshot record")

	undigested := clone(t, parts[2])
	undigested.GetTrailer().Digest = nil
	r = newRig(t, "kinetix.local", "payment", parts[0], parts[1], undigested)
	_, err = r.source.Snapshot(context.Background(), "IDR")
	require.Equal(t, notary.KindSourceRefused, notary.KindOf(err))

	misdigested := clone(t, parts[2])
	misdigested.GetTrailer().Digest[0] ^= 0xff
	r = newRig(t, "kinetix.local", "payment", parts[0], parts[1], misdigested)
	_, err = r.source.Snapshot(context.Background(), "IDR")
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))

	renanoed := clone(t, parts[0])
	renanoed.GetHeader().AsOf.Nanos = 6
	r = newRig(t, "kinetix.local", "payment", renanoed, parts[1], parts[2])
	_, err = r.source.Snapshot(context.Background(), "IDR")
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func clone(t *testing.T, part *paymentv1.StreamLiabilitySnapshotResponse) *paymentv1.StreamLiabilitySnapshotResponse {
	t.Helper()

	copied, ok := proto.Clone(part).(*paymentv1.StreamLiabilitySnapshotResponse)
	require.True(t, ok)

	return copied
}
