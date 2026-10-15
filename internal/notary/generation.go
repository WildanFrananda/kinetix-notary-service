package notary

type Generation struct {
	ChainID       uint64
	Ledger        Address
	Court         Address
	Withholding   Address
	GenesisEpoch  uint64
	GenesisDigest Digest
	Scheme        uint16
	RegisteredAt  UnixTime
}

func (g Generation) LedgerIdentity() LedgerIdentity {
	return LedgerIdentity{ChainID: g.ChainID, Address: g.Ledger}
}
