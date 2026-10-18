package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

var (
	_ notary.EvidenceStore       = (*Store)(nil)
	_ notary.EvidenceProofSource = (*Store)(nil)
)

func (s *Store) Cursor(ctx context.Context, stream notary.EvidenceStream) (notary.LogPosition, error) {
	const op = "postgres.Store.Cursor"

	key, err := streamKey(op, stream)
	if err != nil {
		return notary.LogPosition{}, err
	}

	row, err := s.queries.EvidenceCursor(ctx, key)
	if errors.Is(err, pgx.ErrNoRows) {
		return notary.LogPosition{}, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the database keeps no cursor for the %s stream",
			stream,
		)
	}

	if err != nil {
		return notary.LogPosition{}, unavailable(op, err)
	}

	return positionOf(op, row.ThroughTransaction, row.ThroughSequence)
}

func (s *Store) Collect(
	ctx context.Context,
	lease notary.Lease,
	stream notary.EvidenceStream,
	from notary.LogPosition,
	entries []notary.EvidenceEntry,
	to notary.LogPosition,
) error {
	const op = "postgres.Store.Collect"

	key, err := streamKey(op, stream)
	if err != nil {
		return err
	}

	if len(entries) == 0 && to == from {
		return nil
	}

	if !to.After(from) {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the %s cursor moves forward, not from %d/%d to %d/%d",
			stream,
			from.Transaction,
			from.Sequence,
			to.Transaction,
			to.Sequence,
		)
	}

	rows, err := collectedRowsOf(op, key, stream, entries)
	if err != nil {
		return err
	}

	fromTransaction, fromSequence, err := positionKey(op, from)
	if err != nil {
		return err
	}

	toTransaction, toSequence, err := positionKey(op, to)
	if err != nil {
		return err
	}

	err = s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		moved, err := q.AdvanceEvidenceCursor(ctx, storegen.AdvanceEvidenceCursorParams{
			ToTransaction:   toTransaction,
			ToSequence:      toSequence,
			Stream:          key,
			FromTransaction: fromTransaction,
			FromSequence:    fromSequence,
		})
		if err != nil {
			return err
		}

		if moved != 1 {
			return notary.Errf(
				notary.KindLeaseLost,
				op,
				nil,
				"the %s cursor is no longer at %d/%d",
				stream,
				from.Transaction,
				from.Sequence,
			)
		}

		return rows.insert(ctx, op, q)
	})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return notary.Errf(
			notary.KindDisagreement,
			op,
			err,
			"the database already keeps a %s record this page carries",
			stream,
		)
	}

	return err
}

func (s *Store) Unbatched(
	ctx context.Context,
	stream notary.EvidenceStream,
	limit int32,
) ([]notary.EvidenceEntry, error) {
	const op = "postgres.Store.Unbatched"

	if limit <= 0 {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "a batch takes at least one record")
	}

	stored, err := s.unbatched(ctx, op, stream, limit)
	if err != nil {
		return nil, err
	}

	entries := make([]notary.EvidenceEntry, len(stored))
	for i, record := range stored {
		if record.Batched {
			return nil, notary.Errf(notary.KindStoreUnavailable, op, nil, "an unbatched %s record has a batch", stream)
		}

		entries[i] = record.Entry
	}

	return entries, nil
}

func (s *Store) unbatched(
	ctx context.Context,
	op string,
	stream notary.EvidenceStream,
	limit int32,
) ([]notary.StoredEvidence, error) {
	switch stream {
	case notary.StreamOrderLifecycle:
		rows, err := s.queries.UnbatchedLifecycleRecords(ctx, limit)
		if err != nil {
			return nil, unavailable(op, err)
		}

		return convertAll(op, rows, lifecycleOf)
	case notary.StreamDelivery:
		rows, err := s.queries.UnbatchedDeliveryRecords(ctx, limit)
		if err != nil {
			return nil, unavailable(op, err)
		}

		return convertAll(op, rows, deliveryOf)
	case notary.StreamReview:
		rows, err := s.queries.UnbatchedReviewRecords(ctx, limit)
		if err != nil {
			return nil, unavailable(op, err)
		}

		return convertAll(op, rows, reviewOf)
	default:
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "stream %d is not an evidence stream", stream)
	}
}

