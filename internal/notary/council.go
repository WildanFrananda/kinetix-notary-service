package notary

import "context"

type Council interface {
	Identity(ctx context.Context) (CouncilIdentity, error)
	IsMember(ctx context.Context, member Address) (bool, error)
	Threshold(ctx context.Context) (uint64, error)
	ApprovalDigest(ctx context.Context, number uint64, epoch Digest) (Digest, error)
	AnchorsLedger(ctx context.Context) (bool, error)
}
