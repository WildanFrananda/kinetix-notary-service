package postgres_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/postgres"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const (
	databaseURLKey = "KINETIX_NOTARY_TEST_DATABASE_URL"
	schemaVersion  = int64(3)
)

var ledgerID = notary.LedgerIdentity{ChainID: 84532, Address: notary.Address{0x4e, 0x86}}

type db struct {
	pool  *pgxpool.Pool
	store *postgres.Store
}

func open(t *testing.T) *db {
	t.Helper()

	url := os.Getenv(databaseURLKey)
	if url == "" {
		t.Skip(databaseURLKey + " is not set; these tests need Postgres, and CI must set it")
	}

	ctx := context.Background()
	suffix := make([]byte, 6)
	_, err := rand.Read(suffix)
	require.NoError(t, err)

	schema := "t_" + hex.EncodeToString(suffix)
	admin, err := pgx.Connect(ctx, url)
	require.NoError(t, err)
	_, err = admin.Exec(ctx, "CREATE SCHEMA "+schema)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE")
		_ = admin.Close(context.Background())
	})

	config, err := pgxpool.ParseConfig(url)
	require.NoError(t, err)

	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	before, after, err := postgres.Migrate(ctx, pool)
	require.NoError(t, err)
	require.Equal(t, int64(0), before)
	require.Equal(t, schemaVersion, after)

	return &db{pool: pool, store: postgres.NewStore(pool, keccak.NewScheme())}
}

func (d *db) lease(t *testing.T, instance string) notary.Lease {
	t.Helper()

	lease, held, err := d.store.Lease(context.Background(), instance, time.Minute)
	require.NoError(t, err)
	require.True(t, held)

	return lease
}

func plan(t *testing.T, number uint64, previous notary.Digest, accounts int) notary.EpochPlan {
	t.Helper()
	snapshot := notarytest.Snapshot(t, 1760000000+number, notarytest.Liabilities(t, accounts))
	p, err := notary.PlanEpoch(
		keccak.NewScheme(),
		&notarytest.Salts{Seed: number},
		&notarytest.Secrets{},
		nil,
		number,
		previous,
		snapshot,
	)
	require.NoError(t, err)

	return p
}

func (d *db) count(t *testing.T, query string, args ...any) int64 {
	t.Helper()

	var n int64
	require.NoError(t, d.pool.QueryRow(context.Background(), query, args...).Scan(&n))

	return n
}

func TestMigratingTwiceChangesNothing(t *testing.T) {
	d := open(t)

	before, after, err := postgres.Migrate(context.Background(), d.pool)

	require.NoError(t, err)
	require.Equal(t, schemaVersion, before)
	require.Equal(t, schemaVersion, after)
}

func TestOneWriterHoldsTheLeaseUntilItLapses(t *testing.T) {
	d := open(t)
	ctx := context.Background()

	first, held, err := d.store.Lease(ctx, "a", 300*time.Millisecond)
	require.NoError(t, err)
	require.True(t, held)
	renewed, held, err := d.store.Lease(ctx, "a", 300*time.Millisecond)
	require.NoError(t, err)
	require.True(t, held)
	require.Equal(t, first.Token, renewed.Token)

	_, held, err = d.store.Lease(ctx, "b", time.Minute)
	require.NoError(t, err)
	require.False(t, held)

	time.Sleep(400 * time.Millisecond)

	taken, held, err := d.store.Lease(ctx, "b", time.Minute)
	require.NoError(t, err)
	require.True(t, held)
	require.Greater(t, taken.Token, first.Token)

	_, err = d.store.SaveBuilt(ctx, first, ledgerID, plan(t, 0, notary.Digest{}, 3))
	require.Equal(t, notary.KindLeaseLost, notary.KindOf(err))
	require.Zero(t, d.count(t, "SELECT count(*) FROM epochs"))
	require.Zero(t, d.count(t, "SELECT count(*) FROM accounts"))
}

func TestRacingWritersGetOneLease(t *testing.T) {
	d := open(t)

	var (
		wins sync.Map
		wg   sync.WaitGroup
	)
	for i := range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			instance := string(rune('a' + i))
			if _, held, err := d.store.Lease(context.Background(), instance, time.Minute); err == nil && held {
				wins.Store(instance, true)
			}
		}()
	}

	wg.Wait()

	winners := 0

	wins.Range(func(any, any) bool { winners++; return true })
	require.Equal(t, 1, winners)
}

