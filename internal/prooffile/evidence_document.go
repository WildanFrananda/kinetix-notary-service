package prooffile

type evidenceDocument struct {
	Format   string         `json:"format"`
	Scheme   uint16         `json:"scheme"`
	ChainID  string         `json:"chainId"`
	Registry string         `json:"registry"`
	Batch    evidenceBatch  `json:"batch"`
	Record   evidenceRecord `json:"record"`
	Salt     string         `json:"salt"`
	Index    string         `json:"index"`
	Siblings []string       `json:"siblings"`
}
