package notary

type EpochState uint8

const (
	StateBuilt EpochState = iota + 1
	StateSubmitted
	StateMined
	StateFinalized
	StateAbandoned
)

func ParseEpochState(text string) (EpochState, error) {
	for state := StateBuilt; state <= StateAbandoned; state++ {
		if state.String() == text {
			return state, nil
		}
	}

	return 0, Errf(KindInvalidArgument, "notary.ParseEpochState", nil, "%q is not an epoch state", text)
}

func (s EpochState) String() string {
	switch s {
	case StateBuilt:
		return "built"
	case StateSubmitted:
		return "submitted"
	case StateMined:
		return "mined"
	case StateFinalized:
		return "finalized"
	case StateAbandoned:
		return "abandoned"
	default:
		return "invalid"
	}
}

func (s EpochState) InFlight() bool {
	return s == StateBuilt || s == StateSubmitted || s == StateMined
}

func (s EpochState) CanBecome(next EpochState) bool {
	switch s {
	case StateBuilt:
		return next == StateSubmitted || next == StateAbandoned
	case StateSubmitted:
		return next == StateMined || next == StateFinalized || next == StateAbandoned
	case StateMined:
		return next == StateFinalized || next == StateSubmitted || next == StateAbandoned
	default:
		return false
	}
}
