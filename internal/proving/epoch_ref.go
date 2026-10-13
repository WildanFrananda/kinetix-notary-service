package proving

import (
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

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
	if text == "latest" {
		return Latest(), nil
	}

	number, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return EpochRef{}, notary.Errf(
			notary.KindInvalidArgument,
			"proving.ParseEpochRef",
			nil,
			"%q is neither an epoch number nor latest",
			text,
		)
	}

	return Number(number), nil
}
