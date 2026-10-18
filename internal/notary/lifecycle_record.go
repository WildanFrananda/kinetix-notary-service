package notary

type LifecycleRecord struct {
	At                  LogPosition
	OrderNumber         string
	Kind                uint32
	FromStatus          uint32
	ToStatus            uint32
	ReturnNumber        string
	ReturnStatus        uint32
	ActorKind           uint32
	ActorPrincipalID    string
	BuyerPrincipalID    string
	MerchantPrincipalID string
	OccurredAt          EvidenceTime
}

func (LifecycleRecord) Stream() EvidenceStream {
	return StreamOrderLifecycle
}

func (r LifecycleRecord) Position() LogPosition {
	return r.At
}

func (r LifecycleRecord) Parties() []string {
	return distinctParties(r.BuyerPrincipalID, r.MerchantPrincipalID, r.ActorPrincipalID)
}
