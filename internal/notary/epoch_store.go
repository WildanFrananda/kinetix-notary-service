package notary

import (
	"context"
	"time"
)

type EpochStore interface {
	Lease(ctx context.Context, instance string, ttl time.Duration) (Lease, bool, error)
	InFlight(ctx context.Context) (EpochRecord, bool, error)
	AccountSecrets(ctx context.Context, keys []AccountKey) (map[AccountKey]AccountSecret, error)
	LastFinalized(ctx context.Context) (EpochRecord, bool, error)
	SaveBuilt(ctx context.Context, lease Lease, ledger LedgerIdentity, plan EpochPlan) (EpochRecord, error)
	RecordAttempt(ctx context.Context, lease Lease, attempt Attempt) error
	Attempts(ctx context.Context, epoch uint64) ([]Attempt, error)
	SettleAttempt(ctx context.Context, lease Lease, hash Digest, outcome AttemptOutcome, receipt *AnchorReceipt) error
	Transition(ctx context.Context, lease Lease, epoch uint64, from, to EpochState, reason string) error
	AbandonedCount(ctx context.Context) (int64, error)
	RecentCosts(ctx context.Context, limit int32) ([]uint64, error)
}
