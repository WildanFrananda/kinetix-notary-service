package feeds

import (
	"crypto/sha256"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	commonv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/common/v1"
)

func refused(op, peer, format string, args ...any) error {
	return notary.Errf(notary.KindSourceRefused, op, nil, peer+" sent "+format, args...)
}

func cursorOf(position notary.LogPosition) *commonv1.LogCursor {
	return &commonv1.LogCursor{TransactionId: position.Transaction, Sequence: position.Sequence}
}

func positionOf(op, peer, what string, cursor *commonv1.LogCursor) (notary.LogPosition, error) {
	if cursor == nil {
		return notary.LogPosition{}, refused(op, peer, "%s without a position", what)
	}

	return notary.LogPosition{Transaction: cursor.GetTransactionId(), Sequence: cursor.GetSequence()}, nil
}

func timeOf(op, peer, field string, at *timestamppb.Timestamp) (notary.EvidenceTime, error) {
	if at == nil {
		return notary.EvidenceTime{}, refused(op, peer, "a record without %s", field)
	}

	t, err := notary.EvidenceTimeFrom(at.GetSeconds(), at.GetNanos())
	if err != nil {
		return notary.EvidenceTime{}, notary.Errf(
			notary.KindSourceRefused,
			op,
			err,
			"%s sent %s outside any time",
			peer,
			field,
		)
	}

	return t, nil
}

func hashOf(op, peer, field string, raw []byte, required bool) ([]byte, error) {
	if len(raw) == 0 && !required {
		return nil, nil
	}

	if len(raw) != sha256.Size {
		return nil, refused(op, peer, "a %d-byte %s, not a SHA-256", len(raw), field)
	}

	return append([]byte(nil), raw...), nil
}

func texts(op, peer string, fields ...field) error {
	for _, f := range fields {
		if f.value == "" && !f.optional {
			return refused(op, peer, "a record without %s", f.name)
		}

		if len(f.value) > notary.MaxEvidenceField {
			return refused(
				op,
				peer,
				"%s of %d bytes, more than the scheme's %d",
				f.name,
				len(f.value),
				notary.MaxEvidenceField,
			)
		}

		if strings.ContainsRune(f.value, 0) {
			return refused(op, peer, "%s with a NUL byte", f.name)
		}
	}

	return nil
}

func enumOf(op, peer, field string, value int32, required bool) (uint32, error) {
	if value < 0 || (required && value == 0) {
		return 0, refused(op, peer, "%s %d", field, value)
	}

	return uint32(value), nil
}

func pageOf(
	op, peer string,
	records []notary.EvidenceRecord,
	next *commonv1.LogCursor,
	hasMore bool,
) (notary.EvidencePage, error) {
	position, err := positionOf(op, peer, "a page", next)
	if err != nil {
		return notary.EvidencePage{}, err
	}

	return notary.EvidencePage{Records: records, Next: position, HasMore: hasMore}, nil
}
