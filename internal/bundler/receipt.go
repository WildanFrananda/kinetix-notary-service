package bundler

type receipt struct {
	Success       bool               `json:"success"`
	ActualGasCost string             `json:"actualGasCost"`
	Reason        string             `json:"reason"`
	Receipt       transactionReceipt `json:"receipt"`
}
