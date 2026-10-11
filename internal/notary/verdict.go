package notary

type Verdict uint8

const (
	VerdictIncluded Verdict = iota + 1
	VerdictRejected
)

func (v Verdict) String() string {
	switch v {
	case VerdictIncluded:
		return "included"
	case VerdictRejected:
		return "rejected"
	default:
		return "unset"
	}
}
