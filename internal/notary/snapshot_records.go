package notary

import "context"

type SnapshotRecords interface {
	Record(ctx context.Context, id SnapshotID) (SnapshotRecord, error)
}
