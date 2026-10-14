package registering

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Issued struct {
	Court         notary.CourtIdentity
	Epoch         uint64
	Registrations []notary.IssuedRegistration
}
