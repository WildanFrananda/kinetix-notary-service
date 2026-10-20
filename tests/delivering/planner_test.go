package delivering_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const wib = 7 * time.Hour

var policy = delivering.Policy{
	UTCOffset:     wib,
	EpochsPerTick: 50,
	Batch:         20,
	FirstRetry:    time.Minute,
	LastRetry:     6 * time.Hour,
	OverdueAfter:  6 * time.Hour,
}

func fixed() time.Time {
	return time.Date(2026, 10, 21, 0, 45, 0, 0, time.UTC)
}

func quiet() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func epochAt(t *testing.T, number uint64, at time.Time) notary.EpochRecord {
	t.Helper()

	asOf, err := notary.UnixTimeOf(at)
	require.NoError(t, err)

	return notary.EpochRecord{Header: notary.EpochHeader{Number: number, AsOf: asOf}, State: notary.StateFinalized}
}

func holding(t *testing.T, holder string, kind notary.AccountKind, amount uint64) notary.Account {
	t.Helper()

	return notary.Account{Liability: notary.Liability{Holder: notarytest.Holder(t, holder), Kind: kind, Amount: amount}}
}

func day(year int, month time.Month, d int) time.Time {
	return time.Date(year, month, d, 0, 0, 0, 0, time.UTC)
}

func rowsOf(store *notarytest.StatementStore, state string) map[string]uint64 {
	out := map[string]uint64{}

	for _, row := range store.Snapshot() {
		if row.State == state {
			out[row.Holder.String()] = row.Epoch
		}
	}

	return out
}

func TestTheFirstPlanStartsAtTheNewestFinalEpochAndLeavesSuspenseOut(t *testing.T) {
	store := &notarytest.StatementStore{
		Finalized: []notary.EpochRecord{
			epochAt(t, 180, time.Date(2026, 10, 20, 10, 27, 0, 0, time.UTC)),
			epochAt(t, 181, time.Date(2026, 10, 20, 11, 27, 0, 0, time.UTC)),
			epochAt(t, 182, time.Date(2026, 10, 20, 12, 27, 0, 0, time.UTC)),
		},
		Holdings: map[uint64][]notary.Account{
			180: {holding(t, "principal-z", notary.AccountCustomerBalance, 1)},
			182: {
				holding(t, "principal-a", notary.AccountCustomerBalance, 100),
				holding(t, "principal-a", notary.AccountDriverAvailable, 50),
				holding(t, "principal-b", notary.AccountMerchantAvailable, 10),
				holding(t, "refunds", notary.AccountSuspenseAvailable, 999),
			},
		},
	}

	require.NoError(t, delivering.NewPlanner(store, fixed, quiet(), policy).Plan(context.Background()))

	require.Equal(t, map[string]uint64{"principal-a": 182, "principal-b": 182}, rowsOf(store, "pending"))
	require.Equal(t, notary.StatementProgress{Epoch: 182, Day: day(2026, 10, 20)}, *store.Progress)
}

func TestAHolderIsSentAtMostOneStatementADay(t *testing.T) {
	store := &notarytest.StatementStore{
		Finalized: []notary.EpochRecord{epochAt(t, 182, time.Date(2026, 10, 20, 12, 27, 0, 0, time.UTC))},
		Holdings: map[uint64][]notary.Account{
			182: {holding(t, "principal-a", notary.AccountCustomerBalance, 100)},
			183: {holding(t, "principal-a", notary.AccountCustomerBalance, 250)},
		},
	}
	planner := delivering.NewPlanner(store, fixed, quiet(), policy)
	require.NoError(t, planner.Plan(context.Background()))

	store.Finalized = append(store.Finalized, epochAt(t, 183, time.Date(2026, 10, 20, 16, 27, 0, 0, time.UTC)))
	require.NoError(t, planner.Plan(context.Background()))

	require.Len(t, store.Snapshot(), 1)
	require.Equal(t, notary.StatementProgress{Epoch: 183, Day: day(2026, 10, 20)}, *store.Progress)
}

