package notary

import "context"

type Bundler interface {
	Estimate(ctx context.Context, op UserOperation, entryPoint Address, delegation *Delegation) (OperationGas, error)
	Send(ctx context.Context, op UserOperation, entryPoint Address, delegation *Delegation) (Digest, error)
	Receipt(ctx context.Context, hash Digest) (OperationReceipt, bool, error)
}
