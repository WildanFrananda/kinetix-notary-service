package keccak_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAnAccountIDMatchesCast(t *testing.T) {
	got := keccak.NewScheme().AccountID(
		notarytest.Holder(t, "9f1d4a3e-1c62-4d0a-9a7b-2f5c8e0b41d7"),
		notary.AccountCustomerBalance,
		notary.AccountSecret(repeated(0x77)),
	)

	require.Equal(t, digestOf(t, "0xa5618bc093caf7474bbb829f0edec7502ef3b5cc7aad3df1d9368334ffff79dc"), got)
}

func TestASecondDerivationCommitmentMatchesCast(t *testing.T) {
	got := keccak.NewScheme().AccountCommitment(repeated(0x11), notary.Salt(repeated(0x22)))

	require.Equal(t, digestOf(t, "0xaf293031243d180a098189fc0110b5fadfb2ae95851c586dddfd1011ea27d629"), got)
}

func TestEveryPartOfAnAccountIDChangesIt(t *testing.T) {
	scheme := keccak.NewScheme()
	holder := notarytest.Holder(t, "p-1")
	base := scheme.AccountID(holder, notary.AccountMerchantAvailable, notary.AccountSecret(repeated(1)))

	require.NotEqual(
		t,
		base,
		scheme.AccountID(
			notarytest.Holder(t, "p-2"),
			notary.AccountMerchantAvailable,
			notary.AccountSecret(repeated(1)),
		),
	)
	require.NotEqual(
		t,
		base,
		scheme.AccountID(holder, notary.AccountMerchantPendingEscrow, notary.AccountSecret(repeated(1))),
	)
	require.NotEqual(
		t,
		base,
		scheme.AccountID(holder, notary.AccountMerchantAvailable, notary.AccountSecret(repeated(2))),
	)
}

func TestTheTwoDerivationsDifferForTheSameInputs(t *testing.T) {
	scheme := keccak.NewScheme()
	holder := notarytest.Holder(t, "p-1")
	salt := notary.Salt(repeated(9))

	first := scheme.Commitment(holder, notary.AccountCustomerBalance, salt)
	account := scheme.AccountID(holder, notary.AccountCustomerBalance, notary.AccountSecret(salt))

	require.NotEqual(t, first, account)
	require.NotEqual(t, first, scheme.AccountCommitment(account, salt))
}

func TestARegistrationDigestDependsOnEveryFieldAndTheDomain(t *testing.T) {
	scheme := keccak.NewScheme()
	court := notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0x01}}
	registration := notary.Registration{Holder: notary.Address{0x02}, Account: repeated(0x03), FromEpoch: 4}
	base := scheme.RegistrationDigest(court, registration)

	changed := []notary.Digest{
		scheme.RegistrationDigest(notary.CourtIdentity{ChainID: 1, Address: court.Address}, registration),
		scheme.RegistrationDigest(
			notary.CourtIdentity{ChainID: court.ChainID, Address: notary.Address{0x09}},
			registration,
		),
		scheme.RegistrationDigest(
			court,
			notary.Registration{Holder: notary.Address{0x09}, Account: registration.Account, FromEpoch: 4},
		),
		scheme.RegistrationDigest(
			court,
			notary.Registration{Holder: registration.Holder, Account: repeated(0x09), FromEpoch: 4},
		),
		scheme.RegistrationDigest(
			court,
			notary.Registration{Holder: registration.Holder, Account: registration.Account, FromEpoch: 5},
		),
	}
	for _, other := range changed {
		require.NotEqual(t, base, other)
	}

	statementLike := scheme.StatementDigest(court, notary.Statement{Ledger: registration.Holder})
	require.NotEqual(t, base, statementLike)
}
