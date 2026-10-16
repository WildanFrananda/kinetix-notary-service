package sponsorcli

type report struct {
	Action      string `json:"action"`
	Epoch       string `json:"epoch,omitempty"`
	Account     string `json:"account,omitempty"`
	Verdict     string `json:"verdict"`
	Operation   string `json:"operation,omitempty"`
	Transaction string `json:"transaction,omitempty"`
	GasCost     string `json:"gasCost,omitempty"`
	Reason      string `json:"reason,omitempty"`
}
