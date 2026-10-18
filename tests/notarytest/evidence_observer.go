package notarytest

import (
	"math/big"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type EvidenceObserver struct {
	FeedUp      map[notary.EvidenceStream]bool
	Collected   map[notary.EvidenceStream]int
	Pending     map[notary.EvidenceStream]int64
	Finalized   map[notary.EvidenceStream]notary.EvidenceBatchRecord
	Abandoned   int64
	Balance     *big.Int
	MeanCost    uint64
	Runway      time.Duration
	RunwayKnown bool
	CeilingHits int
}

func NewEvidenceObserver() *EvidenceObserver {
	return &EvidenceObserver{
		FeedUp:    map[notary.EvidenceStream]bool{},
		Collected: map[notary.EvidenceStream]int{},
		Pending:   map[notary.EvidenceStream]int64{},
		Finalized: map[notary.EvidenceStream]notary.EvidenceBatchRecord{},
	}
}

func (o *EvidenceObserver) EvidenceFeed(stream notary.EvidenceStream, healthy bool) {
	o.FeedUp[stream] = healthy
}

func (o *EvidenceObserver) EvidenceCollected(stream notary.EvidenceStream, records int) {
	o.Collected[stream] += records
}

func (o *EvidenceObserver) EvidencePending(stream notary.EvidenceStream, records int64) {
	o.Pending[stream] = records
}

func (o *EvidenceObserver) EvidenceFinalized(record notary.EvidenceBatchRecord) {
	o.Finalized[record.Header.Stream] = record
}

func (o *EvidenceObserver) EvidenceAbandoned(count int64) { o.Abandoned = count }

func (o *EvidenceObserver) EvidenceFunds(balance *big.Int, meanCost uint64, runway time.Duration, known bool) {
	o.Balance, o.MeanCost, o.Runway, o.RunwayKnown = balance, meanCost, runway, known
}

func (o *EvidenceObserver) EvidenceFeeCeilingHit() { o.CeilingHits++ }
