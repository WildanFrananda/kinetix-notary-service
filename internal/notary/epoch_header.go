package notary

type EpochHeader struct {
	Number         uint64
	PreviousDigest Digest
	Root           Digest
	Total          uint64
	LeafCount      uint64
	AsOf           UnixTime
}
