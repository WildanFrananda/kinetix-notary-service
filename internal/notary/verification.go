package notary

type Verification struct {
	Verdict Verdict
	Ledger  LedgerIdentity
	Epoch   AnchoredEpoch
	Reason  string
}
