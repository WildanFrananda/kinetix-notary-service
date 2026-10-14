package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (e *edge) Instrument(route string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		e.httpRequests.WithLabelValues(r.Method, route, strconv.Itoa(recorder.status)).Inc()
		e.httpDuration.WithLabelValues(r.Method, route).Observe(time.Since(started).Seconds())
	})
}
