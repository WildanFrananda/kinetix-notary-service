package plonk_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon2"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/plonk"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestTheCircuitAcceptsEveryHonestEpoch(t *testing.T) {
	scheme := poseidon.NewScheme()

	for depth := 0; depth <= 3; depth++ {
		for count := 0; count <= 1<<depth; count++ {
			entries := notarytest.Entries(uint64(depth*10+count), count)
			assignment := honest(t, scheme, 41, depth, entries)

			err := test.IsSolved(plonk.NewLiabilityCircuit(depth), assignment, ecc.BN254.ScalarField())
			require.NoError(t, err, "depth %d count %d", depth, count)
		}
	}
}

func TestTheCircuitRefusesAnotherTreeRoot(t *testing.T) {
	assignment := honest(t, poseidon.NewScheme(), 3, 2, notarytest.Entries(1, 3))
	assignment.TreeRoot = new(big.Int).Add(bigOf(t, assignment.TreeRoot), big.NewInt(1))

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), assignment, ecc.BN254.ScalarField()))
}

func TestTheCircuitRefusesAnotherTotalCommitment(t *testing.T) {
	assignment := honest(t, poseidon.NewScheme(), 3, 2, notarytest.Entries(1, 3))
	assignment.TotalCommitment = new(big.Int).Add(bigOf(t, assignment.TotalCommitment), big.NewInt(1))

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), assignment, ecc.BN254.ScalarField()))
}

func TestTheCircuitRefusesAnotherEpoch(t *testing.T) {
	assignment := honest(t, poseidon.NewScheme(), 3, 2, notarytest.Entries(1, 3))
	assignment.Epoch = 4

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), assignment, ecc.BN254.ScalarField()))
}

func TestTheCircuitRefusesANegativeAmountThatCancelsAnother(t *testing.T) {
	minusOne := new(big.Int).Sub(fr.Modulus(), big.NewInt(1))
	forged := consistent(t, 3, 2, []*big.Int{minusOne, big.NewInt(2)}, nil)

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), forged, ecc.BN254.ScalarField()))
}

func TestTheCircuitRefusesATotalThatOverflows(t *testing.T) {
	maximum := new(big.Int).SetUint64(math.MaxUint64)
	forged := consistent(t, 3, 2, []*big.Int{maximum, maximum}, nil)

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), forged, ecc.BN254.ScalarField()))
}

func TestTheCircuitRefusesALimbBeyondItsWidth(t *testing.T) {
	wide := new(big.Int).Lsh(big.NewInt(1), 128)
	forged := consistent(t, 3, 2, []*big.Int{big.NewInt(5)}, wide)

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), forged, ecc.BN254.ScalarField()))
}

func TestTheCircuitRefusesAnAmountHiddenInThePadding(t *testing.T) {
	forged := consistent(t, 3, 2, []*big.Int{big.NewInt(5), big.NewInt(6), big.NewInt(7)}, nil)
	forged.Amounts[3] = big.NewInt(9)

	total := big.NewInt(5 + 6 + 7 + 9)
	forged.TotalCommitment = commitmentOf(t, total, blindingOf())

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), forged, ecc.BN254.ScalarField()))
}

func TestTheCircuitRefusesMoreAccountsThanItsCapacity(t *testing.T) {
	forged := consistent(t, 3, 2, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4)}, nil)
	forged.LeafCount = 5

	require.Error(t, test.IsSolved(plonk.NewLiabilityCircuit(2), forged, ecc.BN254.ScalarField()))
}

func TestTheConsistentForgeryBuilderAgreesWithTheCircuit(t *testing.T) {
	built := consistent(t, 3, 2, []*big.Int{big.NewInt(5), big.NewInt(6), big.NewInt(7)}, nil)

	require.NoError(t, test.IsSolved(plonk.NewLiabilityCircuit(2), built, ecc.BN254.ScalarField()))
}

func TestTheCircuitsSizeStaysWithinItsBudgetPerAccount(t *testing.T) {
	previous := 0

	for depth := 2; depth <= 7; depth++ {
		compiled, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, plonk.NewLiabilityCircuit(depth))
		require.NoError(t, err)

		count := compiled.GetNbConstraints()
		if previous > 0 {
			marginal := (count - previous) / (1 << (depth - 1))
			t.Logf("depth %d: %d constraints, %d per added account", depth, count, marginal)
			require.Less(t, marginal, 1500, "depth %d", depth)
		}

		previous = count
	}
}

