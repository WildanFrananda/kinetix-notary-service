package sponsoring_test

import (
	"context"
	"io"
	"log/slog"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type observed struct {
	deposit *big.Int
	stake   *big.Int
	current *bool
}

func (o *observed) PaymasterDeposit(wei *big.Int) {
	o.deposit = wei
}

func (o *observed) PaymasterStake(wei *big.Int) {
	o.stake = wei
}

func (o *observed) PaymasterCurrent(current bool) {
	o.current = &current
}

func notaryCourts() notary.Generation {
	return notary.Generation{ChainID: identity.ChainID, Court: identity.Court, Withholding: identity.Withholding}
}

func watcherFor(paymaster notary.PaymasterIdentity, signer notary.Address, seen *observed) *sponsoring.Watcher {
	return watcherWith(
		paymaster,
		signer,
		notary.PaymasterFunds{Deposit: big.NewInt(42), Stake: big.NewInt(7), Staked: true},
		seen,
	)
}

func watcherWith(
	paymaster notary.PaymasterIdentity,
	signer notary.Address,
	funds notary.PaymasterFunds,
	seen *observed,
) *sponsoring.Watcher {
	return sponsoring.NewWatcher(
		&notarytest.Paymaster{IdentityValue: paymaster, FundsValue: funds},
		&notarytest.Court{Signers: map[uint64]notary.Address{math.MaxUint64: signer}},
		notaryCourts(),
		seen,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		time.Hour,
	)
}

func TestThePaymasterMustSponsorTheseCourtsAndTheirCurrentSigner(t *testing.T) {
	require.NoError(t, watcherFor(identity, identity.Signer, &observed{}).Check(context.Background()))

	for name, change := range map[string]func(*notary.PaymasterIdentity){
		"another chain":             func(p *notary.PaymasterIdentity) { p.ChainID = 1 },
		"another court":             func(p *notary.PaymasterIdentity) { p.Court = notary.Address{0xc1} },
		"another withholding court": func(p *notary.PaymasterIdentity) { p.Withholding = notary.Address{0xd1} },
		"a retired signer":          func(p *notary.PaymasterIdentity) { p.Signer = notary.Address{0x52} },
	} {
		t.Run(name, func(t *testing.T) {
			other := identity
			change(&other)

			err := watcherFor(other, identity.Signer, &observed{}).Check(context.Background())

			require.Equal(t, notary.KindNotSponsored, notary.KindOf(err))
		})
	}
}

func TestEachLookReportsTheDepositTheStakeAndWhetherTheHoldersAreSponsored(t *testing.T) {
	seen := &observed{}
	watcherFor(identity, identity.Signer, seen).Look(context.Background())

	require.Equal(t, big.NewInt(42), seen.deposit)
	require.Equal(t, big.NewInt(7), seen.stake)
	require.True(t, *seen.current)

	unlocking := &observed{}
	watcherWith(
		identity,
		identity.Signer,
		notary.PaymasterFunds{Deposit: big.NewInt(42), Stake: big.NewInt(7), Staked: false},
		unlocking,
	).Look(context.Background())

	require.Equal(t, 0, unlocking.stake.Sign())

	stale := &observed{}
	watcherFor(identity, notary.Address{0x52}, stale).Look(context.Background())

	require.False(t, *stale.current)
}
