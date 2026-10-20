package anchoring

import (
	"math/big"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Observer interface {
	LeaseHeld(held bool)
	LastFinalized(record notary.EpochRecord)
	Abandoned(count int64)
	Funds(balance *big.Int, meanCost uint64, runway time.Duration, known bool)
	FeeCeilingHit()
	Proved(elapsed time.Duration, failed bool)
	Occupancy(accounts, capacity uint64)
}
