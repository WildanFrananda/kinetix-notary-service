package prooffile

type evidenceRecord struct {
	Transaction string           `json:"transaction"`
	Sequence    string           `json:"sequence"`
	Lifecycle   *lifecycleFields `json:"lifecycle,omitempty"`
	Delivery    *deliveryFields  `json:"delivery,omitempty"`
	Review      *reviewFields    `json:"review,omitempty"`
}
