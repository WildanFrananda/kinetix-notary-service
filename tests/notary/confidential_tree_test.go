package notary_test

import (
	"math"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

const confidentialDepth = 4

func TestEveryAccountOfAConfidentialEpochProvesItsInclusion(t *testing.T) {
	scheme := poseidon.NewScheme()

	for count := 1; count <= 1<<confidentialDepth; count++ {
		entries := notarytest.Entries(uint64(count), count)

		tree, err := notary.BuildConfidentialTree(scheme, 9, confidentialDepth, entries)
		require.NoError(t, err)

		epoch := confidentialEpochOf(tree)
		for i, entry := range entries {
			proof, err := tree.Prove(uint64(i))
			require.NoError(t, err)
			require.Len(t, proof.Siblings, confidentialDepth)
			require.NoError(
				t,
				notary.VerifyConfidentialInclusion(scheme, epoch, entry, proof),
				"count %d index %d",
				count,
				i,
			)
		}
	}
}

func TestAConfidentialTreePadsToItsGenerationsCapacity(t *testing.T) {
	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(3, 3)

	tree, err := notary.BuildConfidentialTree(scheme, 2, 2, entries)
	require.NoError(t, err)

	left := nodeOf(t, scheme, scheme.LeafHash(2, entries[0]), scheme.LeafHash(2, entries[1]))
	right := nodeOf(t, scheme, scheme.LeafHash(2, entries[2]), scheme.EmptyHash())

	require.Equal(t, nodeOf(t, scheme, left, right), tree.TreeRoot())
	require.Equal(t, 2, tree.Depth())
	require.Equal(t, uint64(3), tree.LeafCount())
}

func TestAnEmptyConfidentialEpochIsAllPadding(t *testing.T) {
	scheme := poseidon.NewScheme()

	tree, err := notary.BuildConfidentialTree(scheme, 2, 1, nil)
	require.NoError(t, err)

	require.Equal(t, nodeOf(t, scheme, scheme.EmptyHash(), scheme.EmptyHash()), tree.TreeRoot())
	require.Equal(t, uint64(0), tree.Total())
}

func TestAConfidentialTreeKeepsItsTotalPrivateButExact(t *testing.T) {
	entries := notarytest.Entries(5, 11)

	var want uint64
	for _, entry := range entries {
		want += entry.Amount
	}

	tree, err := notary.BuildConfidentialTree(poseidon.NewScheme(), 1, confidentialDepth, entries)
	require.NoError(t, err)
	require.Equal(t, want, tree.Total())
}

func TestAConfidentialTreeRefusesMoreAccountsThanItsCapacity(t *testing.T) {
	_, err := notary.BuildConfidentialTree(poseidon.NewScheme(), 1, 2, notarytest.Entries(1, 5))

	require.Equal(t, notary.KindTooManyLeaves, notary.KindOf(err))
}

func TestAConfidentialTreeRefusesATotalThatDoesNotFit(t *testing.T) {
	entries := []notary.Entry{
		{Commitment: notary.Digest{1}, Amount: math.MaxUint64},
		{Commitment: notary.Digest{2}, Amount: 1},
	}

	_, err := notary.BuildConfidentialTree(poseidon.NewScheme(), 1, 1, entries)

	require.Equal(t, notary.KindOverflow, notary.KindOf(err))
}

func TestAConfidentialTreeRefusesADepthOutsideItsRange(t *testing.T) {
	for _, depth := range []int{-1, notary.MaxConfidentialDepth + 1} {
		_, err := notary.BuildConfidentialTree(poseidon.NewScheme(), 1, depth, nil)

		require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err), "depth %d", depth)
	}
}

