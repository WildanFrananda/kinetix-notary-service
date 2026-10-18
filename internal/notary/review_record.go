package notary

const ReviewKindDriver uint32 = 2

type ReviewRecord struct {
	At                  LogPosition
	ReviewID            string
	Kind                uint32
	OrderID             string
	SubjectID           string
	AuthorPrincipalID   string
	MerchantPrincipalID string
	Rating              uint32
	CommentSHA256       []byte
	CreatedAt           EvidenceTime
}

func (ReviewRecord) Stream() EvidenceStream {
	return StreamReview
}

func (r ReviewRecord) Position() LogPosition {
	return r.At
}

func (r ReviewRecord) Parties() []string {
	if r.Kind == ReviewKindDriver {
		return distinctParties(r.AuthorPrincipalID, r.MerchantPrincipalID, r.SubjectID)
	}

	return distinctParties(r.AuthorPrincipalID, r.MerchantPrincipalID)
}