func (s *Store) UnbatchedCount(ctx context.Context, stream notary.EvidenceStream) (int64, error) {
	const op = "postgres.Store.UnbatchedCount"

	key, err := streamKey(op, stream)
	if err != nil {
		return 0, err
	}

	count, err := s.queries.UnbatchedEvidenceCount(ctx, key)
	if err != nil {
		return 0, unavailable(op, err)
	}

	return count, nil
}

func (s *Store) InFlightBatch(ctx context.Context) (notary.EvidenceBatchRecord, bool, error) {
	const op = "postgres.Store.InFlightBatch"

	row, err := s.queries.InFlightEvidenceBatch(ctx)

	return batchFound(op, row, err)
}

func (s *Store) LastBatch(ctx context.Context, stream notary.EvidenceStream) (notary.EvidenceBatchRecord, bool, error) {
	const op = "postgres.Store.LastBatch"

	key, err := streamKey(op, stream)
	if err != nil {
		return notary.EvidenceBatchRecord{}, false, err
	}

	row, err := s.queries.LastEvidenceBatch(ctx, key)

	return batchFound(op, row, err)
}

func (s *Store) LastFinalizedBatch(
	ctx context.Context,
	stream notary.EvidenceStream,
) (notary.EvidenceBatchRecord, bool, error) {
	const op = "postgres.Store.LastFinalizedBatch"

	key, err := streamKey(op, stream)
	if err != nil {
		return notary.EvidenceBatchRecord{}, false, err
	}

	row, err := s.queries.LastFinalizedEvidenceBatch(ctx, key)

	return batchFound(op, row, err)
}

func (s *Store) EvidenceBatch(
	ctx context.Context,
	stream notary.EvidenceStream,
	number uint64,
) (notary.EvidenceBatchRecord, bool, error) {
	const op = "postgres.Store.EvidenceBatch"

	key, err := streamKey(op, stream)
	if err != nil {
		return notary.EvidenceBatchRecord{}, false, err
	}

	n, err := signed(op, number)
	if err != nil {
		return notary.EvidenceBatchRecord{}, false, err
	}

	row, err := s.queries.EvidenceBatchByNumber(ctx, storegen.EvidenceBatchByNumberParams{Stream: key, Number: n})

	return batchFound(op, row, err)
}

func batchFound(op string, row storegen.EvidenceBatch, err error) (notary.EvidenceBatchRecord, bool, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return notary.EvidenceBatchRecord{}, false, nil
	}

	if err != nil {
		return notary.EvidenceBatchRecord{}, false, unavailable(op, err)
	}

	record, err := batchOf(op, row)
	if err != nil {
		return notary.EvidenceBatchRecord{}, false, err
	}

	return record, true, nil
}

func (s *Store) EvidenceBatchLeaves(
	ctx context.Context,
	stream notary.EvidenceStream,
	number uint64,
) ([]notary.Digest, error) {
	const op = "postgres.Store.EvidenceBatchLeaves"

	key, err := streamKey(op, stream)
	if err != nil {
		return nil, err
	}

	n, err := signed(op, number)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.EvidenceBatchLeaves(ctx, storegen.EvidenceBatchLeavesParams{Stream: key, BatchNumber: &n})
	if err != nil {
		return nil, unavailable(op, err)
	}

	leaves := make([]notary.Digest, len(rows))
	for i, raw := range rows {
		if leaves[i], err = digestOf(op, raw); err != nil {
			return nil, err
		}
	}

	return leaves, nil
}

