package notary_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const (
	vectorSnapshot = "5d0c7a52-8e1b-4f6a-9c33-0b7e2d41a9f8"
	vectorDigest   = "0x89276c0d19d674863ff54989337008621a5774b1a5226f6fbd5029b30bb0bd8f"
	emptyDigest    = "0x16956e78b53678e1767f5d74d09d4ba7884ed4cc1109f6ae47f9113baa1cde49"
)

func vectorLiabilities(t *testing.T) []notary.Liability {
	t.Helper()

	rows := []struct {
		kind   notary.AccountKind
		holder string
		amount uint64
	}{
		{notary.AccountCustomerBalance, "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7", 50_000_000},
		{notary.AccountMerchantAvailable, "3c9a77b1-58de-4a01-8f2e-6d4b19c0a8f3", 0},
		{notary.AccountMerchantPendingEscrow, "3c9a77b1-58de-4a01-8f2e-6d4b19c0a8f3", 13_000_000},
		{notary.AccountDriverAvailable, "b7e2c05f-9a34-4c88-b1d6-0e7a3f52d914", 0},
		{notary.AccountDriverPendingEscrow, "b7e2c05f-9a34-4c88-b1d6-0e7a3f52d914", 2_000_000},
		{notary.AccountSuspenseAvailable, "UNASSIGNED_DRIVER_SHIPPING_FEE", 0},
		{notary.AccountSuspensePendingEscrow, "UNASSIGNED_DRIVER_SHIPPING_FEE", 0},
		{notary.AccountCustomerBalance, "1b2c3d4e-0000-4000-8000-000000000001", 1},
	}

	liabilities := make([]notary.Liability, len(rows))
	for i, row := range rows {
		liabilities[i] = notary.Liability{Holder: notarytest.Holder(t, row.holder), Kind: row.kind, Amount: row.amount}
	}

	return liabilities
}

func TestTheSnapshotDigestIsTheOnePaymentAndAnIndependentImplementationGive(t *testing.T) {
	record := notary.SnapshotRecord{
		AsOf:     notarytest.Unix(t, 1791795600),
		Nanos:    123456000,
		Currency: "IDR",
		Count:    8,
		Total:    65_000_001,
	}
	liabilities := vectorLiabilities(t)

	digest, err := notary.SnapshotDigest(record, liabilities)
	require.NoError(t, err)
	require.Equal(t, vectorDigest, digest.Hex())

	reversed := make([]notary.Liability, len(liabilities))
	for i := range liabilities {
		reversed[i] = liabilities[len(liabilities)-1-i]
	}

	digest, err = notary.SnapshotDigest(record, reversed)
	require.NoError(t, err)
	require.Equal(t, vectorDigest, digest.Hex())

	empty := notary.SnapshotRecord{AsOf: notarytest.Unix(t, 1791795600), Currency: "IDR"}
	digest, err = notary.SnapshotDigest(empty, nil)
	require.NoError(t, err)
	require.Equal(t, emptyDigest, digest.Hex())
}

func TestASnapshotIsRecordedOnlyUnderTheDigestItsRowsGive(t *testing.T) {
	snapshot, err := notary.NewSnapshot(notarytest.Unix(t, 1791795600), "IDR", vectorLiabilities(t))
	require.NoError(t, err)

	id, err := notary.NewSnapshotID(vectorSnapshot)
	require.NoError(t, err)

	digest, err := notary.ParseDigest(vectorDigest)
	require.NoError(t, err)

	recorded, err := snapshot.Recorded(id, 123456000, digest)
	require.NoError(t, err)
	require.Equal(
		t,
		notary.SnapshotRecord{
			ID:       id,
			AsOf:     notarytest.Unix(t, 1791795600),
			Nanos:    123456000,
			Currency: "IDR",
			Count:    8,
			Total:    65_000_001,
			Digest:   digest,
		},
		recorded.Record(),
	)

	_, err = snapshot.Recorded(id, 123456001, digest)
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))

	_, err = snapshot.Recorded("", 123456000, digest)
	require.Equal(t, notary.KindSourceRefused, notary.KindOf(err))

	_, err = snapshot.Recorded(id, 1_000_000_000, digest)
	require.Equal(t, notary.KindSourceRefused, notary.KindOf(err))
}

func TestASnapshotIDIsALowercaseCanonicalUUID(t *testing.T) {
	_, err := notary.NewSnapshotID(vectorSnapshot)
	require.NoError(t, err)

	for _, malformed := range []string{
		"",
		"5D0C7A52-8E1B-4F6A-9C33-0B7E2D41A9F8",
		"5d0c7a528e1b4f6a9c330b7e2d41a9f8",
		"5d0c7a52-8e1b-4f6a-9c33-0b7e2d41a9f",
		"5d0c7a52_8e1b_4f6a_9c33_0b7e2d41a9f8",
		"5d0c7a52-8e1b-4f6a-9c33-0b7e2d41a9fg",
	} {
		_, err := notary.NewSnapshotID(malformed)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), malformed)
	}
}
