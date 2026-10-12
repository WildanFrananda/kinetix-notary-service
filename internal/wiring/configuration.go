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

type ledgerSettingsParams struct {
	dig.In
	Config config.Verifier
}

func provideLedgerSettings(p ledgerSettingsParams) ethereum.Settings {
	return ethereum.Settings{RPCURL: p.Config.RPCURL, Ledger: p.Config.Ledger, Timeout: p.Config.RPCTimeout}
}
