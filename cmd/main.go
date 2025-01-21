package main

import (
	"log"
	"net/http"
	"os"

	"canciones/internal/handlers"
	"canciones/internal/middleware"
	"canciones/internal/storage"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	if err := storage.ConnectMongoDB(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/login", handlers.Login).Methods("POST")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		w.Write([]byte("Welcome, " + user))
	}).Methods("GET")
	protected.HandleFunc("/search", handlers.SearchHandler).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	log.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
