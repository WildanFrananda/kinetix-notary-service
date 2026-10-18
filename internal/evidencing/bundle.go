package evidencing

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Bundle struct {
	Claims  []notary.EvidenceClaim
	Pending int
}
