package registering_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const finalEpoch = 4

var withholdingIdentity = notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xd0, 0x0d}}

type rig struct {
	source    *notarytest.ProofSource
	court     *notarytest.Court
	chain     *notarytest.WithholdingCourt
	store     *notarytest.Registrations
	holderKey *ecdsa.PrivateKey
	holder    notary.Address
	signer    notary.Address
	now       time.Time
	registrar *registering.Registrar
}

func newRig(t *testing.T) *rig {
	t.Helper()

	scheme := keccak.NewScheme()
	snapshot := notarytest.Snapshot(t, 1760000000, notarytest.Liabilities(t, 15))
	plan, err := notary.PlanEpoch(
		scheme,
		&notarytest.Salts{Seed: 1},
		&notarytest.Secrets{Seed: 1},
		nil,
		finalEpoch,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)

	record := notary.EpochRecord{
		Ledger:   notarytest.LedgerIdentity,
		Header:   plan.Header,
		Digest:   plan.Digest,
		Currency: "IDR",
		State:    notary.StateFinalized,
	}
	signerKey, signer := notarytest.StatementKey(t, "registration-signer")
	holderKey, holder := notarytest.StatementKey(t, "registration-holder")

	r := &rig{
		source: &notarytest.ProofSource{
			Records: map[uint64]notary.EpochRecord{finalEpoch: record},
			Stored:  map[uint64][]notary.Account{finalEpoch: plan.Accounts},
		},
		court:     &notarytest.Court{IdentityValue: notarytest.CourtIdentity, Statements: scheme, Signer: signer},
		chain:     notarytest.NewWithholdingCourt(withholdingIdentity, scheme),
		store:     notarytest.NewRegistrations(),
		holderKey: holderKey,
		holder:    holder,
		signer:    signer,
		now:       time.Unix(1760000600, 0).UTC(),
	}
	r.registrar = registering.NewRegistrar(
		r.source,
		scheme,
		scheme,
		r.court,
		r.chain,
		ethereum.NewStatementSigner([]*ecdsa.PrivateKey{signerKey}),
		ethereum.NewRecoverer(),
		r.store,
		withholdingIdentity,
		func() time.Time { return r.now },
	)

	return r
}

func (r *rig) request(t *testing.T, principal string, issuedAt time.Time) notary.RegistrationRequest {
	t.Helper()

	return notary.RegistrationRequest{
		Holder:    r.holder,
		Principal: notarytest.Holder(t, principal),
		IssuedAt:  notarytest.Unix(t, uint64(issuedAt.Unix())),
	}
}

func sign(t *testing.T, key *ecdsa.PrivateKey, request notary.RegistrationRequest) notary.Signature {
	t.Helper()

	digest := keccak.NewScheme().RegistrationRequestDigest(withholdingIdentity, request)
	raw, err := crypto.Sign(digest[:], key)
	require.NoError(t, err)

	var signature notary.Signature
	copy(signature[:], raw)
	signature[64] += 27

	return signature
}

func (r *rig) register(t *testing.T, principal string) (registering.Issued, error) {
	t.Helper()

	request := r.request(t, principal, r.now.Add(-time.Minute))

	return r.registrar.Register(context.Background(), request, sign(t, r.holderKey, request))
}

func TestEachAccountOfThePrincipalGetsARegistrationTheCourtsSignerSigned(t *testing.T) {
	r := newRig(t)
	scheme := keccak.NewScheme()

	issued, err := r.register(t, "principal-0001")

	require.NoError(t, err)
	require.Equal(t, withholdingIdentity, issued.Court)
	require.Equal(t, uint64(finalEpoch), issued.Epoch)
	require.Len(t, issued.Registrations, 7)

	kinds := map[notary.AccountKind]bool{}

	for _, one := range issued.Registrations {
		require.Equal(t, r.holder, one.Registration.Holder)
		require.Equal(t, uint64(finalEpoch), one.Registration.FromEpoch)
		require.Equal(t, scheme.RegistrationDigest(withholdingIdentity, one.Registration), one.Digest)

		recovered, err := ethereum.NewRecoverer().Recover(one.Digest, one.Signature)
		require.NoError(t, err)
		require.Equal(t, r.signer, recovered)

		var account notary.Account

		for _, stored := range r.source.Stored[finalEpoch] {
			if stored.ID == one.Registration.Account {
				account = stored
			}
		}

		require.Equal(t, "principal-0001", account.Liability.Holder.String())
		require.Equal(t, account.Liability.Kind, one.Kind)

		kinds[one.Kind] = true

		require.Equal(t, one, r.store.Saved[one.Digest])
		require.Equal(t, "principal-0001", r.store.Principals[one.Digest].String())
	}

	require.Len(t, kinds, 7)
}

func TestARequestTheHoldersKeyDidNotSignIsRefused(t *testing.T) {
	r := newRig(t)
	request := r.request(t, "principal-0001", r.now)
	otherKey, _ := notarytest.StatementKey(t, "someone-else")

	_, err := r.registrar.Register(context.Background(), request, sign(t, otherKey, request))

	require.Equal(t, notary.KindRegistrationRefused, notary.KindOf(err))
	require.Empty(t, r.store.Saved)
}

