package answering

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type requestKey struct {
	venue        int
	registration notary.Digest
	epoch        uint64
}
