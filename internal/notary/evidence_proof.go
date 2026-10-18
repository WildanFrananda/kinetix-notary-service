package notary

type EvidenceProof struct {
	Index    uint64
	Siblings []Digest
}