func TestABuiltEpochIsStoredWithEveryAccountOrNotAtAll(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	p := plan(t, 0, notary.Digest{}, 40)

	saved, err := d.store.SaveBuilt(ctx, lease, ledgerID, p)
	require.NoError(t, err)

	inFlight, found, err := d.store.InFlight(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, saved, inFlight)
	require.Equal(
		t,
		notary.EpochRecord{
			Ledger:   ledgerID,
			Header:   p.Header,
			Digest:   p.Digest,
			Currency: "IDR",
			State:    notary.StateBuilt,
		},
		inFlight,
	)
	require.Equal(t, int64(40), d.count(t, "SELECT count(*) FROM accounts WHERE epoch_number = 0"))

	for _, account := range p.Accounts {
		require.Equal(
			t,
			int64(1),
			d.count(
				t,
				"SELECT count(*) FROM accounts WHERE position = $1 AND holder = $2 AND kind = $3 AND amount = $4 AND salt = $5 AND commitment = $6",
				int64(account.Position),
				account.Liability.Holder.String(),
				int16(account.Liability.Kind),
				int64(account.Liability.Amount),
				account.Salt[:],
				account.Commitment[:],
			),
		)
	}

	_, err = d.store.SaveBuilt(ctx, lease, ledgerID, plan(t, 1, p.Digest, 2))
	require.Equal(t, notary.KindLedgerDiverged, notary.KindOf(err))
	require.Equal(t, int64(40), d.count(t, "SELECT count(*) FROM accounts"))
}

func TestAnEpochWalksToFinalWithItsAttempts(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	p := plan(t, 0, notary.Digest{}, 5)
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, p)
	require.NoError(t, err)

	sent := time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)
	cheap := notary.Attempt{
		Epoch: 0,
		Anchor: notary.PreparedAnchor{
			Hash:     notary.Digest{1},
			Raw:      []byte{0xaa},
			Nonce:    7,
			GasLimit: 180_000,
			Fees:     notary.FeeQuote{MaxFeePerGas: 1_000, MaxPriorityFeePerGas: 10},
		},
		SentAt: sent,
	}
	dear := cheap
	dear.Anchor.Hash = notary.Digest{2}
	dear.Anchor.Fees = notary.FeeQuote{MaxFeePerGas: 1_200, MaxPriorityFeePerGas: 12}
	dear.SentAt = sent.Add(5 * time.Minute)

	require.NoError(t, d.store.RecordAttempt(ctx, lease, cheap))
	require.NoError(t, d.store.RecordAttempt(ctx, lease, dear))

	inFlight, _, err := d.store.InFlight(ctx)
	require.NoError(t, err)
	require.Equal(t, notary.StateSubmitted, inFlight.State)

	receipt := notary.AnchorReceipt{
		Hash:              dear.Anchor.Hash,
		Succeeded:         true,
		BlockNumber:       48_100_000,
		GasUsed:           144_000,
		EffectiveGasPrice: 1_100,
		L1FeeWei:          9,
	}
	require.NoError(t, d.store.SettleAttempt(ctx, lease, dear.Anchor.Hash, notary.AttemptIncluded, &receipt))
	require.NoError(t, d.store.SettleAttempt(ctx, lease, cheap.Anchor.Hash, notary.AttemptSuperseded, nil))
	require.NoError(t, d.store.Transition(ctx, lease, 0, notary.StateSubmitted, notary.StateMined, ""))
	require.NoError(t, d.store.Transition(ctx, lease, 0, notary.StateMined, notary.StateFinalized, ""))

	attempts, err := d.store.Attempts(ctx, 0)
	require.NoError(t, err)
	require.Len(t, attempts, 2)
	require.Equal(t, notary.AttemptSuperseded, attempts[0].Outcome)
	require.Nil(t, attempts[0].Receipt)
	require.Equal(t, cheap.Anchor, attempts[0].Anchor)
	require.Equal(t, notary.AttemptIncluded, attempts[1].Outcome)
	require.Equal(t, &receipt, attempts[1].Receipt)
	require.True(t, attempts[1].SentAt.Equal(dear.SentAt))

	_, found, err := d.store.InFlight(ctx)
	require.NoError(t, err)
	require.False(t, found)
	last, found, err := d.store.LastFinalized(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, p.Digest, last.Digest)

	costs, err := d.store.RecentCosts(ctx, 24)
	require.NoError(t, err)
	require.Equal(t, []uint64{144_000*1_100 + 9}, costs)

	_, err = d.store.SaveBuilt(ctx, lease, ledgerID, plan(t, 1, p.Digest, 3))
	require.NoError(t, err)
}

func TestTransitionsAreCheckedTwice(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, plan(t, 0, notary.Digest{}, 1))
	require.NoError(t, err)

	err = d.store.Transition(ctx, lease, 0, notary.StateBuilt, notary.StateFinalized, "")
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	err = d.store.Transition(ctx, lease, 0, notary.StateSubmitted, notary.StateMined, "")
	require.Equal(t, notary.KindLeaseLost, notary.KindOf(err))
	err = d.store.Transition(ctx, lease, 0, notary.StateBuilt, notary.StateAbandoned, "")
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	err = d.store.SettleAttempt(ctx, lease, notary.Digest{9}, notary.AttemptIncluded, nil)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	require.NoError(
		t,
		d.store.Transition(ctx, lease, 0, notary.StateBuilt, notary.StateAbandoned, "EpochConflict(0, 0xaa, 0xbb)"),
	)
	abandoned, err := d.store.AbandonedCount(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), abandoned)

	_, found, err := d.store.InFlight(ctx)
	require.NoError(t, err)
	require.False(t, found)
}
