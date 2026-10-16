package sponsoring_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var identity = notary.PaymasterIdentity{
	ChainID:     84532,
	Address:     notary.Address{0xfa},
	EntryPoint:  notary.Address{0x43},
	Account:     notary.Address{0x77},
	Court:       notary.Address{0xc0},
	Withholding: notary.Address{0xd0},
	Signer:      notary.Address{0x51},
	Ceiling:     big.NewInt(100_000_000_000_000),
}

type rig struct {
	chain     *notarytest.SponsorChain
	paymaster *notarytest.Paymaster
	bundler   *notarytest.Bundler
	signer    *ethereum.HolderSigner
	sponsor   *sponsoring.Sponsor
}

func newRig(t *testing.T) *rig {
	t.Helper()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	scheme := keccak.NewScheme()
	entryPoint := notary.EntryPointIdentity{ChainID: identity.ChainID, Address: identity.EntryPoint}
	r := &rig{
		chain: &notarytest.SponsorChain{
			AccountNonces:   3,
			OperationNonces: 9,
			MaxFee:          7_700_000,
			Priority:        1_100_000,
			Now:             at(t, 1_800_000_000),
		},
		paymaster: &notarytest.Paymaster{IdentityValue: identity},
		signer:    ethereum.NewHolderSigner(key),
		bundler: &notarytest.Bundler{
			Gas: notary.OperationGas{
				PreVerificationGas:            90_000,
				VerificationGasLimit:          150_000,
				CallGasLimit:                  300_000,
				PaymasterVerificationGasLimit: 100_000,
			},
			Hasher: func(op notary.UserOperation) notary.Digest { return scheme.OperationHash(entryPoint, op) },
			Landed: &notary.OperationReceipt{
				Success:         true,
				TransactionHash: notary.Digest{0x7a},
				ActualGasCost:   big.NewInt(2_000_000_000),
			},
		},
	}
	r.sponsor = sponsoring.NewSponsor(
		r.chain,
		r.paymaster,
		r.bundler,
		r.signer,
		scheme,
		ethereum.NewSponsoredCalls(),
		time.Millisecond,
		50*time.Millisecond,
	)

	return r
}

func at(t *testing.T, seconds uint64) notary.UnixTime {
	t.Helper()

	moment, err := notary.UnixTimeFromSeconds(seconds)
	require.NoError(t, err)

	return moment
}

func call() notary.SponsoredCall {
	return notary.SponsoredCall{Target: identity.Withholding, Data: []byte{0x01, 0x02}, PaymasterData: []byte{0x0a}}
}

func TestAHolderWhoNeverDelegatedAuthorizesThePinnedAccountInTheSameOperation(t *testing.T) {
	r := newRig(t)

	receipt, err := r.sponsor.Submit(context.Background(), identity, call())

	require.NoError(t, err)
	require.True(t, receipt.Success)
	require.Len(t, r.bundler.Sent, 1)

	sent := r.bundler.Sent[0]
	require.Equal(t, r.signer.Holder(), sent.Sender)
	require.Equal(t, identity.Account, sent.Delegate)
	require.Equal(t, uint64(9), sent.Nonce)
	require.Equal(t, identity.Address, sent.Paymaster)
	require.Equal(t, []byte{0x0a}, sent.PaymasterData)
	require.Equal(t, uint64(300_000), sent.CallGasLimit)
	require.Equal(t, uint64(100_000), sent.PaymasterVerificationGasLimit)
	require.Equal(t, uint64(7_700_000), sent.MaxFeePerGas)

	delegation := r.bundler.Delegations[0]
	require.NotNil(t, delegation)
	require.Equal(t, identity.Account, delegation.Delegate)
	require.Equal(t, uint64(3), delegation.Nonce)
	require.Equal(t, identity.ChainID, delegation.ChainID)

	hash := keccak.NewScheme().OperationHash(
		notary.EntryPointIdentity{ChainID: identity.ChainID, Address: identity.EntryPoint},
		sent,
	)
	recovered, err := crypto.SigToPub(hash[:], append(append([]byte{}, sent.Signature[:64]...), sent.Signature[64]-27))
	require.NoError(t, err)
	require.Equal(t, r.signer.Holder(), notary.Address(crypto.PubkeyToAddress(*recovered)))
	require.Equal(t, hash, receipt.Hash)
}

func TestAHolderAlreadyDelegatedToThePinnedAccountSendsNoAuthorization(t *testing.T) {
	r := newRig(t)
	r.chain.Delegate = identity.Account

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.NoError(t, err)
	require.True(t, r.bundler.Sent[0].Delegate.IsZero())
	require.Nil(t, r.bundler.Delegations[0])
}

func TestAHolderDelegatedToAnyOtherAccountIsNeverSent(t *testing.T) {
	r := newRig(t)
	r.chain.Delegate = notary.Address{0x99}

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindNotSponsored, notary.KindOf(err))
	require.Empty(t, r.bundler.Estimated)
	require.Empty(t, r.bundler.Sent)
}

func TestAnOperationThatMayCostMoreThanTheCeilingIsNeverSent(t *testing.T) {
	r := newRig(t)
	r.chain.MaxFee = 1_000_000_000_000

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindNotSponsored, notary.KindOf(err))
	require.ErrorContains(t, err, "ceiling")
	require.Empty(t, r.bundler.Sent)
}

func TestABundlerThatHashesTheOperationDifferentlyIsADisagreement(t *testing.T) {
	r := newRig(t)
	other := notary.Digest{0xee}
	r.bundler.Answer = &other

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestAnOperationWithoutAReceiptInTimeIsUnknownNotFailed(t *testing.T) {
	r := newRig(t)
	r.bundler.Landed = nil

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindSourceUnavailable, notary.KindOf(err))
	require.ErrorContains(t, err, "may still land")
}

func TestTheBundlersRefusalIsPassedOnAsItIs(t *testing.T) {
	r := newRig(t)
	r.bundler.EstimateErr = notary.Errf(notary.KindNotSponsored, "test", nil, "AA33 reverted")

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindNotSponsored, notary.KindOf(err))
	require.Empty(t, r.bundler.Sent)
}

func TestTheHolderAsksThePaymasterBeforeSigningAnything(t *testing.T) {
	r := newRig(t)

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.NoError(t, err)
	require.Equal(t, []notary.Address{r.signer.Holder()}, r.paymaster.Asked)
}

func TestACallTheCourtWouldRefuseIsNeverSignedOrSent(t *testing.T) {
	r := newRig(t)
	r.paymaster.Refusal = notary.Errf(notary.KindCallRefused, "test", nil, "already requested")

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindCallRefused, notary.KindOf(err))
	require.Empty(t, r.bundler.Estimated)
	require.Empty(t, r.bundler.Sent)
}

func TestACallBeforeItsTimeIsNotDueAndOneAfterItsWindowIsRefused(t *testing.T) {
	r := newRig(t)
	r.paymaster.Window = notary.SponsorshipWindow{After: at(t, 1_800_000_000)}

	_, err := r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindNotDue, notary.KindOf(err))
	require.ErrorContains(t, err, "2027-01-15T08:00:00Z")
	require.Empty(t, r.bundler.Estimated)

	r.paymaster.Window = notary.SponsorshipWindow{Until: at(t, 1_799_999_999)}

	_, err = r.sponsor.Submit(context.Background(), identity, call())

	require.Equal(t, notary.KindCallRefused, notary.KindOf(err))
	require.Empty(t, r.bundler.Estimated)
}
