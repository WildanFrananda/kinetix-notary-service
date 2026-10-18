package notary

type EvidenceStreamHead struct {
	Known        bool
	Sealed       bool
	Head         Digest
	GenesisBatch uint64
	NextBatch    uint64
	LastClosedAt uint64
	Through      LogPosition
	BlockTime    UnixTime
}

func (h EvidenceStreamHead) AtGenesis() bool {
	return h.NextBatch == h.GenesisBatch
}
