package notary

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

const snapshotDigestTag = "kinetix.payment.liability-snapshot.v1"

func SnapshotDigest(record SnapshotRecord, liabilities []Liability) (Digest, error) {
	ordered := make([]Liability, len(liabilities))
	copy(ordered, liabilities)
	sort.Slice(ordered, func(a, b int) bool {
		if ordered[a].Kind != ordered[b].Kind {
			return ordered[a].Kind < ordered[b].Kind
		}

		return ordered[a].Holder.String() < ordered[b].Holder.String()
	})

	packed := make([]byte, 0, len(snapshotDigestTag)+8+4+len(record.Currency)+8+8+len(ordered)*(1+1+36+8))
	packed = append(packed, snapshotDigestTag...)
	packed = binary.BigEndian.AppendUint64(packed, record.AsOf.Seconds())
	packed = binary.BigEndian.AppendUint32(packed, record.Nanos)
	packed = append(packed, record.Currency...)
	packed = binary.BigEndian.AppendUint64(packed, record.Count)
	packed = binary.BigEndian.AppendUint64(packed, record.Total)

	for _, liability := range ordered {
		name := liability.Holder.String()

		length := len(name)
		if length > MaxHolderBytes {
			return Digest{}, Errf(
				KindInvalidArgument,
				"notary.SnapshotDigest",
				nil,
				"a holder of %d bytes does not fit the snapshot digest",
				length,
			)
		}

		packed = append(packed, byte(liability.Kind), byte(length))
		packed = append(packed, name...)
		packed = binary.BigEndian.AppendUint64(packed, liability.Amount)
	}

	return sha256.Sum256(packed), nil
}
