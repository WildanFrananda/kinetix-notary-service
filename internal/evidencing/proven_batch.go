package evidencing

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type provenBatch struct {
	record notary.EvidenceBatchRecord
	tree   *notary.EvidenceTree
	leaves []notary.Digest
}