func TestARequestTheHolderSignedForAnotherPrincipalIsRefused(t *testing.T) {
	r := newRig(t)
	signed := r.request(t, "principal-0001", r.now)
	asked := r.request(t, "principal-0000", r.now)

	_, err := r.registrar.Register(context.Background(), asked, sign(t, r.holderKey, signed))

	require.Equal(t, notary.KindRegistrationRefused, notary.KindOf(err))
	require.Contains(t, notary.MessageOf(err), "not by "+r.holder.Hex())
	require.Empty(t, r.store.Saved)

	_, err = r.registrar.Register(context.Background(), signed, sign(t, r.holderKey, signed))
	require.NoError(t, err)
}

func TestTheDomainIsTheWithholdingCourtsOwn(t *testing.T) {
	r := newRig(t)

	require.Equal(
		t,
		notary.TypedDataDomain{
			Name:              keccak.WithholdingCourtName,
			Version:           keccak.WithholdingCourtVersion,
			ChainID:           withholdingIdentity.ChainID,
			VerifyingContract: withholdingIdentity.Address,
		},
		r.registrar.Domain(),
	)
}

func TestARequestOutsideItsFiveMinutesIsRefused(t *testing.T) {
	r := newRig(t)

	for name, issuedAt := range map[string]time.Time{
		"older than five minutes": r.now.Add(-registering.MaxRequestAge - time.Second),
		"from the future":         r.now.Add(registering.MaxRequestSkew + time.Second),
	} {
		t.Run(name, func(t *testing.T) {
			request := r.request(t, "principal-0001", issuedAt)

			_, err := r.registrar.Register(context.Background(), request, sign(t, r.holderKey, request))

			require.Equal(t, notary.KindRegistrationRefused, notary.KindOf(err))
		})
	}

	request := r.request(t, "principal-0001", r.now.Add(-registering.MaxRequestAge))
	_, err := r.registrar.Register(context.Background(), request, sign(t, r.holderKey, request))
	require.NoError(t, err)
}

func TestAPrincipalWithoutAnAccountGetsNoRegistration(t *testing.T) {
	r := newRig(t)

	_, err := r.register(t, "principal-9999")

	require.Equal(t, notary.KindNoAccount, notary.KindOf(err))
}

func TestRegistrationWaitsForAFinalEpochWithSecrets(t *testing.T) {
	r := newRig(t)
	r.source.Records = map[uint64]notary.EpochRecord{}

	_, err := r.register(t, "principal-0001")
	require.Equal(t, notary.KindRegistrationNotOpen, notary.KindOf(err))

	r = newRig(t)
	for i := range r.source.Stored[finalEpoch] {
		r.source.Stored[finalEpoch][i].Secret = notary.AccountSecret{}
	}

	_, err = r.register(t, "principal-0001")
	require.Equal(t, notary.KindRegistrationNotOpen, notary.KindOf(err))
	require.Empty(t, r.store.Saved)
}

func TestAnAccountWhoseIDDoesNotDeriveFromItsSecretIsNeverRegistered(t *testing.T) {
	r := newRig(t)
	for i := range r.source.Stored[finalEpoch] {
		if r.source.Stored[finalEpoch][i].Liability.Holder.String() == "principal-0001" {
			r.source.Stored[finalEpoch][i].Secret[0] ^= 0xff
		}
	}

	_, err := r.register(t, "principal-0001")

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
	require.Empty(t, r.store.Saved)
}

func TestNoSignerMeansNoRegistration(t *testing.T) {
	r := newRig(t)
	r.court.Signer = notary.Address{}

	_, err := r.register(t, "principal-0001")
	require.Equal(t, notary.KindSignerUnavailable, notary.KindOf(err))

	r = newRig(t)
	r.court.Signer = notary.Address{0x99}

	_, err = r.register(t, "principal-0001")
	require.Equal(t, notary.KindSignerUnavailable, notary.KindOf(err))
	require.Empty(t, r.store.Saved)
}

func TestARegistrationThatCannotBeKeptIsNotIssued(t *testing.T) {
	r := newRig(t)
	r.store.Err = notary.Errf(notary.KindStoreUnavailable, "test", errors.New("down"), "")

	issued, err := r.register(t, "principal-0001")

	require.Equal(t, notary.KindStoreUnavailable, notary.KindOf(err))
	require.Empty(t, issued.Registrations)
}

func TestTheWithholdingCourtIsCheckedBeforeTheNotaryRegisters(t *testing.T) {
	r := newRig(t)
	require.NoError(t, r.registrar.CheckWithholding(context.Background()))

	r.chain.IdentityValue = notary.CourtIdentity{ChainID: 84532, Address: notary.Address{0xee}}
	require.Equal(t, notary.KindWrongCourt, notary.KindOf(r.registrar.CheckWithholding(context.Background())))

	r = newRig(t)
	r.chain.Disagree = true
	require.Equal(t, notary.KindDisagreement, notary.KindOf(r.registrar.CheckWithholding(context.Background())))

	r = newRig(t)
	r.chain.IdentityErr = notary.Errf(notary.KindLedgerUnavailable, "test", nil, "")
	require.Equal(t, notary.KindLedgerUnavailable, notary.KindOf(r.registrar.CheckWithholding(context.Background())))
}
