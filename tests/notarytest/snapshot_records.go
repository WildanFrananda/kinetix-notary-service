package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type SnapshotRecords struct {
	Kept map[notary.SnapshotID]notary.SnapshotRecord
	Err  error
}

var _ notary.SnapshotRecords = (*SnapshotRecords)(nil)

func (s *SnapshotRecords) Record(_ context.Context, id notary.SnapshotID) (notary.SnapshotRecord, error) {
	if s.Err != nil {
		return notary.SnapshotRecord{}, s.Err
	}

	record, ok := s.Kept[id]
	if !ok {
		return notary.SnapshotRecord{}, notary.Errf(
			notary.KindSnapshotUnknown,
			"notarytest.SnapshotRecords",
			nil,
			"%s",
			id,
		)
	}

	return record, nil
}
