package notary_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAStatementFingerprintIgnoresTheOrderOfAHoldersAccounts(t *testing.T) {
	holder := notarytest.Holder(t, "principal-1")
	wallet := notary.Liability{Holder: holder, Kind: notary.AccountCustomerBalance, Amount: 1_500}
	earnings := notary.Liability{Holder: holder, Kind: notary.AccountDriverAvailable, Amount: 70_000}

	require.Equal(
		t,
		notary.StatementFingerprint([]notary.Liability{wallet, earnings}),
		notary.StatementFingerprint([]notary.Liability{earnings, wallet}),
	)
}

func TestAStatementFingerprintChangesWithAnyAmountOrKind(t *testing.T) {
	holder := notarytest.Holder(t, "principal-1")
	base := []notary.Liability{
		{Holder: holder, Kind: notary.AccountCustomerBalance, Amount: 1_500},
		{Holder: holder, Kind: notary.AccountDriverAvailable, Amount: 70_000},
	}
	fingerprint := notary.StatementFingerprint(base)

	changedAmount := []notary.Liability{base[0], {Holder: holder, Kind: notary.AccountDriverAvailable, Amount: 70_001}}
	movedKind := []notary.Liability{base[0], {Holder: holder, Kind: notary.AccountDriverPendingEscrow, Amount: 70_000}}
	gainedAccount := append(
		[]notary.Liability{{Holder: holder, Kind: notary.AccountMerchantAvailable, Amount: 0}},
		base...,
	)

	require.NotEqual(t, fingerprint, notary.StatementFingerprint(changedAmount))
	require.NotEqual(t, fingerprint, notary.StatementFingerprint(movedKind))
	require.NotEqual(t, fingerprint, notary.StatementFingerprint(gainedAccount))
	require.NotEqual(t, fingerprint, notary.StatementFingerprint(nil))
}

func TestOnlyAPrincipalsAccountsAreHeldByAPrincipal(t *testing.T) {
	for kind := notary.AccountCustomerBalance; kind <= notary.AccountDriverPendingEscrow; kind++ {
		require.True(t, kind.HeldByPrincipal(), kind.String())
	}

	require.False(t, notary.AccountSuspenseAvailable.HeldByPrincipal())
	require.False(t, notary.AccountSuspensePendingEscrow.HeldByPrincipal())
	require.False(t, notary.AccountKind(0).HeldByPrincipal())
	require.False(t, notary.AccountKind(8).HeldByPrincipal())
}
