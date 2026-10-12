package notarytest

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Salts struct {
	Seed  uint64
	Drawn int
	Err   error
	Zero  bool
}

var _ notary.SaltSource = (*Salts)(nil)

func (s *Salts) Salt() (notary.Salt, error) {
	if s.Err != nil {
		return notary.Salt{}, s.Err
	}

	if s.Zero {
		return notary.Salt{}, nil
	}

	s.Drawn++
	seed := binary.BigEndian.AppendUint64(nil, s.Seed)

	return notary.Salt(sha256.Sum256(binary.BigEndian.AppendUint64(seed, uint64(s.Drawn)))), nil
}
