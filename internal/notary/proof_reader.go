package notary

import "context"

type ProofReader interface {
	Depth(ctx context.Context) (int, error)
	IsProven(ctx context.Context, number uint64, root Digest, leafCount uint64) (bool, error)
}
