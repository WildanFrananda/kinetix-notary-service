package notary

type EvidenceBatchHeader struct {
	Stream         EvidenceStream
	Number         uint64
	PreviousDigest Digest
	Root           Digest
	LeafCount      uint64
	Through        LogPosition
	ClosedAt       UnixTime
}
