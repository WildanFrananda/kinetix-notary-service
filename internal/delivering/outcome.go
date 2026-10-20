package delivering

type Outcome uint8

const (
	OutcomeSent Outcome = iota + 1
	OutcomeRefused
	OutcomePostponed
)

var Outcomes = []Outcome{OutcomeSent, OutcomeRefused, OutcomePostponed}

func (o Outcome) String() string {
	switch o {
	case OutcomeSent:
		return "sent"
	case OutcomeRefused:
		return "refused"
	case OutcomePostponed:
		return "postponed"
	default:
		return "unknown"
	}
}
