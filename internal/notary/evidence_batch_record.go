package notary

type EvidenceBatchRecord struct {
	Registry        RegistryIdentity
	Scheme          uint16
	Header          EvidenceBatchHeader
	Digest          Digest
	State           EpochState
	AbandonedReason string
}
