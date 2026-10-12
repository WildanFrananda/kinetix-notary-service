package notary

type EpochRecord struct {
	Ledger          LedgerIdentity
	Header          EpochHeader
	Digest          Digest
	Currency        string
	State           EpochState
	AbandonedReason string
}
