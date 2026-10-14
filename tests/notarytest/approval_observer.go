package notarytest

type ApprovalObserver struct {
	Outcomes []string
}

func (o *ApprovalObserver) Approved() { o.Outcomes = append(o.Outcomes, "approved") }

func (o *ApprovalObserver) Refused() { o.Outcomes = append(o.Outcomes, "refused") }

func (o *ApprovalObserver) Failed() { o.Outcomes = append(o.Outcomes, "failed") }
