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

	return snapshot
}
