package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/chaim/hetic-cdn/handlers"
)

// SetupRoutes configure toutes les routes et retourne un http.Handler
func SetupRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.HandleFunc("/metrics", handlers.MetricsHandler).Methods("GET")
	r.HandleFunc("/files", handlers.GetFiles).Methods("GET")
	r.HandleFunc("/upload", handlers.UploadFile).Methods("POST")
	r.HandleFunc("/download/{id:[0-9]+}", handlers.DownloadFile).Methods("GET")
	r.HandleFunc("/delete/{id:[0-9]+}", handlers.DeleteFile).Methods("DELETE")

	// 🌍 Middleware CORS pour autoriser le frontend React
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"}, // Change si besoin
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "X-Secret-Key"},
		AllowCredentials: true,
	}).Handler(r)

	return corsHandler // ✅ Retourne un http.Handler
}
