package notary

type EvidenceScheme interface {
	EvidenceVersion() uint16
	RecordDigest(record EvidenceRecord) (Digest, error)
	EvidenceLeaf(stream EvidenceStream, salt Salt, record Digest) Digest
	EvidenceNode(left, right Digest) Digest
	EvidenceEmpty() Digest
	BatchDigest(header EvidenceBatchHeader) Digest
}
