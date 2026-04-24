package middleware

import "strings"

// isSkipPath reports whether the request path should be excluded from metrics
// and logging (e.g. high-frequency probes or telemetry endpoints).
func isSkipPath(path string) bool {
	return strings.HasPrefix(path, "/healthz/") || path == "/metrics"
}
