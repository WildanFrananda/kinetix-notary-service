package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type daemonConfigParams struct {
	dig.In
	Lookup config.Lookup
}

func provideDaemonConfig(p daemonConfigParams) (config.Daemon, error) {
	cfg, err := config.LoadDaemon(p.Lookup)
	if err != nil {
		return config.Daemon{}, notary.Errf(notary.KindInvalidArgument, "wiring.provideDaemonConfig", err, "")
	}

	return cfg, nil
}

type daemonLedgerSettingsParams struct {
	dig.In
	Config config.Daemon
}

func provideDaemonLedgerSettings(p daemonLedgerSettingsParams) ethereum.Settings {
	chain := p.Config.Chain

	return ethereum.Settings{
		RPCURL:  chain.RPCURL,
		Ledger:  chain.Ledger,
		Court:   chain.Court,
		Timeout: chain.RPCTimeout,
	}
}
