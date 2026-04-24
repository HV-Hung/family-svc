package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LogRequest is a structured access-log middleware. Each completed request is
// logged as a single slog record with method, path, status, latency, and the
// remote address. Health-check probes (/healthz/*) are silently passed
// through to avoid flooding logs with high-frequency Kubernetes probe noise.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for telemetry and health-check probes.
		if isSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		start := time.Now()
		next.ServeHTTP(rw, r)
		elapsed := time.Since(start)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"latency_ms", elapsed.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}
