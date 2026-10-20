package poseidon

import (
	"encoding/binary"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"golang.org/x/crypto/sha3"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	version    uint16 = 2
	leafTag    byte   = 1
	emptyTag   byte   = 2
	totalTag   byte   = 5
	epochTag   byte   = 0x03
	bindingTag byte   = 0x04
	tagByte           = 15
)

type Scheme struct {
	permutation permutation
}

var _ notary.ConfidentialScheme = Scheme{}

func NewScheme() Scheme {
	return Scheme{permutation: defaultPermutation()}
}

func (Scheme) Version() uint16 {
	return version
}

func (s Scheme) LeafHash(epoch uint64, entry notary.Entry) notary.Digest {
	var head, high, low notary.Digest

	head[tagByte] = leafTag
	binary.BigEndian.PutUint64(head[16:24], epoch)
	binary.BigEndian.PutUint64(head[24:32], entry.Amount)
	copy(high[16:], entry.Commitment[:16])
	copy(low[16:], entry.Commitment[16:])

	return s.hash(head, high, low)
}

func (s Scheme) NodeHash(left, right notary.Digest) (notary.Digest, error) {
	const op = "poseidon.Scheme.NodeHash"

	var x, y fr.Element
	if err := x.SetBytesCanonical(left[:]); err != nil {
		return notary.Digest{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"the left child is not a field element",
		)
	}

	if err := y.SetBytesCanonical(right[:]); err != nil {
		return notary.Digest{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"the right child is not a field element",
		)
	}

	return digestOf(s.permutation.compress(x, y)), nil
}

func (s Scheme) EmptyHash() notary.Digest {
	var block notary.Digest

	block[tagByte] = emptyTag

	return s.hash(block)
}

func (s Scheme) TotalCommitment(total uint64, blinding notary.Blinding) notary.Digest {
	var head, tail notary.Digest

	head[tagByte] = totalTag
	binary.BigEndian.PutUint64(head[24:32], total)
	copy(tail[1:], blinding[:])

	return s.hash(head, tail)
}

func (Scheme) BindingRoot(treeRoot, totalCommitment notary.Digest) notary.Digest {
	packed := make([]byte, 0, 1+32+32)
	packed = append(packed, bindingTag)
	packed = append(packed, treeRoot[:]...)
	packed = append(packed, totalCommitment[:]...)

	return keccak(packed)
}

func (Scheme) EpochDigest(header notary.EpochHeader) notary.Digest {
	packed := make([]byte, 0, 1+2+8+32+32+8+8+8)
	packed = append(packed, epochTag)
	packed = binary.BigEndian.AppendUint16(packed, version)
	packed = binary.BigEndian.AppendUint64(packed, header.Number)
	packed = append(packed, header.PreviousDigest[:]...)
	packed = append(packed, header.Root[:]...)
	packed = binary.BigEndian.AppendUint64(packed, header.Total)
	packed = binary.BigEndian.AppendUint64(packed, header.LeafCount)
	packed = binary.BigEndian.AppendUint64(packed, header.AsOf.Seconds())

	return keccak(packed)
}

func (s Scheme) hash(blocks ...notary.Digest) notary.Digest {
	var state fr.Element

	for _, block := range blocks {
		var element fr.Element
		element.SetBytes(block[:])
		state = s.permutation.compress(state, element)
	}

	return digestOf(state)
}

func digestOf(element fr.Element) notary.Digest {
	return notary.Digest(element.Bytes())
}

func keccak(packed []byte) notary.Digest {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(packed)

	var out notary.Digest
	copy(out[:], hasher.Sum(nil))

	return out
}
