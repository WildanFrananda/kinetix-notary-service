package wiring

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/poseidon"
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
	Scheme       notary.Scheme
	Confidential notary.ConfidentialScheme
	Statements   notary.StatementScheme
	Withholding  notary.WithholdingScheme
	Approvals    notary.ApprovalScheme
	Operations   notary.OperationScheme
}

func provideScheme() schemeResult {
	scheme := keccak.NewScheme()

	return schemeResult{
		Scheme:       scheme,
		Confidential: poseidon.NewScheme(),
		Statements:   scheme,
		Withholding:  scheme,
		Approvals:    scheme,
		Operations:   scheme,
	}
}

type ledgerParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
	Config   config.Daemon
}

type ledgerResult struct {
	dig.Out
	Ledger notary.Ledger
}

func provideLedger(p ledgerParams) ledgerResult {
	return ledgerResult{Ledger: ethereum.NewLedgerOfScheme(p.Client, p.Settings, p.Config.Scheme)}
}
