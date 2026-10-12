package notary_test

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAHolderIsAPrincipalOrAPurposeAndNothingLooser(t *testing.T) {
	for _, valid := range []string{"9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7", "UNASSIGNED_DRIVER_SHIPPING_FEE", strings.Repeat("x", 255)} {
		holder, err := notary.NewHolder(valid)
		require.NoError(t, err, valid)
		require.Equal(t, valid, holder.String())
	}

	for _, invalid := range []string{"", strings.Repeat("x", 256), "two words", " leading", "tab\there", "nul\x00", "\xff\xfe"} {
		_, err := notary.NewHolder(invalid)
		require.Error(t, err, "%q", invalid)
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
	}

	require.True(t, notary.Holder{}.IsZero())
}

func TestAccountKindsKeepPaymentsNumbers(t *testing.T) {
	names := map[uint8]string{
		1: "customer_balance",
		2: "merchant_available",
		3: "merchant_pending_escrow",
		4: "driver_available",
		5: "driver_pending_escrow",
		6: "suspense_available",
		7: "suspense_pending_escrow",
	}
	for code, name := range names {
		kind, err := notary.AccountKindOf(code)
		require.NoError(t, err)
		require.Equal(t, name, kind.String())
	}

	for _, code := range []uint8{0, 8, 255} {
		_, err := notary.AccountKindOf(code)
		require.Error(t, err)
		require.Equal(t, "invalid", notary.AccountKind(code).String())
	}
}

func TestASnapshotRefusesWhatCannotBeCommitted(t *testing.T) {
	asOf := notarytest.Unix(t, 1760000000)
	holder := notarytest.Holder(t, "p-1")
	one := notary.Liability{Holder: holder, Kind: notary.AccountCustomerBalance, Amount: 5}

	cases := map[string]struct {
		asOf        notary.UnixTime
		currency    string
		liabilities []notary.Liability
		kind        notary.Kind
	}{
		"no time":     {currency: "IDR", liabilities: []notary.Liability{one}, kind: notary.KindInvalidArgument},
		"lowercase":   {asOf: asOf, currency: "idr", kind: notary.KindInvalidArgument},
		"two letters": {asOf: asOf, currency: "ID", kind: notary.KindInvalidArgument},
		"same account": {
			asOf:        asOf,
			currency:    "IDR",
			liabilities: []notary.Liability{one, one},
			kind:        notary.KindInvalidArgument,
		},
		"no holder": {
			asOf:        asOf,
			currency:    "IDR",
			liabilities: []notary.Liability{{Kind: notary.AccountCustomerBalance}},
			kind:        notary.KindInvalidArgument,
		},
		"no kind": {
			asOf:        asOf,
			currency:    "IDR",
			liabilities: []notary.Liability{{Holder: holder}},
			kind:        notary.KindInvalidArgument,
		},
		"total past int64": {
			asOf:     asOf,
			currency: "IDR",
			liabilities: []notary.Liability{
				{Holder: holder, Kind: 1, Amount: math.MaxInt64},
				{Holder: holder, Kind: 2, Amount: 1},
			},
			kind: notary.KindOverflow,
		},
		"one past int64": {
			asOf:        asOf,
			currency:    "IDR",
			liabilities: []notary.Liability{{Holder: holder, Kind: 1, Amount: math.MaxInt64 + 1}},
			kind:        notary.KindOverflow,
		},
	}
	for name, c := range cases {
		_, err := notary.NewSnapshot(c.asOf, c.currency, c.liabilities)
		require.Error(t, err, name)
		require.Equal(t, c.kind, notary.KindOf(err), name)
	}
}

func TestASnapshotTotalsAndKeepsItsOwnCopy(t *testing.T) {
	holder := notarytest.Holder(t, "p-1")
	liabilities := []notary.Liability{
		{Holder: holder, Kind: notary.AccountMerchantAvailable, Amount: math.MaxInt64 - 10},
		{Holder: holder, Kind: notary.AccountMerchantPendingEscrow, Amount: 10},
		{Holder: notarytest.Holder(t, "p-2"), Kind: notary.AccountMerchantAvailable, Amount: 0},
	}

	snapshot := notarytest.Snapshot(t, 1760000000, liabilities)
	liabilities[0].Amount = 1
	listed := snapshot.Liabilities()
	listed[1].Amount = 99

	require.Equal(t, uint64(math.MaxInt64), snapshot.Total())
	require.Equal(t, uint64(math.MaxInt64-10), snapshot.Liabilities()[0].Amount)
	require.Equal(t, uint64(10), snapshot.Liabilities()[1].Amount)
	require.Equal(t, "IDR", snapshot.Currency())
}
