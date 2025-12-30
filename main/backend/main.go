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

	// Public API routes (register BEFORE root route)
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
	http.HandleFunc("/api/profile/avatar", enableCORS(middleware.AuthMiddleware(handlers.UploadAvatarHandler)))
	http.HandleFunc("/api/profile/avatar/delete", enableCORS(middleware.AuthMiddleware(handlers.DeleteAvatarHandler)))
	http.HandleFunc("/api/profile/cover", enableCORS(middleware.AuthMiddleware(handlers.UploadCoverPhotoHandler)))
	http.HandleFunc("/api/profile/cover/delete", enableCORS(middleware.AuthMiddleware(handlers.DeleteCoverPhotoHandler)))
	http.HandleFunc("/api/posts/drafts", enableCORS(middleware.AuthMiddleware(handlers.DraftsHandler)))
	http.HandleFunc("/api/posts", enableCORS(middleware.AuthMiddleware(handlers.PostsHandler)))
	http.HandleFunc("/api/posts/", enableCORS(middleware.AuthMiddleware(handlers.PostHandler)))
	http.HandleFunc("/api/posts/user/", enableCORS(middleware.AuthMiddleware(handlers.UserPostsHandler)))
	http.HandleFunc("/api/posts/pet/", enableCORS(middleware.AuthMiddleware(handlers.PetPostsHandler)))

	// Comments
	http.HandleFunc("/api/comments/post/", enableCORS(middleware.AuthMiddleware(handlers.CommentsHandler)))
	http.HandleFunc("/api/comments/", enableCORS(middleware.AuthMiddleware(handlers.DeleteCommentHandler)))

	// Polls
	http.HandleFunc("/api/polls/", enableCORS(middleware.AuthMiddleware(handlers.VoteHandler)))

	// Pets
	http.HandleFunc("/api/pets", enableCORS(middleware.AuthMiddleware(handlers.PetsHandler)))
	http.HandleFunc("/api/pets/", enableCORS(middleware.AuthMiddleware(handlers.PetHandler)))
	http.HandleFunc("/api/pets/user/", enableCORS(middleware.AuthMiddleware(handlers.UserPetsHandler)))

	// Pet Announcements
	http.HandleFunc("/api/announcements", enableCORS(middleware.AuthMiddleware(handlers.AnnouncementsHandler)))
	http.HandleFunc("/api/announcements/", enableCORS(middleware.AuthMiddleware(handlers.AnnouncementHandler)))
	http.HandleFunc("/api/announcements/posts/", enableCORS(middleware.AuthMiddleware(handlers.AnnouncementPostsHandler)))
	http.HandleFunc("/api/announcements/donations/", enableCORS(middleware.AuthMiddleware(handlers.AnnouncementDonationsHandler)))

	// Friends
	http.HandleFunc("/api/friends", enableCORS(middleware.AuthMiddleware(handlers.GetFriendsHandler)))
	http.HandleFunc("/api/friends/requests", enableCORS(middleware.AuthMiddleware(handlers.GetFriendRequestsHandler)))
	http.HandleFunc("/api/friends/send", enableCORS(middleware.AuthMiddleware(handlers.SendFriendRequestHandler)))
	http.HandleFunc("/api/friends/accept", enableCORS(middleware.AuthMiddleware(handlers.AcceptFriendRequestHandler)))
	http.HandleFunc("/api/friends/reject", enableCORS(middleware.AuthMiddleware(handlers.RejectFriendRequestHandler)))
	http.HandleFunc("/api/friends/remove", enableCORS(middleware.AuthMiddleware(handlers.RemoveFriendHandler)))
	http.HandleFunc("/api/friends/status", enableCORS(middleware.AuthMiddleware(handlers.GetFriendshipStatusHandler)))

	// Organizations
	http.HandleFunc("/api/organizations/all", enableCORS(handlers.GetAllOrganizationsHandler)) // Публичный endpoint
	http.HandleFunc("/api/organizations", enableCORS(middleware.AuthMiddleware(handlers.CreateOrganizationHandler)))
	http.HandleFunc("/api/organizations/", enableCORS(handlers.GetOrganizationHandler)) // Публичный для просмотра
	http.HandleFunc("/api/organizations/user/", enableCORS(middleware.AuthMiddleware(handlers.GetUserOrganizationsHandler)))
	http.HandleFunc("/api/organizations/members/", enableCORS(middleware.AuthMiddleware(handlers.GetOrganizationMembersHandler)))

	// Media - более специфичные роуты должны быть первыми
	mediaHandler := handlers.NewMediaHandler(database.DB)
	http.HandleFunc("/api/media/upload", enableCORS(middleware.AuthMiddleware(mediaHandler.UploadMedia)))
	http.HandleFunc("/api/media/stats", enableCORS(middleware.AuthMiddleware(mediaHandler.GetMediaStats)))
	http.HandleFunc("/api/media/user/", enableCORS(middleware.AuthMiddleware(mediaHandler.GetUserMedia)))
	http.HandleFunc("/api/media/file/", enableCORS(mediaHandler.GetMediaFile)) // Public для отображения
	http.HandleFunc("/api/media/delete/", enableCORS(middleware.AuthMiddleware(mediaHandler.DeleteMedia)))

	// Chunked Upload
	chunkedHandler := handlers.NewChunkedUploadHandler(database.DB)
	http.HandleFunc("/api/media/chunked/initiate", enableCORS(middleware.AuthMiddleware(chunkedHandler.InitiateUpload)))
	http.HandleFunc("/api/media/chunked/upload", enableCORS(middleware.AuthMiddleware(chunkedHandler.UploadChunk)))
	http.HandleFunc("/api/media/chunked/complete", enableCORS(middleware.AuthMiddleware(chunkedHandler.CompleteUpload)))

	// Static files - serve uploads directory from project root
	fs := http.FileServer(http.Dir("../.."))
	http.Handle("/uploads/", enableCORS(http.StripPrefix("/", fs).ServeHTTP))

	// Root route MUST be registered LAST
	http.HandleFunc("/", enableCORS(handleRoot))

	port := ":8000"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Только для точного пути "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Welcome to the API"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok"}`)
}
