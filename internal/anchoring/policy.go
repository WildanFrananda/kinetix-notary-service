package anchoring

import (
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Policy struct {
	Instance      string
	Ledger        notary.LedgerIdentity
	Currency      string
	EpochInterval time.Duration
	ReplaceAfter  time.Duration
	LeaseTerm     time.Duration
	FeeCeiling    uint64
	CostWindow    int32
}
