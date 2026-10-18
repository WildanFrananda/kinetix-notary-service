package notary

type AnchoredBatch struct {
	Stream     EvidenceStream
	Number     uint64
	Root       Digest
	Digest     Digest
	LeafCount  uint64
	Through    LogPosition
	ClosedAt   UnixTime
	AnchoredAt UnixTime
}
