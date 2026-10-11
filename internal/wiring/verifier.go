package wiring

import (
	"fmt"

	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
)

func Verifier(lookup config.Lookup) (*dig.Container, error) {
	container := dig.New()
	providers := []struct {
		constructor any
		options     []dig.ProvideOption
	}{
		{constructor: func() config.Lookup { return lookup }},
		{constructor: config.LoadVerifier},
		{constructor: ethereumSettings},
		{constructor: ethereum.Dial},
		{constructor: keccak.NewScheme, options: []dig.ProvideOption{dig.As(new(notary.Scheme))}},
		{constructor: ethereum.NewLedger, options: []dig.ProvideOption{dig.As(new(notary.Ledger))}},
		{constructor: verifying.NewVerifier},
	}
	for _, provider := range providers {
		if err := container.Provide(provider.constructor, provider.options...); err != nil {
			return nil, fmt.Errorf("provide %T: %w", provider.constructor, err)
		}
	}
	return container, nil
}

func ethereumSettings(cfg config.Verifier) ethereum.Settings {
	return ethereum.Settings{RPCURL: cfg.RPCURL, Ledger: cfg.Ledger, Timeout: cfg.RPCTimeout}
}
