package wiring

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type chainClientParams struct {
	dig.In
	Settings ethereum.Settings
}

type chainClientResult struct {
	dig.Out
	Client *ethclient.Client
}

func provideChainClient(p chainClientParams) (chainClientResult, error) {
	client, err := ethereum.Dial(p.Settings)
	if err != nil {
		return chainClientResult{}, err
	}

	return chainClientResult{Client: client}, nil
}

type schemeResult struct {
	dig.Out
	Scheme      notary.Scheme
	Statements  notary.StatementScheme
	Withholding notary.WithholdingScheme
}

func provideScheme() schemeResult {
	scheme := keccak.NewScheme()

	return schemeResult{Scheme: scheme, Statements: scheme, Withholding: scheme}
}

type ledgerParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
	Scheme   notary.Scheme
}

type ledgerResult struct {
	dig.Out
	Ledger notary.Ledger
}

func provideLedger(p ledgerParams) ledgerResult {
	return ledgerResult{Ledger: ethereum.NewLedger(p.Client, p.Settings, p.Scheme)}
}
