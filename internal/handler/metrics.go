package handler

import (
	"net/http"

	"github.com/HV-Hung/family-svc/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler serves Prometheus metrics from the application registry.
// Intended to be scraped by a Prometheus ServiceMonitor on GET /metrics.
func MetricsHandler(reg *telemetry.Registry) http.HandlerFunc {
	h := promhttp.HandlerFor(
		reg.Prometheus,
		promhttp.HandlerOpts{EnableOpenMetrics: true},
	)
	return h.ServeHTTP
}
