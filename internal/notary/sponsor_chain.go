package notary

import "context"

type SponsorChain interface {
	OperationNonce(ctx context.Context, entryPoint Address, sender Address) (uint64, error)
	AccountNonce(ctx context.Context, account Address) (uint64, error)
	DelegateOf(ctx context.Context, account Address) (Address, error)
	Fees(ctx context.Context) (uint64, uint64, error)
	BlockTime(ctx context.Context) (UnixTime, error)
}
