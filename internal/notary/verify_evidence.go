package notary

func VerifyEvidence(
	scheme EvidenceScheme,
	batch AnchoredBatch,
	salt Salt,
	record Digest,
	proof EvidenceProof,
) error {
	const op = "notary.VerifyEvidence"
	if proof.Index >= batch.LeafCount {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"index %d is outside the %d records of %s batch %d",
			proof.Index,
			batch.LeafCount,
			batch.Stream,
			batch.Number,
		)
	}

	if want := DepthFor(batch.LeafCount); len(proof.Siblings) != want {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"%s batch %d needs %d siblings, the proof carries %d",
			batch.Stream,
			batch.Number,
			want,
			len(proof.Siblings),
		)
	}

	node := scheme.EvidenceLeaf(batch.Stream, salt, record)

	position := proof.Index
	for _, sibling := range proof.Siblings {
		if position&1 == 0 {
			node = scheme.EvidenceNode(node, sibling)
		} else {
			node = scheme.EvidenceNode(sibling, node)
		}

		position >>= 1
	}

	if node != batch.Root {
		return Errf(
			KindProofRejected,
			op,
			nil,
			"the proof reaches %s, %s batch %d holds %s",
			node.Hex(),
			batch.Stream,
			batch.Number,
			batch.Root.Hex(),
		)
	}

	return nil
}
