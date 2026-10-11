package notary

type Claim struct {
	Ledger LedgerIdentity
	Epoch  uint64
	Entry  Entry
	Proof  InclusionProof
}
