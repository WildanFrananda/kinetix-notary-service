package notary_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func TestADigestRoundTripsThroughHex(t *testing.T) {
	text := "0x9f4e54b0d366f7dfbe07c3641e4137ff738e612d4df1d01fd35dac7622e7915b"

	d, err := notary.ParseDigest(text)

	require.NoError(t, err)
	require.Equal(t, text, d.Hex())
	require.False(t, d.IsZero())
	require.True(t, notary.Digest{}.IsZero())
}

func TestMalformedDigestsAreRefused(t *testing.T) {
	for _, text := range []string{
		"",
		"9f4e54b0d366f7dfbe07c3641e4137ff738e612d4df1d01fd35dac7622e7915b",
		"0x9f4e",
		"0x9f4e54b0d366f7dfbe07c3641e4137ff738e612d4df1d01fd35dac7622e7915bff",
		"0xzz4e54b0d366f7dfbe07c3641e4137ff738e612d4df1d01fd35dac7622e7915b",
	} {
		_, err := notary.ParseDigest(text)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), "%q", text)
	}
}

func TestAnAddressAcceptsEitherCase(t *testing.T) {
	lower, err := notary.ParseAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3")
	require.NoError(t, err)
	mixed, err := notary.ParseAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	require.NoError(t, err)

	require.Equal(t, lower, mixed)
	require.Equal(t, "0x5fbdb2315678afecb367f032d93f642f64180aa3", mixed.Hex())

	_, err = notary.ParseAddress("0x5fbdb2315678afecb367f032d93f642f64180a")
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestASnapshotTimeIsAWholeSecondAfterTheEpoch(t *testing.T) {
	at := time.Date(2026, 10, 11, 9, 0, 0, 0, time.UTC)

	u, err := notary.UnixTimeOf(at)

	require.NoError(t, err)
	require.Equal(t, uint64(at.Unix()), u.Seconds())
	require.True(t, at.Equal(u.Time()))
	_, err = notary.UnixTimeOf(at.Add(time.Millisecond))
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	_, err = notary.UnixTimeOf(time.Unix(0, 0))
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	_, err = notary.UnixTimeOf(time.Unix(-5, 0))
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestSecondsFromTheChainMustFitATime(t *testing.T) {
	_, err := notary.UnixTimeFromSeconds(0)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	_, err = notary.UnixTimeFromSeconds(math.MaxInt64 + 1)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	u, err := notary.UnixTimeFromSeconds(math.MaxInt64)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt64), u.Time().Unix())
	require.True(t, notary.UnixTime{}.Time().IsZero())
}

func TestAKindSurvivesWrapping(t *testing.T) {
	inner := notary.Errf(notary.KindLedgerUnavailable, "ethereum.Epoch", errors.New("dial tcp: refused"), "")

	wrapped := notary.Errf(notary.KindLedgerUnavailable, "verifying.Verify", inner, "epoch %d", 3)

	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(wrapped))
	require.Equal(t, notary.KindUnknown, notary.KindOf(errors.New("plain")))
	require.Contains(t, wrapped.Error(), "dial tcp: refused")
}
