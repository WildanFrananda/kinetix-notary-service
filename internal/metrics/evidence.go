package metrics

import (
	"math"
	"math/big"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type evidence struct {
	feedUp      *prometheus.GaugeVec
	collected   *prometheus.CounterVec
	pending     *prometheus.GaugeVec
	lastBatch   *prometheus.GaugeVec
	lastClosed  *prometheus.GaugeVec
	abandoned   prometheus.Gauge
	balance     prometheus.Gauge
	cost        prometheus.Gauge
	runway      prometheus.Gauge
	ceilingHits prometheus.Counter
}

func newEvidence() evidence {
	streams := []string{"stream"}

	return evidence{
		feedUp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "kinetix_notary_evidence_feed_up",
				Help: "1 when the last read of a stream's evidence feed succeeded, 0 when it failed; NaN before the first read.",
			},
			streams,
		),
		collected: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kinetix_notary_evidence_collected_total",
				Help: "Evidence records collected from each stream's feed and kept with a salt.",
			},
			streams,
		),
		pending: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "kinetix_notary_evidence_pending_records",
				Help: "Evidence records collected and not yet in a batch.",
			},
			streams,
		),
		lastBatch: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "kinetix_notary_evidence_last_finalized_batch",
				Help: "Number of the newest batch of each stream the registry holds at a finalized block; NaN until one is.",
			},
			streams,
		),
		lastClosed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "kinetix_notary_evidence_last_finalized_closed_at_seconds",
				Help: "Close time of the newest final batch of each stream, in Unix seconds; NaN until one is final.",
			},
			streams,
		),
		abandoned: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_abandoned_evidence_batches",
			Help: "Evidence batches the chain refused for good; each stops its stream until a person acts.",
		}),
		balance: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_evidence_anchorer_balance_wei",
			Help: "Balance of the evidence anchorer key in wei; NaN until it is read.",
		}),
		cost: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_evidence_anchor_cost_wei",
			Help: "Mean cost of recent final evidence batches in wei, execution plus L1 data fee.",
		}),
		runway: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_evidence_runway_seconds",
			Help: "How long the evidence anchorer balance lasts at one batch per stream per interval; NaN until a batch is final.",
		}),
		ceilingHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kinetix_notary_evidence_fee_ceiling_hits_total",
			Help: "Times the fee ceiling held an evidence anchor back.",
		}),
	}
}

func (e *evidence) register(factory func(prometheus.Collector)) {
	for _, c := range []prometheus.Collector{
		e.feedUp,
		e.collected,
		e.pending,
		e.lastBatch,
		e.lastClosed,
		e.abandoned,
		e.balance,
		e.cost,
		e.runway,
		e.ceilingHits,
	} {
		factory(c)
	}

	for _, stream := range notary.EvidenceStreams() {
		name := stream.String()
		e.feedUp.WithLabelValues(name).Set(math.NaN())
		e.collected.WithLabelValues(name)
		e.pending.WithLabelValues(name)
		e.lastBatch.WithLabelValues(name).Set(math.NaN())
		e.lastClosed.WithLabelValues(name).Set(math.NaN())
	}

	e.balance.Set(math.NaN())
	e.runway.Set(math.NaN())
}

func (e *evidence) EvidenceFeed(stream notary.EvidenceStream, healthy bool) {
	if healthy {
		e.feedUp.WithLabelValues(stream.String()).Set(1)
		return
	}

	e.feedUp.WithLabelValues(stream.String()).Set(0)
}

func (e *evidence) EvidenceCollected(stream notary.EvidenceStream, records int) {
	e.collected.WithLabelValues(stream.String()).Add(float64(records))
}

func (e *evidence) EvidencePending(stream notary.EvidenceStream, records int64) {
	e.pending.WithLabelValues(stream.String()).Set(float64(records))
}

func (e *evidence) EvidenceFinalized(record notary.EvidenceBatchRecord) {
	name := record.Header.Stream.String()
	e.lastBatch.WithLabelValues(name).Set(float64(record.Header.Number))
	e.lastClosed.WithLabelValues(name).Set(float64(record.Header.ClosedAt.Seconds()))
}

func (e *evidence) EvidenceAbandoned(count int64) {
	e.abandoned.Set(float64(count))
}

func (e *evidence) EvidenceFunds(balance *big.Int, meanCost uint64, runway time.Duration, known bool) {
	if balance != nil {
		wei, _ := new(big.Float).SetInt(balance).Float64()
		e.balance.Set(wei)
	}

	e.cost.Set(float64(meanCost))

	if !known {
		e.runway.Set(math.NaN())
		return
	}

	e.runway.Set(runway.Seconds())
}

func (e *evidence) EvidenceFeeCeilingHit() {
	e.ceilingHits.Inc()
}
