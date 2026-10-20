package postgres

import (
	"math"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/storegen"
)

func confidentialParams(op string, plan *notary.ConfidentialPlan, params *storegen.InsertEpochParams) error {
	if plan == nil {
		return nil
	}

	if plan.Depth < 0 || plan.Depth > notary.MaxConfidentialDepth || len(plan.Proof) == 0 {
		return notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"a scheme-%d epoch needs a depth in 0..%d and a proof",
			plan.Scheme,
			notary.MaxConfidentialDepth,
		)
	}

	depth := int16(plan.Depth)
	params.Scheme = int32(plan.Scheme)
	params.TreeRoot = plan.TreeRoot[:]
	params.TotalCommitment = plan.TotalCommitment[:]
	params.Blinding = plan.Blinding[:]
	params.Depth = &depth
	params.Proof = plan.Proof

	return nil
}

func confidentialOf(op string, row storegen.Epoch) (*notary.ConfidentialPlan, error) {
	if row.TreeRoot == nil {
		return nil, nil
	}

	treeRoot, err := digestOf(op, row.TreeRoot)
	if err != nil {
		return nil, err
	}

	commitment, err := digestOf(op, row.TotalCommitment)
	if err != nil {
		return nil, err
	}

	var blinding notary.Blinding
	if len(row.Blinding) != len(blinding) || row.Depth == nil {
		return nil, notary.Errf(
			notary.KindStoreUnavailable,
			op,
			nil,
			"epoch %d holds a broken scheme-2 record",
			row.Number,
		)
	}

	copy(blinding[:], row.Blinding)

	if row.Scheme < 1 || row.Scheme > math.MaxUint16 {
		return nil, notary.Errf(notary.KindStoreUnavailable, op, nil, "epoch %d has scheme %d", row.Number, row.Scheme)
	}

	return &notary.ConfidentialPlan{
		Scheme:          uint16(row.Scheme),
		TreeRoot:        treeRoot,
		TotalCommitment: commitment,
		Blinding:        blinding,
		Depth:           int(*row.Depth),
		Proof:           row.Proof,
	}, nil
}