func TestAConfidentialTreeNeverProvesPadding(t *testing.T) {
	tree, err := notary.BuildConfidentialTree(poseidon.NewScheme(), 1, 2, notarytest.Entries(1, 3))
	require.NoError(t, err)

	_, err = tree.Prove(3)

	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestAConfidentialProofIsRejectedWhenAnythingDiffers(t *testing.T) {
	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(7, 6)

	tree, err := notary.BuildConfidentialTree(scheme, 4, confidentialDepth, entries)
	require.NoError(t, err)

	epoch := confidentialEpochOf(tree)

	proof, err := tree.Prove(2)
	require.NoError(t, err)

	modulus := notary.Digest(fr.Modulus().FillBytes(make([]byte, 32)))

	cases := map[string]func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof){
		"another amount": func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof) {
			entry := entries[2]
			entry.Amount++

			return epoch, entry, proof
		},
		"another commitment": func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof) {
			return epoch, entries[3], proof
		},
		"another epoch": func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof) {
			other := epoch
			other.Number++

			return other, entries[2], proof
		},
		"a padding index": func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof) {
			moved := proof
			moved.Index = epoch.LeafCount

			return epoch, entries[2], moved
		},
		"a sibling too few": func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof) {
			short := proof
			short.Siblings = proof.Siblings[:confidentialDepth-1]

			return epoch, entries[2], short
		},
		"a sibling carrying a sum": func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof) {
			summed := notary.InclusionProof{Index: proof.Index, Siblings: append([]notary.Node(nil), proof.Siblings...)}
			summed.Siblings[0].Sum = 1

			return epoch, entries[2], summed
		},
		"a sibling that is not a field element": func() (notary.ConfidentialEpoch, notary.Entry, notary.InclusionProof) {
			forged := notary.InclusionProof{Index: proof.Index, Siblings: append([]notary.Node(nil), proof.Siblings...)}
			forged.Siblings[1].Hash = modulus

			return epoch, entries[2], forged
		},
	}

	for name, build := range cases {
		e, entry, p := build()

		err := notary.VerifyConfidentialInclusion(scheme, e, entry, p)
		require.Equal(t, notary.KindProofRejected, notary.KindOf(err), name)
	}
}

func TestAConfidentialHeaderBindsTheTreeRootAndTheCommitment(t *testing.T) {
	scheme := poseidon.NewScheme()
	entries := notarytest.Entries(8, 5)

	tree, err := notary.BuildConfidentialTree(scheme, 6, confidentialDepth, entries)
	require.NoError(t, err)

	previous := notary.Digest{0x42}
	asOf := notarytest.Unix(t, 1_790_000_000)
	blinding := notary.Blinding{0x07}

	header, err := tree.Header(scheme, previous, asOf, blinding)
	require.NoError(t, err)

	commitment := scheme.TotalCommitment(tree.Total(), blinding)

	require.Equal(t, tree.TreeRoot(), header.TreeRoot)
	require.Equal(t, commitment, header.TotalCommitment)
	require.Equal(t, scheme.BindingRoot(tree.TreeRoot(), commitment), header.Header.Root)
	require.Equal(t, uint64(0), header.Header.Total)
	require.Equal(t, uint64(5), header.Header.LeafCount)
	require.Equal(t, previous, header.Header.PreviousDigest)
	require.Equal(t, uint64(6), header.Header.Number)

	other, err := tree.Header(scheme, previous, asOf, notary.Blinding{0x08})
	require.NoError(t, err)
	require.NotEqual(t, header.TotalCommitment, other.TotalCommitment)
	require.NotEqual(t, header.Header.Root, other.Header.Root)
}

func TestAConfidentialHeaderNeedsATimeAndABlinding(t *testing.T) {
	scheme := poseidon.NewScheme()

	tree, err := notary.BuildConfidentialTree(scheme, 6, 1, notarytest.Entries(1, 1))
	require.NoError(t, err)

	_, err = tree.Header(scheme, notary.Digest{}, notary.UnixTime{}, notary.Blinding{0x01})
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	_, err = tree.Header(scheme, notary.Digest{}, notarytest.Unix(t, 1), notary.Blinding{})
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func confidentialEpochOf(tree *notary.ConfidentialTree) notary.ConfidentialEpoch {
	return notary.ConfidentialEpoch{
		Number:    tree.Epoch(),
		TreeRoot:  tree.TreeRoot(),
		LeafCount: tree.LeafCount(),
		Depth:     tree.Depth(),
	}
}

func nodeOf(t *testing.T, scheme notary.ConfidentialScheme, left, right notary.Digest) notary.Digest {
	t.Helper()

	node, err := scheme.NodeHash(left, right)
	require.NoError(t, err)

	return node
}
