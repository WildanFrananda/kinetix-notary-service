package evidencing_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestEveryStreamIsCollectedBatchedAnchoredAndFinalizedInTurn(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	require.Equal(t, 4, r.observer.Collected[notary.StreamOrderLifecycle])
	require.Equal(t, 3, r.observer.Collected[notary.StreamDelivery])
	require.Equal(t, 2, r.observer.Collected[notary.StreamReview])
	require.Equal(t, "submitted", r.state(t, notary.StreamOrderLifecycle, 0))
	require.Empty(t, r.state(t, notary.StreamDelivery, 0), "one anchor in flight at a time")
	require.Len(t, r.chain.Broadcasts, 1)

	r.chain.Mine()
	r.tick(t)
	require.Equal(t, "mined", r.state(t, notary.StreamOrderLifecycle, 0))

	r.chain.Finalize()
	r.tick(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamOrderLifecycle, 0))

	r.tick(t)
	require.Equal(t, "submitted", r.state(t, notary.StreamDelivery, 0))
	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamDelivery, 0))

	r.tick(t)
	require.Equal(t, "submitted", r.state(t, notary.StreamReview, 0))
	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamReview, 0))

	for _, stream := range notary.EvidenceStreams() {
		head, err := r.chain.StreamHead(context.Background(), stream)
		require.NoError(t, err)

		last, found, err := r.store.LastFinalizedBatch(context.Background(), stream)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, head.Head, last.Digest)
		require.Equal(t, uint64(1), head.NextBatch)
		require.Equal(t, last, r.observer.Finalized[stream])
		require.Zero(t, r.observer.Pending[stream])
	}

	require.Equal(t, 4, r.count(t, "SELECT count(*) FROM evidence_leaves WHERE stream = 1 AND batch_number = 0"))
	require.True(t, r.observer.RunwayKnown)
	require.Positive(t, r.observer.MeanCost)
}

func TestTheNextBatchWaitsForTheIntervalAndChainsToTheLastOne(t *testing.T) {
	r := newRig(t)
	r.feeds[notary.StreamDelivery].Records = nil
	r.feeds[notary.StreamReview].Records = nil

	r.tick(t)
	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamOrderLifecycle, 0))

	r.feeds[notary.StreamOrderLifecycle].Records = append(
		r.feeds[notary.StreamOrderLifecycle].Records,
		notarytest.Records(notarytest.Lifecycle(t, "ORD-2", "buyer-2", "merchant-2", 200, 2))...,
	)

	r.advance(time.Hour)
	r.tick(t)
	require.Equal(t, int64(2), r.observer.Pending[notary.StreamOrderLifecycle])
	require.Empty(t, r.state(t, notary.StreamOrderLifecycle, 1), "a batch before its interval")

	r.advance(24 * time.Hour)
	r.tick(t)
	require.Equal(t, "submitted", r.state(t, notary.StreamOrderLifecycle, 1))

	first, _, err := r.store.LastFinalizedBatch(context.Background(), notary.StreamOrderLifecycle)
	require.NoError(t, err)

	second, found, err := r.store.InFlightBatch(context.Background())
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, first.Digest, second.Header.PreviousDigest)
	require.Equal(t, uint64(2), second.Header.LeafCount)
	require.Equal(t, notary.LogPosition{Transaction: 201, Sequence: 2}, second.Header.Through)
}

func TestAFullBatchDoesNotWaitForTheInterval(t *testing.T) {
	r := newRig(t)
	r.policy.MaxBatch = 3
	r.build()
	r.feeds[notary.StreamOrderLifecycle].Records = notarytest.Records(
		notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 100, 7),
	)
	r.feeds[notary.StreamDelivery].Records = nil
	r.feeds[notary.StreamReview].Records = nil

	r.tick(t)
	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamOrderLifecycle, 0))

	r.advance(time.Minute)
	r.tick(t)
	require.Equal(t, "submitted", r.state(t, notary.StreamOrderLifecycle, 1))
	r.settle(t)

	r.advance(time.Minute)
	r.tick(t)
	require.Empty(t, r.state(t, notary.StreamOrderLifecycle, 2), "the remaining record waits for the interval")
	require.Equal(t, int64(1), r.observer.Pending[notary.StreamOrderLifecycle])
}

func TestAnOutageOfOneOwnerStallsOnlyItsOwnStream(t *testing.T) {
	r := newRig(t)
	r.feeds[notary.StreamOrderLifecycle].Err = notary.Errf(notary.KindSourceUnavailable, "test", nil, "order is down")

	r.tick(t)
	require.False(t, r.observer.FeedUp[notary.StreamOrderLifecycle])
	require.True(t, r.observer.FeedUp[notary.StreamDelivery])
	require.Zero(t, r.observer.Collected[notary.StreamOrderLifecycle])
	require.Equal(t, "submitted", r.state(t, notary.StreamDelivery, 0))
	require.Contains(t, r.logs.String(), "evidence collection will retry")

	r.feeds[notary.StreamOrderLifecycle].Err = nil
	r.settle(t)
	r.tick(t)
	require.Equal(t, "submitted", r.state(t, notary.StreamOrderLifecycle, 0))
}

func TestAPageOutOfOrderIsKeptNowhere(t *testing.T) {
	r := newRig(t)
	r.feeds[notary.StreamReview].Tamper = func(page *notary.EvidencePage) {
		if len(page.Records) == 2 {
			page.Records[0], page.Records[1] = page.Records[1], page.Records[0]
		}
	}

	r.tick(t)
	require.False(t, r.observer.FeedUp[notary.StreamReview])

	cursor, err := r.store.Cursor(context.Background(), notary.StreamReview)
	require.NoError(t, err)
	require.Equal(t, notary.LogPosition{}, cursor)
	require.Zero(t, r.count(t, "SELECT count(*) FROM evidence_leaves WHERE stream = 3"))
	require.Contains(t, r.logs.String(), "evidence collection stopped short")
}

