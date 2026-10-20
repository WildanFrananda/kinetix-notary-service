package notary

type ConfidentialEpoch struct {
	Number    uint64
	TreeRoot  Digest
	LeafCount uint64
	Depth     int
}
