package httpapi

import (
	"net/http"
)

type Instrument func(route string, next http.Handler) http.Handler

type Check func(r *http.Request) error

func NewRouter(
	instrument Instrument,
	metrics http.Handler,
	ready Check,
	proofs *ProofHandler,
	registrations *RegistrationHandler,
	orderEvidence *EvidenceHandler,
	reviewEvidence *EvidenceHandler,
) http.Handler {
	routes := newRouteTable(instrument)

	routes.handle(http.MethodGet, ProofsRoute, proofs)
	routes.handle(http.MethodGet, RegistrationsRoute, http.HandlerFunc(registrations.ServeTemplate))
	routes.handle(http.MethodPost, RegistrationsRoute, registrations)
	routes.handle(http.MethodGet, OrderEvidenceRoute, orderEvidence)
	routes.handle(http.MethodGet, ReviewEvidenceRoute, reviewEvidence)
	operational(routes, metrics, ready)

	return routes.handler()
}

func NewWatchRouter(instrument Instrument, metrics http.Handler, ready Check) http.Handler {
	routes := newRouteTable(instrument)
	operational(routes, metrics, ready)

	return routes.handler()
}

func operational(routes *routeTable, metrics http.Handler, ready Check) {
	routes.handle(
		http.MethodGet,
		"/health/live",
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"status": "live"})
		}),
	)

	routes.handle(
		http.MethodGet,
		"/health/ready",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := ready(r); err != nil {
				writeJSON(
					w,
					http.StatusServiceUnavailable,
					map[string]string{"status": "not_ready", "reason": err.Error()},
				)

				return
			}

			writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
		}),
	)

	routes.handle(http.MethodGet, "/metrics", metrics)
}
