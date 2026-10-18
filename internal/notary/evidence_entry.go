package notary

type EvidenceEntry struct {
	Record       EvidenceRecord
	Salt         Salt
	RecordDigest Digest
	Leaf         Digest
}

func EnterEvidence(scheme EvidenceScheme, salts SaltSource, record EvidenceRecord) (EvidenceEntry, error) {
	const op = "notary.EnterEvidence"

	if !record.Stream().Valid() {
		return EvidenceEntry{}, Errf(
			KindInvalidArgument,
			op,
			nil,
			"stream %d is not an evidence stream",
			record.Stream(),
		)
	}

	digest, err := scheme.RecordDigest(record)
	if err != nil {
		return EvidenceEntry{}, err
	}

	salt, err := salts.Salt()
	if err != nil {
		return EvidenceEntry{}, Errf(KindUnknown, op, err, "drawing a salt")
	}

	if salt.IsZero() {
		return EvidenceEntry{}, Errf(KindInvalidArgument, op, nil, "the salt source returned an all-zero salt")
	}

	return EvidenceEntry{
		Record:       record,
		Salt:         salt,
		RecordDigest: digest,
		Leaf:         scheme.EvidenceLeaf(record.Stream(), salt, digest),
	}, nil
}

func (e EvidenceEntry) Holds(scheme EvidenceScheme) error {
	const op = "notary.EvidenceEntry.Holds"

	position := e.Record.Position()

	digest, err := scheme.RecordDigest(e.Record)
	if err != nil {
		return err
	}

	if digest != e.RecordDigest {
		return Errf(
			KindDisagreement,
			op,
			nil,
			"the %s record at %d/%d hashes to %s, the database keeps %s",
			e.Record.Stream(),
			position.Transaction,
			position.Sequence,
			digest.Hex(),
			e.RecordDigest.Hex(),
		)
	}

	if leaf := scheme.EvidenceLeaf(e.Record.Stream(), e.Salt, digest); leaf != e.Leaf {
		return Errf(
			KindDisagreement,
			op,
			nil,
			"the %s record at %d/%d leads to leaf %s, the database keeps %s",
			e.Record.Stream(),
			position.Transaction,
			position.Sequence,
			leaf.Hex(),
			e.Leaf.Hex(),
		)
	}

	return nil
}
