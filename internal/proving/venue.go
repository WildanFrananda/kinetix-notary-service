package proving

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Venue struct {
	Ledger   notary.LedgerIdentity
	Court    notary.Court
	Identity notary.CourtIdentity
}
