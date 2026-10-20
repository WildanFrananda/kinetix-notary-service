package notary

func RebuildEpochTree(
	scheme Scheme,
	confidential ConfidentialScheme,
	record EpochRecord,
	entries []Entry,
) (EpochTree, error) {
	const op = "notary.RebuildEpochTree"

	number := record.Header.Number
	if record.Confidential == nil {
		tree, err := BuildTree(scheme, number, entries)
		if err != nil {
			return nil, err
		}

		if root := tree.Root(); root.Hash != record.Header.Root || root.Sum != record.Header.Total {
			return nil, Errf(
				KindDisagreement,
				op,
				nil,
				"the accounts stored for epoch %d do not rebuild the root it anchored",
				number,
			)
		}

		return tree, nil
	}

	if confidential == nil || confidential.Version() != record.Confidential.Scheme {
		return nil, Errf(
			KindWrongLedger,
			op,
			nil,
			"epoch %d is scheme %d, which this notary does not build",
			number,
			record.Confidential.Scheme,
		)
	}

	tree, err := BuildConfidentialTree(confidential, number, record.Confidential.Depth, entries)
	if err != nil {
		return nil, err
	}

	commitment := confidential.TotalCommitment(tree.Total(), record.Confidential.Blinding)
	if tree.TreeRoot() != record.Confidential.TreeRoot ||
		commitment != record.Confidential.TotalCommitment ||
		confidential.BindingRoot(tree.TreeRoot(), commitment) != record.Header.Root {
		return nil, Errf(
			KindDisagreement,
			op,
			nil,
			"the accounts stored for epoch %d do not rebuild the statement it anchored",
			number,
		)
	}

	return tree, nil
}
