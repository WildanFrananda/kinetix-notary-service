package plonk

import (
	"math/big"

	"github.com/consensys/gnark/frontend"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func NewAssignment(
	epoch uint64,
	depth int,
	entries []notary.Entry,
	header notary.ConfidentialHeader,
	blinding notary.Blinding,
) (*LiabilityCircuit, error) {
	const op = "plonk.NewAssignment"
	if depth < 0 || depth > notary.MaxConfidentialDepth {
		return nil, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"depth %d is outside 0..%d",
			depth,
			notary.MaxConfidentialDepth,
		)
	}

	capacity := 1 << depth
	if len(entries) > capacity {
		return nil, notary.Errf(
			notary.KindTooManyLeaves,
			op,
			nil,
			"%d accounts exceed the circuit's %d",
			len(entries),
			capacity,
		)
	}

	assignment := NewLiabilityCircuit(depth)
	assignment.Epoch = epoch
	assignment.TreeRoot = new(big.Int).SetBytes(header.TreeRoot[:])
	assignment.TotalCommitment = new(big.Int).SetBytes(header.TotalCommitment[:])
	assignment.LeafCount = len(entries)
	assignment.Blinding = new(big.Int).SetBytes(blinding[:])

	for i := range capacity {
		if i >= len(entries) {
			assignment.High[i] = 0
			assignment.Low[i] = 0
			assignment.Amounts[i] = 0

			continue
		}

		commitment := entries[i].Commitment
		assignment.High[i] = new(big.Int).SetBytes(commitment[:16])
		assignment.Low[i] = new(big.Int).SetBytes(commitment[16:])
		assignment.Amounts[i] = new(big.Int).SetUint64(entries[i].Amount)
	}

	return assignment, nil
}

var _ frontend.Circuit = (*LiabilityCircuit)(nil)
