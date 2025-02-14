package main

import (
	"log"
	"net/http"

	"github.com/chaim/hetic-cdn/db"
	"github.com/chaim/hetic-cdn/metrics"
	"github.com/chaim/hetic-cdn/middleware"
	"github.com/chaim/hetic-cdn/routes"
)

func main() {
	// ✅ Initialisation
	db.InitDB()
	db.CreateTables()
	db.InitRedis() // 🔥 Ajout de Redis
	metrics.RegisterMetrics()
	// ✅ Récupération du routeur correctement configuré
	router := routes.SetupRoutes()

	// ✅ Application des middlewares
	router = middleware.BlockIPMiddleware(router)  // 🔥 Bloque les IP bannies
	router = middleware.RateLimitMiddleware(router) // ✅ Applique le Rate Limiting

	log.Println("🚀 Serveur CDN démarré sur le port 8080")
	// log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "cert.key", router)) // ✅ Serveur HTTPS

	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "cert.key", router)) // ✅ Serveur HTTPS
}
