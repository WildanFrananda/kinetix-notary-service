package notary

import "time"

type Attempt struct {
	Epoch   uint64
	Anchor  PreparedAnchor
	Outcome AttemptOutcome
	SentAt  time.Time
	Receipt *AnchorReceipt
}
