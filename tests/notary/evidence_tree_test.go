package notary_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func evidenceLeaves(t *testing.T, count int) ([]notary.Salt, []notary.Digest, []notary.Digest) {
	t.Helper()

	scheme := keccak.NewEvidence()
	salts := make([]notary.Salt, count)
	records := make([]notary.Digest, count)
	leaves := make([]notary.Digest, count)

	for i, record := range notarytest.Deliveries(t, count) {
		digest, err := scheme.RecordDigest(record)
		require.NoError(t, err)

		salts[i] = notary.Salt{byte(i), 0x5a}
		records[i] = digest
		leaves[i] = scheme.EvidenceLeaf(notary.StreamDelivery, salts[i], digest)
	}

	return salts, records, leaves
}

func anchoredFrom(tree *notary.EvidenceTree) notary.AnchoredBatch {
	return notary.AnchoredBatch{
		Stream:    notary.StreamDelivery,
		Number:    3,
		Root:      tree.Root(),
		LeafCount: tree.LeafCount(),
	}
}

func TestEveryRecordOfAnEvidenceTreeVerifies(t *testing.T) {
	scheme := keccak.NewEvidence()

	for _, count := range []int{1, 2, 3, 5, 8, 9, 33} {
		salts, records, leaves := evidenceLeaves(t, count)
		tree, err := notary.BuildEvidenceTree(scheme, notary.StreamDelivery, leaves)
		require.NoError(t, err)
		require.Equal(t, notary.DepthFor(uint64(count)), tree.Depth())

		for i := range leaves {
			proof, err := tree.Prove(uint64(i))
			require.NoError(t, err)
			require.NoError(t, notary.VerifyEvidence(scheme, anchoredFrom(tree), salts[i], records[i], proof))
		}
	}
}

func TestAnEvidenceProofThatProvesNothingIsRejected(t *testing.T) {
	scheme := keccak.NewEvidence()
	salts, records, leaves := evidenceLeaves(t, 5)
	tree, err := notary.BuildEvidenceTree(scheme, notary.StreamDelivery, leaves)
	require.NoError(t, err)

	batch := anchoredFrom(tree)
	proof, err := tree.Prove(2)
	require.NoError(t, err)

	cases := map[string]func() error{
		"another salt": func() error {
			return notary.VerifyEvidence(scheme, batch, notary.Salt{0xff}, records[2], proof)
		},
		"another record": func() error {
			return notary.VerifyEvidence(scheme, batch, salts[2], records[3], proof)
		},
		"another stream": func() error {
			other := batch
			other.Stream = notary.StreamReview

			return notary.VerifyEvidence(scheme, other, salts[2], records[2], proof)
		},
		"another index": func() error {
			moved := notary.EvidenceProof{Index: 3, Siblings: proof.Siblings}

			return notary.VerifyEvidence(scheme, batch, salts[2], records[2], moved)
		},
		"an index in the padding": func() error {
			padded := notary.EvidenceProof{Index: 6, Siblings: proof.Siblings}

			return notary.VerifyEvidence(scheme, batch, salts[2], records[2], padded)
		},
		"one level short": func() error {
			short := notary.EvidenceProof{Index: 2, Siblings: proof.Siblings[:len(proof.Siblings)-1]}

			return notary.VerifyEvidence(scheme, batch, salts[2], records[2], short)
		},
		"a forged sibling": func() error {
			forged := append([]notary.Digest(nil), proof.Siblings...)
			forged[1][4] ^= 1

			return notary.VerifyEvidence(
				scheme,
				batch,
				salts[2],
				records[2],
				notary.EvidenceProof{Index: 2, Siblings: forged},
			)
		},
	}
	for name, attempt := range cases {
		require.Equal(t, notary.KindProofRejected, notary.KindOf(attempt()), name)
	}
}

func TestAnEvidenceTreeOfNothingIsRefused(t *testing.T) {
	_, err := notary.BuildEvidenceTree(keccak.NewEvidence(), notary.StreamReview, nil)

	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestPartiesAreTheRecordsPrincipalsOnceEach(t *testing.T) {
	lifecycle := notary.LifecycleRecord{BuyerPrincipalID: "b", MerchantPrincipalID: "m", ActorPrincipalID: "b"}
	require.Equal(t, []string{"b", "m"}, lifecycle.Parties())

	driverRating := notary.ReviewRecord{Kind: notary.ReviewKindDriver, AuthorPrincipalID: "a", SubjectID: "d"}
	require.Equal(t, []string{"a", "d"}, driverRating.Parties())

	productReview := notary.ReviewRecord{Kind: 1, AuthorPrincipalID: "a", MerchantPrincipalID: "m", SubjectID: "SKU-1"}
	require.Equal(t, []string{"a", "m"}, productReview.Parties())

	require.Equal(t, []string{"d"}, notary.DeliveryRecord{DriverPrincipalID: "d"}.Parties())
}

func TestAnEvidenceTimeIsAfterTheEpochAndWithinASecond(t *testing.T) {
	for _, bad := range [][2]int64{{0, 0}, {-1, 0}, {1, -1}, {1, 1_000_000_000}} {
		_, err := notary.EvidenceTimeFrom(bad[0], int32(bad[1]))
		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), bad)
	}

	at := notarytest.EvidenceTime(t, 1_760_000_000, 999_999_999)
	require.Equal(t, int64(1_760_000_000), at.Time().Unix())
	require.Equal(t, 999_999_999, at.Time().Nanosecond())
}

func TestALogPositionOrdersByTransactionThenSequence(t *testing.T) {
	require.True(
		t,
		notary.LogPosition{Transaction: 2, Sequence: 1}.After(notary.LogPosition{Transaction: 1, Sequence: 9}),
	)
	require.True(
		t,
		notary.LogPosition{Transaction: 1, Sequence: 2}.After(notary.LogPosition{Transaction: 1, Sequence: 1}),
	)
	require.False(
		t,
		notary.LogPosition{Transaction: 1, Sequence: 1}.After(notary.LogPosition{Transaction: 1, Sequence: 1}),
	)
	require.False(
		t,
		notary.LogPosition{Transaction: 0, Sequence: 99}.After(notary.LogPosition{Transaction: 1, Sequence: 0}),
	)
}
