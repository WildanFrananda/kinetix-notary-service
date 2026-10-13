package metrics

import (
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/watching"
)

type Build struct {
	Service string
	Version string
}

type Metrics struct {
	registry     *prometheus.Registry
	httpRequests *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
	grpcCalls    *prometheus.CounterVec
	lastEpoch    prometheus.Gauge
	lastAsOf     prometheus.Gauge
	abandoned    prometheus.Gauge
	balance      prometheus.Gauge
	cost         prometheus.Gauge
	runway       prometheus.Gauge
	ceilingHits  prometheus.Counter
	leaseHeld    prometheus.Gauge
	faults       prometheus.Gauge
}

var (
	_ anchoring.Observer = (*Metrics)(nil)
	_ watching.Observer  = (*Metrics)(nil)
)

func New(build Build) *Metrics {
	registry := prometheus.NewRegistry()
	factory := func(c prometheus.Collector) { registry.MustRegister(c) }
	factory(collectors.NewGoCollector())
	factory(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	info := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kinetix_build_info",
			Help: "Always 1. Carries the service name and build version as labels.",
		},
		[]string{"service", "version"},
	)
	factory(info)
	info.WithLabelValues(build.Service, build.Version).Set(1)

	m := &Metrics{
		registry: registry,
		httpRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kinetix_http_requests_total",
				Help: "HTTP requests served, by method, route template and status.",
			},
			[]string{"method", "route", "status"},
		),
		httpDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "kinetix_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route"},
		),
		grpcCalls: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kinetix_grpc_client_calls_total",
				Help: "gRPC calls this service made, by peer, method and final code.",
			},
			[]string{"peer", "grpc_method", "grpc_code"},
		),
		lastEpoch: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_last_finalized_epoch",
			Help: "Number of the newest epoch the ledger holds at a finalized block.",
		}),
		lastAsOf: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_last_finalized_as_of_seconds",
			Help: "Snapshot time of the newest final epoch, in Unix seconds.",
		}),
		abandoned: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_abandoned_epochs",
			Help: "Epochs the chain refused for good; each one needs a person.",
		}),
		balance: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_anchorer_balance_wei",
			Help: "Balance of the anchorer key in wei.",
		}),
		cost: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_anchor_cost_wei",
			Help: "Mean cost of recent final epochs in wei, execution plus L1 data fee.",
		}),
		runway: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_runway_seconds",
			Help: "How long the anchorer balance lasts at the mean cost; NaN until an epoch is final.",
		}),
		ceilingHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kinetix_notary_fee_ceiling_hits_total",
			Help: "Times the fee ceiling held an anchor back.",
		}),
		leaseHeld: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_lease_held",
			Help: "1 while this process is the one writer, 0 while it waits.",
		}),
		faults: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_proven_faults",
			Help: "Faults the court holds against this notary's statements; NaN until the court is read.",
		}),
	}
	for _, c := range []prometheus.Collector{
		m.httpRequests,
		m.httpDuration,
		m.grpcCalls,
		m.lastEpoch,
		m.lastAsOf,
		m.abandoned,
		m.balance,
		m.cost,
		m.runway,
		m.ceilingHits,
		m.leaseHeld,
		m.faults,
	} {
		factory(c)
	}

	m.runway.Set(math.NaN())
	m.faults.Set(math.NaN())

	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registry})
}

func (m *Metrics) SeedRoute(method, route string) {
	m.httpRequests.WithLabelValues(method, route, strconv.Itoa(http.StatusOK))
	m.httpDuration.WithLabelValues(method, route)
}

func (m *Metrics) SeedCall(peer, method string) {
	m.grpcCalls.WithLabelValues(peer, method, "OK")
}

func (m *Metrics) LeaseHeld(held bool) {
	if held {
		m.leaseHeld.Set(1)
		return
	}

	m.leaseHeld.Set(0)
}

func (m *Metrics) LastFinalized(record notary.EpochRecord) {
	m.lastEpoch.Set(float64(record.Header.Number))
	m.lastAsOf.Set(float64(record.Header.AsOf.Seconds()))
}

func (m *Metrics) Abandoned(count int64) {
	m.abandoned.Set(float64(count))
}

func (m *Metrics) Funds(balance *big.Int, meanCost uint64, runway time.Duration, known bool) {
	if balance != nil {
		wei, _ := new(big.Float).SetInt(balance).Float64()
		m.balance.Set(wei)
	}

	m.cost.Set(float64(meanCost))

	if !known {
		m.runway.Set(math.NaN())
		return
	}

	m.runway.Set(runway.Seconds())
}

func (m *Metrics) FeeCeilingHit() {
	m.ceilingHits.Inc()
}

func (m *Metrics) ProvenFaults(total uint64) {
	m.faults.Set(float64(total))
}
