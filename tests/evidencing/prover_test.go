package evidencing_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func provingRig(t *testing.T) *rig {
	t.Helper()

	r := newRig(t)
	events := notarytest.Lifecycle(t, "ORD-1", "buyer-1", "merchant-1", 100, 4)
	events = append(events, notary.LifecycleRecord{
		At:                  notary.LogPosition{Transaction: 104, Sequence: 5},
		OrderNumber:         "ORD-1",
		Kind:                1,
		FromStatus:          5,
		ToStatus:            6,
		ActorKind:           4,
		ActorPrincipalID:    "admin-9",
		BuyerPrincipalID:    "buyer-1",
		MerchantPrincipalID: "merchant-1",
		OccurredAt:          notarytest.EvidenceTime(t, 1_760_000_500, 0),
	})
	r.feeds[notary.StreamOrderLifecycle].Records = notarytest.Records(events)

	deliveries := notarytest.Deliveries(t, 2)
	deliveries[0].OrderNumber, deliveries[0].DriverPrincipalID = "ORD-1", "courier-1"
	deliveries[1].OrderNumber, deliveries[1].DriverPrincipalID = "ORD-9", "courier-2"
	r.feeds[notary.StreamDelivery].Records = notarytest.Records(deliveries)

	reviews := notarytest.Reviews(t, 2)
	reviews[0].ReviewID, reviews[0].AuthorPrincipalID, reviews[0].MerchantPrincipalID = "review-p", "buyer-1", "merchant-1"
	reviews[1].ReviewID, reviews[1].Kind, reviews[1].SubjectID = "review-d", notary.ReviewKindDriver, "courier-1"
	reviews[1].AuthorPrincipalID, reviews[1].MerchantPrincipalID = "buyer-1", "merchant-1"
	r.feeds[notary.StreamReview].Records = notarytest.Records(reviews)

	return r
}

func (r *rig) finalizeAll(t *testing.T) {
	t.Helper()

	for range 3 {
		r.tick(t)
		r.settle(t)
	}

	for _, stream := range notary.EvidenceStreams() {
		require.Equal(t, "finalized", r.state(t, stream, 0))
	}
}

func requireProven(t *testing.T, r *rig, bundle evidencing.Bundle) {
	t.Helper()

	scheme := keccak.NewEvidence()

	for _, claim := range bundle.Claims {
		anchored, found, err := r.chain.BatchAtFinalized(context.Background(), claim.Header.Stream, claim.Header.Number)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, anchored.Digest, claim.Digest)
		require.Equal(t, claim.Digest, scheme.BatchDigest(claim.Header))
		require.Equal(t, registryID, claim.Registry)
		require.Equal(t, scheme.EvidenceVersion(), claim.Scheme)

		record, err := scheme.RecordDigest(claim.Record)
		require.NoError(t, err)
		require.NoError(t, notary.VerifyEvidence(scheme, anchored, claim.Salt, record, claim.Proof))
	}
}

func streamsOf(bundle evidencing.Bundle) []notary.EvidenceStream {
	streams := make([]notary.EvidenceStream, len(bundle.Claims))
	for i, claim := range bundle.Claims {
		streams[i] = claim.Record.Stream()
	}

	return streams
}

func TestAnOrdersBuyerAndMerchantProveEveryEventAndItsDelivery(t *testing.T) {
	r := provingRig(t)
	r.finalizeAll(t)
	prover := evidencing.NewProver(keccak.NewEvidence(), r.store)

	for _, party := range []string{"buyer-1", "merchant-1"} {
		bundle, err := prover.OrderProofs(context.Background(), party, "ORD-1")
		require.NoError(t, err)
		require.Len(t, bundle.Claims, 6)
		require.Zero(t, bundle.Pending)
		require.Equal(t, notary.StreamDelivery, streamsOf(bundle)[5])
		requireProven(t, r, bundle)
	}
}

