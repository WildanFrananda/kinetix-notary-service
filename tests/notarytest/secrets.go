package notarytest

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Secrets struct {
	Seed  uint64
	Drawn int
	Err   error
	Zero  bool
}

var _ notary.SecretSource = (*Secrets)(nil)

func (s *Secrets) AccountSecret() (notary.AccountSecret, error) {
	if s.Err != nil {
		return notary.AccountSecret{}, s.Err
	}

	if s.Zero {
		return notary.AccountSecret{}, nil
	}

	s.Drawn++
	seed := binary.BigEndian.AppendUint64([]byte("secret"), s.Seed)

	return notary.AccountSecret(sha256.Sum256(binary.BigEndian.AppendUint64(seed, uint64(s.Drawn)))), nil
}
