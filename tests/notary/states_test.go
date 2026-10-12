package notary_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func TestAnEpochMovesOnlyAlongTheStateMachine(t *testing.T) {
	allowed := map[notary.EpochState][]notary.EpochState{
		notary.StateBuilt:     {notary.StateSubmitted, notary.StateAbandoned},
		notary.StateSubmitted: {notary.StateMined, notary.StateFinalized, notary.StateAbandoned},
		notary.StateMined:     {notary.StateFinalized, notary.StateSubmitted, notary.StateAbandoned},
	}
	for from := notary.StateBuilt; from <= notary.StateAbandoned; from++ {
		for to := notary.StateBuilt; to <= notary.StateAbandoned; to++ {
			want := false
			for _, next := range allowed[from] {
				want = want || next == to
			}

			require.Equal(t, want, from.CanBecome(to), "%s -> %s", from, to)
		}

		parsed, err := notary.ParseEpochState(from.String())
		require.NoError(t, err)
		require.Equal(t, from, parsed)
	}

	require.True(t, notary.StateMined.InFlight())
	require.False(t, notary.StateFinalized.InFlight())
	require.False(t, notary.StateAbandoned.InFlight())

	_, err := notary.ParseEpochState("planned")
	require.Error(t, err)
}

func TestAttemptOutcomesRoundTrip(t *testing.T) {
	for outcome := notary.AttemptPending; outcome <= notary.AttemptReverted; outcome++ {
		parsed, err := notary.ParseAttemptOutcome(outcome.String())
		require.NoError(t, err)
		require.Equal(t, outcome, parsed)
	}

	_, err := notary.ParseAttemptOutcome("dropped")
	require.Error(t, err)
}

func TestFeesAreQuotedUnderTheCeiling(t *testing.T) {
	quote, err := notary.QuoteFees(100, 10, 1_000)
	require.NoError(t, err)
	require.Equal(t, notary.FeeQuote{MaxFeePerGas: 210, MaxPriorityFeePerGas: 10}, quote)

	quote, err = notary.QuoteFees(100, 10, 150)
	require.NoError(t, err)
	require.Equal(t, notary.FeeQuote{MaxFeePerGas: 150, MaxPriorityFeePerGas: 10}, quote)

	quote, err = notary.QuoteFees(100, 80, 150)
	require.NoError(t, err)
	require.Equal(t, notary.FeeQuote{MaxFeePerGas: 150, MaxPriorityFeePerGas: 50}, quote)

	_, err = notary.QuoteFees(151, 1, 150)
	require.Equal(t, notary.KindAnchorDeferred, notary.KindOf(err))

	_, err = notary.QuoteFees(1, 1, 0)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	quote, err = notary.QuoteFees(math.MaxUint64/2, math.MaxUint64/2, math.MaxUint64)
	require.NoError(t, err)
	require.Equal(t, uint64(math.MaxUint64), quote.MaxFeePerGas)
}

func TestAReplacementOutbidsByMoreThanTheNodesTenPercent(t *testing.T) {
	quote := notary.FeeQuote{MaxFeePerGas: 1_000, MaxPriorityFeePerGas: 7}

	next, raised := quote.Bumped(10_000)
	require.True(t, raised)
	require.GreaterOrEqual(t, next.MaxFeePerGas*10, quote.MaxFeePerGas*11)
	require.GreaterOrEqual(t, next.MaxPriorityFeePerGas*10, quote.MaxPriorityFeePerGas*11)

	same, raised := quote.Bumped(1_100)
	require.False(t, raised)
	require.Equal(t, quote, same)

	zeroTip, raised := notary.FeeQuote{MaxFeePerGas: 10}.Bumped(100)
	require.True(t, raised)
	require.Positive(t, zeroTip.MaxPriorityFeePerGas)

	top, _ := notary.FeeQuote{MaxFeePerGas: math.MaxUint64, MaxPriorityFeePerGas: math.MaxUint64}.Bumped(math.MaxUint64)
	require.Equal(t, uint64(math.MaxUint64), top.MaxFeePerGas)
}

func TestRunwayIsBalanceOverCostTimesTheInterval(t *testing.T) {
	runway, known := notary.Runway(big.NewInt(1_000), 100, time.Hour)
	require.True(t, known)
	require.Equal(t, 10*time.Hour, runway)

	runway, known = notary.Runway(big.NewInt(99), 100, time.Hour)
	require.True(t, known)
	require.Zero(t, runway)

	_, known = notary.Runway(big.NewInt(1_000), 0, time.Hour)
	require.False(t, known)
	_, known = notary.Runway(nil, 1, time.Hour)
	require.False(t, known)

	huge := new(big.Int).Lsh(big.NewInt(1), 200)
	runway, known = notary.Runway(huge, 1, time.Hour)
	require.True(t, known)
	require.Equal(t, time.Duration(math.MaxInt64), runway)
}

func TestAReceiptCostsExecutionPlusTheL1DataFee(t *testing.T) {
	cost, err := notary.AnchorReceipt{GasUsed: 144_000, EffectiveGasPrice: 1_000_000, L1FeeWei: 5}.CostWei()
	require.NoError(t, err)
	require.Equal(t, uint64(144_000_000_005), cost)

	_, err = notary.AnchorReceipt{GasUsed: math.MaxUint64, EffectiveGasPrice: 2}.CostWei()
	require.Equal(t, notary.KindOverflow, notary.KindOf(err))
	_, err = notary.AnchorReceipt{GasUsed: 1, EffectiveGasPrice: math.MaxUint64, L1FeeWei: 1}.CostWei()
	require.Equal(t, notary.KindOverflow, notary.KindOf(err))
}

func TestALedgerAtGenesisHasNothingAnchored(t *testing.T) {
	require.True(
		t,
		notary.LedgerHead{
			NextEpoch:     4,
			GenesisEpoch:  4,
			Head:          notary.Digest{1},
			GenesisDigest: notary.Digest{1},
		}.AtGenesis(),
	)
	require.False(
		t,
		notary.LedgerHead{
			NextEpoch:     5,
			GenesisEpoch:  4,
			Head:          notary.Digest{1},
			GenesisDigest: notary.Digest{1},
		}.AtGenesis(),
	)
}
