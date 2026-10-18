package wiring

import (
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/feeds"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	evidencePageSize     = 500
	evidencePagesPerTick = 20
	evidenceMaxBatch     = 1 << 16
)

type feedsParams struct {
	dig.In
	Config  config.Daemon
	Metrics *metrics.Metrics
}

type feedsResult struct {
	dig.Out
	Lifecycle  *feeds.LifecycleFeed
	Deliveries *feeds.DeliveryFeed
	Reviews    *feeds.ReviewFeed
	Feeds      []notary.EvidenceFeed
}

func provideFeeds(p feedsParams) (feedsResult, error) {
	identity := mesh.Identity{PKIDir: p.Config.PKIDir, TrustDomains: p.Config.TrustDomains}
	settings := func(endpoint string) feeds.Settings {
		return feeds.Settings{Endpoint: endpoint, Identity: identity, Timeout: p.Config.FeedTimeout}
	}

	lifecycle, err := feeds.DialLifecycle(settings(p.Config.OrderAddress), p.Metrics.ClientInterceptors("order")...)
	if err != nil {
		return feedsResult{}, err
	}

	deliveries, err := feeds.DialDeliveries(
		settings(p.Config.MatchingAddress),
		p.Metrics.ClientInterceptors("matching")...,
	)
	if err != nil {
		_ = lifecycle.Close()
		return feedsResult{}, err
	}

	reviews, err := feeds.DialReviews(settings(p.Config.ReviewAddress), p.Metrics.ClientInterceptors("review")...)
	if err != nil {
		_ = lifecycle.Close()
		_ = deliveries.Close()

		return feedsResult{}, err
	}

	return feedsResult{
		Lifecycle:  lifecycle,
		Deliveries: deliveries,
		Reviews:    reviews,
		Feeds:      []notary.EvidenceFeed{lifecycle, deliveries, reviews},
	}, nil
}

type evidenceChainParams struct {
	dig.In
	Client *ethclient.Client
	Config config.Daemon
	Scheme notary.EvidenceScheme
}

type evidenceChainResult struct {
	dig.Out
	Chain    notary.EvidenceChain
	Registry notary.EvidenceRegistry
	Identity notary.RegistryIdentity
}

func provideEvidenceChain(p evidenceChainParams) (evidenceChainResult, error) {
	key, err := ethereum.ParseEvidenceAnchorerKey(p.Config.EvidenceKey.Reveal())
	if err != nil {
		return evidenceChainResult{}, err
	}

	identity := notary.RegistryIdentity{ChainID: p.Config.Chain.Ledger.ChainID, Address: p.Config.Registry}
	timeout := p.Config.Chain.RPCTimeout

	return evidenceChainResult{
		Chain:    ethereum.NewEvidenceWriter(p.Client, identity, timeout, key),
		Registry: ethereum.NewEvidenceRegistry(p.Client, identity, p.Scheme, timeout),
		Identity: identity,
	}, nil
}

type evidencePolicyParams struct {
	dig.In
	Anchoring anchoring.Policy
	Identity  notary.RegistryIdentity
	Config    config.Daemon
}

func provideEvidencePolicy(p evidencePolicyParams) evidencing.Policy {
	return evidencing.Policy{
		Instance:      p.Anchoring.Instance,
		Registry:      p.Identity,
		BatchInterval: p.Config.EvidenceBatch,
		MaxBatch:      evidenceMaxBatch,
		PageSize:      evidencePageSize,
		PagesPerTick:  evidencePagesPerTick,
		ReplaceAfter:  replaceAfter,
		LeaseTerm:     leaseTerm,
		FeeCeiling:    p.Config.MaxFeePerGasWei,
		CostWindow:    costWindow,
	}
}

type evidencerParams struct {
	dig.In
	Scheme   notary.EvidenceScheme
	Salts    notary.SaltSource
	Feeds    []notary.EvidenceFeed
	Store    notary.EvidenceStore
	Chain    notary.EvidenceChain
	Observer evidencing.Observer
	Clock    func() time.Time
	Log      *slog.Logger
	Policy   evidencing.Policy
	Config   config.Daemon
}

type evidencerResult struct {
	dig.Out
	Evidencer *evidencing.Evidencer
	Runner    *evidencing.Runner
}

func provideEvidencer(p evidencerParams) evidencerResult {
	collector := evidencing.NewCollector(p.Scheme, p.Salts, p.Feeds, p.Store, p.Observer, p.Log, p.Policy)
	closer := evidencing.NewCloser(p.Scheme, p.Store, p.Chain, p.Observer, p.Clock, p.Log, p.Policy)
	evidencer := evidencing.NewEvidencer(p.Store, collector, closer, p.Log, p.Policy)

	return evidencerResult{
		Evidencer: evidencer,
		Runner:    evidencing.NewRunner(evidencer, p.Config.Tick, p.Clock, p.Log),
	}
}

type registryCheckParams struct {
	dig.In
	Registry    notary.EvidenceRegistry
	Chain       notary.EvidenceChain
	Identity    notary.RegistryIdentity
	Liabilities notary.AnchorChain
}

func provideRegistryCheck(p registryCheckParams) *evidencing.RegistryCheck {
	return evidencing.NewRegistryCheck(p.Registry, p.Chain, p.Identity, p.Liabilities.Anchorer())
}

type evidenceHandlersParams struct {
	dig.In
	Authenticator httpapi.Authenticator
	Scheme        notary.EvidenceScheme
	Source        notary.EvidenceProofSource
	Log           *slog.Logger
}

type evidenceHandlersResult struct {
	dig.Out
	Orders  *httpapi.EvidenceHandler `name:"order_evidence"`
	Reviews *httpapi.EvidenceHandler `name:"review_evidence"`
}

func provideEvidenceHandlers(p evidenceHandlersParams) evidenceHandlersResult {
	prover := evidencing.NewProver(p.Scheme, p.Source)

	return evidenceHandlersResult{
		Orders:  httpapi.NewOrderEvidenceHandler(p.Authenticator, prover, p.Log),
		Reviews: httpapi.NewReviewEvidenceHandler(p.Authenticator, prover, p.Log),
	}
}
