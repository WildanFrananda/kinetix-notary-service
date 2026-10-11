package prooffile

type document struct {
	Format     string    `json:"format"`
	Scheme     uint16    `json:"scheme"`
	ChainID    string    `json:"chainId"`
	Ledger     string    `json:"ledger"`
	Epoch      string    `json:"epoch"`
	Commitment string    `json:"commitment"`
	Amount     string    `json:"amount"`
	Index      string    `json:"index"`
	Siblings   []sibling `json:"siblings"`
}

type sibling struct {
	Hash string `json:"hash"`
	Sum  string `json:"sum"`
}