func TestTheFirstEpochOfADayCarriesOnlyTheHoldersWhoseBalanceChanged(t *testing.T) {
	store := &notarytest.StatementStore{
		Finalized: []notary.EpochRecord{epochAt(t, 182, time.Date(2026, 10, 20, 12, 27, 0, 0, time.UTC))},
		Holdings: map[uint64][]notary.Account{
			182: {
				holding(t, "principal-a", notary.AccountCustomerBalance, 100),
				holding(t, "principal-b", notary.AccountMerchantAvailable, 10),
			},
			187: {
				holding(t, "principal-a", notary.AccountCustomerBalance, 100),
				holding(t, "principal-b", notary.AccountMerchantAvailable, 11),
				holding(t, "principal-c", notary.AccountDriverAvailable, 5),
			},
			188: {holding(t, "principal-a", notary.AccountCustomerBalance, 999)},
		},
	}
	planner := delivering.NewPlanner(store, fixed, quiet(), policy)
	require.NoError(t, planner.Plan(context.Background()))

	store.Finalized = append(
		store.Finalized,
		epochAt(t, 187, time.Date(2026, 10, 20, 17, 27, 0, 0, time.UTC)),
		epochAt(t, 188, time.Date(2026, 10, 20, 18, 27, 0, 0, time.UTC)),
	)
	require.NoError(t, planner.Plan(context.Background()))

	require.Equal(
		t,
		map[string]uint64{"principal-a": 182, "principal-b": 187, "principal-c": 187},
		rowsOf(store, "pending"),
	)
	require.Equal(t, map[string]uint64{"principal-b": 182}, rowsOf(store, "superseded"))
	require.Equal(t, notary.StatementProgress{Epoch: 188, Day: day(2026, 10, 21)}, *store.Progress)
}

func TestTheDayTurnsAtMidnightInTheConfiguredOffset(t *testing.T) {
	store := &notarytest.StatementStore{
		Finalized: []notary.EpochRecord{epochAt(t, 10, time.Date(2026, 10, 20, 16, 59, 59, 0, time.UTC))},
		Holdings: map[uint64][]notary.Account{
			10: {holding(t, "principal-a", notary.AccountCustomerBalance, 1)},
			11: {holding(t, "principal-a", notary.AccountCustomerBalance, 2)},
		},
	}
	planner := delivering.NewPlanner(store, fixed, quiet(), policy)
	require.NoError(t, planner.Plan(context.Background()))
	require.Equal(t, day(2026, 10, 20), store.Progress.Day)

	store.Finalized = append(store.Finalized, epochAt(t, 11, time.Date(2026, 10, 20, 17, 0, 0, 0, time.UTC)))
	require.NoError(t, planner.Plan(context.Background()))

	require.Equal(t, day(2026, 10, 21), store.Progress.Day)
	require.Equal(t, map[string]uint64{"principal-a": 11}, rowsOf(store, "pending"))
}

func TestAHolderWithNowhereToSendIsAskedAgainOnlyOnceTheirBalanceChanges(t *testing.T) {
	store := &notarytest.StatementStore{
		Finalized: []notary.EpochRecord{epochAt(t, 1, time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC))},
		Holdings: map[uint64][]notary.Account{
			1: {holding(t, "principal-a", notary.AccountCustomerBalance, 100)},
			2: {holding(t, "principal-a", notary.AccountCustomerBalance, 100)},
			3: {holding(t, "principal-a", notary.AccountCustomerBalance, 101)},
		},
	}
	planner := delivering.NewPlanner(store, fixed, quiet(), policy)
	require.NoError(t, planner.Plan(context.Background()))
	require.NoError(
		t,
		store.StatementRefused(
			context.Background(),
			notary.StatementDelivery{Holder: notarytest.Holder(t, "principal-a"), Epoch: 1},
			"RECIPIENT_UNKNOWN",
		),
	)

	store.Finalized = append(store.Finalized, epochAt(t, 2, time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)))
	require.NoError(t, planner.Plan(context.Background()))
	require.Empty(t, rowsOf(store, "pending"))

	store.Finalized = append(store.Finalized, epochAt(t, 3, time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)))
	require.NoError(t, planner.Plan(context.Background()))
	require.Equal(t, map[string]uint64{"principal-a": 3}, rowsOf(store, "pending"))
}

func TestAPassStopsWhenAnotherInstanceAdvancedFirst(t *testing.T) {
	store := &notarytest.StatementStore{
		Finalized: []notary.EpochRecord{epochAt(t, 1, time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC))},
		Holdings:  map[uint64][]notary.Account{1: {holding(t, "principal-a", notary.AccountCustomerBalance, 1)}},
		Raced:     true,
	}

	require.NoError(t, delivering.NewPlanner(store, fixed, quiet(), policy).Plan(context.Background()))
	require.Empty(t, store.Snapshot())
	require.Nil(t, store.Progress)
}

func TestNothingIsPlannedBeforeAnyEpochIsFinal(t *testing.T) {
	store := &notarytest.StatementStore{}

	require.NoError(t, delivering.NewPlanner(store, fixed, quiet(), policy).Plan(context.Background()))
	require.Nil(t, store.Progress)
}
