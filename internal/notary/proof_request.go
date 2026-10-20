package notary

type ProofRequest struct {
	Number          uint64
	Depth           int
	Entries         []Entry
	TreeRoot        Digest
	TotalCommitment Digest
	Blinding        Blinding
}
