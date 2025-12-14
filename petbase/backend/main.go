package main

import (
	"database"
	"fmt"
	"log"
	"net/http"
	"petbase/handlers"

	"github.com/joho/godotenv"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:4100"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Create tables
	if err := createTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Routes
	http.HandleFunc("/", enableCORS(handleRoot))
	http.HandleFunc("/api/health", enableCORS(handleHealth))

	// Species routes
	http.HandleFunc("/api/species", enableCORS(handlers.SpeciesHandler))
	http.HandleFunc("/api/species/", enableCORS(handlers.SpeciesDetailHandler))

	// Breeds routes
	http.HandleFunc("/api/breeds", enableCORS(handlers.BreedsHandler))
	http.HandleFunc("/api/breeds/", enableCORS(handlers.BreedDetailHandler))
	http.HandleFunc("/api/breeds/species/", enableCORS(handlers.BreedsBySpeciesHandler))

	// Cards routes
	http.HandleFunc("/api/cards", enableCORS(handlers.CardsHandler))
	http.HandleFunc("/api/cards/breed/", enableCORS(handlers.CardsByBreedHandler))
	http.HandleFunc("/api/cards/", enableCORS(handlers.CardDetailHandler))

	port := ":8100"
	fmt.Printf("🐾 ЗооБаза API starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "ЗооБаза API", "version": "1.0.0"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "service": "petbase"}`)
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS species (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_en TEXT NOT NULL,
			description TEXT,
			icon TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS breeds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			species_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			name_en TEXT,
			description TEXT,
			origin TEXT,
			size TEXT,
			weight_min REAL,
			weight_max REAL,
			lifespan_min INTEGER,
			lifespan_max INTEGER,
			temperament TEXT,
			care_level TEXT,
			photo TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (species_id) REFERENCES species(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pet_cards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			breed_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			characteristics TEXT,
			care_tips TEXT,
			health_info TEXT,
			nutrition TEXT,
			photos TEXT,
			is_published BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (breed_id) REFERENCES breeds(id)
		)`,
	}

	for _, query := range queries {
		if _, err := database.DB.Exec(query); err != nil {
			return err
		}
	}

	log.Println("✅ ЗооБаза tables created successfully")
	return nil
}
