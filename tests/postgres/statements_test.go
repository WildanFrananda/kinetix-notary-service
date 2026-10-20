package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func (d *db) finalized(t *testing.T, epochs ...uint64) {
	t.Helper()

	ctx := context.Background()
	lease := d.lease(t, "writer")
	previous := notary.Digest{}

	for _, number := range epochs {
		p := plan(t, number, previous, 3)
		_, err := d.store.SaveBuilt(ctx, lease, ledgerID, p)
		require.NoError(t, err)
		_, err = d.pool.Exec(ctx, "UPDATE epochs SET state = 'finalized' WHERE number = $1", int64(number))
		require.NoError(t, err)

		previous = p.Digest
	}
}

func (d *db) statement(t *testing.T, holder string, epoch uint64) (string, int32, string) {
	t.Helper()

	var (
		state        string
		attempts     int32
		notification *string
	)

	row := d.pool.QueryRow(
		context.Background(),
		"SELECT state, attempts, notification_id FROM statement_deliveries WHERE holder = $1 AND epoch_number = $2",
		holder,
		int64(epoch),
	)
	require.NoError(t, row.Scan(&state, &attempts, &notification))

	if notification == nil {
		return state, attempts, ""
	}

	return state, attempts, *notification
}

func planned(t *testing.T, fingerprint byte, holders ...string) []notary.PlannedStatement {
	t.Helper()

	out := make([]notary.PlannedStatement, len(holders))
	for i, holder := range holders {
		out[i] = notary.PlannedStatement{Holder: notarytest.Holder(t, holder), Fingerprint: notary.Digest{fingerprint}}
	}

	return out
}

func TestStatementProgressOnlyMovesForward(t *testing.T) {
	d := open(t)
	ctx := context.Background()

	d.finalized(t, 0, 1, 2)

	_, started, err := d.store.StatementProgress(ctx)
	require.NoError(t, err)
	require.False(t, started)

	day := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	advanced, err := d.store.AdvanceStatements(
		ctx,
		notary.StatementProgress{Epoch: 1, Day: day},
		planned(t, 1, "a"),
		time.Now(),
	)
	require.NoError(t, err)
	require.True(t, advanced)

	advanced, err = d.store.AdvanceStatements(
		ctx,
		notary.StatementProgress{Epoch: 1, Day: day},
		planned(t, 2, "b"),
		time.Now(),
	)
	require.NoError(t, err)
	require.False(t, advanced, "a second pass over the same epoch changes nothing")
	require.Equal(t, int64(1), d.count(t, "SELECT count(*) FROM statement_deliveries"))

	progress, started, err := d.store.StatementProgress(ctx)
	require.NoError(t, err)
	require.True(t, started)
	require.Equal(t, notary.StatementProgress{Epoch: 1, Day: day}, progress)

	after, err := d.store.FinalizedAfter(ctx, progress.Epoch, 10)
	require.NoError(t, err)
	require.Len(t, after, 1)
	require.Equal(t, uint64(2), after[0].Header.Number)
}

func TestANewerStatementSupersedesAnUnsentOneAndOnlyLiveOnesAreCompared(t *testing.T) {
	d := open(t)
	ctx := context.Background()

	d.finalized(t, 0, 1, 2)

	first := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	_, err := d.store.AdvanceStatements(
		ctx,
		notary.StatementProgress{Epoch: 1, Day: first},
		planned(t, 1, "a", "b"),
		time.Now(),
	)
	require.NoError(t, err)

	holderA := notary.StatementDelivery{Holder: notarytest.Holder(t, "a"), Epoch: 1}
	require.NoError(t, d.store.StatementSent(ctx, holderA, "NTF-1"))

	_, err = d.store.AdvanceStatements(
		ctx,
		notary.StatementProgress{Epoch: 2, Day: first.AddDate(0, 0, 1)},
		planned(t, 2, "b"),
		time.Now(),
	)
	require.NoError(t, err)

	state, _, _ := d.statement(t, "b", 1)
	require.Equal(t, "superseded", state)

	state, attempts, notification := d.statement(t, "a", 1)
	require.Equal(t, "sent", state)
	require.Equal(t, int32(1), attempts)
	require.Equal(t, "NTF-1", notification)

	fingerprints, err := d.store.StatementFingerprints(
		ctx,
		[]notary.Holder{notarytest.Holder(t, "a"), notarytest.Holder(t, "b"), notarytest.Holder(t, "c")},
	)
	require.NoError(t, err)
	require.Equal(
		t,
		map[notary.Holder]notary.Digest{
			notarytest.Holder(t, "a"): {1},
			notarytest.Holder(t, "b"): {2},
		},
		fingerprints,
	)
}

func TestDueStatementsWaitForTheirNextAttemptAndCountTowardsTheBacklog(t *testing.T) {
	d := open(t)
	ctx := context.Background()

	d.finalized(t, 0, 1)

	day := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	_, err := d.store.AdvanceStatements(
		ctx,
		notary.StatementProgress{Epoch: 1, Day: day},
		planned(t, 1, "a", "b", "c"),
		time.Now(),
	)
	require.NoError(t, err)

	now := time.Now()
	due, err := d.store.DueStatements(ctx, now, 10)
	require.NoError(t, err)
	require.Len(t, due, 3)

	require.NoError(t, d.store.StatementPostponed(ctx, due[0], now.Add(time.Hour)))
	require.NoError(t, d.store.StatementRefused(ctx, due[1], "NO_CHANNEL_AVAILABLE: no address"))

	again, err := d.store.DueStatements(ctx, now, 10)
	require.NoError(t, err)
	require.Equal(t, []notary.StatementDelivery{due[2]}, again)

	later, err := d.store.DueStatements(ctx, now.Add(2*time.Hour), 10)
	require.NoError(t, err)
	require.Len(t, later, 2)
	require.Contains(t, later, notary.StatementDelivery{Holder: due[0].Holder, Epoch: 1, Attempts: 1})

	backlog, err := d.store.StatementBacklog(ctx, now.Add(-time.Hour))
	require.NoError(t, err)
	require.Equal(t, notary.StatementBacklog{Pending: 2, Overdue: 0}, backlog)

	backlog, err = d.store.StatementBacklog(ctx, now.Add(time.Hour))
	require.NoError(t, err)
	require.Equal(t, notary.StatementBacklog{Pending: 2, Overdue: 2}, backlog)
}

func TestAHolderGetsOneStatementADayEvenIfPlannedTwice(t *testing.T) {
	d := open(t)
	ctx := context.Background()

	d.finalized(t, 0, 1, 2)

	day := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	_, err := d.store.AdvanceStatements(
		ctx,
		notary.StatementProgress{Epoch: 1, Day: day},
		planned(t, 1, "a"),
		time.Now(),
	)
	require.NoError(t, err)
	_, err = d.store.AdvanceStatements(
		ctx,
		notary.StatementProgress{Epoch: 2, Day: day},
		planned(t, 2, "a"),
		time.Now(),
	)
	require.NoError(t, err)

	require.Equal(
		t,
		int64(1),
		d.count(t, "SELECT count(*) FROM statement_deliveries WHERE holder = 'a' AND state = 'pending'"),
	)
}
