package notarytest

import (
	"crypto/sha256"
	"encoding/binary"
	"math"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func Entries(seed uint64, count int) []notary.Entry {
	entries := make([]notary.Entry, count)
	for i := range entries {
		entries[i] = notary.Entry{
			Commitment: notary.Digest(sha256.Sum256(label(seed, "commitment", i))),
			Amount:     amount(seed, i, count),
		}
	}
	return entries
}

func Uint64(seed uint64, purpose string, i int) uint64 {
	sum := sha256.Sum256(label(seed, purpose, i))
	return binary.BigEndian.Uint64(sum[:8])
}

func amount(seed uint64, i, count int) uint64 {
	raw := Uint64(seed, "amount", i)
	if count <= 1 {
		return raw
	}
	return raw % (math.MaxUint64/uint64(count) + 1)
}

func label(seed uint64, purpose string, i int) []byte {
	out := binary.BigEndian.AppendUint64(nil, seed)
	out = append(out, purpose...)
	return binary.BigEndian.AppendUint64(out, uint64(i))
}