func TestACourierAndAnAdminSeeOnlyWhatTheyTookPartIn(t *testing.T) {
	r := provingRig(t)
	r.finalizeAll(t)
	prover := evidencing.NewProver(keccak.NewEvidence(), r.store)

	courier, err := prover.OrderProofs(context.Background(), "courier-1", "ORD-1")
	require.NoError(t, err)
	require.Equal(t, []notary.EvidenceStream{notary.StreamDelivery}, streamsOf(courier))
	requireProven(t, r, courier)

	admin, err := prover.OrderProofs(context.Background(), "admin-9", "ORD-1")
	require.NoError(t, err)
	require.Len(t, admin.Claims, 1)

	acted, ok := admin.Claims[0].Record.(notary.LifecycleRecord)
	require.True(t, ok)
	require.Equal(t, "admin-9", acted.ActorPrincipalID)

	for _, stranger := range []string{"stranger", "courier-2"} {
		_, err := prover.OrderProofs(context.Background(), stranger, "ORD-1")
		require.Equal(t, notary.KindNoEvidence, notary.KindOf(err), stranger)
	}

	_, err = prover.OrderProofs(context.Background(), "buyer-1", "ORD-404")
	require.Equal(t, notary.KindNoEvidence, notary.KindOf(err))
}

func TestAReviewIsProvenToItsAuthorItsMerchantAndARatedCourier(t *testing.T) {
	r := provingRig(t)
	r.finalizeAll(t)
	prover := evidencing.NewProver(keccak.NewEvidence(), r.store)

	for _, ask := range []struct{ party, review string }{
		{"buyer-1", "review-p"},
		{"merchant-1", "review-p"},
		{"buyer-1", "review-d"},
		{"courier-1", "review-d"},
	} {
		bundle, err := prover.ReviewProofs(context.Background(), ask.party, ask.review)
		require.NoError(t, err, ask)
		require.Len(t, bundle.Claims, 1, ask)
		requireProven(t, r, bundle)
	}

	_, err := prover.ReviewProofs(context.Background(), "courier-1", "review-p")
	require.Equal(t, notary.KindNoEvidence, notary.KindOf(err))
}

func TestRecordsWhoseBatchIsNotFinalArePendingNotProven(t *testing.T) {
	r := provingRig(t)
	r.tick(t)
	prover := evidencing.NewProver(keccak.NewEvidence(), r.store)

	bundle, err := prover.OrderProofs(context.Background(), "buyer-1", "ORD-1")
	require.NoError(t, err)
	require.Empty(t, bundle.Claims)
	require.Equal(t, 6, bundle.Pending)
}

func TestARecordChangedInTheDatabaseIsNeverProven(t *testing.T) {
	r := provingRig(t)
	r.finalizeAll(t)
	_, err := r.pool.Exec(
		context.Background(),
		"UPDATE lifecycle_records SET to_status = 99 WHERE transaction_id = 101",
	)
	require.NoError(t, err)

	_, err = evidencing.NewProver(keccak.NewEvidence(), r.store).OrderProofs(context.Background(), "buyer-1", "ORD-1")
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestARegistryCheckRefusesAWrongRegistryTheLiabilitiesKeyAndAKeyWithoutTheRole(t *testing.T) {
	r := newRig(t)
	registry := &notarytest.EvidenceRegistry{Chain: r.chain, Answers: registryID}
	liabilities := notary.Address{0x51, 0xb2}

	require.NoError(
		t,
		evidencing.NewRegistryCheck(registry, r.chain, registryID, liabilities).Check(context.Background()),
	)

	other := registryID
	other.Address = notary.Address{0xe0, 0x02}
	err := evidencing.NewRegistryCheck(registry, r.chain, other, liabilities).Check(context.Background())
	require.Equal(t, notary.KindWrongLedger, notary.KindOf(err))

	err = evidencing.NewRegistryCheck(registry, r.chain, registryID, r.chain.From).Check(context.Background())
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	r.chain.Authorized = false
	err = evidencing.NewRegistryCheck(registry, r.chain, registryID, liabilities).Check(context.Background())
	require.Equal(t, notary.KindNotAnchorer, notary.KindOf(err))
}
