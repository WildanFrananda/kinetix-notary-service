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
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET "+ProofsRoute, instrument(ProofsRoute, proofs))
	mux.Handle("POST "+RegistrationsRoute, instrument(RegistrationsRoute, registrations))
	operational(mux, instrument, metrics, ready)

	return mux
}

func NewWatchRouter(instrument Instrument, metrics http.Handler, ready Check) http.Handler {
	mux := http.NewServeMux()
	operational(mux, instrument, metrics, ready)

	return mux
}

func operational(mux *http.ServeMux, instrument Instrument, metrics http.Handler, ready Check) {
	mux.Handle(
		"GET /health/live",
		instrument(
			"/health/live",
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, http.StatusOK, map[string]string{"status": "live"})
			}),
		),
	)

	mux.Handle(
		"GET /health/ready",
		instrument(
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
		),
	)

	mux.Handle("GET /metrics", instrument("/metrics", metrics))
	mux.Handle("GET /", instrument("unmatched", http.NotFoundHandler()))
}
