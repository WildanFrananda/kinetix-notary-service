package notary

type AttemptKind uint8

const (
	AttemptForAnchor AttemptKind = iota
	AttemptForProof
)

func ParseAttemptKind(text string) (AttemptKind, error) {
	for kind := AttemptForAnchor; kind <= AttemptForProof; kind++ {
		if kind.String() == text {
			return kind, nil
		}
	}

	return 0, Errf(KindInvalidArgument, "notary.ParseAttemptKind", nil, "%q is not an attempt kind", text)
}

func (k AttemptKind) String() string {
	switch k {
	case AttemptForAnchor:
		return "anchor"
	case AttemptForProof:
		return "proof"
	default:
		return "invalid"
	}
}
