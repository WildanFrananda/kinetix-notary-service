package notary

import (
	"cmp"
	"crypto/sha256"
	"encoding/binary"
	"slices"
)

const statementFingerprintDomain = "kinetix.notary.statement-fingerprint.v1"

func StatementFingerprint(liabilities []Liability) Digest {
	sorted := slices.Clone(liabilities)
	slices.SortFunc(sorted, func(a, b Liability) int {
		return cmp.Compare(a.Kind, b.Kind)
	})

	hash := sha256.New()
	_, _ = hash.Write([]byte(statementFingerprintDomain))

	var entry [9]byte

	for _, liability := range sorted {
		entry[0] = byte(liability.Kind)
		binary.BigEndian.PutUint64(entry[1:], liability.Amount)
		_, _ = hash.Write(entry[:])
	}

	return Digest(hash.Sum(nil))
}
