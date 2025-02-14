package handlers

import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler handles the /metrics endpoint
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
    promhttp.Handler().ServeHTTP(w, r)
}