package notarytest

import (
	"math/big"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Observer struct {
	Held          bool
	Last          notary.EpochRecord
	AbandonedSeen int64
	Balance       *big.Int
	MeanCost      uint64
	Runway        time.Duration
	RunwayKnown   bool
	CeilingHits   int
}

func (o *Observer) LeaseHeld(held bool)                     { o.Held = held }
func (o *Observer) LastFinalized(record notary.EpochRecord) { o.Last = record }
func (o *Observer) Abandoned(count int64)                   { o.AbandonedSeen = count }
func (o *Observer) FeeCeilingHit()                          { o.CeilingHits++ }

func (o *Observer) Funds(balance *big.Int, meanCost uint64, runway time.Duration, known bool) {
	o.Balance, o.MeanCost, o.Runway, o.RunwayKnown = balance, meanCost, runway, known
}
