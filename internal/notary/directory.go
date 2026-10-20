package notary

import "context"

type Directory interface {
	Identity(ctx context.Context) (DirectoryIdentity, error)
	Current(ctx context.Context) (Generation, error)
	GenerationOf(ctx context.Context, epoch uint64) (Generation, error)
	Lineage(ctx context.Context) ([]Generation, error)
}
