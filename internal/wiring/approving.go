package wiring

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/approving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/payment"
	notaryv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notary/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

type watchMetricsParams struct {
	dig.In
	Build metrics.Build
}

type watchMetricsResult struct {
	dig.Out
	Metrics  *metrics.Watch
	Observer approving.Observer
	RPC      ethereum.RPCObserver
}

func provideWatchMetrics(p watchMetricsParams) watchMetricsResult {
	m := metrics.NewWatch(p.Build)

	for _, route := range httpapi.WatchRoutes {
		m.SeedRoute("GET", route)
	}

	m.SeedCall("payment", paymentv1.LiabilityRecordService_GetLiabilitySnapshotRecord_FullMethodName)
	m.SeedServed(notaryv1.EpochWatchService_ApproveEpoch_FullMethodName)

	return watchMetricsResult{Metrics: m, Observer: m, RPC: m}
}

type watchSignerParams struct {
	dig.In
	Config config.Watch
}

type watchSignerResult struct {
	dig.Out
	Signer notary.ApprovalSigner
}

func provideWatchSigner(p watchSignerParams) (watchSignerResult, error) {
	key, err := ethereum.ParseWatchKey(p.Config.WatchKey.Reveal())
	if err != nil {
		return watchSignerResult{}, err
	}

	return watchSignerResult{Signer: ethereum.NewApprovalSigner(key)}, nil
}

type snapshotRecordsParams struct {
	dig.In
	Config  config.Watch
	Metrics *metrics.Watch
}

type snapshotRecordsResult struct {
	dig.Out
	Payment *payment.Records
	Records notary.SnapshotRecords
}

func provideSnapshotRecords(p snapshotRecordsParams) (snapshotRecordsResult, error) {
	records, err := payment.DialRecords(
		payment.Settings{
			Endpoint: p.Config.PaymentAddress,
			Identity: mesh.Identity{PKIDir: p.Config.PKIDir, TrustDomains: p.Config.TrustDomains},
			Timeout:  p.Config.Chain.RPCTimeout,
		},
		p.Metrics.ClientInterceptors("payment")...,
	)
	if err != nil {
		return snapshotRecordsResult{}, err
	}

	return snapshotRecordsResult{Payment: records, Records: records}, nil
}

type approverParams struct {
	dig.In
	Config       config.Watch
	Records      notary.SnapshotRecords
	Heads        notary.LedgerHeads
	Council      notary.Council
	Signer       notary.ApprovalSigner
	Scheme       notary.Scheme
	Confidential notary.ConfidentialScheme
	Proofs       notary.ProofReader
	Approvals    notary.ApprovalScheme
	Observer     approving.Observer
	Log          *slog.Logger
}

func provideApprover(p approverParams) *approving.Approver {
	return approving.NewApprover(
		p.Records,
		p.Heads,
		p.Council,
		p.Signer,
		p.Scheme,
		p.Confidential,
		p.Proofs,
		p.Approvals,
		notary.CouncilIdentity{ChainID: p.Config.Chain.Ledger.ChainID, Address: p.Config.Council},
		p.Config.Chain.Ledger.Address,
		p.Observer,
		p.Log,
	)
}

type proofReaderParams struct {
	dig.In
	Client   *ethclient.Client
	Settings ethereum.Settings
}

type proofReaderResult struct {
	dig.Out
	Proofs notary.ProofReader
}

func provideProofReader(p proofReaderParams) proofReaderResult {
	return proofReaderResult{Proofs: ethereum.NewProofReader(p.Client, p.Settings)}
}
