package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
)

type metricsParams struct {
	dig.In
	Build metrics.Build
}

type metricsResult struct {
	dig.Out
	Metrics  *metrics.Metrics
	Observer anchoring.Observer
}

func provideMetrics(p metricsParams) metricsResult {
	m := metrics.New(p.Build)
	for _, route := range httpapi.Routes {
		m.SeedRoute("GET", route)
	}

	m.SeedCall("payment", paymentv1.LiabilityService_StreamLiabilitySnapshot_FullMethodName)

	return metricsResult{Metrics: m, Observer: m}
}
