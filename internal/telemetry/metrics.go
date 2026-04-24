package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Registry holds the custom Prometheus registry and all application metrics.
type Registry struct {
	Prometheus *prometheus.Registry

	// HTTPRequestsTotal counts every completed HTTP request.
	HTTPRequestsTotal *prometheus.CounterVec

	// HTTPRequestDuration tracks latency per route.
	HTTPRequestDuration *prometheus.HistogramVec

	// HTTPRequestsInFlight tracks concurrent requests being handled.
	HTTPRequestsInFlight prometheus.Gauge
}

// NewRegistry creates an isolated Prometheus registry pre-populated with
// the standard Go runtime / process collectors and the application HTTP metrics.
func NewRegistry() *Registry {
	reg := prometheus.NewRegistry()

	// Standard runtime collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// HTTP request counter.
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by method, path, and status code.",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP request duration histogram.
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds, partitioned by method and path.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)

	// In-flight gauge.
	httpRequestsInFlight := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed.",
		},
	)

	reg.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInFlight,
	)

	return &Registry{
		Prometheus:           reg,
		HTTPRequestsTotal:    httpRequestsTotal,
		HTTPRequestDuration:  httpRequestDuration,
		HTTPRequestsInFlight: httpRequestsInFlight,
	}
}
