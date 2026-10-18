package keccak

import (
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	evidenceVersion uint16 = 1
	evidenceLeafTag byte   = 0x10
	evidenceNodeTag byte   = 0x11
	evidenceEmpty   byte   = 0x12
	batchTag        byte   = 0x13
	recordTag       byte   = 0x20
)

type Evidence struct{}

var _ notary.EvidenceScheme = Evidence{}

func NewEvidence() Evidence {
	return Evidence{}
}

func (Evidence) EvidenceVersion() uint16 {
	return evidenceVersion
}

func (Evidence) EvidenceLeaf(stream notary.EvidenceStream, salt notary.Salt, record notary.Digest) notary.Digest {
	packed := make([]byte, 0, 1+1+32+32)
	packed = append(packed, evidenceLeafTag, byte(stream))
	packed = append(packed, salt[:]...)
	packed = append(packed, record[:]...)

	return digest(packed)
}

func (Evidence) EvidenceNode(left, right notary.Digest) notary.Digest {
	packed := make([]byte, 0, 1+32+32)
	packed = append(packed, evidenceNodeTag)
	packed = append(packed, left[:]...)
	packed = append(packed, right[:]...)

	return digest(packed)
}

func (Evidence) EvidenceEmpty() notary.Digest {
	return digest([]byte{evidenceEmpty})
}

func (Evidence) BatchDigest(header notary.EvidenceBatchHeader) notary.Digest {
	packed := make([]byte, 0, 1+2+1+8+32+32+8+8+8+8)
	packed = append(packed, batchTag)
	packed = binary.BigEndian.AppendUint16(packed, evidenceVersion)
	packed = append(packed, byte(header.Stream))
	packed = binary.BigEndian.AppendUint64(packed, header.Number)
	packed = append(packed, header.PreviousDigest[:]...)
	packed = append(packed, header.Root[:]...)
	packed = binary.BigEndian.AppendUint64(packed, header.LeafCount)
	packed = binary.BigEndian.AppendUint64(packed, header.Through.Transaction)
	packed = binary.BigEndian.AppendUint64(packed, header.Through.Sequence)
	packed = binary.BigEndian.AppendUint64(packed, header.ClosedAt.Seconds())

	return digest(packed)
}

func (Evidence) RecordDigest(record notary.EvidenceRecord) (notary.Digest, error) {
	const op = "keccak.Evidence.RecordDigest"

	fields := newRecordFields(record)

	switch r := record.(type) {
	case notary.LifecycleRecord:
		fields.text(r.OrderNumber)
		fields.number(r.Kind)
		fields.number(r.FromStatus)
		fields.number(r.ToStatus)
		fields.text(r.ReturnNumber)
		fields.number(r.ReturnStatus)
		fields.number(r.ActorKind)
		fields.text(r.ActorPrincipalID)
		fields.text(r.BuyerPrincipalID)
		fields.text(r.MerchantPrincipalID)
		fields.time(r.OccurredAt)
	case notary.DeliveryRecord:
		fields.text(r.OrderNumber)
		fields.text(r.DriverPrincipalID)
		fields.bytes(r.PhotoSHA256)
		fields.bytes(r.SignatureSHA256)
		fields.time(r.DeliveredAt)
	case notary.ReviewRecord:
		fields.text(r.ReviewID)
		fields.number(r.Kind)
		fields.text(r.OrderID)
		fields.text(r.SubjectID)
		fields.text(r.AuthorPrincipalID)
		fields.text(r.MerchantPrincipalID)
		fields.number(r.Rating)
		fields.bytes(r.CommentSHA256)
		fields.time(r.CreatedAt)
	default:
		return notary.Digest{}, notary.Errf(notary.KindInvalidArgument, op, nil, "no encoding for a %T", record)
	}

	if fields.tooLong {
		return notary.Digest{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"a field of the %s record at %d/%d is longer than %d bytes",
			record.Stream(),
			record.Position().Transaction,
			record.Position().Sequence,
			notary.MaxEvidenceField,
		)
	}

	return digest(fields.packed), nil
}
