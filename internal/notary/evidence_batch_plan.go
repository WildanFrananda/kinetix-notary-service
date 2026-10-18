package notary

type EvidenceBatchPlan struct {
	Scheme  uint16
	Header  EvidenceBatchHeader
	Digest  Digest
	Entries []EvidenceEntry
}

func PlanEvidenceBatch(
	scheme EvidenceScheme,
	stream EvidenceStream,
	number uint64,
	previous Digest,
	closedAt UnixTime,
	entries []EvidenceEntry,
) (EvidenceBatchPlan, error) {
	const op = "notary.PlanEvidenceBatch"

	if closedAt.IsZero() {
		return EvidenceBatchPlan{}, Errf(KindInvalidArgument, op, nil, "a batch closes at a time after the Unix epoch")
	}

	leaves := make([]Digest, len(entries))

	var through LogPosition

	for i, entry := range entries {
		if entry.Record.Stream() != stream {
			return EvidenceBatchPlan{}, Errf(
				KindInvalidArgument,
				op,
				nil,
				"entry %d belongs to the %s stream, not %s",
				i,
				entry.Record.Stream(),
				stream,
			)
		}

		position := entry.Record.Position()
		if i > 0 && !position.After(through) {
			return EvidenceBatchPlan{}, Errf(KindInvalidArgument, op, nil, "entry %d is out of feed order", i)
		}

		if err := entry.Holds(scheme); err != nil {
			return EvidenceBatchPlan{}, err
		}

		through = position
		leaves[i] = entry.Leaf
	}

	tree, err := BuildEvidenceTree(scheme, stream, leaves)
	if err != nil {
		return EvidenceBatchPlan{}, err
	}

	header := EvidenceBatchHeader{
		Stream:         stream,
		Number:         number,
		PreviousDigest: previous,
		Root:           tree.Root(),
		LeafCount:      tree.LeafCount(),
		Through:        through,
		ClosedAt:       closedAt,
	}

	return EvidenceBatchPlan{
		Scheme:  scheme.EvidenceVersion(),
		Header:  header,
		Digest:  scheme.BatchDigest(header),
		Entries: entries,
	}, nil
}
