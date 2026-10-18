package notary

type EvidenceVerification struct {
	Verdict  Verdict
	Registry RegistryIdentity
	Batch    AnchoredBatch
	Record   EvidenceRecord
	Reason   string
}
