package wiring

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"
	"time"

	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	replaceAfter = 5 * time.Minute
	leaseTerm    = 30 * time.Second
	costWindow   = 24
)

type policyParams struct {
	dig.In
	Config config.Daemon
}

func providePolicy(p policyParams) (anchoring.Policy, error) {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "notaryd"
	}

	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		return anchoring.Policy{}, notary.Errf(notary.KindUnknown, "wiring.providePolicy", err, "naming this instance")
	}

	return anchoring.Policy{
		Instance:      host + "-" + hex.EncodeToString(suffix),
		Ledger:        p.Config.Chain.Ledger,
		Council:       notary.CouncilIdentity{ChainID: p.Config.Chain.Ledger.ChainID, Address: p.Config.Council},
		Currency:      p.Config.Currency,
		EpochInterval: p.Config.EpochInterval,
		ReplaceAfter:  replaceAfter,
		LeaseTerm:     leaseTerm,
		FeeCeiling:    p.Config.MaxFeePerGasWei,
		CostWindow:    costWindow,
	}, nil
}

type anchorerParams struct {
	dig.In
	Scheme   notary.Scheme
	Salts    notary.SaltSource
	Secrets  notary.SecretSource
	Source   notary.SnapshotSource
	Store    notary.EpochStore
	Chain    notary.AnchorChain
	Quorum   anchoring.Quorum
	Observer anchoring.Observer
	Clock    func() time.Time
	Log      *slog.Logger
	Policy   anchoring.Policy
	Two      *anchoring.SchemeTwo
}

func provideAnchorer(p anchorerParams) *anchoring.Anchorer {
	return anchoring.NewAnchorer(
		p.Scheme,
		p.Salts,
		p.Secrets,
		p.Source,
		p.Store,
		p.Chain,
		p.Quorum,
		p.Observer,
		p.Clock,
		p.Log,
		p.Policy,
		p.Two,
	)
}

type runnerParams struct {
	dig.In
	Anchorer *anchoring.Anchorer
	Config   config.Daemon
	Clock    func() time.Time
	Log      *slog.Logger
}

func provideRunner(p runnerParams) *anchoring.Runner {
	return anchoring.NewRunner(p.Anchorer, p.Config.Tick, p.Clock, p.Log)
}
