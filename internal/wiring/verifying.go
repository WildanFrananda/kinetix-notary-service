package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
)

type verifierParams struct {
	dig.In
	Ledger     notary.Ledger
	Scheme     notary.Scheme
	Statements notary.StatementScheme
	Court      notary.Court
}

func provideVerifier(p verifierParams) *verifying.Verifier {
	return verifying.NewVerifier(p.Ledger, p.Scheme, p.Statements, p.Court)
}
