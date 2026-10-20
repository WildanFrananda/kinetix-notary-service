package notary

func VerifyConfidentialInclusion(
	scheme ConfidentialScheme,
	epoch ConfidentialEpoch,
	entry Entry,
	proof InclusionProof,
) error {
	const op = "notary.VerifyConfidentialInclusion"
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

	if len(proof.Siblings) != epoch.Depth {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"epoch %d needs %d siblings, the proof carries %d",
			epoch.Number,
			epoch.Depth,
			len(proof.Siblings),
		)
	}

	node := scheme.LeafHash(epoch.Number, entry)

	position := proof.Index
	for level, sibling := range proof.Siblings {
		if sibling.Sum != 0 {
			return Errf(
				KindProofRejected,
				op,
				nil,
				"sibling %d carries a sum, which scheme %d never publishes",
				level,
				scheme.Version(),
			)
		}

		var err error
		if position&1 == 0 {
			node, err = scheme.NodeHash(node, sibling.Hash)
		} else {
			node, err = scheme.NodeHash(sibling.Hash, node)
		}

		if err != nil {
			return Errf(KindProofRejected, op, err, "sibling %d", level)
		}

		position >>= 1
	}

	if node != epoch.TreeRoot {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"the proof reaches %s, epoch %d is anchored at %s",
			node.Hex(),
			epoch.Number,
			epoch.TreeRoot.Hex(),
		)
	}

	return nil
}
