package utils

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter est une structure pour stocker les limiteurs par IP
type RateLimiter struct {
	limiters    map[string]*rate.Limiter
	mu          sync.Mutex
	rateLimit   rate.Limit
	burstLimit  int
	cleanupTime time.Duration
}

// NewRateLimiter crée un gestionnaire de rate limit
func NewRateLimiter(r rate.Limit, b int, cleanup time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limiters:    make(map[string]*rate.Limiter),
		rateLimit:   r,
		burstLimit:  b,
		cleanupTime: cleanup,
	}

	// Lancer un cleanup automatique des IPs inactives
	go rl.cleanupLoop()

	return rl
}

// GetLimiter retourne le rate limiter pour une IP donnée
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rateLimit, rl.burstLimit)
		rl.limiters[ip] = limiter
	}
	return limiter
}

// Supprime les IPs inactives après un certain temps
func (rl *RateLimiter) cleanupLoop() {
	for {
		time.Sleep(rl.cleanupTime)
		rl.mu.Lock()
		for ip, limiter := range rl.limiters {
			if limiter.AllowN(time.Now(), rl.burstLimit) {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}
	
