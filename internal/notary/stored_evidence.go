package notary

type StoredEvidence struct {
	Entry   EvidenceEntry
	Batched bool
	Batch   uint64
	Index   uint64
}
