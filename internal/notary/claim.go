package notary

type Claim struct {
	Ledger    LedgerIdentity
	Scheme    uint16
	Epoch     uint64
	Entry     Entry
	Proof     InclusionProof
	Opening   *Opening
	Statement *StatementSignature
}
