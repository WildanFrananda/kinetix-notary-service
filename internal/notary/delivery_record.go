package notary

type DeliveryRecord struct {
	At                LogPosition
	OrderNumber       string
	DriverPrincipalID string
	PhotoSHA256       []byte
	SignatureSHA256   []byte
	DeliveredAt       EvidenceTime
}

func (DeliveryRecord) Stream() EvidenceStream {
	return StreamDelivery
}

func (r DeliveryRecord) Position() LogPosition {
	return r.At
}

func (r DeliveryRecord) Parties() []string {
	return distinctParties(r.DriverPrincipalID)
}
