package main

import (
	"backend/handlers"
	"backend/middleware"
	"database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:3000"
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
		log.Println("Warning: .env file not found, using default values")
	}

	// Debug: проверяем JWT_SECRET
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("⚠️  JWT_SECRET not set, using default")
	} else {
		log.Printf("✅ JWT_SECRET loaded: %s...\n", secret[:10])
	}

	// Initialize JWT Secret AFTER loading .env
	middleware.InitJWTSecret()

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Public routes
	http.HandleFunc("/", enableCORS(handleRoot))
	http.HandleFunc("/api/health", enableCORS(handleHealth))
	http.HandleFunc("/api/auth/register", enableCORS(handlers.RegisterHandler))
	http.HandleFunc("/api/auth/login", enableCORS(handlers.LoginHandler))
	http.HandleFunc("/api/auth/logout", enableCORS(handlers.LogoutHandler))
	http.HandleFunc("/api/auth/me", enableCORS(handlers.MeHandler))
	http.HandleFunc("/api/auth/verify", enableCORS(handlers.VerifyTokenHandler))

	// Protected routes
	http.HandleFunc("/api/users", enableCORS(middleware.AuthMiddleware(handlers.UsersHandler)))
	http.HandleFunc("/api/users/", enableCORS(middleware.AuthMiddleware(handlers.UserHandler)))
	http.HandleFunc("/api/profile", enableCORS(middleware.AuthMiddleware(handlers.UpdateProfileHandler)))
	http.HandleFunc("/api/posts/drafts", enableCORS(middleware.AuthMiddleware(handlers.DraftsHandler)))
	http.HandleFunc("/api/posts", enableCORS(middleware.AuthMiddleware(handlers.PostsHandler)))
	http.HandleFunc("/api/posts/", enableCORS(middleware.AuthMiddleware(handlers.PostHandler)))
	http.HandleFunc("/api/posts/user/", enableCORS(middleware.AuthMiddleware(handlers.UserPostsHandler)))

	// Comments
	http.HandleFunc("/api/comments/post/", enableCORS(middleware.AuthMiddleware(handlers.CommentsHandler)))
	http.HandleFunc("/api/comments/", enableCORS(middleware.AuthMiddleware(handlers.DeleteCommentHandler)))

	// Polls
	http.HandleFunc("/api/polls/", enableCORS(middleware.AuthMiddleware(handlers.VoteHandler)))

	// Pets
	http.HandleFunc("/api/pets", enableCORS(middleware.AuthMiddleware(handlers.PetsHandler)))
	http.HandleFunc("/api/pets/", enableCORS(middleware.AuthMiddleware(handlers.PetHandler)))
	http.HandleFunc("/api/pets/user/", enableCORS(middleware.AuthMiddleware(handlers.UserPetsHandler)))

	// Media
	mediaHandler := handlers.NewMediaHandler(database.DB)
	http.HandleFunc("/api/media/upload", enableCORS(middleware.AuthMiddleware(mediaHandler.UploadMedia)))
	http.HandleFunc("/api/media/user/", enableCORS(middleware.AuthMiddleware(mediaHandler.GetUserMedia)))
	http.HandleFunc("/api/media/file/", enableCORS(mediaHandler.GetMediaFile)) // Public для отображения
	http.HandleFunc("/api/media/stats", enableCORS(middleware.AuthMiddleware(mediaHandler.GetMediaStats)))
	http.HandleFunc("/api/media/", enableCORS(middleware.AuthMiddleware(mediaHandler.DeleteMedia)))

	port := ":8000"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Welcome to the API"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok"}`)
}
