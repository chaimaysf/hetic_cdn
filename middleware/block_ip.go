package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"
)

// BlockIPMiddleware interdit l'accès aux IP bannies
func BlockIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getRealIP(r)

		// 🔍 Affiche l'IP détectée
		log.Println("🔍 Vérification de l'IP :", ip)

		// ✅ Vérifie si l'IP est bannie
		if IsIPBanned(ip) {
			log.Println("🚫 IP bannie détectée :", ip)
			http.Error(w, "🚫 Accès interdit - IP bannie", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getRealIP récupère la vraie IP de l'utilisateur (IPv4 ou IPv6 propre)
func getRealIP(r *http.Request) string {
	// 1️⃣ Vérifie si l'IP est dans X-Forwarded-For (proxy)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		ip := strings.TrimSpace(ips[0])
		return cleanIP(ip)
	}

	// 2️⃣ Sinon, récupère l'IP depuis RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Si erreur, retourne tel quel
	}

	return cleanIP(host)
}

// cleanIP formate correctement l'IP (évite [::1] au lieu de 127.0.0.1)
func cleanIP(ip string) string {
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}
