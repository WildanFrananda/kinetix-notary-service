package keccak

import (
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type recordFields struct {
	packed  []byte
	tooLong bool
}

func newRecordFields(record notary.EvidenceRecord) *recordFields {
	position := record.Position()

	packed := make([]byte, 0, 256)
	packed = append(packed, recordTag, byte(record.Stream()))
	packed = binary.BigEndian.AppendUint64(packed, position.Transaction)
	packed = binary.BigEndian.AppendUint64(packed, position.Sequence)

	return &recordFields{packed: packed}
}

func (f *recordFields) text(value string) {
	f.bytes([]byte(value))
}

func (f *recordFields) bytes(value []byte) {
	length := len(value)
	if length > notary.MaxEvidenceField {
		f.tooLong = true

		return
	}

	f.packed = binary.BigEndian.AppendUint32(f.packed, uint32(length))
	f.packed = append(f.packed, value...)
}

func (f *recordFields) number(value uint32) {
	f.packed = binary.BigEndian.AppendUint32(f.packed, value)
}

func (f *recordFields) time(value notary.EvidenceTime) {
	f.packed = binary.BigEndian.AppendUint64(f.packed, value.Seconds())
	f.packed = binary.BigEndian.AppendUint32(f.packed, value.Nanos())
}
