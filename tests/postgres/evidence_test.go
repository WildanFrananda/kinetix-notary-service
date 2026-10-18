package postgres_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var evidenceRegistry = notary.RegistryIdentity{ChainID: 84532, Address: notary.Address{0xe0, 0x12}}

func entered(t *testing.T, records []notary.EvidenceRecord) []notary.EvidenceEntry {
	t.Helper()

	salts := &notarytest.Salts{Seed: 23}
	entries := make([]notary.EvidenceEntry, len(records))

	for i, record := range records {
		entry, err := notary.EnterEvidence(keccak.NewEvidence(), salts, record)
		require.NoError(t, err)

		entries[i] = entry
	}

	return entries
}

func last(entries []notary.EvidenceEntry) notary.LogPosition {
	return entries[len(entries)-1].Record.Position()
}

func (d *db) collect(t *testing.T, lease notary.Lease, stream notary.EvidenceStream, entries []notary.EvidenceEntry) {
	t.Helper()

	from, err := d.store.Cursor(context.Background(), stream)
	require.NoError(t, err)
	require.NoError(t, d.store.Collect(context.Background(), lease, stream, from, entries, last(entries)))
}

func batchPlan(
	t *testing.T,
	number uint64,
	previous notary.Digest,
	entries []notary.EvidenceEntry,
) notary.EvidenceBatchPlan {
	t.Helper()

	closedAt, err := notary.UnixTimeFromSeconds(1_760_000_000 + number)
	require.NoError(t, err)

	plan, err := notary.PlanEvidenceBatch(
		keccak.NewEvidence(),
		entries[0].Record.Stream(),
		number,
		previous,
		closedAt,
		entries,
	)
	require.NoError(t, err)

	return plan
}

func TestEveryEvidenceFieldComesBackAsItWasKept(t *testing.T) {
	d := open(t)
	lease := d.lease(t, "notaryd-1")

	lifecycle := notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 3)
	lifecycle[1].ReturnNumber, lifecycle[1].ReturnStatus, lifecycle[1].ActorPrincipalID = "RET-1", 2, ""
	lifecycle[2].OccurredAt = notarytest.EvidenceTime(t, 1_760_000_123, 999_999_900)
	deliveries := notarytest.Deliveries(t, 2)
	deliveries[1].SignatureSHA256 = bytes.Repeat([]byte{0x5e}, 32)
	reviews := notarytest.Reviews(t, 2)
	reviews[0].CommentSHA256, reviews[0].MerchantPrincipalID = nil, ""

	for _, records := range [][]notary.EvidenceRecord{
		notarytest.Records(lifecycle),
		notarytest.Records(deliveries),
		notarytest.Records(reviews),
	} {
		entries := entered(t, records)
		stream := records[0].Stream()
		d.collect(t, lease, stream, entries)

		kept, err := d.store.Unbatched(context.Background(), stream, 100)
		require.NoError(t, err)
		require.Equal(t, entries, kept)

		cursor, err := d.store.Cursor(context.Background(), stream)
		require.NoError(t, err)
		require.Equal(t, last(entries), cursor)

		pending, err := d.store.UnbatchedCount(context.Background(), stream)
		require.NoError(t, err)
		require.Equal(t, int64(len(entries)), pending)
	}
}

