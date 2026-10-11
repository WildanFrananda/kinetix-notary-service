package notary

type AnchoredEpoch struct {
	Number     uint64
	Root       Digest
	Digest     Digest
	Total      uint64
	LeafCount  uint64
	AsOf       UnixTime
	AnchoredAt UnixTime
}
