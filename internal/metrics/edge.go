package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type edge struct {
	httpRequests *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
	grpcCalls    *prometheus.CounterVec
}

func newEdge() edge {
	return edge{
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
	}
}

func (e *edge) register(factory func(prometheus.Collector)) {
	factory(e.httpRequests)
	factory(e.httpDuration)
	factory(e.grpcCalls)
}

func (e *edge) SeedRoute(method, route string) {
	e.httpRequests.WithLabelValues(method, route, strconv.Itoa(http.StatusOK))
	e.httpDuration.WithLabelValues(method, route)
}

func (e *edge) SeedCall(peer, method string) {
	e.grpcCalls.WithLabelValues(peer, method, "OK")
}
