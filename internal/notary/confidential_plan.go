package notary

type ConfidentialPlan struct {
	Scheme          uint16
	TreeRoot        Digest
	TotalCommitment Digest
	Blinding        Blinding
	Depth           int
	Proof           []byte
}
