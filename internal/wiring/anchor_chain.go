package wiring

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type anchorChainParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
	Config   config.Daemon
}

type anchorChainResult struct {
	dig.Out
	Chain notary.AnchorChain
}

func provideAnchorChain(p anchorChainParams) (anchorChainResult, error) {
	key, err := ethereum.ParseAnchorerKey(p.Config.AnchorerKey.Reveal())
	if err != nil {
		return anchorChainResult{}, err
	}

	return anchorChainResult{Chain: ethereum.NewWriter(p.Client, p.Settings, key)}, nil
}
