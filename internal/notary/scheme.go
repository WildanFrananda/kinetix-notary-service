package notary

type Scheme interface {
	Version() uint16
	Commitment(holder Holder, kind AccountKind, salt Salt) Digest
	LeafHash(epoch uint64, entry Entry) Digest
	NodeHash(left, right Node) Digest
	EmptyHash() Digest
	EpochDigest(header EpochHeader) Digest
}
