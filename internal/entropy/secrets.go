package entropy

import (
	"crypto/rand"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Secrets struct{}

var _ notary.SecretSource = Secrets{}

func NewSecrets() Secrets {
	return Secrets{}
}

func (Secrets) AccountSecret() (notary.AccountSecret, error) {
	var secret notary.AccountSecret
	if _, err := rand.Read(secret[:]); err != nil {
		return notary.AccountSecret{}, notary.Errf(
			notary.KindUnknown,
			"entropy.Secrets.AccountSecret",
			err,
			"reading the system's randomness",
		)
	}

	return secret, nil
}
