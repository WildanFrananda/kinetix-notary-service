package notary

type EvidenceClaim struct {
	Registry RegistryIdentity
	Scheme   uint16
	Header   EvidenceBatchHeader
	Digest   Digest
	Record   EvidenceRecord
	Salt     Salt
	Proof    EvidenceProof
}
