package notary

import "context"

type Court interface {
	Identity(ctx context.Context) (CourtIdentity, error)
	SignerOf(ctx context.Context, epoch uint64) (Address, error)
	StatementDigest(ctx context.Context, statement Statement) (Digest, error)
	RecoverSigner(ctx context.Context, statement Statement, signature Signature) (Address, error)
	FaultCount(ctx context.Context, epoch uint64) (uint64, error)
	TotalFaults(ctx context.Context) (uint64, error)
}
