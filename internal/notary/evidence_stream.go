package notary

type EvidenceStream uint8

const (
	StreamOrderLifecycle EvidenceStream = 1
	StreamDelivery       EvidenceStream = 2
	StreamReview         EvidenceStream = 3
)

func EvidenceStreams() []EvidenceStream {
	return []EvidenceStream{
		StreamOrderLifecycle,
		StreamDelivery,
		StreamReview,
	}
}

func (s EvidenceStream) Valid() bool {
	return s >= StreamOrderLifecycle && s <= StreamReview
}

func (s EvidenceStream) String() string {
	switch s {
	case StreamOrderLifecycle:
		return "order_lifecycle"
	case StreamDelivery:
		return "delivery"
	case StreamReview:
		return "review"
	default:
		return "unknown"
	}
}

func EvidenceStreamNamed(name string) (EvidenceStream, error) {
	for _, stream := range EvidenceStreams() {
		if stream.String() == name {
			return stream, nil
		}
	}

	return 0, Errf(KindInvalidArgument, "notary.EvidenceStreamNamed", nil, "%q is not an evidence stream", name)
}
