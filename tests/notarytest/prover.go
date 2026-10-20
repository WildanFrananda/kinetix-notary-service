package notarytest

import (
	"context"
	"crypto/sha256"
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Prover struct {
	Err      error
	During   func()
	Requests []notary.ProofRequest
}

var _ notary.EpochProver = (*Prover)(nil)

func (p *Prover) Prove(_ context.Context, request notary.ProofRequest) ([]byte, error) {
	p.Requests = append(p.Requests, request)
	if p.During != nil {
		p.During()
	}

	if p.Err != nil {
		return nil, p.Err
	}

	seed := binary.BigEndian.AppendUint64(append([]byte{}, request.TreeRoot[:]...), request.Number)
	proof := sha256.Sum256(append(seed, request.TotalCommitment[:]...))

	return proof[:], nil
}
