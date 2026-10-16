package notary_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func moment(t *testing.T, seconds uint64) notary.UnixTime {
	t.Helper()

	at, err := notary.UnixTimeFromSeconds(seconds)
	require.NoError(t, err)

	return at
}

func TestAWindowAdmitsExactlyWhatTheEntryPointWould(t *testing.T) {
	window := notary.SponsorshipWindow{After: moment(t, 100), Until: moment(t, 200)}

	for _, c := range []struct {
		now  uint64
		kind notary.Kind
	}{
		{99, notary.KindNotDue},
		{100, notary.KindNotDue},
		{101, notary.KindUnknown},
		{200, notary.KindUnknown},
		{201, notary.KindCallRefused},
	} {
		err := window.Admit(moment(t, c.now))

		if c.kind == notary.KindUnknown {
			require.NoError(t, err, c.now)
			continue
		}

		require.Equal(t, c.kind, notary.KindOf(err), c.now)
	}
}

func TestAnUnboundedWindowAdmitsAnyTime(t *testing.T) {
	require.NoError(t, notary.SponsorshipWindow{}.Admit(moment(t, 1)))
	require.NoError(t, notary.SponsorshipWindow{After: moment(t, 5)}.Admit(moment(t, 6)))
	require.NoError(t, notary.SponsorshipWindow{Until: moment(t, 5)}.Admit(moment(t, 5)))
}

func TestTheTwoNewVerdictsHaveNames(t *testing.T) {
	require.Equal(t, "not_due", notary.KindNotDue.String())
	require.Equal(t, "call_refused", notary.KindCallRefused.String())
}

func TestALockedStakeIsZeroOnceUnlocked(t *testing.T) {
	require.Equal(t, 0, notary.PaymasterFunds{}.LockedStake().Sign())
}