func (s *Store) OrderEvidence(ctx context.Context, orderNumber string) ([]notary.StoredEvidence, error) {
	const op = "postgres.Store.OrderEvidence"

	if orderNumber == "" {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "an order is named by its number")
	}

	lifecycle, err := s.queries.LifecycleRecordsOfOrder(ctx, orderNumber)
	if err != nil {
		return nil, unavailable(op, err)
	}

	deliveries, err := s.queries.DeliveryRecordsOfOrder(ctx, orderNumber)
	if err != nil {
		return nil, unavailable(op, err)
	}

	stored := make([]notary.StoredEvidence, 0, len(lifecycle)+len(deliveries))

	for _, row := range lifecycle {
		record, err := lifecycleOf(op, storegen.UnbatchedLifecycleRecordsRow(row))
		if err != nil {
			return nil, err
		}

		stored = append(stored, record)
	}

	for _, row := range deliveries {
		record, err := deliveryOf(op, storegen.UnbatchedDeliveryRecordsRow(row))
		if err != nil {
			return nil, err
		}

		stored = append(stored, record)
	}

	return stored, nil
}

func (s *Store) ReviewEvidence(ctx context.Context, reviewID string) ([]notary.StoredEvidence, error) {
	const op = "postgres.Store.ReviewEvidence"

	if reviewID == "" {
		return nil, notary.Errf(notary.KindInvalidArgument, op, nil, "a review is named by its id")
	}

	rows, err := s.queries.ReviewRecordsOfReview(ctx, reviewID)
	if err != nil {
		return nil, unavailable(op, err)
	}

	stored := make([]notary.StoredEvidence, 0, len(rows))

	for _, row := range rows {
		record, err := reviewOf(op, storegen.UnbatchedReviewRecordsRow(row))
		if err != nil {
			return nil, err
		}

		stored = append(stored, record)
	}

	return stored, nil
}

func (s *Store) SaveBatch(
	ctx context.Context,
	lease notary.Lease,
	registry notary.RegistryIdentity,
	plan notary.EvidenceBatchPlan,
) (notary.EvidenceBatchRecord, error) {
	const op = "postgres.Store.SaveBatch"

	header := plan.Header

	key, err := streamKey(op, header.Stream)
	if err != nil {
		return notary.EvidenceBatchRecord{}, err
	}

	if uint64(len(plan.Entries)) != header.LeafCount {
		return notary.EvidenceBatchRecord{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"batch %d counts %d records and carries %d",
			header.Number,
			header.LeafCount,
			len(plan.Entries),
		)
	}

	ints := make([]int64, 0, 6)

	for _, v := range []uint64{
		header.Number,
		registry.ChainID,
		header.LeafCount,
		header.Through.Transaction,
		header.Through.Sequence,
		header.ClosedAt.Seconds(),
	} {
		i, err := signed(op, v)
		if err != nil {
			return notary.EvidenceBatchRecord{}, err
		}

		ints = append(ints, i)
	}

	assign := storegen.AssignEvidenceLeavesParams{
		BatchNumber:    ints[0],
		Stream:         key,
		TransactionIds: make([]int64, len(plan.Entries)),
		Sequences:      make([]int64, len(plan.Entries)),
		LeafIndexes:    make([]int64, len(plan.Entries)),
		Leaves:         make([][]byte, len(plan.Entries)),
	}

	for i, entry := range plan.Entries {
		transaction, sequence, err := positionKey(op, entry.Record.Position())
		if err != nil {
			return notary.EvidenceBatchRecord{}, err
		}

		assign.TransactionIds[i], assign.Sequences[i], assign.LeafIndexes[i] = transaction, sequence, int64(i)
		assign.Leaves[i] = entry.Leaf[:]
	}

	err = s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		if err := q.InsertEvidenceBatch(ctx, storegen.InsertEvidenceBatchParams{
			Stream:             key,
			Number:             ints[0],
			ChainID:            ints[1],
			RegistryAddress:    registry.Address[:],
			Scheme:             int32(plan.Scheme),
			PreviousDigest:     header.PreviousDigest[:],
			Root:               header.Root[:],
			Digest:             plan.Digest[:],
			LeafCount:          ints[2],
			ThroughTransaction: ints[3],
			ThroughSequence:    ints[4],
			ClosedAt:           ints[5],
		}); err != nil {
			return err
		}

		assigned, err := q.AssignEvidenceLeaves(ctx, assign)
		if err != nil {
			return err
		}

		if assigned != int64(len(plan.Entries)) {
			return notary.Errf(
				notary.KindDisagreement,
				op,
				nil,
				"batch %d of the %s stream names %d records, and the database keeps %d of them unbatched",
				header.Number,
				header.Stream,
				len(plan.Entries),
				assigned,
			)
		}

		return nil
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return notary.EvidenceBatchRecord{}, notary.Errf(
				notary.KindLedgerDiverged,
				op,
				err,
				"batch %d of the %s stream collides with one this database already holds",
				header.Number,
				header.Stream,
			)
		}

		return notary.EvidenceBatchRecord{}, err
	}

	return notary.EvidenceBatchRecord{
		Registry: registry,
		Scheme:   plan.Scheme,
		Header:   header,
		Digest:   plan.Digest,
		State:    notary.StateBuilt,
	}, nil
}

