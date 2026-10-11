package notary_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type anchoredTree struct {
	epoch   notary.AnchoredEpoch
	entries []notary.Entry
	tree    *notary.Tree
}

func anchor(t *testing.T, seed uint64, count int) anchoredTree {
	t.Helper()
	entries := notarytest.Entries(seed, count)
	tree, err := notary.BuildTree(keccak.NewScheme(), 11, entries)
	require.NoError(t, err)
	root := tree.Root()
	return anchoredTree{
		epoch: notary.AnchoredEpoch{
			Number:    11,
			Root:      root.Hash,
			Total:     root.Sum,
			LeafCount: uint64(count),
		},
		entries: entries,
		tree:    tree,
	}
}

func (a anchoredTree) proof(t *testing.T, index uint64) notary.InclusionProof {
	t.Helper()
	proof, err := a.tree.Prove(index)
	require.NoError(t, err)
	return proof
}

func requireRejected(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	require.Equal(t, notary.KindProofRejected, notary.KindOf(err), "%v", err)
}

func TestEveryAccountProvesItsInclusion(t *testing.T) {
	for count := 1; count <= 40; count++ {
		a := anchor(t, uint64(count), count)
		for i := range a.entries {
			err := notary.VerifyInclusion(keccak.NewScheme(), a.epoch, a.entries[i], a.proof(t, uint64(i)))
			require.NoError(t, err, "count %d index %d", count, i)
		}
	}
}

func TestAnInflatedAmountIsRejected(t *testing.T) {
	a := anchor(t, 1, 5)
	entry := a.entries[2]
	entry.Amount++

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, entry, a.proof(t, 2)))
}

func TestAForeignCommitmentIsRejected(t *testing.T) {
	a := anchor(t, 1, 5)
	entry := a.entries[2]
	entry.Commitment[0] ^= 0xff

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, entry, a.proof(t, 2)))
}

func TestValueMovedFromASiblingOntoTheAmountIsRejected(t *testing.T) {
	a := anchor(t, 1, 5)
	proof := a.proof(t, 2)
	entry := a.entries[2]
	require.Positive(t, proof.Siblings[0].Sum)
	entry.Amount += proof.Siblings[0].Sum
	proof.Siblings[0].Sum = 0

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, entry, proof))
}

func TestAForgedSiblingHashIsRejected(t *testing.T) {
	a := anchor(t, 1, 5)
	proof := a.proof(t, 2)
	proof.Siblings[1].Hash[5] ^= 0x01

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, a.entries[2], proof))
}

func TestTheWrongIndexIsRejected(t *testing.T) {
	a := anchor(t, 1, 5)
	proof := a.proof(t, 2)
	proof.Index = 3

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, a.entries[2], proof))
}

func TestAnIndexInThePaddingIsRejected(t *testing.T) {
	a := anchor(t, 1, 5)
	proof := a.proof(t, 4)
	proof.Index = 6

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, notary.Entry{}, proof))
}

func TestAProofOfTheWrongLengthIsRejected(t *testing.T) {
	a := anchor(t, 1, 5)
	proof := a.proof(t, 2)
	proof.Siblings = proof.Siblings[:len(proof.Siblings)-1]

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, a.entries[2], proof))
}

func TestSumsThatOverflowAreRejectedRatherThanWrapped(t *testing.T) {
	a := anchor(t, 1, 5)
	proof := a.proof(t, 2)
	for i := range proof.Siblings {
		proof.Siblings[i].Sum = math.MaxUint64
	}

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, a.entries[2], proof))
}

func TestAProofIsBoundToItsEpochNumber(t *testing.T) {
	a := anchor(t, 1, 5)
	moved := a.epoch
	moved.Number = 12

	requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), moved, a.entries[2], a.proof(t, 2)))
}

func FuzzAnyChangeToAClaimIsRejected(f *testing.F) {
	f.Add(uint64(1), uint8(5), uint8(2), uint8(0), uint8(0), uint64(1))
	f.Add(uint64(7), uint8(33), uint8(31), uint8(4), uint8(3), uint64(math.MaxUint64))
	f.Add(uint64(9), uint8(2), uint8(1), uint8(0), uint8(4), uint64(77))
	f.Fuzz(func(t *testing.T, seed uint64, rawCount, rawIndex, rawLevel, tamper uint8, delta uint64) {
		count := int(rawCount%64) + 2
		index := uint64(rawIndex) % uint64(count)
		if delta == 0 {
			delta = 1
		}
		a := anchor(t, seed, count)
		proof := a.proof(t, index)
		entry := a.entries[index]
		level := int(rawLevel) % len(proof.Siblings)

		switch tamper % 5 {
		case 0:
			entry.Commitment[delta%32] ^= byte(delta>>8) | 1
		case 1:
			entry.Amount += delta
		case 2:
			proof.Siblings[level].Hash[delta%32] ^= byte(delta>>8) | 1
		case 3:
			proof.Siblings[level].Sum += delta
		default:
			proof.Index ^= 1 << level
		}

		requireRejected(t, notary.VerifyInclusion(keccak.NewScheme(), a.epoch, entry, proof))
	})
}
