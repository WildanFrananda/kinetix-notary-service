package notary

import "context"

type SnapshotSource interface {
	Snapshot(ctx context.Context, currency string) (Snapshot, error)
}
