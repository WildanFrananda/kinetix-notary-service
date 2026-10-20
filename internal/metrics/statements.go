package metrics

import (
	"math"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/WildanFrananda/kinetix-notary-service/internal/delivering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type statements struct {
	delivered *prometheus.CounterVec
	pending   prometheus.Gauge
	overdue   prometheus.Gauge
}

func newStatements() statements {
	return statements{
		delivered: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kinetix_notary_statements_total",
				Help: "Statement deliveries handed to notification, by outcome: sent, refused (no address), or postponed.",
			},
			[]string{"outcome"},
		),
		pending: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_statements_pending",
			Help: "Statement deliveries planned and not yet handed over; NaN until it is read.",
		}),
		overdue: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "kinetix_notary_statements_overdue",
			Help: "Statement deliveries still pending six hours after they were planned; NaN until it is read.",
		}),
	}
}

func (s statements) register(factory func(prometheus.Collector)) {
	factory(s.delivered)
	factory(s.pending)
	factory(s.overdue)

	for _, outcome := range delivering.Outcomes {
		s.delivered.WithLabelValues(outcome.String())
	}

	s.pending.Set(math.NaN())
	s.overdue.Set(math.NaN())
}

func (s statements) Statement(outcome delivering.Outcome) {
	s.delivered.WithLabelValues(outcome.String()).Inc()
}

func (s statements) StatementBacklog(backlog notary.StatementBacklog) {
	s.pending.Set(float64(backlog.Pending))
	s.overdue.Set(float64(backlog.Overdue))
}
