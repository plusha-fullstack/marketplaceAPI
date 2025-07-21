package main

import (
	"log"
	"net/http"

	"github.com/plusha-fullstack/marketplaceAPI/internal/database"
	"github.com/plusha-fullstack/marketplaceAPI/internal/handlers"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	authHandler := handlers.NewAuthHandler(db)
	adHandler := handlers.NewAdHandler(db)

	http.HandleFunc("/api/register", corsMiddleware(authHandler.Register))
	http.HandleFunc("/api/login", corsMiddleware(authHandler.Login))
	http.HandleFunc("/api/ads", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AuthMiddleware(adHandler.CreateAd)(w, r)
		case http.MethodGet:
			handlers.OptionalAuthMiddleware(adHandler.GetAds)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
