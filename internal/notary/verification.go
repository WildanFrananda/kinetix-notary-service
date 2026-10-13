package notary

type Verification struct {
	Verdict  Verdict
	Ledger   LedgerIdentity
	Epoch    AnchoredEpoch
	Reason   string
	Opening  *Opening
	Court    CourtIdentity
	Check    StatementCheck
	Signer   Address
	Faults   uint64
	Provable bool
}
