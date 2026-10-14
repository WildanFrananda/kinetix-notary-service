package notary

type RequestState uint8

const (
	RequestNone RequestState = iota
	RequestOpen
	RequestAnswered
	RequestUnanswered
)

func (s RequestState) String() string {
	switch s {
	case RequestNone:
		return "none"
	case RequestOpen:
		return "open"
	case RequestAnswered:
		return "answered"
	case RequestUnanswered:
		return "unanswered"
	default:
		return "unknown"
	}
}