func (s *Store) RecordBatchAttempt(ctx context.Context, lease notary.Lease, attempt notary.BatchAttempt) error {
	const op = "postgres.Store.RecordBatchAttempt"

	key, err := streamKey(op, attempt.Stream)
	if err != nil {
		return err
	}

	anchor := attempt.Anchor
	ints := make([]int64, 0, 5)

	for _, v := range []uint64{
		attempt.Batch,
		anchor.Nonce,
		anchor.GasLimit,
		anchor.Fees.MaxFeePerGas,
		anchor.Fees.MaxPriorityFeePerGas,
	} {
		i, err := signed(op, v)
		if err != nil {
			return err
		}

		ints = append(ints, i)
	}

	return s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		if err := q.InsertEvidenceAttempt(ctx, storegen.InsertEvidenceAttemptParams{
			TxHash:               anchor.Hash[:],
			Stream:               key,
			BatchNumber:          ints[0],
			Nonce:                ints[1],
			RawTx:                anchor.Raw,
			GasLimit:             ints[2],
			MaxFeePerGas:         ints[3],
			MaxPriorityFeePerGas: ints[4],
			SentAt:               pgtype.Timestamptz{Time: attempt.SentAt, Valid: true},
		}); err != nil {
			return err
		}

		_, err := q.TransitionEvidenceBatch(ctx, storegen.TransitionEvidenceBatchParams{
			ToState:   notary.StateSubmitted.String(),
			Stream:    key,
			Number:    ints[0],
			FromState: notary.StateBuilt.String(),
		})

		return err
	})
}

func (s *Store) BatchAttempts(
	ctx context.Context,
	stream notary.EvidenceStream,
	batch uint64,
) ([]notary.BatchAttempt, error) {
	const op = "postgres.Store.BatchAttempts"

	key, err := streamKey(op, stream)
	if err != nil {
		return nil, err
	}

	number, err := signed(op, batch)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.EvidenceBatchAttempts(
		ctx,
		storegen.EvidenceBatchAttemptsParams{Stream: key, BatchNumber: number},
	)
	if err != nil {
		return nil, unavailable(op, err)
	}

	attempts := make([]notary.BatchAttempt, 0, len(rows))
	for _, row := range rows {
		attempt, err := batchAttemptOf(op, row)
		if err != nil {
			return nil, err
		}

		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

func (s *Store) SettleBatchAttempt(
	ctx context.Context,
	lease notary.Lease,
	hash notary.Digest,
	outcome notary.AttemptOutcome,
	receipt *notary.AnchorReceipt,
) error {
	const op = "postgres.Store.SettleBatchAttempt"

	settled := outcome == notary.AttemptIncluded || outcome == notary.AttemptReverted
	if settled != (receipt != nil) {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"an %s attempt carries a receipt exactly when it was mined",
			outcome,
		)
	}

	params := storegen.SettleEvidenceAttemptParams{Outcome: outcome.String(), TxHash: hash[:]}

	if receipt != nil {
		fields := make([]*int64, 0, 4)

		for _, v := range []uint64{receipt.BlockNumber, receipt.GasUsed, receipt.EffectiveGasPrice, receipt.L1FeeWei} {
			i, err := signed(op, v)
			if err != nil {
				return err
			}

			fields = append(fields, &i)
		}

		params.BlockNumber, params.GasUsed, params.EffectiveGasPrice, params.L1FeeWei = fields[0], fields[1], fields[2], fields[3]
	}

	return s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		changed, err := q.SettleEvidenceAttempt(ctx, params)
		if err != nil {
			return err
		}

		if changed != 1 {
			return notary.Errf(notary.KindInvalidArgument, op, nil, "no batch attempt %s is recorded", hash.Hex())
		}

		return nil
	})
}

