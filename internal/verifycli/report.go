package verifycli

type report struct {
	Verdict    string `json:"verdict"`
	ChainID    string `json:"chainId,omitempty"`
	Ledger     string `json:"ledger,omitempty"`
	Epoch      string `json:"epoch,omitempty"`
	Total      string `json:"total,omitempty"`
	AsOf       string `json:"asOf,omitempty"`
	AnchoredAt string `json:"anchoredAt,omitempty"`
	Holder     string `json:"holder,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Amount     string `json:"amount,omitempty"`
	Court      string `json:"court,omitempty"`
	Statement  string `json:"statement,omitempty"`
	Faults     string `json:"faults,omitempty"`
	Provable   string `json:"provable,omitempty"`
	Reason     string `json:"reason,omitempty"`
}
