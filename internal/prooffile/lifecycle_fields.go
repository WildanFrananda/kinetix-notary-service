package prooffile

type lifecycleFields struct {
	OrderNumber         string `json:"orderNumber"`
	Kind                uint32 `json:"kind"`
	FromStatus          uint32 `json:"fromStatus"`
	ToStatus            uint32 `json:"toStatus"`
	ReturnNumber        string `json:"returnNumber"`
	ReturnStatus        uint32 `json:"returnStatus"`
	ActorKind           uint32 `json:"actorKind"`
	ActorPrincipalID    string `json:"actorPrincipalId"`
	BuyerPrincipalID    string `json:"buyerPrincipalId"`
	MerchantPrincipalID string `json:"merchantPrincipalId"`
	OccurredAt          string `json:"occurredAt"`
}
