package middleware

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// 🔥 Déclaration d'un compteur Prometheus pour suivre les attaques bloquées
var blockedRequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "blocked_requests_total",
		Help: "Nombre total de requêtes bloquées par le WAF",
	},
	[]string{"type"},
)

func init() {
	prometheus.MustRegister(blockedRequestsCounter)
}

// WebApplicationFirewall bloque les requêtes malveillantes
func WebApplicationFirewall(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		badPatterns := []string{"<script>", "DROP TABLE", "SELECT * FROM", "1=1", "--", "'"}

		for _, pattern := range badPatterns {
			if strings.Contains(r.URL.RawQuery, pattern) || strings.Contains(r.PostForm.Encode(), pattern) {
				blockedRequestsCounter.WithLabelValues("SQL_XSS").Inc()
				http.Error(w, "🚨 Requête malveillante détectée !", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
