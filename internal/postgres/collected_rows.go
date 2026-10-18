package postgres

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

type collectedRows struct {
	leaves     []storegen.InsertEvidenceLeavesParams
	lifecycle  []storegen.InsertLifecycleRecordsParams
	deliveries []storegen.InsertDeliveryRecordsParams
	reviews    []storegen.InsertReviewRecordsParams
}

func collectedRowsOf(
	op string,
	key int16,
	stream notary.EvidenceStream,
	entries []notary.EvidenceEntry,
) (*collectedRows, error) {
	rows := &collectedRows{leaves: make([]storegen.InsertEvidenceLeavesParams, 0, len(entries))}

	for i, entry := range entries {
		if entry.Record.Stream() != stream {
			return nil, notary.Errf(
				notary.KindInvalidArgument,
				op,
				nil,
				"entry %d belongs to the %s stream, not %s",
				i,
				entry.Record.Stream(),
				stream,
			)
		}

		transaction, sequence, err := positionKey(op, entry.Record.Position())
		if err != nil {
			return nil, err
		}

		rows.leaves = append(rows.leaves, storegen.InsertEvidenceLeavesParams{
			Stream:        key,
			TransactionID: transaction,
			Sequence:      sequence,
			Salt:          entry.Salt[:],
			RecordDigest:  entry.RecordDigest[:],
			Leaf:          entry.Leaf[:],
		})

		if err := rows.keep(op, transaction, sequence, entry.Record); err != nil {
			return nil, err
		}
	}

	return rows, nil
}

func (c *collectedRows) keep(op string, transaction, sequence int64, record notary.EvidenceRecord) error {
	switch r := record.(type) {
	case notary.LifecycleRecord:
		seconds, nanos, err := timeKey(op, r.OccurredAt)
		if err != nil {
			return err
		}

		c.lifecycle = append(c.lifecycle, storegen.InsertLifecycleRecordsParams{
			TransactionID:       transaction,
			Sequence:            sequence,
			OrderNumber:         r.OrderNumber,
			Kind:                int64(r.Kind),
			FromStatus:          int64(r.FromStatus),
			ToStatus:            int64(r.ToStatus),
			ReturnNumber:        r.ReturnNumber,
			ReturnStatus:        int64(r.ReturnStatus),
			ActorKind:           int64(r.ActorKind),
			ActorPrincipalID:    r.ActorPrincipalID,
			BuyerPrincipalID:    r.BuyerPrincipalID,
			MerchantPrincipalID: r.MerchantPrincipalID,
			OccurredSeconds:     seconds,
			OccurredNanos:       nanos,
		})
	case notary.DeliveryRecord:
		seconds, nanos, err := timeKey(op, r.DeliveredAt)
		if err != nil {
			return err
		}

		c.deliveries = append(c.deliveries, storegen.InsertDeliveryRecordsParams{
			TransactionID:     transaction,
			Sequence:          sequence,
			OrderNumber:       r.OrderNumber,
			DriverPrincipalID: r.DriverPrincipalID,
			PhotoSha256:       r.PhotoSHA256,
			SignatureSha256:   r.SignatureSHA256,
			DeliveredSeconds:  seconds,
			DeliveredNanos:    nanos,
		})
	case notary.ReviewRecord:
		seconds, nanos, err := timeKey(op, r.CreatedAt)
		if err != nil {
			return err
		}

		c.reviews = append(c.reviews, storegen.InsertReviewRecordsParams{
			TransactionID:       transaction,
			Sequence:            sequence,
			ReviewID:            r.ReviewID,
			Kind:                int64(r.Kind),
			OrderID:             r.OrderID,
			SubjectID:           r.SubjectID,
			AuthorPrincipalID:   r.AuthorPrincipalID,
			MerchantPrincipalID: r.MerchantPrincipalID,
			Rating:              int64(r.Rating),
			CommentSha256:       r.CommentSHA256,
			CreatedSeconds:      seconds,
			CreatedNanos:        nanos,
		})
	default:
		return notary.Errf(notary.KindInvalidArgument, op, nil, "no table keeps a %T", record)
	}

	return nil
}

func (c *collectedRows) insert(ctx context.Context, op string, q *storegen.Queries) error {
	counts := []struct {
		table  string
		wanted int
		copy   func() (int64, error)
	}{
		{"evidence_leaves", len(c.leaves), func() (int64, error) { return q.InsertEvidenceLeaves(ctx, c.leaves) }},
		{
			"lifecycle_records",
			len(c.lifecycle),
			func() (int64, error) { return q.InsertLifecycleRecords(ctx, c.lifecycle) },
		},
		{
			"delivery_records",
			len(c.deliveries),
			func() (int64, error) { return q.InsertDeliveryRecords(ctx, c.deliveries) },
		},
		{"review_records", len(c.reviews), func() (int64, error) { return q.InsertReviewRecords(ctx, c.reviews) }},
	}

	for _, count := range counts {
		if count.wanted == 0 {
			continue
		}

		copied, err := count.copy()
		if err != nil {
			return err
		}

		if copied != int64(count.wanted) {
			return notary.Errf(
				notary.KindStoreUnavailable,
				op,
				nil,
				"wrote %d of %d rows into %s",
				copied,
				count.wanted,
				count.table,
			)
		}
	}

	return nil
}
