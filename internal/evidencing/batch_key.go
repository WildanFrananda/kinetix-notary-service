package evidencing

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type batchKey struct {
	stream notary.EvidenceStream
	number uint64
}
