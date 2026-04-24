package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HV-Hung/family-svc/internal/telemetry"
)

// responseWriter is a thin wrapper around http.ResponseWriter that captures
// the status code written by the downstream handler.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// InstrumentHandler wraps next with Prometheus instrumentation using the
// metrics defined in the provided Registry. Health-check probes (/healthz/*)
// are passed through without recording any metrics.
func InstrumentHandler(reg *telemetry.Registry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip metric recording for telemetry and health-check probes.
		if isSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		reg.HTTPRequestsInFlight.Inc()
		defer reg.HTTPRequestsInFlight.Dec()

		start := time.Now()
		next.ServeHTTP(rw, r)
		elapsed := time.Since(start).Seconds()

		method := r.Method
		path := r.URL.Path
		status := fmt.Sprintf("%d", rw.status)

		reg.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		reg.HTTPRequestDuration.WithLabelValues(method, path).Observe(elapsed)
	})
}
