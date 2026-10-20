package poseidon_test

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon2"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestTheNodeHashIsGnarkCryptosCompression(t *testing.T) {
	scheme := poseidon.NewScheme()
	reference := poseidon2.NewDefaultPermutation()

	for range 64 {
		left, right := randomElement(t), randomElement(t)

		got, err := scheme.NodeHash(left, right)
		require.NoError(t, err)

		want, err := reference.Compress(left[:], right[:])
		require.NoError(t, err)
		require.Equal(t, notary.Digest(want), got)
	}
}

func TestALeafIsTheMerkleDamgardHashOfItsThreeBlocks(t *testing.T) {
	scheme := poseidon.NewScheme()

	for range 32 {
		commitment := randomBytes(t)
		entry := notary.Entry{Commitment: commitment, Amount: randomUint64(t)}
		epoch := randomUint64(t)

		var head, high, low notary.Digest

		head[15] = 1
		binary.BigEndian.PutUint64(head[16:24], epoch)
		binary.BigEndian.PutUint64(head[24:32], entry.Amount)
		copy(high[16:], commitment[:16])
		copy(low[16:], commitment[16:])

		require.Equal(t, merkleDamgard(t, head, high, low), scheme.LeafHash(epoch, entry))
	}
}

func TestTheEmptyLeafIsTheHashOfItsTag(t *testing.T) {
	var block notary.Digest

	block[15] = 2

	require.Equal(t, merkleDamgard(t, block), poseidon.NewScheme().EmptyHash())
}

func TestTheTotalCommitmentHashesTheTotalAndTheBlinding(t *testing.T) {
	var blinding notary.Blinding

	_, err := rand.Read(blinding[:])
	require.NoError(t, err)

	var head, tail notary.Digest

	head[15] = 5
	binary.BigEndian.PutUint64(head[24:32], 123_456_789)
	copy(tail[1:], blinding[:])

	require.Equal(t, merkleDamgard(t, head, tail), poseidon.NewScheme().TotalCommitment(123_456_789, blinding))
}

func TestTheBindingRootIsKeccakOfTheTreeRootAndTheCommitment(t *testing.T) {
	treeRoot, commitment := randomBytes(t), randomBytes(t)

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte{0x04})
	hasher.Write(treeRoot[:])
	hasher.Write(commitment[:])

	require.Equal(t, notary.Digest(hasher.Sum(nil)), poseidon.NewScheme().BindingRoot(treeRoot, commitment))
}

func TestTheEpochDigestIsSchemeTwosShapeOfVersionOne(t *testing.T) {
	header := notary.EpochHeader{
		Number:         7,
		PreviousDigest: randomBytes(t),
		Root:           randomBytes(t),
		LeafCount:      3,
		AsOf:           notarytest.Unix(t, 1_790_000_000),
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte{0x03, 0x00, 0x02})
	hasher.Write(binary.BigEndian.AppendUint64(nil, header.Number))
	hasher.Write(header.PreviousDigest[:])
	hasher.Write(header.Root[:])
	hasher.Write(binary.BigEndian.AppendUint64(nil, 0))
	hasher.Write(binary.BigEndian.AppendUint64(nil, header.LeafCount))
	hasher.Write(binary.BigEndian.AppendUint64(nil, 1_790_000_000))

	require.Equal(t, notary.Digest(hasher.Sum(nil)), poseidon.NewScheme().EpochDigest(header))
}

func TestANodeOfANonCanonicalChildIsRefused(t *testing.T) {
	modulus := notary.Digest(fr.Modulus().FillBytes(make([]byte, 32)))

	_, err := poseidon.NewScheme().NodeHash(modulus, notary.Digest{})
	require.Error(t, err)

	_, err = poseidon.NewScheme().NodeHash(notary.Digest{}, modulus)
	require.Error(t, err)
}

func TestTheSchemeIsVersionTwo(t *testing.T) {
	require.Equal(t, uint16(2), poseidon.NewScheme().Version())
}

func TestKnownAnswersPinTheParameters(t *testing.T) {
	scheme := poseidon.NewScheme()

	node, err := scheme.NodeHash(notary.Digest{}, notary.Digest{})
	require.NoError(t, err)

	require.Equal(t, knownNodeOfZeros, node.Hex())
	require.Equal(t, knownEmpty, scheme.EmptyHash().Hex())
}

func merkleDamgard(t *testing.T, blocks ...notary.Digest) notary.Digest {
	t.Helper()

	hasher := poseidon2.NewMerkleDamgardHasher()
	for _, block := range blocks {
		_, err := hasher.Write(block[:])
		require.NoError(t, err)
	}

	return notary.Digest(hasher.Sum(nil))
}

func randomElement(t *testing.T) notary.Digest {
	t.Helper()

	value, err := rand.Int(rand.Reader, fr.Modulus())
	require.NoError(t, err)

	return notary.Digest(value.FillBytes(make([]byte, 32)))
}

func randomBytes(t *testing.T) notary.Digest {
	t.Helper()

	var out notary.Digest

	_, err := rand.Read(out[:])
	require.NoError(t, err)

	return out
}

func randomUint64(t *testing.T) uint64 {
	t.Helper()

	value, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 64))
	require.NoError(t, err)

	return value.Uint64()
}
