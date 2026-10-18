package prooffile

type deliveryFields struct {
	OrderNumber       string `json:"orderNumber"`
	DriverPrincipalID string `json:"driverPrincipalId"`
	PhotoSHA256       string `json:"photoSha256"`
	SignatureSHA256   string `json:"signatureSha256,omitempty"`
	DeliveredAt       string `json:"deliveredAt"`
}
