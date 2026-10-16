package sponsorcli

import "github.com/WildanFrananda/kinetix-notary-service/internal/notary"

type Dependencies struct {
	Sponsor    Sponsorship
	Directory  notary.Directory
	Verifier   ClaimVerifier
	Statements notary.StatementScheme
	Scheme     uint16
}
