package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
)

type verifierParams struct {
	dig.In
	Directory  notary.Directory
	Readers    notary.GenerationReaders
	Scheme     notary.Scheme
	Statements notary.StatementScheme
}

func provideVerifier(p verifierParams) *verifying.DirectoryVerifier {
	return verifying.NewDirectoryVerifier(p.Directory, p.Readers, p.Scheme, p.Statements)
}
