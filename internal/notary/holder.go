package notary

import (
	"unicode"
	"unicode/utf8"
)

const MaxHolderBytes = 255

type Holder struct {
	value string
}

func NewHolder(value string) (Holder, error) {
	const op = "notary.NewHolder"
	if value == "" || len(value) > MaxHolderBytes {
		return Holder{}, Errf(
			KindInvalidArgument,
			op,
			nil,
			"a holder is 1 to %d bytes, not %d",
			MaxHolderBytes,
			len(value),
		)
	}

	if !utf8.ValidString(value) {
		return Holder{}, Errf(KindInvalidArgument, op, nil, "a holder must be valid UTF-8")
	}

	for _, r := range value {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return Holder{}, Errf(KindInvalidArgument, op, nil, "a holder carries no spaces or control characters")
		}
	}

	return Holder{value: value}, nil
}

func (h Holder) String() string {
	return h.value
}

func (h Holder) IsZero() bool {
	return h.value == ""
}
