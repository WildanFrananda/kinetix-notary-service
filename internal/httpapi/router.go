package httpapi

import (
	"encoding/json"
	"net/http"
)

type Instrument func(route string, next http.Handler) http.Handler

type Check func(r *http.Request) error

func NewRouter(instrument Instrument, metrics http.Handler, ready Check) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(
		"GET /health/live",
		instrument(
			"/health/live",
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				write(w, http.StatusOK, map[string]string{"status": "live"})
			}),
		),
	)
	mux.Handle(
		"GET /health/ready",
		instrument(
			"/health/ready",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := ready(r); err != nil {
					write(
						w,
						http.StatusServiceUnavailable,
						map[string]string{"status": "not_ready", "reason": err.Error()},
					)

					return
				}

				write(w, http.StatusOK, map[string]string{"status": "ready"})
			}),
		),
	)
	mux.Handle("GET /metrics", instrument("/metrics", metrics))
	mux.Handle("/", instrument("unmatched", http.NotFoundHandler()))

	return mux
}

func write(w http.ResponseWriter, status int, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
