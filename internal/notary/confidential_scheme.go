package notary

type ConfidentialScheme interface {
	Version() uint16
	LeafHash(epoch uint64, entry Entry) Digest
	NodeHash(left, right Digest) (Digest, error)
	EmptyHash() Digest
	TotalCommitment(total uint64, blinding Blinding) Digest
	BindingRoot(treeRoot, totalCommitment Digest) Digest
	EpochDigest(header EpochHeader) Digest
}