func TestACollectorResumesFromItsCursorWithoutKeepingARecordTwice(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	r.build()
	r.tick(t)

	require.Equal(t, 4, r.count(t, "SELECT count(*) FROM lifecycle_records"))
	require.Equal(t, 3, r.count(t, "SELECT count(*) FROM delivery_records"))
	require.Equal(t, 2, r.count(t, "SELECT count(*) FROM review_records"))

	asked := r.feeds[notary.StreamDelivery].Asked
	require.Equal(t, notary.LogPosition{Transaction: 900, Sequence: 3}, asked[len(asked)-1])
}

func TestAWriterWithoutTheLeaseCollectsNothing(t *testing.T) {
	r := newRig(t)
	_, held, err := r.store.Lease(context.Background(), "notaryd-2", time.Minute)
	require.NoError(t, err)
	require.True(t, held)

	r.tick(t)

	require.Empty(t, r.feeds[notary.StreamOrderLifecycle].Asked)
	require.Zero(t, r.count(t, "SELECT count(*) FROM evidence_leaves"))
}

func TestAForeignBatchAbandonsOursAndStopsOnlyThatStream(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	ours, found, err := r.store.InFlightBatch(context.Background())
	require.NoError(t, err)
	require.True(t, found)

	foreign := ours.Header
	foreign.Root = notary.Digest{0xba, 0xd0}
	r.chain.Plant(foreign)
	r.chain.Mine()

	r.tick(t)
	require.Equal(t, "abandoned", r.state(t, notary.StreamOrderLifecycle, 0))
	require.Contains(t, r.logs.String(), "evidence batch abandoned")

	r.tick(t)
	require.Equal(t, "submitted", r.state(t, notary.StreamDelivery, 0))
	require.Equal(t, int64(1), r.observer.Abandoned)

	r.settle(t)
	r.tick(t)
	require.Contains(t, r.logs.String(), "a person must decide what happens next")
	require.Empty(t, r.state(t, notary.StreamOrderLifecycle, 1))
	require.Equal(t, "submitted", r.state(t, notary.StreamReview, 0))
}

func TestAStuckEvidenceAnchorIsOutbidAtTheSameNonce(t *testing.T) {
	r := newRig(t)

	r.tick(t)
	r.advance(6 * time.Minute)
	r.tick(t)

	attempts, err := r.store.BatchAttempts(context.Background(), notary.StreamOrderLifecycle, 0)
	require.NoError(t, err)
	require.Len(t, attempts, 2)
	require.Equal(t, attempts[0].Anchor.Nonce, attempts[1].Anchor.Nonce)
	require.GreaterOrEqual(t, attempts[1].Anchor.Fees.MaxFeePerGas*8, attempts[0].Anchor.Fees.MaxFeePerGas*9)

	r.settle(t)
	require.Equal(t, "finalized", r.state(t, notary.StreamOrderLifecycle, 0))

	attempts, err = r.store.BatchAttempts(context.Background(), notary.StreamOrderLifecycle, 0)
	require.NoError(t, err)
	require.Equal(t, notary.AttemptSuperseded, attempts[0].Outcome)
	require.Equal(t, notary.AttemptIncluded, attempts[1].Outcome)
}

func TestTheFeeCeilingHoldsAnEvidenceAnchorBack(t *testing.T) {
	r := newRig(t)
	r.chain.BaseFee = 20_000

	r.tick(t)

	require.Equal(t, "built", r.state(t, notary.StreamOrderLifecycle, 0))
	require.Empty(t, r.state(t, notary.StreamDelivery, 0))
	require.NotContains(t, r.logs.String(), "collides", "a second batch was attempted beside the one in flight")
	require.Empty(t, r.chain.Broadcasts)
	require.Equal(t, 1, r.observer.CeilingHits)
}

func TestASealedRegistryGetsNoBatch(t *testing.T) {
	r := newRig(t)
	r.chain.Sealed = true

	r.tick(t)

	require.Zero(t, r.count(t, "SELECT count(*) FROM evidence_batches"))
	require.Contains(t, r.logs.String(), "the evidence registry is sealed")
}

func TestARegistryThatKeepsNoStreamIsNeverWrittenTo(t *testing.T) {
	r := newRig(t)
	r.chain.Unknown[notary.StreamOrderLifecycle] = true

	r.tick(t)

	require.Empty(t, r.state(t, notary.StreamOrderLifecycle, 0))
	require.Equal(t, "submitted", r.state(t, notary.StreamDelivery, 0))
}

func TestAChainThatAlreadyHoldsBatchesThisDatabaseNeverBuiltIsRefused(t *testing.T) {
	r := newRig(t)
	scheme := r.chain.Scheme
	closedAt, err := notary.UnixTimeFromSeconds(r.chain.BlockTime - 10)
	require.NoError(t, err)

	r.chain.Plant(notary.EvidenceBatchHeader{
		Stream:    notary.StreamOrderLifecycle,
		Number:    0,
		Root:      scheme.EvidenceEmpty(),
		LeafCount: 1,
		Through:   notary.LogPosition{Transaction: 1, Sequence: 1},
		ClosedAt:  closedAt,
	})

	r.tick(t)

	require.Empty(t, r.state(t, notary.StreamOrderLifecycle, 0))
	require.Contains(t, r.logs.String(), "this database built none of them")
}
