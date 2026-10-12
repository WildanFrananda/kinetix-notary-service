package wiring

import (
	"fmt"

	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
)

func Verifier(lookup config.Lookup) (*dig.Container, error) {
	container := dig.New()

	providers := []any{
		func() config.Lookup { return lookup },
		provideVerifierConfig,
		provideLedgerSettings,
		provideChainClient,
		provideScheme,
		provideLedger,
		provideVerifier,
	}
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, fmt.Errorf("provide %T: %w", provider, err)
		}
	}

	return container, nil
}
