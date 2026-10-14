package notarytest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func Holder(t *testing.T, value string) notary.Holder {
	t.Helper()

	holder, err := notary.NewHolder(value)
	require.NoError(t, err)

	return holder
}

func Liabilities(t *testing.T, count int) []notary.Liability {
	t.Helper()

	out := make([]notary.Liability, count)
	for i := range out {
		out[i] = notary.Liability{
			Holder: Holder(t, fmt.Sprintf("principal-%04d", i/7)),
			Kind:   notary.AccountKind(i%7 + 1),
			Amount: uint64(i) * 1000,
		}
	}

	return out
}

func Snapshot(t *testing.T, asOf uint64, liabilities []notary.Liability) notary.Snapshot {
	t.Helper()
	snapshot, err := notary.NewSnapshot(Unix(t, asOf), "IDR", liabilities)
	require.NoError(t, err)

	record := notary.SnapshotRecord{
		AsOf:     snapshot.AsOf(),
		Currency: snapshot.Currency(),
		Count:    uint64(len(liabilities)),
		Total:    snapshot.Total(),
	}
	digest, err := notary.SnapshotDigest(record, liabilities)
	require.NoError(t, err)

	recorded, err := snapshot.Recorded(SnapshotID(t, asOf), 0, digest)
	require.NoError(t, err)

	return recorded
}

func SnapshotID(t *testing.T, seed uint64) notary.SnapshotID {
	t.Helper()

	id, err := notary.NewSnapshotID(fmt.Sprintf("00000000-0000-4000-8000-%012x", seed&0xffffffffffff))
	require.NoError(t, err)

	return id
}
