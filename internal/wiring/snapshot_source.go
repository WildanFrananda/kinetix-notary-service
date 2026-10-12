package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/entropy"
	"github.com/WildanFrananda/kinetix-notary-service/internal/mesh"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/payment"
)

type snapshotSourceParams struct {
	dig.In
	Config  config.Daemon
	Metrics *metrics.Metrics
}

type snapshotSourceResult struct {
	dig.Out
	Payment *payment.Source
	Source  notary.SnapshotSource
}

func provideSnapshotSource(p snapshotSourceParams) (snapshotSourceResult, error) {
	source, err := payment.Dial(
		payment.Settings{
			Endpoint: p.Config.PaymentAddress,
			Identity: mesh.Identity{PKIDir: p.Config.PKIDir, TrustDomains: p.Config.TrustDomains},
			Timeout:  p.Config.SnapshotTimeout,
		},
		p.Metrics.ClientInterceptors("payment")...,
	)
	if err != nil {
		return snapshotSourceResult{}, err
	}

	return snapshotSourceResult{Payment: source, Source: source}, nil
}

type saltsResult struct {
	dig.Out
	Salts notary.SaltSource
}

func provideSalts() saltsResult {
	return saltsResult{Salts: entropy.NewSalts()}
}
