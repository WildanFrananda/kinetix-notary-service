package notary

type ConfidentialHeader struct {
	Header          EpochHeader
	TreeRoot        Digest
	TotalCommitment Digest
}
