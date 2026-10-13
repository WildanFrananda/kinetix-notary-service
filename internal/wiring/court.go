package wiring

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type courtParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
}

type courtResult struct {
	dig.Out
	Court    notary.Court
	Identity notary.CourtIdentity
}

func provideCourt(p courtParams) courtResult {
	return courtResult{
		Court:    ethereum.NewCourt(p.Client, p.Settings),
		Identity: notary.CourtIdentity{ChainID: p.Settings.Ledger.ChainID, Address: p.Settings.Court},
	}
}
