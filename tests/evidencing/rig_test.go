package evidencing_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/postgres"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var registryID = notary.RegistryIdentity{ChainID: 84532, Address: notary.Address{0xe0, 0x01}}

type rig struct {
	now      time.Time
	pool     *pgxpool.Pool
	store    *postgres.Store
	chain    *notarytest.EvidenceChain
	feeds    map[notary.EvidenceStream]*notarytest.EvidenceFeed
	observer *notarytest.EvidenceObserver
	salts    *notarytest.Salts
	logs     *bytes.Buffer
	policy   evidencing.Policy
	ticker   *evidencing.Evidencer
}

func newRig(t *testing.T) *rig {
	t.Helper()

	pool := notarytest.Database(t)
	r := &rig{
		now:      time.Unix(1760000000, 0).UTC(),
		pool:     pool,
		store:    postgres.NewStore(pool, keccak.NewScheme()),
		chain:    notarytest.NewEvidenceChain(keccak.NewEvidence()),
		observer: notarytest.NewEvidenceObserver(),
		salts:    &notarytest.Salts{Seed: 11},
		logs:     &bytes.Buffer{},
		feeds: map[notary.EvidenceStream]*notarytest.EvidenceFeed{
			notary.StreamOrderLifecycle: {
				Kept:    notary.StreamOrderLifecycle,
				Records: notarytest.Records(notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 100, 4)),
			},
			notary.StreamDelivery: {
				Kept:    notary.StreamDelivery,
				Records: notarytest.Records(notarytest.Deliveries(t, 3)),
			},
			notary.StreamReview: {
				Kept:    notary.StreamReview,
				Records: notarytest.Records(notarytest.Reviews(t, 2)),
			},
		},
	}
	r.chain.BlockTime = uint64(r.now.Unix())
	r.policy = evidencing.Policy{
		Instance:      "notaryd-1",
		Registry:      registryID,
		BatchInterval: 24 * time.Hour,
		MaxBatch:      100,
		PageSize:      2,
		PagesPerTick:  10,
		ReplaceAfter:  5 * time.Minute,
		LeaseTerm:     30 * time.Second,
		FeeCeiling:    10_000,
		CostWindow:    24,
	}
	r.build()

	return r
}

func (r *rig) build() {
	log := slog.New(slog.NewJSONHandler(r.logs, nil))
	feeds := []notary.EvidenceFeed{
		r.feeds[notary.StreamOrderLifecycle],
		r.feeds[notary.StreamDelivery],
		r.feeds[notary.StreamReview],
	}
	scheme := keccak.NewEvidence()
	collector := evidencing.NewCollector(scheme, r.salts, feeds, r.store, r.observer, log, r.policy)
	closer := evidencing.NewCloser(scheme, r.store, r.chain, r.observer, r.clock, log, r.policy)
	r.ticker = evidencing.NewEvidencer(r.store, collector, closer, log, r.policy)
}

func (r *rig) clock() time.Time { return r.now }

func (r *rig) advance(d time.Duration) {
	r.now = r.now.Add(d)
	r.chain.BlockTime = uint64(r.now.Unix())
}

func (r *rig) tick(t *testing.T) {
	t.Helper()
	require.NoError(t, r.ticker.Tick(context.Background()))
}

func (r *rig) settle(t *testing.T) {
	t.Helper()

	r.chain.Mine()
	r.chain.Finalize()
	r.tick(t)
}

func (r *rig) state(t *testing.T, stream notary.EvidenceStream, number uint64) string {
	t.Helper()

	var state string

	err := r.pool.QueryRow(
		context.Background(),
		"SELECT state FROM evidence_batches WHERE stream = $1 AND number = $2",
		int16(stream),
		int64(number),
	).Scan(&state)
	if errors.Is(err, pgx.ErrNoRows) {
		return ""
	}

	require.NoError(t, err)

	return state
}

func (r *rig) count(t *testing.T, query string, args ...any) int {
	t.Helper()

	var n int
	require.NoError(t, r.pool.QueryRow(context.Background(), query, args...).Scan(&n))

	return n
}
