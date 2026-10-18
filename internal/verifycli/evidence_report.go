package verifycli

type evidenceReport struct {
	Verdict     string `json:"verdict"`
	ChainID     string `json:"chainId,omitempty"`
	Registry    string `json:"registry,omitempty"`
	Stream      string `json:"stream,omitempty"`
	Batch       string `json:"batch,omitempty"`
	ClosedAt    string `json:"closedAt,omitempty"`
	AnchoredAt  string `json:"anchoredAt,omitempty"`
	Transaction string `json:"transaction,omitempty"`
	Sequence    string `json:"sequence,omitempty"`
	Reason      string `json:"reason,omitempty"`
}
