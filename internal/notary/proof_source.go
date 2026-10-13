package notary

import "context"

type ProofSource interface {
	Epoch(ctx context.Context, number uint64) (EpochRecord, bool, error)
	LastFinalized(ctx context.Context) (EpochRecord, bool, error)
	Accounts(ctx context.Context, epoch uint64) ([]Account, error)
}
