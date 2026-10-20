package notarytest

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Proofs struct {
	DepthValue int
	Proven     bool
	Err        error
	Asked      int
}

var _ notary.ProofReader = (*Proofs)(nil)

func (p *Proofs) Depth(context.Context) (int, error) {
	return p.DepthValue, p.Err
}

func (p *Proofs) IsProven(context.Context, uint64, notary.Digest, uint64) (bool, error) {
	p.Asked++

	return p.Proven, p.Err
}
