package verifying_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

var evidenceRegistry = notary.RegistryIdentity{ChainID: 84532, Address: notary.Address{0xe0, 0x0f}}

type evidenceRig struct {
	chain    *notarytest.EvidenceChain
	registry *notarytest.EvidenceRegistry
	batch    notarytest.EvidenceBatch
	verifier *verifying.EvidenceVerifier
}

func newEvidenceRig(t *testing.T) *evidenceRig {
	t.Helper()

	scheme := keccak.NewEvidence()
	chain := notarytest.NewEvidenceChain(scheme)
	batch := notarytest.BatchOf(
		t,
		scheme,
		evidenceRegistry,
		0,
		notary.Digest{},
		chain.BlockTime-60,
		notarytest.Records(notarytest.Deliveries(t, 5)),
	)
	registry := &notarytest.EvidenceRegistry{Chain: chain, Answers: evidenceRegistry}

	return &evidenceRig{
		chain:    chain,
		registry: registry,
		batch:    batch,
		verifier: verifying.NewEvidenceVerifier(registry, scheme),
	}
}

func (r *evidenceRig) anchor() {
	r.chain.Plant(r.batch.Header)
	r.chain.Finalize()
}

func TestEvidenceInAFinalBatchIsIncluded(t *testing.T) {
	r := newEvidenceRig(t)
	r.anchor()

	for _, claim := range r.batch.Claims {
		verification, err := r.verifier.Verify(context.Background(), claim)
		require.NoError(t, err)
		require.Equal(t, notary.VerdictIncluded, verification.Verdict, verification.Reason)
		require.Equal(t, evidenceRegistry, verification.Registry)
		require.Equal(t, r.batch.Digest, verification.Batch.Digest)
	}
}

func TestEvidenceThatWasChangedIsRejected(t *testing.T) {
	r := newEvidenceRig(t)
	r.anchor()

	changed := r.batch.Claims[2]
	record, ok := changed.Record.(notary.DeliveryRecord)
	require.True(t, ok)

	record.DriverPrincipalID = "someone-else"
	changed.Record = record

	verification, err := r.verifier.Verify(context.Background(), changed)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, verification.Verdict)
	require.Contains(t, verification.Reason, "the proof reaches")

	moved := r.batch.Claims[2]
	moved.Proof.Index = 3

	verification, err = r.verifier.Verify(context.Background(), moved)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, verification.Verdict)
}

func TestABatchHeaderThatDoesNotHashToItsDigestIsRejected(t *testing.T) {
	r := newEvidenceRig(t)
	r.anchor()

	claim := r.batch.Claims[0]
	claim.Header.Through.Sequence++

	verification, err := r.verifier.Verify(context.Background(), claim)
	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, verification.Verdict)
	require.Contains(t, verification.Reason, "does not hash")
}

func TestAProofNamingAnotherBatchThanTheRegistryHoldsIsRejected(t *testing.T) {
	r := newEvidenceRig(t)
	other := r.batch.Header
	other.Root = notary.Digest{0xba, 0xd0}
	r.chain.Plant(other)
	r.chain.Finalize()

	verification, err := r.verifier.Verify(context.Background(), r.batch.Claims[0])
	require.NoError(t, err)
	require.Equal(t, notary.VerdictRejected, verification.Verdict)
	require.Contains(t, verification.Reason, "the registry holds")
}

func TestABatchNotFinalOrNotAnchoredIsNeverAVerdict(t *testing.T) {
	r := newEvidenceRig(t)

	_, err := r.verifier.Verify(context.Background(), r.batch.Claims[0])
	require.Equal(t, notary.KindEpochNotAnchored, notary.KindOf(err))

	r.chain.Plant(r.batch.Header)

	_, err = r.verifier.Verify(context.Background(), r.batch.Claims[0])
	require.Equal(t, notary.KindEpochNotFinal, notary.KindOf(err))
}

func TestAnotherRegistryIsNeverTrusted(t *testing.T) {
	r := newEvidenceRig(t)
	r.anchor()
	r.registry.Answers = notary.RegistryIdentity{ChainID: 84532, Address: notary.Address{0xe0, 0x10}}

	_, err := r.verifier.Verify(context.Background(), r.batch.Claims[0])
	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))
}

func TestARegistryThatDisagreesWithTheSchemeGivesNoVerdict(t *testing.T) {
	r := newEvidenceRig(t)
	r.anchor()
	r.registry.Disagree = true

	_, err := r.verifier.Verify(context.Background(), r.batch.Claims[0])
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestAnotherEvidenceSchemeIsRefused(t *testing.T) {
	r := newEvidenceRig(t)
	r.anchor()

	claim := r.batch.Claims[0]
	claim.Scheme = 2

	_, err := r.verifier.Verify(context.Background(), claim)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}
