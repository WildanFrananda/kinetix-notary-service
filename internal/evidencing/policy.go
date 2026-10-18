package evidencing

import (
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Policy struct {
	Instance      string
	Registry      notary.RegistryIdentity
	BatchInterval time.Duration
	MaxBatch      int32
	PageSize      int32
	PagesPerTick  int
	ReplaceAfter  time.Duration
	LeaseTerm     time.Duration
	FeeCeiling    uint64
	CostWindow    int32
}
