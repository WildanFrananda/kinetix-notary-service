package proving

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Bundle struct {
	Epoch  notary.EpochRecord
	Claims []notary.Claim
}