func honest(
	t *testing.T,
	scheme notary.ConfidentialScheme,
	epoch uint64,
	depth int,
	entries []notary.Entry,
) *plonk.LiabilityCircuit {
	t.Helper()

	tree, err := notary.BuildConfidentialTree(scheme, epoch, depth, entries)
	require.NoError(t, err)

	blinding := notary.Blinding{0x5a, 0x01}

	header, err := tree.Header(scheme, notary.Digest{}, notarytest.Unix(t, 1_790_000_000), blinding)
	require.NoError(t, err)

	assignment, err := plonk.NewAssignment(epoch, depth, entries, header, blinding)
	require.NoError(t, err)

	return assignment
}

func consistent(t *testing.T, epoch uint64, depth int, amounts []*big.Int, firstHigh *big.Int) *plonk.LiabilityCircuit {
	t.Helper()

	capacity := 1 << depth
	assignment := plonk.NewLiabilityCircuit(depth)
	leaves := make([]*big.Int, capacity)
	total := new(big.Int)

	for i := range capacity {
		high, low, amount := big.NewInt(int64(i+1)), big.NewInt(int64(i+100)), big.NewInt(0)
		if i == 0 && firstHigh != nil {
			high = firstHigh
		}

		if i < len(amounts) {
			amount = amounts[i]
			leaves[i] = leafOf(t, epoch, high, low, amount)
			total.Add(total, amount)
		} else {
			high, low = big.NewInt(0), big.NewInt(0)
			leaves[i] = emptyLeaf(t)
		}

		assignment.High[i] = high
		assignment.Low[i] = low
		assignment.Amounts[i] = amount
	}

	assignment.Epoch = epoch
	assignment.LeafCount = len(amounts)
	assignment.TreeRoot = rootOf(t, leaves)
	assignment.TotalCommitment = commitmentOf(t, total, blindingOf())
	assignment.Blinding = blindingOf()

	return assignment
}

func leafOf(t *testing.T, epoch uint64, high, low, amount *big.Int) *big.Int {
	t.Helper()

	head := new(big.Int).Lsh(big.NewInt(1), 128)
	head.Add(head, new(big.Int).Lsh(new(big.Int).SetUint64(epoch), 64))
	head.Add(head, amount)

	return hash(t, head, high, low)
}

func emptyLeaf(t *testing.T) *big.Int {
	t.Helper()

	return hash(t, new(big.Int).Lsh(big.NewInt(2), 128))
}

func commitmentOf(t *testing.T, total, blinding *big.Int) *big.Int {
	t.Helper()

	head := new(big.Int).Lsh(big.NewInt(5), 128)

	return hash(t, head.Add(head, total), blinding)
}

func rootOf(t *testing.T, leaves []*big.Int) *big.Int {
	t.Helper()

	permutation := poseidon2.NewDefaultPermutation()
	level := leaves

	for len(level) > 1 {
		above := make([]*big.Int, len(level)/2)
		for i := range above {
			out, err := permutation.Compress(canonical(level[2*i]), canonical(level[2*i+1]))
			require.NoError(t, err)

			above[i] = new(big.Int).SetBytes(out)
		}

		level = above
	}

	return level[0]
}

func hash(t *testing.T, blocks ...*big.Int) *big.Int {
	t.Helper()

	hasher := poseidon2.NewMerkleDamgardHasher()
	for _, block := range blocks {
		_, err := hasher.Write(canonical(block))
		require.NoError(t, err)
	}

	return new(big.Int).SetBytes(hasher.Sum(nil))
}

func canonical(value *big.Int) []byte {
	reduced := new(big.Int).Mod(value, fr.Modulus())

	return reduced.FillBytes(make([]byte, 32))
}

func bigOf(t *testing.T, value frontend.Variable) *big.Int {
	t.Helper()

	out, ok := value.(*big.Int)
	require.True(t, ok)

	return out
}

func blindingOf() *big.Int {
	return big.NewInt(0x5a01)
}
