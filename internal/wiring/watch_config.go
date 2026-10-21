package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type watchConfigParams struct {
	dig.In
	Lookup config.Lookup
}

func provideWatchConfig(p watchConfigParams) (config.Watch, error) {
	cfg, err := config.LoadWatch(p.Lookup)
	if err != nil {
		return config.Watch{}, notary.Errf(notary.KindInvalidArgument, "wiring.provideWatchConfig", err, "")
	}

	return cfg, nil
}

type watchLedgerSettingsParams struct {
	dig.In
	Config config.Watch
}

func provideWatchLedgerSettings(p watchLedgerSettingsParams) ethereum.Settings {
	chain := p.Config.Chain

	return ethereum.Settings{
		RPCURL:  chain.RPCURL,
		Ledger:  chain.Ledger,
		Court:   chain.Court,
		Council: p.Config.Council,
		Timeout: chain.RPCTimeout,
		Rate:    chain.RPCRate,
	}
}
