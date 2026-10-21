package wiring

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type councilParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
}

type councilResult struct {
	dig.Out
	Council notary.Council
}

func provideCouncil(p councilParams) councilResult {
	return councilResult{Council: ethereum.NewCouncil(p.Client, p.Settings)}
}

type headReaderResult struct {
	dig.Out
	Heads  notary.LedgerHeads
	Reader *ethereum.HeadReader
}

func provideHeadReader(p councilParams) headReaderResult {
	reader := ethereum.NewHeadReader(p.Client, p.Settings)

	return headReaderResult{Heads: reader, Reader: reader}
}
