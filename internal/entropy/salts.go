package entropy

import (
	"crypto/rand"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Salts struct{}

var _ notary.SaltSource = Salts{}

func NewSalts() Salts {
	return Salts{}
}

func (Salts) Salt() (notary.Salt, error) {
	var salt notary.Salt
	if _, err := rand.Read(salt[:]); err != nil {
		return notary.Salt{}, notary.Errf(
			notary.KindUnknown,
			"entropy.Salts.Salt",
			err,
			"reading the system's randomness",
		)
	}

	return salt, nil
}
