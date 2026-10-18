package notary

import "time"

type BatchAttempt struct {
	Stream  EvidenceStream
	Batch   uint64
	Anchor  PreparedAnchor
	Outcome AttemptOutcome
	SentAt  time.Time
	Receipt *AnchorReceipt
}
