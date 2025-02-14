package metrics

import "github.com/prometheus/client_golang/prometheus"

// Déclaration des métriques
var (
	TotalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cdn_total_requests",
			Help: "Total des requêtes reçues par le CDN",
		})

	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cdn_request_duration_seconds",
			Help:    "Durée des requêtes au CDN",
			Buckets: prometheus.DefBuckets,
		})
)

// Enregistrement des métriques
func RegisterMetrics() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)
}