func TestTheCursorMovesOnlyFromWhereItStands(t *testing.T) {
	d := open(t)
	lease := d.lease(t, "notaryd-1")
	entries := entered(t, notarytest.Records(notarytest.Deliveries(t, 3)))

	err := d.store.Collect(
		context.Background(),
		lease,
		notary.StreamDelivery,
		notary.LogPosition{Transaction: 1, Sequence: 1},
		entries,
		last(entries),
	)
	require.Equal(t, notary.KindLeaseLost, notary.KindOf(err))
	require.Zero(t, d.count(t, "SELECT count(*) FROM evidence_leaves"))

	err = d.store.Collect(
		context.Background(),
		lease,
		notary.StreamDelivery,
		notary.LogPosition{},
		nil,
		notary.LogPosition{},
	)
	require.NoError(t, err)

	err = d.store.Collect(
		context.Background(),
		lease,
		notary.StreamDelivery,
		last(entries),
		nil,
		notary.LogPosition{},
	)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestAWriterThatLostTheLeaseKeepsNoEvidence(t *testing.T) {
	d := open(t)
	lease := d.lease(t, "notaryd-1")
	entries := entered(t, notarytest.Records(notarytest.Deliveries(t, 2)))
	lease.Token++

	err := d.store.Collect(
		context.Background(),
		lease,
		notary.StreamDelivery,
		notary.LogPosition{},
		entries,
		last(entries),
	)

	require.Equal(t, notary.KindLeaseLost, notary.KindOf(err))
	require.Zero(t, d.count(t, "SELECT count(*) FROM delivery_records"))
}

func TestARecordKeptTwiceIsADisagreement(t *testing.T) {
	d := open(t)
	lease := d.lease(t, "notaryd-1")
	entries := entered(t, notarytest.Records(notarytest.Reviews(t, 2)))
	d.collect(t, lease, notary.StreamReview, entries)

	_, err := d.pool.Exec(
		context.Background(),
		"UPDATE evidence_cursors SET through_transaction = 0, through_sequence = 0 WHERE stream = 3",
	)
	require.NoError(t, err)

	err = d.store.Collect(
		context.Background(),
		lease,
		notary.StreamReview,
		notary.LogPosition{},
		entries,
		last(entries),
	)

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
	require.Equal(t, int64(2), d.count(t, "SELECT count(*) FROM review_records"))
}

func TestABatchTakesExactlyTheUnbatchedLeavesItNames(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "notaryd-1")
	entries := entered(t, notarytest.Records(notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 5)))
	d.collect(t, lease, notary.StreamOrderLifecycle, entries)

	plan := batchPlan(t, 0, notary.Digest{}, entries[:3])
	record, err := d.store.SaveBatch(ctx, lease, evidenceRegistry, plan)
	require.NoError(t, err)
	require.Equal(t, notary.StateBuilt, record.State)

	stored, found, err := d.store.EvidenceBatch(ctx, notary.StreamOrderLifecycle, 0)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, record, stored)

	leaves, err := d.store.EvidenceBatchLeaves(ctx, notary.StreamOrderLifecycle, 0)
	require.NoError(t, err)
	require.Equal(t, []notary.Digest{entries[0].Leaf, entries[1].Leaf, entries[2].Leaf}, leaves)

	left, err := d.store.Unbatched(ctx, notary.StreamOrderLifecycle, 100)
	require.NoError(t, err)
	require.Equal(t, entries[3:], left)

	require.NoError(
		t,
		d.store.TransitionBatch(
			ctx,
			lease,
			notary.StreamOrderLifecycle,
			0,
			notary.StateBuilt,
			notary.StateAbandoned,
			"test",
		),
	)

	again := batchPlan(t, 1, plan.Digest, entries[2:4])
	_, err = d.store.SaveBatch(ctx, lease, evidenceRegistry, again)
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))

	_, found, err = d.store.EvidenceBatch(ctx, notary.StreamOrderLifecycle, 1)
	require.NoError(t, err)
	require.False(t, found, "a refused batch left its row behind")
}

func TestOnlyOneEvidenceBatchIsInFlightAcrossEveryStream(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "notaryd-1")
	lifecycle := entered(t, notarytest.Records(notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 2)))
	deliveries := entered(t, notarytest.Records(notarytest.Deliveries(t, 2)))
	d.collect(t, lease, notary.StreamOrderLifecycle, lifecycle)
	d.collect(t, lease, notary.StreamDelivery, deliveries)

	_, err := d.store.SaveBatch(ctx, lease, evidenceRegistry, batchPlan(t, 0, notary.Digest{}, lifecycle))
	require.NoError(t, err)

	_, err = d.store.SaveBatch(ctx, lease, evidenceRegistry, batchPlan(t, 0, notary.Digest{}, deliveries))
	require.Equal(t, notary.KindLedgerDiverged, notary.KindOf(err))

	inFlight, found, err := d.store.InFlightBatch(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, notary.StreamOrderLifecycle, inFlight.Header.Stream)
}

