package evidencing

import (
	"math/big"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Observer interface {
	EvidenceFeed(stream notary.EvidenceStream, healthy bool)
	EvidenceCollected(stream notary.EvidenceStream, records int)
	EvidencePending(stream notary.EvidenceStream, records int64)
	EvidenceFinalized(record notary.EvidenceBatchRecord)
	EvidenceAbandoned(count int64)
	EvidenceFunds(balance *big.Int, meanCost uint64, runway time.Duration, known bool)
	EvidenceFeeCeilingHit()
}
