package notarytest

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Blindings struct {
	Seed  uint64
	Drawn int
	Err   error
	Zero  bool
}

var _ notary.BlindingSource = (*Blindings)(nil)

func (b *Blindings) Blinding() (notary.Blinding, error) {
	if b.Err != nil {
		return notary.Blinding{}, b.Err
	}

	if b.Zero {
		return notary.Blinding{}, nil
	}

	b.Drawn++
	seed := binary.BigEndian.AppendUint64(nil, b.Seed)
	sum := sha256.Sum256(binary.BigEndian.AppendUint64(seed, uint64(b.Drawn)))

	var blinding notary.Blinding
	copy(blinding[:], sum[:])

	return blinding, nil
}
