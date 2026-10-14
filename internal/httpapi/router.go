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

	return mux
}
