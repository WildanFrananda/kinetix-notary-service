package notary

type LedgerHead struct {
	NextEpoch     uint64
	Head          Digest
	LastAsOf      uint64
	GenesisEpoch  uint64
	GenesisDigest Digest
	Sealed        bool
}

func (h LedgerHead) AtGenesis() bool {
	return h.NextEpoch == h.GenesisEpoch && h.Head == h.GenesisDigest
}
