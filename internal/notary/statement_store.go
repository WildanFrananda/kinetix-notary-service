package notary

import (
	"context"
	"time"
)

type StatementStore interface {
	StatementProgress(ctx context.Context) (StatementProgress, bool, error)
	LastFinalized(ctx context.Context) (EpochRecord, bool, error)
	FinalizedAfter(ctx context.Context, after uint64, limit int32) ([]EpochRecord, error)
	Accounts(ctx context.Context, epoch uint64) ([]Account, error)
	StatementFingerprints(ctx context.Context, holders []Holder) (map[Holder]Digest, error)
	AdvanceStatements(
		ctx context.Context,
		progress StatementProgress,
		planned []PlannedStatement,
		at time.Time,
	) (bool, error)
	DueStatements(ctx context.Context, now time.Time, limit int32) ([]StatementDelivery, error)
	StatementSent(ctx context.Context, delivery StatementDelivery, notification string) error
	StatementRefused(ctx context.Context, delivery StatementDelivery, reason string) error
	StatementPostponed(ctx context.Context, delivery StatementDelivery, next time.Time) error
	StatementBacklog(ctx context.Context, plannedBefore time.Time) (StatementBacklog, error)
}
