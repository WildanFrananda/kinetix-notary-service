package notarytest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func Unix(t testing.TB, seconds uint64) notary.UnixTime {
	t.Helper()

	u, err := notary.UnixTimeFromSeconds(seconds)
	require.NoError(t, err)

	return u
}
