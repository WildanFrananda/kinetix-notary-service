package notary

import "time"

type Attempt struct {
	Epoch   uint64
	Kind    AttemptKind
	Anchor  PreparedAnchor
	Outcome AttemptOutcome
	SentAt  time.Time
	Receipt *AnchorReceipt
}
