package plonk

import (
	"context"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EpochProver struct {
	prover *Prover
}

var _ notary.EpochProver = (*EpochProver)(nil)

func NewEpochProver(prover *Prover) *EpochProver {
	return &EpochProver{prover: prover}
}

func (e *EpochProver) Prove(ctx context.Context, request notary.ProofRequest) ([]byte, error) {
	const op = "plonk.EpochProver.Prove"
	if err := ctx.Err(); err != nil {
		return nil, notary.Errf(notary.KindUnknown, op, err, "epoch %d", request.Number)
	}

	if request.Depth != e.prover.Depth() {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"epoch %d asks for depth %d, this prover's circuit has depth %d",
			request.Number,
			request.Depth,
			e.prover.Depth(),
		)
	}

	header := notary.ConfidentialHeader{TreeRoot: request.TreeRoot, TotalCommitment: request.TotalCommitment}

	assignment, err := NewAssignment(request.Number, request.Depth, request.Entries, header, request.Blinding)
	if err != nil {
		return nil, err
	}

	return e.prover.Prove(assignment)
}
