package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/epochwatch"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type epochWatchParams struct {
	dig.In
	Config  config.Daemon
	Metrics *metrics.Metrics
}

type epochWatchResult struct {
	dig.Out
	Client *epochwatch.Client
	Watch  notary.EpochWatch
}

func provideEpochWatch(p epochWatchParams) (epochWatchResult, error) {
	client, err := epochwatch.Dial(
		epochwatch.Settings{
			Endpoint: p.Config.WatchAddress,
			Identity: mesh.Identity{PKIDir: p.Config.PKIDir, TrustDomains: p.Config.TrustDomains},
			Timeout:  p.Config.SnapshotTimeout,
		},
		p.Config.Currency,
		p.Metrics.ClientInterceptors(epochwatch.Service)...,
	)
	if err != nil {
		return epochWatchResult{}, err
	}

	return epochWatchResult{Client: client, Watch: client}, nil
}

type quorumParams struct {
	dig.In
	Config    config.Daemon
	Council   notary.Council
	Watch     notary.EpochWatch
	Approvals notary.ApprovalScheme
}

func provideQuorum(p quorumParams) (anchoring.Quorum, error) {
	key, err := ethereum.ParseAnchorerKey(p.Config.AnchorerKey.Reveal())
	if err != nil {
		return anchoring.Quorum{}, err
	}

	return anchoring.Quorum{
		Council: p.Council,
		Watch:   p.Watch,
		Signer:  ethereum.NewApprovalSigner(key),
		Scheme:  p.Approvals,
	}, nil
}
