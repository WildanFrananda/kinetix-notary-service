package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func withholdingRequest(t *testing.T, tag byte, epoch, dueAt uint64) notary.WithholdingRequest {
	t.Helper()

	return notary.WithholdingRequest{
		Registration: notary.Digest{0xdd, tag},
		Account:      notary.Digest{0xac, tag},
		Epoch:        epoch,
		Holder:       notary.Address{0x40, tag},
		DueAt:        notarytest.Unix(t, dueAt),
	}
}

func TestAWithholdingScanResumesWhereItWasRecordedAndNeverMovesBack(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	court := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xa3}}
	other := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xe0}}

	_, found, err := d.store.ScanFrom(ctx, court)
	require.NoError(t, err)
	require.False(t, found)

	require.NoError(t, d.store.RecordScan(ctx, court, 48_365_427, nil))
	require.NoError(t, d.store.RecordScan(ctx, court, 48_365_227, nil))

	next, found, err := d.store.ScanFrom(ctx, court)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint64(48_365_427), next)

	_, found, err = d.store.ScanFrom(ctx, other)
	require.NoError(t, err)
	require.False(t, found)
}

func TestAnOpenWithholdingRequestIsKeptOnceUntilItIsClosed(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	court := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xa3}}
	other := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xe0}}
	later := withholdingRequest(t, 1, 186, 1_761_000_000)
	sooner := withholdingRequest(t, 2, 187, 1_760_900_000)

	require.NoError(t, d.store.RecordScan(ctx, court, 100, []notary.WithholdingRequest{later}))
	require.NoError(t, d.store.RecordScan(ctx, court, 300, []notary.WithholdingRequest{later, sooner}))
	require.NoError(t, d.store.RecordScan(ctx, other, 50, nil))
	require.Equal(t, int64(2), d.count(t, "SELECT count(*) FROM withholding_requests"))

	open, err := d.store.OpenRequests(ctx, court)
	require.NoError(t, err)
	require.Equal(t, []notary.WithholdingRequest{sooner, later}, open)

	open, err = d.store.OpenRequests(ctx, other)
	require.NoError(t, err)
	require.Empty(t, open)

	require.NoError(t, d.store.CloseRequest(ctx, court, sooner.Registration, sooner.Epoch))
	require.NoError(t, d.store.CloseRequest(ctx, court, sooner.Registration, sooner.Epoch))

	open, err = d.store.OpenRequests(ctx, court)
	require.NoError(t, err)
	require.Equal(t, []notary.WithholdingRequest{later}, open)
}
