package proving

import (
	"errors"
	"math"
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const MaxEpochNumber = uint64(math.MaxInt64)

type EpochRef struct {
	latest bool
	number uint64
}

func Latest() EpochRef {
	return EpochRef{latest: true}
}

func Number(number uint64) EpochRef {
	return EpochRef{number: number}
}

func ParseEpochRef(text string) (EpochRef, error) {
	const op = "proving.ParseEpochRef"

	if text == "latest" {
		return Latest(), nil
	}

	number, err := strconv.ParseUint(text, 10, 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return EpochRef{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%q is neither an epoch number nor latest",
			text,
		)
	}

	if err != nil || number > MaxEpochNumber {
		return EpochRef{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%q is above %d, the highest epoch number",
			text,
			MaxEpochNumber,
		)
	}

	if len(text) > 1 && text[0] == '0' {
		return EpochRef{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"%q is written with a leading zero; an epoch number has none",
			text,
		)
	}

	return Number(number), nil
}
