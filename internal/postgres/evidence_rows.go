package postgres

import (
	"math"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

func streamKey(op string, stream notary.EvidenceStream) (int16, error) {
	if !stream.Valid() {
		return 0, notary.Errf(notary.KindInvalidArgument, op, nil, "stream %d is not an evidence stream", stream)
	}

	return int16(stream), nil
}

func streamOf(op string, key int16) (notary.EvidenceStream, error) {
	if key <= 0 || key > math.MaxUint8 {
		return 0, notary.Errf(notary.KindStoreUnavailable, op, nil, "the database holds stream %d", key)
	}

	stream := notary.EvidenceStream(uint8(key))
	if !stream.Valid() {
		return 0, notary.Errf(notary.KindStoreUnavailable, op, nil, "the database holds stream %d", key)
	}

	return stream, nil
}

func positionKey(op string, position notary.LogPosition) (int64, int64, error) {
	transaction, err := signed(op, position.Transaction)
	if err != nil {
		return 0, 0, err
	}

	sequence, err := signed(op, position.Sequence)
	if err != nil {
		return 0, 0, err
	}

	return transaction, sequence, nil
}

func positionOf(op string, transaction, sequence int64) (notary.LogPosition, error) {
	t, err := unsigned(op, transaction)
	if err != nil {
		return notary.LogPosition{}, err
	}

	s, err := unsigned(op, sequence)
	if err != nil {
		return notary.LogPosition{}, err
	}

	return notary.LogPosition{Transaction: t, Sequence: s}, nil
}

func word(op string, value int64) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the database holds %d where a 32-bit word belongs",
			value,
		)
	}

	return uint32(value), nil
}

func words(op string, values ...int64) ([]uint32, error) {
	out := make([]uint32, len(values))

	for i, value := range values {
		w, err := word(op, value)
		if err != nil {
			return nil, err
		}

		out[i] = w
	}

	return out, nil
}

