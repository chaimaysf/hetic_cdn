package middleware

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// Map pour stocker les limiteurs par IP
var (
	ipLimiters = make(map[string]*rate.Limiter)
)

// getLimiter retourne un rate limiter spécifique à une IP
func getLimiter(ip string) *rate.Limiter {
	if IsIPBanned(ip) { // ✅ Vérification immédiate
		return nil
	}

	limiter, exists := ipLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(2, 5) // 2 requêtes/s, burst de 5
		ipLimiters[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware applique une limite par IP et banne les abus
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ips := strings.Split(forwarded, ",")
			ip = strings.TrimSpace(ips[0])
		}

		// ✅ Vérifie si l'IP est bannie avant de limiter
		if IsIPBanned(ip) {
			http.Error(w, "🚫 Accès interdit - IP bannie", http.StatusForbidden)
			return
		}

		limiter := getLimiter(ip)
		if limiter == nil || !limiter.Allow() {
			// ✅ Bannir l'IP si elle dépasse trop souvent la limite
			BanIP(ip, 10*time.Minute)
			http.Error(w, "🚫 Trop de requêtes. Vous êtes temporairement banni.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
