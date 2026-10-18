package notary

import (
	"context"
	"time"
)

type EvidenceStore interface {
	Lease(ctx context.Context, instance string, ttl time.Duration) (Lease, bool, error)
	Cursor(ctx context.Context, stream EvidenceStream) (LogPosition, error)
	Collect(
		ctx context.Context,
		lease Lease,
		stream EvidenceStream,
		from LogPosition,
		entries []EvidenceEntry,
		to LogPosition,
	) error
	Unbatched(ctx context.Context, stream EvidenceStream, limit int32) ([]EvidenceEntry, error)
	UnbatchedCount(ctx context.Context, stream EvidenceStream) (int64, error)
	InFlightBatch(ctx context.Context) (EvidenceBatchRecord, bool, error)
	LastBatch(ctx context.Context, stream EvidenceStream) (EvidenceBatchRecord, bool, error)
	LastFinalizedBatch(ctx context.Context, stream EvidenceStream) (EvidenceBatchRecord, bool, error)
	SaveBatch(
		ctx context.Context,
		lease Lease,
		registry RegistryIdentity,
		plan EvidenceBatchPlan,
	) (EvidenceBatchRecord, error)
	RecordBatchAttempt(ctx context.Context, lease Lease, attempt BatchAttempt) error
	BatchAttempts(ctx context.Context, stream EvidenceStream, batch uint64) ([]BatchAttempt, error)
	SettleBatchAttempt(
		ctx context.Context,
		lease Lease,
		hash Digest,
		outcome AttemptOutcome,
		receipt *AnchorReceipt,
	) error
	TransitionBatch(
		ctx context.Context,
		lease Lease,
		stream EvidenceStream,
		batch uint64,
		from, to EpochState,
		reason string,
	) error
	AbandonedBatches(ctx context.Context) (int64, error)
	RecentBatchCosts(ctx context.Context, limit int32) ([]uint64, error)
}
