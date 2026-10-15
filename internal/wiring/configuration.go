package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
)

type verifierConfigParams struct {
	dig.In
	Lookup config.Lookup
}

func provideVerifierConfig(p verifierConfigParams) (config.Verifier, error) {
	return config.LoadVerifier(p.Lookup)
}

type verifierSettingsParams struct {
	dig.In
	Config config.Verifier
}

func provideVerifierSettings(p verifierSettingsParams) ethereum.Settings {
	return ethereum.Settings{
		RPCURL:    p.Config.RPCURL,
		Directory: p.Config.Directory,
		Timeout:   p.Config.RPCTimeout,
	}
}
