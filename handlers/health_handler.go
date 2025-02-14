package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler gère la route /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "CDN Go is running")
}
