package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/WildanFrananda/kinetix-notary-service/internal/approving"
)

const (
	approvalApproved = "approved"
	approvalRefused  = "refused"
	approvalFailed   = "failed"
)

type Watch struct {
	edge

	registry  *prometheus.Registry
	approvals *prometheus.CounterVec
	served    *prometheus.CounterVec
}

var _ approving.Observer = (*Watch)(nil)

func NewWatch(build Build) *Watch {
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

	w := &Watch{
		edge:     newEdge(),
		registry: registry,
		approvals: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kinetix_notary_watch_approvals_total",
				Help: "Epochs the watch was asked to approve, by outcome: approved, refused, or failed and left for a retry.",
			},
			[]string{"outcome"},
		),
		served: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kinetix_grpc_server_calls_total",
				Help: "gRPC calls this service answered, by method and the code it returned.",
			},
			[]string{"grpc_method", "grpc_code"},
		),
	}
	w.register(factory)
	factory(w.approvals)
	factory(w.served)

	for _, outcome := range []string{approvalApproved, approvalRefused, approvalFailed} {
		w.approvals.WithLabelValues(outcome)
	}

	return w
}

func (w *Watch) Handler() http.Handler {
	return promhttp.HandlerFor(w.registry, promhttp.HandlerOpts{Registry: w.registry})
}

func (w *Watch) ServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(func(
			ctx context.Context,
			req any,
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (any, error) {
			reply, err := handler(ctx, req)
			w.served.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()

			return reply, err
		}),
		grpc.ChainStreamInterceptor(func(
			srv any,
			stream grpc.ServerStream,
			info *grpc.StreamServerInfo,
			handler grpc.StreamHandler,
		) error {
			err := handler(srv, stream)
			w.served.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()

			return err
		}),
	}
}

func (w *Watch) SeedServed(method string) {
	w.served.WithLabelValues(method, "OK")
}

func (w *Watch) Approved() {
	w.approvals.WithLabelValues(approvalApproved).Inc()
}

func (w *Watch) Refused() {
	w.approvals.WithLabelValues(approvalRefused).Inc()
}

func (w *Watch) Failed() {
	w.approvals.WithLabelValues(approvalFailed).Inc()
}
