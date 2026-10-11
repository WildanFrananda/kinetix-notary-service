package keccak

import (
	"encoding/binary"

	"golang.org/x/crypto/sha3"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	version  uint16 = 1
	leafTag  byte   = 0x00
	nodeTag  byte   = 0x01
	emptyTag byte   = 0x02
	epochTag byte   = 0x03
)

type Scheme struct{}

var _ notary.Scheme = Scheme{}

func NewScheme() Scheme {
	return Scheme{}
}

func (Scheme) Version() uint16 {
	return version
}

func (Scheme) LeafHash(epoch uint64, entry notary.Entry) notary.Digest {
	packed := make([]byte, 0, 1+8+32+8)
	packed = append(packed, leafTag)
	packed = binary.BigEndian.AppendUint64(packed, epoch)
	packed = append(packed, entry.Commitment[:]...)
	packed = binary.BigEndian.AppendUint64(packed, entry.Amount)
	return digest(packed)
}

func (Scheme) NodeHash(left, right notary.Node) notary.Digest {
	packed := make([]byte, 0, 1+32+8+32+8)
	packed = append(packed, nodeTag)
	packed = append(packed, left.Hash[:]...)
	packed = binary.BigEndian.AppendUint64(packed, left.Sum)
	packed = append(packed, right.Hash[:]...)
	packed = binary.BigEndian.AppendUint64(packed, right.Sum)
	return digest(packed)
}

func (Scheme) EmptyHash() notary.Digest {
	return digest([]byte{emptyTag})
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
	return digest(packed)
}

func digest(packed []byte) notary.Digest {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(packed)
	var out notary.Digest
	copy(out[:], hasher.Sum(nil))
	return out
}
