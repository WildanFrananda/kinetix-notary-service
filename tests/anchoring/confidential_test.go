package anchoring_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const depth = 4

func newConfidentialRig(t *testing.T) (*rig, *notarytest.Prover) {
	t.Helper()

	r := newRig(t)
	scheme := poseidon.NewScheme()
	prover := &notarytest.Prover{}
	r.chain.Confidential = scheme
	r.chain.DepthValue = depth
	r.watch.Scheme = scheme
	r.two = &anchoring.SchemeTwo{
		Scheme:    scheme,
		Depth:     depth,
		Blindings: &notarytest.Blindings{Seed: 9},
		Prover:    prover,
		Chain:     r.chain,
	}
	r.build()

	return r, prover
}

func (r *rig) attemptsOf(epoch uint64, kind notary.AttemptKind) []notary.Attempt {
	var kept []notary.Attempt

	for _, attempt := range r.store.AttemptsOf(epoch) {
		if attempt.Kind == kind {
			kept = append(kept, attempt)
		}
	}

	return kept
}

func TestASchemeTwoEpochIsProvenThenApprovedThenAnchored(t *testing.T) {
	r, prover := newConfidentialRig(t)

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateBuilt, r.state(t, 0))
	require.Len(t, prover.Requests, 1)
	require.Len(t, prover.Requests[0].Entries, 12)
	require.Equal(t, depth, prover.Requests[0].Depth)
	require.Len(t, r.attemptsOf(0, notary.AttemptForProof), 1)
	require.Empty(t, r.watch.Proposals, "the watch is asked only once the proof has landed")

	record := r.store.Records[0]
	require.Equal(t, uint64(0), record.Header.Total)
	require.NotNil(t, record.Confidential)
	require.NotEmpty(t, record.Confidential.Proof)
	require.Equal(t, poseidon.NewScheme().EpochDigest(record.Header), record.Digest)
	require.Equal(t, 1, r.observer.Proofs)
	require.Equal(t, uint64(12), r.observer.Accounts)
	require.Equal(t, uint64(16), r.observer.Capacity)

	r.chain.Mine()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
	require.Len(t, r.watch.Proposals, 1)
	require.Len(t, r.attemptsOf(0, notary.AttemptForAnchor), 1)

	r.chain.Mine()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateMined, r.state(t, 0))

	r.chain.Finalize()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateFinalized, r.state(t, 0))

	for _, attempt := range r.store.AttemptsOf(0) {
		require.Equal(t, notary.AttemptIncluded, attempt.Outcome, attempt.Kind.String())
	}

	proof := r.attemptsOf(0, notary.AttemptForProof)[0].Receipt
	anchor := r.attemptsOf(0, notary.AttemptForAnchor)[0].Receipt
	proofCost, _ := proof.CostWei()
	anchorCost, _ := anchor.CostWei()
	require.Equal(t, proofCost+anchorCost, r.observer.MeanCost, "an epoch costs its proof and its anchor")
}

func TestAStuckProofIsOutbidAtTheSameNonce(t *testing.T) {
	r, _ := newConfidentialRig(t)

	require.NoError(t, r.tick(t))
	first := r.attemptsOf(0, notary.AttemptForProof)[0]

	r.advance(r.policy.ReplaceAfter)
	require.NoError(t, r.tick(t))

	proofs := r.attemptsOf(0, notary.AttemptForProof)
	require.Len(t, proofs, 2)
	require.Equal(t, first.Anchor.Nonce, proofs[1].Anchor.Nonce)
	require.Greater(t, proofs[1].Anchor.Fees.MaxFeePerGas, first.Anchor.Fees.MaxFeePerGas)
	require.Equal(t, notary.StateBuilt, r.state(t, 0))

	r.chain.Mine()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
}

func TestAProofTheChainRefusesAbandonsTheEpoch(t *testing.T) {
	r, _ := newConfidentialRig(t)
	r.chain.RefuseProofs = true

	err := r.tick(t)
	require.Equal(t, notary.KindAnchorRefused, notary.KindOf(err))
	require.Equal(t, notary.StateAbandoned, r.state(t, 0))
	require.Contains(t, r.store.Records[0].AbandonedReason, "ProofRejected")
	require.Empty(t, r.watch.Proposals)
}

func TestAProverThatFailsBuildsNothingAndIsCounted(t *testing.T) {
	r, prover := newConfidentialRig(t)
	killed := errors.New("signal: killed")
	prover.Err = killed

	require.ErrorIs(t, r.tick(t), killed)
	require.Empty(t, r.store.Records)
	require.Equal(t, 1, r.observer.ProofFailures)

	prover.Err = nil

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateBuilt, r.state(t, 0))
}

func TestAProofThatLeftTheChainIsSentAgainBeforeAnyApproval(t *testing.T) {
	r, _ := newConfidentialRig(t)

	require.NoError(t, r.tick(t))
	r.chain.Mine()
	r.chain.Reorg()
	r.chain.Forget()

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateBuilt, r.state(t, 0))
	require.Empty(t, r.watch.Proposals)

	r.chain.Mine()
	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateSubmitted, r.state(t, 0))
}

func TestAProofLongerThanTheLeaseIsStillSavedByItsWriter(t *testing.T) {
	r, prover := newConfidentialRig(t)
	prover.During = func() { r.advance(2 * r.policy.LeaseTerm) }

	require.NoError(t, r.tick(t))
	require.Equal(t, notary.StateBuilt, r.state(t, 0))
}

func TestAWriterThatLostTheLeaseWhileProvingSavesNothing(t *testing.T) {
	r, prover := newConfidentialRig(t)
	prover.During = func() {
		r.advance(2 * r.policy.LeaseTerm)

		_, held, err := r.store.Lease(context.Background(), "notaryd-2", time.Minute)
		require.NoError(t, err)
		require.True(t, held)
	}

	require.Equal(t, notary.KindLeaseLost, notary.KindOf(r.tick(t)))
	require.Empty(t, r.store.Records)
}
