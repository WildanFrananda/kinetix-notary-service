package verifycli

type report struct {
	Verdict    string `json:"verdict"`
	ChainID    string `json:"chainId,omitempty"`
	Ledger     string `json:"ledger,omitempty"`
	Epoch      string `json:"epoch,omitempty"`
	Total      string `json:"total,omitempty"`
	AsOf       string `json:"asOf,omitempty"`
	AnchoredAt string `json:"anchoredAt,omitempty"`
	Reason     string `json:"reason,omitempty"`
}
