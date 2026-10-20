package metrics

import (
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/watching"
)

const (
	passServed = "served"
	passFailed = "failed"
)

type Build struct {
	Service string
	Version string
}

type Metrics struct {
	edge
	evidence
	statements

	registry     *prometheus.Registry
	lastEpoch    prometheus.Gauge
	lastAsOf     prometheus.Gauge
	abandoned    prometheus.Gauge
	balance      prometheus.Gauge
	cost         prometheus.Gauge
	runway       prometheus.Gauge
	ceilingHits  prometheus.Counter
	proving      prometheus.Gauge
	provingFails prometheus.Counter
	accounts     prometheus.Gauge
	capacity     prometheus.Gauge
	leaseHeld    prometheus.Gauge
	faults       prometheus.Gauge
	silences     prometheus.Gauge
	openRequests prometheus.Gauge
	unanswerable prometheus.Gauge
	answerer     prometheus.Gauge
	passes       *prometheus.CounterVec
	deposit      prometheus.Gauge
	stake        prometheus.Gauge
	current      prometheus.Gauge
}

var (
	_ anchoring.Observer  = (*Metrics)(nil)
	_ watching.Observer   = (*Metrics)(nil)
	_ answering.Observer  = (*Metrics)(nil)
	_ sponsoring.Observer = (*Metrics)(nil)
	_ evidencing.Observer = (*Metrics)(nil)
	_ delivering.Observer = (*Metrics)(nil)
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
		edge:       newEdge(),
		evidence:   newEvidence(),
		statements: newStatements(),
		registry:   registry,
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
		proving: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_proving_seconds",
			Help: "How long the last scheme-2 proof took to make, failed or not.",
		}),
		provingFails: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kinetix_notary_proving_failures_total",
			Help: "Scheme-2 proofs the prover did not make.",
		}),
		accounts: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_accounts",
			Help: "Accounts in the newest snapshot an epoch was planned from.",
		}),
		capacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_account_capacity",
			Help: "Accounts one epoch of the current generation can hold; 0 for scheme 1, which has no fixed capacity.",
		}),
		leaseHeld: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_lease_held",
			Help: "1 while this process is the one writer, 0 while it waits.",
		}),
		faults: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_proven_faults",
			Help: "Faults the court holds against this notary's statements; NaN until the court is read.",
		}),
		silences: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_withholding_silences",
			Help: "Requests the withholding court recorded as unanswered; NaN until the court is read.",
		}),
		openRequests: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_withholding_open_requests",
			Help: "Statement requests this notary has seen open and not yet answered.",
		}),
		unanswerable: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_withholding_unanswerable_requests",
			Help: "Open statement requests this notary cannot answer from its own records.",
		}),
		answerer: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_answerer_balance_wei",
			Help: "The answerer's balance; NaN until it is read.",
		}),
		passes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kinetix_notary_withholding_passes_total",
				Help: "Passes over the withholding court, by outcome: served, or failed and retried next tick.",
			},
			[]string{"outcome"},
		),
		deposit: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_paymaster_deposit_wei",
			Help: "The paymaster's deposit in the EntryPoint, which pays holders' sponsored calls; NaN until it is read.",
		}),
		stake: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_paymaster_stake_wei",
			Help: "The paymaster's locked stake in the EntryPoint, 0 once it is unlocked; NaN until it is read.",
		}),
		current: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_paymaster_current",
			Help: "1 while the paymaster sponsors this notary's courts and their current signer, 0 when it does not.",
		}),
	}
	m.edge.register(factory)
	m.evidence.register(factory)
	m.statements.register(factory)

	for _, c := range []prometheus.Collector{
		m.lastEpoch,
		m.lastAsOf,
		m.abandoned,
		m.balance,
		m.cost,
		m.runway,
		m.ceilingHits,
		m.proving,
		m.provingFails,
		m.accounts,
		m.capacity,
		m.leaseHeld,
		m.faults,
		m.silences,
		m.openRequests,
		m.unanswerable,
		m.answerer,
		m.passes,
		m.deposit,
		m.stake,
		m.current,
	} {
		factory(c)
	}

	m.runway.Set(math.NaN())
	m.faults.Set(math.NaN())
	m.silences.Set(math.NaN())
	m.answerer.Set(math.NaN())
	m.deposit.Set(math.NaN())
	m.stake.Set(math.NaN())
	m.current.Set(math.NaN())
	m.passes.WithLabelValues(passServed)
	m.passes.WithLabelValues(passFailed)

	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registry})
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

func (m *Metrics) Proved(elapsed time.Duration, failed bool) {
	m.proving.Set(elapsed.Seconds())

	if failed {
		m.provingFails.Inc()
	}
}

func (m *Metrics) Occupancy(accounts, capacity uint64) {
	m.accounts.Set(float64(accounts))
	m.capacity.Set(float64(capacity))
}

func (m *Metrics) ProvenFaults(total uint64) {
	m.faults.Set(float64(total))
}

func (m *Metrics) Silences(total uint64) {
	m.silences.Set(float64(total))
}

func (m *Metrics) OpenRequests(count int) {
	m.openRequests.Set(float64(count))
}

func (m *Metrics) Unanswerable(count int) {
	m.unanswerable.Set(float64(count))
}

func (m *Metrics) PaymasterDeposit(wei *big.Int) {
	if wei == nil {
		return
	}

	value, _ := new(big.Float).SetInt(wei).Float64()
	m.deposit.Set(value)
}

func (m *Metrics) PaymasterStake(wei *big.Int) {
	if wei == nil {
		return
	}

	value, _ := new(big.Float).SetInt(wei).Float64()
	m.stake.Set(value)
}

func (m *Metrics) PaymasterCurrent(current bool) {
	if current {
		m.current.Set(1)
		return
	}

	m.current.Set(0)
}

func (m *Metrics) AnswererBalance(wei *big.Int) {
	if wei == nil {
		return
	}

	value, _ := new(big.Float).SetInt(wei).Float64()
	m.answerer.Set(value)
}

func (m *Metrics) Pass(served bool) {
	if served {
		m.passes.WithLabelValues(passServed).Inc()
		return
	}

	m.passes.WithLabelValues(passFailed).Inc()
}
