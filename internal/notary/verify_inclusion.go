package notary

import "math/bits"

func VerifyInclusion(scheme Scheme, epoch AnchoredEpoch, entry Entry, proof InclusionProof) error {
	const op = "notary.VerifyInclusion"
	if proof.Index >= epoch.LeafCount {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"index %d is outside the %d accounts of epoch %d",
			proof.Index,
			epoch.LeafCount,
			epoch.Number,
		)
	}

	if want := DepthFor(epoch.LeafCount); len(proof.Siblings) != want {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"epoch %d needs %d siblings, the proof carries %d",
			epoch.Number,
			want,
			len(proof.Siblings),
		)
	}

	total := entry.Amount

	for _, sibling := range proof.Siblings {
		var carry uint64

		total, carry = bits.Add64(total, sibling.Sum, 0)
		if carry != 0 {
			return Errf(KindProofRejected, op, nil, "the proof's sums do not fit in 64 bits")
		}
	}

	if total != epoch.Total {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"the proof sums to %d, epoch %d holds %d",
			total,
			epoch.Number,
			epoch.Total,
		)
	}

	node := Node{Hash: scheme.LeafHash(epoch.Number, entry), Sum: entry.Amount}

	position := proof.Index
	for _, sibling := range proof.Siblings {
		var err error
		if position&1 == 0 {
			node, err = Parent(scheme, node, sibling)
		} else {
			node, err = Parent(scheme, sibling, node)
		}

		if err != nil {
			return Errf(KindProofRejected, op, err, "")
		}

		position >>= 1
	}

	if node.Hash != epoch.Root {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"the proof reaches %s, epoch %d is anchored at %s",
			node.Hash.Hex(),
			epoch.Number,
			epoch.Root.Hex(),
		)
	}

	return nil
}
