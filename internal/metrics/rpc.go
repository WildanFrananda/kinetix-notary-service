package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type rpc struct {
	requests  prometheus.Counter
	throttled prometheus.Counter
	paced     prometheus.Counter
}

func newRPC() rpc {
	return rpc{
		requests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kinetix_notary_rpc_requests_total",
			Help: "HTTP requests sent to the chain's RPC endpoint, each retry counted again.",
		}),
		throttled: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kinetix_notary_rpc_throttled_total",
			Help: "RPC responses refused with 429 Too Many Requests, retried or not.",
		}),
		paced: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kinetix_notary_rpc_paced_seconds_total",
			Help: "Seconds RPC requests waited to keep to KINETIX_NOTARY_RPC_RATE.",
		}),
	}
}

func (r rpc) register(factory func(prometheus.Collector)) {
	factory(r.requests)
	factory(r.throttled)
	factory(r.paced)
}

func (r rpc) RPCRequest() {
	r.requests.Inc()
}

func (r rpc) RPCThrottled() {
	r.throttled.Inc()
}

func (r rpc) RPCPaced(waited time.Duration) {
	r.paced.Add(waited.Seconds())
}
