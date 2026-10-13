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
	Account    *account  `json:"account,omitempty"`
	Statement  *seal     `json:"statement,omitempty"`
}

type sibling struct {
	Hash string `json:"hash"`
	Sum  string `json:"sum"`
}

type account struct {
	Holder string `json:"holder"`
	Kind   string `json:"kind"`
	Salt   string `json:"salt"`
}

type seal struct {
	Court     string `json:"court"`
	Signature string `json:"signature"`
}
