package notary

type AttemptOutcome uint8

const (
	AttemptPending AttemptOutcome = iota + 1
	AttemptIncluded
	AttemptSuperseded
	AttemptReverted
)

func ParseAttemptOutcome(text string) (AttemptOutcome, error) {
	for outcome := AttemptPending; outcome <= AttemptReverted; outcome++ {
		if outcome.String() == text {
			return outcome, nil
		}
	}

	return 0, Errf(KindInvalidArgument, "notary.ParseAttemptOutcome", nil, "%q is not an attempt outcome", text)
}

func (o AttemptOutcome) String() string {
	switch o {
	case AttemptPending:
		return "pending"
	case AttemptIncluded:
		return "included"
	case AttemptSuperseded:
		return "superseded"
	case AttemptReverted:
		return "reverted"
	default:
		return "invalid"
	}
}
