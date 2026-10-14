package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/epochwatch"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	notaryv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notary/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
	"github.com/WildanFrananda/kinetix-notary-service/internal/watching"
)

type metricsParams struct {
	dig.In
	Build metrics.Build
}

type metricsResult struct {
	dig.Out
	Metrics   *metrics.Metrics
	Observer  anchoring.Observer
	Watching  watching.Observer
	Answering answering.Observer
}

func provideMetrics(p metricsParams) metricsResult {
	m := metrics.New(p.Build)

	for _, route := range httpapi.Routes {
		method := "GET"
		if route == httpapi.RegistrationsRoute {
			method = "POST"
		}

		m.SeedRoute(method, route)
	}

	m.SeedCall("payment", paymentv1.LiabilityService_StreamLiabilitySnapshot_FullMethodName)
	m.SeedCall(epochwatch.Service, notaryv1.EpochWatchService_ApproveEpoch_FullMethodName)

	return metricsResult{Metrics: m, Observer: m, Watching: m, Answering: m}
}