func (s *Store) TransitionBatch(
	ctx context.Context,
	lease notary.Lease,
	stream notary.EvidenceStream,
	batch uint64,
	from, to notary.EpochState,
	reason string,
) error {
	const op = "postgres.Store.TransitionBatch"

	if !from.CanBecome(to) {
		return notary.Errf(notary.KindInvalidArgument, op, nil, "a batch does not go from %s to %s", from, to)
	}

	if (to == notary.StateAbandoned) != (reason != "") {
		return notary.Errf(notary.KindInvalidArgument, op, nil, "a reason is given exactly when a batch is abandoned")
	}

	key, err := streamKey(op, stream)
	if err != nil {
		return err
	}

	number, err := signed(op, batch)
	if err != nil {
		return err
	}

	var why *string
	if reason != "" {
		why = &reason
	}

	return s.fenced(ctx, op, lease, func(_ pgx.Tx, q *storegen.Queries) error {
		changed, err := q.TransitionEvidenceBatch(ctx, storegen.TransitionEvidenceBatchParams{
			ToState:   to.String(),
			Reason:    why,
			Stream:    key,
			Number:    number,
			FromState: from.String(),
		})
		if err != nil {
			return err
		}

		if changed != 1 {
			return notary.Errf(
				notary.KindLeaseLost,
				op,
				nil,
				"batch %d of the %s stream is no longer %s",
				batch,
				stream,
				from,
			)
		}

		return nil
	})
}

func (s *Store) AbandonedBatches(ctx context.Context) (int64, error) {
	count, err := s.queries.AbandonedEvidenceBatches(ctx)
	if err != nil {
		return 0, unavailable("postgres.Store.AbandonedBatches", err)
	}

	return count, nil
}

func (s *Store) RecentBatchCosts(ctx context.Context, limit int32) ([]uint64, error) {
	const op = "postgres.Store.RecentBatchCosts"

	rows, err := s.queries.RecentEvidenceReceipts(ctx, limit)
	if err != nil {
		return nil, unavailable(op, err)
	}

	costs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		hash, err := digestOf(op, row.TxHash)
		if err != nil {
			return nil, err
		}

		receipt, err := receiptOf(op, hash, row.BlockNumber, row.GasUsed, row.EffectiveGasPrice, row.L1FeeWei)
		if err != nil || receipt == nil {
			return nil, notary.Errf(notary.KindStoreUnavailable, op, err, "an included batch attempt has no receipt")
		}

		cost, err := receipt.CostWei()
		if err != nil {
			return nil, err
		}

		costs = append(costs, cost)
	}

	return costs, nil
}

func convertAll[R any](
	op string,
	rows []R,
	convert func(string, R) (notary.StoredEvidence, error),
) ([]notary.StoredEvidence, error) {
	stored := make([]notary.StoredEvidence, len(rows))

	for i, row := range rows {
		record, err := convert(op, row)
		if err != nil {
			return nil, err
		}

		stored[i] = record
	}

	return stored, nil
}
