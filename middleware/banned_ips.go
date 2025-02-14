package middleware

import (
	
	"sync"
	"time"
)

// 🔥 Liste des IP bannies (accessible globalement)
var (
	banMu     sync.Mutex
	bannedIPs = make(map[string]time.Time)
)

// BanIP ajoute une IP à la liste des bannies pour une durée donnée
func BanIP(ip string, duration time.Duration) {
	banMu.Lock()
	defer banMu.Unlock()
	bannedIPs[ip] = time.Now().Add(duration) // 🔥 Ban temporaire
}

// IsIPBanned vérifie si une IP est actuellement bannie
func IsIPBanned(ip string) bool {
	banMu.Lock()
	defer banMu.Unlock()

	banTime, exists := bannedIPs[ip]
	if !exists {
		return false
	}

	// 🔥 Si le ban est expiré, on le supprime
	if time.Now().After(banTime) {
		delete(bannedIPs, ip)
		return false
	}
	return true
}
func init() {
    bannedIPs["192.168.1.1"] = time.Now().Add(10 * time.Minute) // 🔥 Ban de test pour 10 min
}