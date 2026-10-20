package wiring

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
)

type directoryParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
}

type directoryResult struct {
	dig.Out
	Directory notary.Directory
}

func provideDirectory(p directoryParams) directoryResult {
	return directoryResult{Directory: ethereum.NewDirectory(p.Client, p.Settings)}
}

type generationReadersParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
}

type generationReadersResult struct {
	dig.Out
	Readers notary.GenerationReaders
}

func provideGenerationReaders(p generationReadersParams) generationReadersResult {
	return generationReadersResult{Readers: ethereum.NewReaders(p.Client, p.Settings)}
}

type earlierGenerationsParams struct {
	dig.In
	Directory notary.Directory
	Settings  ethereum.Settings
}

type earlierGenerationsResult struct {
	dig.Out
	Earlier []notary.Generation `name:"earlier"`
}

func provideEarlierGenerations(p earlierGenerationsParams) (earlierGenerationsResult, error) {
	lineage, err := p.Directory.Lineage(context.Background())
	if err != nil {
		return earlierGenerationsResult{}, err
	}

	earlier := make([]notary.Generation, 0, len(lineage))

	for _, generation := range lineage {
		if generation.LedgerIdentity() == p.Settings.Ledger {
			break
		}

		if generation.ChainID == p.Settings.Ledger.ChainID {
			earlier = append(earlier, generation)
		}
	}

	return earlierGenerationsResult{Earlier: earlier}, nil
}

func settingsOf(settings ethereum.Settings, generation notary.Generation) ethereum.Settings {
	settings.Ledger = generation.LedgerIdentity()
	settings.Court = generation.Court
	settings.Withholding = generation.Withholding

	return settings
}

type directoryCheckParams struct {
	dig.In
	Directory notary.Directory
	Settings  ethereum.Settings
}

func provideDirectoryCheck(p directoryCheckParams) *proving.DirectoryCheck {
	return proving.NewDirectoryCheck(
		p.Directory,
		notary.Generation{
			ChainID:     p.Settings.Ledger.ChainID,
			Ledger:      p.Settings.Ledger.Address,
			Court:       p.Settings.Court,
			Withholding: p.Settings.Withholding,
		},
	)
}

type sponsorshipWatcherParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
	Config   config.Daemon
	Court    notary.Court
	Observer sponsoring.Observer
	Log      *slog.Logger
}

func provideSponsorshipWatcher(p sponsorshipWatcherParams) *sponsoring.Watcher {
	return sponsoring.NewWatcher(
		ethereum.NewPaymaster(p.Client, p.Settings.Ledger.ChainID, p.Config.Paymaster, p.Settings.Timeout),
		p.Court,
		notary.Generation{
			ChainID:     p.Settings.Ledger.ChainID,
			Court:       p.Settings.Court,
			Withholding: p.Settings.Withholding,
		},
		p.Observer,
		p.Log,
		p.Config.Tick,
	)
}
