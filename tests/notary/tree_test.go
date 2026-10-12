package notary_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestAnEmptyEpochCommitsToTheEmptyLeaf(t *testing.T) {
	scheme := keccak.NewScheme()

	tree, err := notary.BuildTree(scheme, 3, nil)

	require.NoError(t, err)
	require.Equal(t, notary.Node{Hash: scheme.EmptyHash()}, tree.Root())
	require.Equal(t, 0, tree.Depth())
	require.Equal(t, uint64(0), tree.LeafCount())
}

func TestASingleAccountIsItsOwnRoot(t *testing.T) {
	scheme := keccak.NewScheme()
	entries := notarytest.Entries(1, 1)

	tree, err := notary.BuildTree(scheme, 3, entries)

	require.NoError(t, err)
	require.Equal(t, notary.Node{Hash: scheme.LeafHash(3, entries[0]), Sum: entries[0].Amount}, tree.Root())
	require.Equal(t, 0, tree.Depth())
}

func TestTheRootSumIsTheTotalOfEveryAccount(t *testing.T) {
	for count := 0; count <= 40; count++ {
		entries := notarytest.Entries(uint64(count), count)

		var want uint64
		for _, e := range entries {
			want += e.Amount
		}

		tree, err := notary.BuildTree(keccak.NewScheme(), 1, entries)

		require.NoError(t, err)
		require.Equal(t, want, tree.Root().Sum, "count %d", count)
		require.Equal(t, notary.DepthFor(uint64(count)), tree.Depth(), "count %d", count)
	}
}

func TestATotalThatDoesNotFitIsRefused(t *testing.T) {
	entries := []notary.Entry{{Amount: math.MaxUint64}, {Amount: 1}}

	_, err := notary.BuildTree(keccak.NewScheme(), 1, entries)

	require.Equal(t, notary.KindOverflow, notary.KindOf(err))
}

func TestATotalOfExactlyTheMaximumIsAccepted(t *testing.T) {
	entries := []notary.Entry{{Amount: math.MaxUint64 - 5}, {Amount: 5}}

	tree, err := notary.BuildTree(keccak.NewScheme(), 1, entries)

	require.NoError(t, err)
	require.Equal(t, uint64(math.MaxUint64), tree.Root().Sum)
}

func TestTheSameAccountsInAnotherEpochHaveAnotherRoot(t *testing.T) {
	entries := notarytest.Entries(9, 6)

	first, err := notary.BuildTree(keccak.NewScheme(), 1, entries)
	require.NoError(t, err)
	second, err := notary.BuildTree(keccak.NewScheme(), 2, entries)
	require.NoError(t, err)

	require.NotEqual(t, first.Root().Hash, second.Root().Hash)
	require.Equal(t, first.Root().Sum, second.Root().Sum)
}

func TestAProofIsRefusedForAnIndexOutsideTheAccounts(t *testing.T) {
	tree, err := notary.BuildTree(keccak.NewScheme(), 1, notarytest.Entries(1, 5))
	require.NoError(t, err)

	_, err = tree.Prove(5)

	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestAHeaderCarriesTheRootAndTheSnapshotTime(t *testing.T) {
	tree, err := notary.BuildTree(keccak.NewScheme(), 4, notarytest.Entries(2, 3))
	require.NoError(t, err)

	previous := notary.Digest{1}

	header, err := tree.Header(previous, notarytest.Unix(t, 1760000000))

	require.NoError(t, err)
	require.Equal(
		t,
		notary.EpochHeader{
			Number:         4,
			PreviousDigest: previous,
			Root:           tree.Root().Hash,
			Total:          tree.Root().Sum,
			LeafCount:      3,
			AsOf:           notarytest.Unix(t, 1760000000),
		},
		header,
	)
}

func TestAHeaderWithoutASnapshotTimeIsRefused(t *testing.T) {
	tree, err := notary.BuildTree(keccak.NewScheme(), 4, notarytest.Entries(2, 3))
	require.NoError(t, err)

	_, err = tree.Header(notary.Digest{}, notary.UnixTime{})

	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}

func TestDepthIsTheSmallestThatHoldsEveryLeaf(t *testing.T) {
	cases := map[uint64]int{0: 0, 1: 0, 2: 1, 3: 2, 4: 2, 5: 3, 8: 3, 9: 4, notary.MaxLeaves: notary.MaxDepth}
	for leaves, want := range cases {
		require.Equal(t, want, notary.DepthFor(leaves), "leaves %d", leaves)
	}
}
