package prooffile

type reviewFields struct {
	ReviewID            string `json:"reviewId"`
	Kind                uint32 `json:"kind"`
	OrderID             string `json:"orderId"`
	SubjectID           string `json:"subjectId"`
	AuthorPrincipalID   string `json:"authorPrincipalId"`
	MerchantPrincipalID string `json:"merchantPrincipalId"`
	Rating              uint32 `json:"rating"`
	CommentSHA256       string `json:"commentSha256,omitempty"`
	CreatedAt           string `json:"createdAt"`
}
