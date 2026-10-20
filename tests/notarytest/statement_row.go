package notarytest

import (
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type StatementRow struct {
	Holder       notary.Holder
	Epoch        uint64
	Day          time.Time
	Fingerprint  notary.Digest
	State        string
	Attempts     uint32
	Next         time.Time
	Notification string
	Reason       string
	Planned      time.Time
}
