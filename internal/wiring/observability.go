package wiring

import (
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/epochwatch"
	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/httpapi"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	fulfillmentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/fulfillment/v1"
	notaryv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notary/v1"
	notificationv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notification/v1"
	orderv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/order/v1"
	paymentv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/payment/v1"
	reviewv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/review/v1"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
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
	Sponsors  sponsoring.Observer
	Evidence  evidencing.Observer
	Delivery  delivering.Observer
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
	m.SeedCall("order", orderv1.OrderLifecycleService_LifecycleEventsSince_FullMethodName)
	m.SeedCall("matching", fulfillmentv1.DeliveryEvidenceService_DeliveriesSince_FullMethodName)
	m.SeedCall("review", reviewv1.ReviewEvidenceService_ReviewsSince_FullMethodName)
	m.SeedCall("notification", notificationv1.NotificationService_Notify_FullMethodName)

	return metricsResult{Metrics: m, Observer: m, Watching: m, Answering: m, Sponsors: m, Evidence: m, Delivery: m}
}