func TestABatchIsAttemptedSettledFinalizedAndCosted(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "notaryd-1")
	entries := entered(t, notarytest.Records(notarytest.Deliveries(t, 2)))
	d.collect(t, lease, notary.StreamDelivery, entries)
	plan := batchPlan(t, 0, notary.Digest{}, entries)
	_, err := d.store.SaveBatch(ctx, lease, evidenceRegistry, plan)
	require.NoError(t, err)

	attempt := notary.BatchAttempt{
		Stream: notary.StreamDelivery,
		Batch:  0,
		Anchor: notary.PreparedAnchor{
			Hash:     notary.Digest{0xaa},
			Raw:      []byte{0x02, 0x01},
			Nonce:    4,
			GasLimit: 90_000,
			Fees:     notary.FeeQuote{MaxFeePerGas: 2_100, MaxPriorityFeePerGas: 100},
		},
		Outcome: notary.AttemptPending,
		SentAt:  time.Unix(1_760_000_000, 0).UTC(),
	}
	require.NoError(t, d.store.RecordBatchAttempt(ctx, lease, attempt))

	submitted, found, err := d.store.InFlightBatch(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, notary.StateSubmitted, submitted.State)

	receipt := notary.AnchorReceipt{
		Hash:              attempt.Anchor.Hash,
		BlockNumber:       9,
		GasUsed:           80_000,
		EffectiveGasPrice: 1_100,
		L1FeeWei:          5,
	}
	require.NoError(t, d.store.SettleBatchAttempt(ctx, lease, attempt.Anchor.Hash, notary.AttemptIncluded, &receipt))
	require.Equal(
		t,
		notary.KindInvalidArgument,
		notary.KindOf(d.store.SettleBatchAttempt(ctx, lease, attempt.Anchor.Hash, notary.AttemptIncluded, nil)),
	)
	require.NoError(
		t,
		d.store.TransitionBatch(ctx, lease, notary.StreamDelivery, 0, notary.StateSubmitted, notary.StateFinalized, ""),
	)

	attempts, err := d.store.BatchAttempts(ctx, notary.StreamDelivery, 0)
	require.NoError(t, err)
	require.Len(t, attempts, 1)
	require.Equal(t, notary.AttemptIncluded, attempts[0].Outcome)
	require.Equal(t, attempt.Anchor, attempts[0].Anchor)

	costs, err := d.store.RecentBatchCosts(ctx, 24)
	require.NoError(t, err)
	require.Equal(t, []uint64{80_000*1_100 + 5}, costs)

	finalized, found, err := d.store.LastFinalizedBatch(ctx, notary.StreamDelivery)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, plan.Digest, finalized.Digest)

	abandoned, err := d.store.AbandonedBatches(ctx)
	require.NoError(t, err)
	require.Zero(t, abandoned)
}

func TestAnOrdersAndAReviewsEvidenceIsFoundByName(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "notaryd-1")
	lifecycle := entered(t, notarytest.Records(notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 10, 2)))
	deliveryRecords := notarytest.Deliveries(t, 2)
	deliveryRecords[0].OrderNumber = "ORD-1"
	deliveries := entered(t, notarytest.Records(deliveryRecords))
	reviews := entered(t, notarytest.Records(notarytest.Reviews(t, 2)))
	d.collect(t, lease, notary.StreamOrderLifecycle, lifecycle)
	d.collect(t, lease, notary.StreamDelivery, deliveries)
	d.collect(t, lease, notary.StreamReview, reviews)
	_, err := d.store.SaveBatch(ctx, lease, evidenceRegistry, batchPlan(t, 0, notary.Digest{}, lifecycle))
	require.NoError(t, err)

	order, err := d.store.OrderEvidence(ctx, "ORD-1")
	require.NoError(t, err)
	require.Len(t, order, 3)
	require.True(t, order[1].Batched)
	require.Equal(t, uint64(1), order[1].Index)
	require.Equal(t, lifecycle[1], order[1].Entry)
	require.False(t, order[2].Batched)
	require.Equal(t, deliveries[0], order[2].Entry)

	review, err := d.store.ReviewEvidence(ctx, "review-b")
	require.NoError(t, err)
	require.Len(t, review, 1)
	require.Equal(t, reviews[1], review[0].Entry)

	none, err := d.store.OrderEvidence(ctx, "ORD-404")
	require.NoError(t, err)
	require.Empty(t, none)

	_, err = d.store.OrderEvidence(ctx, "")
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestTheEvidenceTablesRefuseWhatTheSchemeCannotHold(t *testing.T) {
	d := open(t)

	for _, statement := range []string{
		"INSERT INTO evidence_cursors (stream, through_transaction, through_sequence) VALUES (4, 0, 0)",
		"INSERT INTO evidence_leaves (stream, transaction_id, sequence, salt, record_digest, leaf) " +
			"VALUES (1, 1, 1, '\\x00', '\\x00', '\\x00')",
		"INSERT INTO delivery_records (transaction_id, sequence, order_number, driver_principal_id, photo_sha256, " +
			"delivered_seconds, delivered_nanos) VALUES (1, 1, '', 'driver', '\\x00', 1, 1000000000)",
	} {
		_, err := d.pool.Exec(context.Background(), statement)
		require.Error(t, err, statement)
	}
}