func evidenceTimeOf(op string, seconds int64, nanos int32) (notary.EvidenceTime, error) {
	t, err := notary.EvidenceTimeFrom(seconds, nanos)
	if err != nil {
		return notary.EvidenceTime{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	return t, nil
}

func timeKey(op string, t notary.EvidenceTime) (int64, int32, error) {
	seconds, err := signed(op, t.Seconds())
	if err != nil {
		return 0, 0, err
	}

	nanos := t.Nanos()
	if t.IsZero() || nanos > math.MaxInt32 {
		return 0, 0, notary.Errf(notary.KindInvalidArgument, op, nil, "a record's time is not after the Unix epoch")
	}

	return seconds, int32(nanos), nil
}

func optionalHash(op string, raw []byte) ([]byte, error) {
	if raw == nil {
		return nil, nil
	}

	if len(raw) != len(notary.Digest{}) {
		return nil, notary.Errf(notary.KindStoreUnavailable, op, nil, "the database holds a %d-byte hash", len(raw))
	}

	return raw, nil
}

func storedOf(
	op string,
	record notary.EvidenceRecord,
	salt, recordDigest, leaf []byte,
	batch, index *int64,
) (notary.StoredEvidence, error) {
	digests := make([]notary.Digest, 0, 3)

	for _, raw := range [][]byte{salt, recordDigest, leaf} {
		d, err := digestOf(op, raw)
		if err != nil {
			return notary.StoredEvidence{}, err
		}

		digests = append(digests, d)
	}

	stored := notary.StoredEvidence{
		Entry: notary.EvidenceEntry{
			Record:       record,
			Salt:         notary.Salt(digests[0]),
			RecordDigest: digests[1],
			Leaf:         digests[2],
		},
	}

	if (batch == nil) != (index == nil) {
		return notary.StoredEvidence{}, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"a %s record is in a batch without a place in it",
			record.Stream(),
		)
	}

	if batch == nil {
		return stored, nil
	}

	number, err := unsigned(op, *batch)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	place, err := unsigned(op, *index)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	stored.Batched, stored.Batch, stored.Index = true, number, place

	return stored, nil
}

func lifecycleOf(op string, row storegen.UnbatchedLifecycleRecordsRow) (notary.StoredEvidence, error) {
	at, err := positionOf(op, row.TransactionID, row.Sequence)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	w, err := words(op, row.Kind, row.FromStatus, row.ToStatus, row.ReturnStatus, row.ActorKind)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	occurredAt, err := evidenceTimeOf(op, row.OccurredSeconds, row.OccurredNanos)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	record := notary.LifecycleRecord{
		At:                  at,
		OrderNumber:         row.OrderNumber,
		Kind:                w[0],
		FromStatus:          w[1],
		ToStatus:            w[2],
		ReturnNumber:        row.ReturnNumber,
		ReturnStatus:        w[3],
		ActorKind:           w[4],
		ActorPrincipalID:    row.ActorPrincipalID,
		BuyerPrincipalID:    row.BuyerPrincipalID,
		MerchantPrincipalID: row.MerchantPrincipalID,
		OccurredAt:          occurredAt,
	}

	return storedOf(op, record, row.Salt, row.RecordDigest, row.Leaf, row.BatchNumber, row.LeafIndex)
}

func deliveryOf(op string, row storegen.UnbatchedDeliveryRecordsRow) (notary.StoredEvidence, error) {
	at, err := positionOf(op, row.TransactionID, row.Sequence)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	photo, err := digestOf(op, row.PhotoSha256)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	signature, err := optionalHash(op, row.SignatureSha256)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	deliveredAt, err := evidenceTimeOf(op, row.DeliveredSeconds, row.DeliveredNanos)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	record := notary.DeliveryRecord{
		At:                at,
		OrderNumber:       row.OrderNumber,
		DriverPrincipalID: row.DriverPrincipalID,
		PhotoSHA256:       photo[:],
		SignatureSHA256:   signature,
		DeliveredAt:       deliveredAt,
	}

	return storedOf(op, record, row.Salt, row.RecordDigest, row.Leaf, row.BatchNumber, row.LeafIndex)
}

func reviewOf(op string, row storegen.UnbatchedReviewRecordsRow) (notary.StoredEvidence, error) {
	at, err := positionOf(op, row.TransactionID, row.Sequence)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	w, err := words(op, row.Kind, row.Rating)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	comment, err := optionalHash(op, row.CommentSha256)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	createdAt, err := evidenceTimeOf(op, row.CreatedSeconds, row.CreatedNanos)
	if err != nil {
		return notary.StoredEvidence{}, err
	}

	record := notary.ReviewRecord{
		At:                  at,
		ReviewID:            row.ReviewID,
		Kind:                w[0],
		OrderID:             row.OrderID,
		SubjectID:           row.SubjectID,
		AuthorPrincipalID:   row.AuthorPrincipalID,
		MerchantPrincipalID: row.MerchantPrincipalID,
		Rating:              w[1],
		CommentSHA256:       comment,
		CreatedAt:           createdAt,
	}

	return storedOf(op, record, row.Salt, row.RecordDigest, row.Leaf, row.BatchNumber, row.LeafIndex)
}

func batchOf(op string, row storegen.EvidenceBatch) (notary.EvidenceBatchRecord, error) {
	state, err := notary.ParseEpochState(row.State)
	if err != nil {
		return notary.EvidenceBatchRecord{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	stream, err := streamOf(op, row.Stream)
	if err != nil {
		return notary.EvidenceBatchRecord{}, err
	}

	if row.Scheme <= 0 || row.Scheme > math.MaxUint16 {
		return notary.EvidenceBatchRecord{}, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"the database holds scheme %d",
			row.Scheme,
		)
	}

	values := make([]uint64, 0, 5)

	for _, v := range []int64{row.Number, row.ChainID, row.LeafCount, row.ClosedAt} {
		u, err := unsigned(op, v)
		if err != nil {
			return notary.EvidenceBatchRecord{}, err
		}

		values = append(values, u)
	}

	through, err := positionOf(op, row.ThroughTransaction, row.ThroughSequence)
	if err != nil {
		return notary.EvidenceBatchRecord{}, err
	}

	closedAt, err := notary.UnixTimeFromSeconds(values[3])
	if err != nil {
		return notary.EvidenceBatchRecord{}, notary.Errf(notary.KindStoreUnavailable, op, err, "")
	}

	address, err := addressOf(op, row.RegistryAddress)
	if err != nil {
		return notary.EvidenceBatchRecord{}, err
	}

	digests := make([]notary.Digest, 0, 3)

	for _, raw := range [][]byte{row.PreviousDigest, row.Root, row.Digest} {
		d, err := digestOf(op, raw)
		if err != nil {
			return notary.EvidenceBatchRecord{}, err
		}

		digests = append(digests, d)
	}

	reason := ""
	if row.AbandonedReason != nil {
		reason = *row.AbandonedReason
	}

	return notary.EvidenceBatchRecord{
		Registry: notary.RegistryIdentity{ChainID: values[1], Address: address},
		Scheme:   uint16(row.Scheme),
		Header: notary.EvidenceBatchHeader{
			Stream:         stream,
			Number:         values[0],
			PreviousDigest: digests[0],
			Root:           digests[1],
			LeafCount:      values[2],
			Through:        through,
			ClosedAt:       closedAt,
		},
		Digest:          digests[2],
		State:           state,
		AbandonedReason: reason,
	}, nil
}

func batchAttemptOf(op string, row storegen.EvidenceAttempt) (notary.BatchAttempt, error) {
	stream, err := streamOf(op, row.Stream)
	if err != nil {
		return notary.BatchAttempt{}, err
	}

	batch, err := unsigned(op, row.BatchNumber)
	if err != nil {
		return notary.BatchAttempt{}, err
	}

	attempt, err := attemptOf(op, batch, storegen.Attempt{
		TxHash:               row.TxHash,
		EpochNumber:          row.BatchNumber,
		Nonce:                row.Nonce,
		RawTx:                row.RawTx,
		GasLimit:             row.GasLimit,
		MaxFeePerGas:         row.MaxFeePerGas,
		MaxPriorityFeePerGas: row.MaxPriorityFeePerGas,
		Outcome:              row.Outcome,
		SentAt:               row.SentAt,
		BlockNumber:          row.BlockNumber,
		GasUsed:              row.GasUsed,
		EffectiveGasPrice:    row.EffectiveGasPrice,
		L1FeeWei:             row.L1FeeWei,
		Kind:                 notary.AttemptForAnchor.String(),
	})
	if err != nil {
		return notary.BatchAttempt{}, err
	}

	return notary.BatchAttempt{
		Stream:  stream,
		Batch:   batch,
		Anchor:  attempt.Anchor,
		Outcome: attempt.Outcome,
		SentAt:  attempt.SentAt,
		Receipt: attempt.Receipt,
	}, nil
}
