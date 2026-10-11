package notary

import "context"

type Ledger interface {
	Identity(ctx context.Context) (LedgerIdentity, error)
	Epoch(ctx context.Context, number uint64) (AnchoredEpoch, error)
	ConfirmsInclusion(ctx context.Context, claim Claim) (bool, error)
}
