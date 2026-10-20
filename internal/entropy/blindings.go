package entropy

import (
	"crypto/rand"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Blindings struct{}

var _ notary.BlindingSource = Blindings{}

func NewBlindings() Blindings {
	return Blindings{}
}

func (Blindings) Blinding() (notary.Blinding, error) {
	var blinding notary.Blinding
	if _, err := rand.Read(blinding[:]); err != nil {
		return notary.Blinding{}, notary.Errf(
			notary.KindUnknown,
			"entropy.Blindings.Blinding",
			err,
			"reading the system's randomness",
		)
	}

	return blinding, nil
}
