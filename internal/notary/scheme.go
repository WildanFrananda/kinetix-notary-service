package notary

type Scheme interface {
	Version() uint16
	LeafHash(epoch uint64, entry Entry) Digest
	NodeHash(left, right Node) Digest
	EmptyHash() Digest
	EpochDigest(header EpochHeader) Digest
}
