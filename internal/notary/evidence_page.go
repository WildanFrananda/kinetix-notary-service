package notary

type EvidencePage struct {
	Records []EvidenceRecord
	Next    LogPosition
	HasMore bool
}

func (p EvidencePage) Follows(stream EvidenceStream, after LogPosition, limit int32) error {
	const op = "notary.EvidencePage.Follows"

	if limit <= 0 || int64(len(p.Records)) > int64(limit) {
		return Errf(
			KindSourceRefused,
			op,
			nil,
			"the %s feed answered %d records to a request for %d",
			stream,
			len(p.Records),
			limit,
		)
	}

	previous := after

	for i, record := range p.Records {
		if record.Stream() != stream {
			return Errf(
				KindSourceRefused,
				op,
				nil,
				"record %d of the %s feed belongs to the %s stream",
				i,
				stream,
				record.Stream(),
			)
		}

		position := record.Position()
		if !position.After(previous) {
			return Errf(
				KindSourceRefused,
				op,
				nil,
				"record %d of the %s feed at %d/%d does not follow %d/%d",
				i,
				stream,
				position.Transaction,
				position.Sequence,
				previous.Transaction,
				previous.Sequence,
			)
		}

		previous = position
	}

	if p.Next != previous {
		return Errf(
			KindSourceRefused,
			op,
			nil,
			"the %s feed's page ends at %d/%d but names %d/%d as the next position",
			stream,
			previous.Transaction,
			previous.Sequence,
			p.Next.Transaction,
			p.Next.Sequence,
		)
	}

	if p.HasMore && len(p.Records) == 0 {
		return Errf(KindSourceRefused, op, nil, "the %s feed sent an empty page and said more follows", stream)
	}

	return nil
}
