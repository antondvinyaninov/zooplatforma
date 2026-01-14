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

		// Разрешённые origins
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true, // Main frontend
			"http://localhost:4000": true, // Admin frontend
			"http://localhost:4100": true, // PetBase frontend
			"http://localhost:6100": true, // Owner frontend
		}

		// Если origin не указан или не в списке разрешённых, используем дефолтный
		if origin == "" || !allowedOrigins[origin] {
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

	// Public user profile endpoint
	http.HandleFunc("/api/users/", enableCORS(handlers.UserHandler)) // Публичный просмотр профилей пользователей

	// Protected routes
	http.HandleFunc("/api/users", enableCORS(middleware.AuthMiddleware(handlers.UsersHandler)))
	http.HandleFunc("/api/profile", enableCORS(middleware.AuthMiddleware(handlers.UpdateProfileHandler)))
	http.HandleFunc("/api/profile/avatar", enableCORS(middleware.AuthMiddleware(handlers.UploadAvatarHandler)))
	http.HandleFunc("/api/profile/avatar/delete", enableCORS(middleware.AuthMiddleware(handlers.DeleteAvatarHandler)))
	http.HandleFunc("/api/profile/cover", enableCORS(middleware.AuthMiddleware(handlers.UploadCoverPhotoHandler)))
	http.HandleFunc("/api/profile/cover/delete", enableCORS(middleware.AuthMiddleware(handlers.DeleteCoverPhotoHandler)))
	http.HandleFunc("/api/posts/drafts", enableCORS(middleware.AuthMiddleware(handlers.DraftsHandler)))
	http.HandleFunc("/api/posts/user/", enableCORS(handlers.UserPostsHandler))                 // Публичный endpoint для просмотра постов пользователя
	http.HandleFunc("/api/posts/pet/", enableCORS(handlers.PetPostsHandler))                   // Публичный endpoint для просмотра постов питомца
	http.HandleFunc("/api/posts/organization/", enableCORS(handlers.OrganizationPostsHandler)) // Публичный endpoint для просмотра постов организации
	http.HandleFunc("/api/posts", enableCORS(handlers.PostsHandler))                           // GET публичный, POST требует авторизации (проверка внутри handler)
	http.HandleFunc("/api/posts/", enableCORS(middleware.AuthMiddleware(handlers.PostHandler)))

	// Comments
	http.HandleFunc("/api/comments/post/", enableCORS(middleware.AuthMiddleware(handlers.CommentsHandler)))
	http.HandleFunc("/api/comments/", enableCORS(middleware.AuthMiddleware(handlers.DeleteCommentHandler)))

	// Polls
	http.HandleFunc("/api/polls/", enableCORS(middleware.AuthMiddleware(handlers.VoteHandler)))

	// Pets
	http.HandleFunc("/api/pets", enableCORS(middleware.AuthMiddleware(handlers.PetsHandler)))
	http.HandleFunc("/api/pets/user/", enableCORS(handlers.UserPetsHandler)) // Публичный endpoint для просмотра питомцев
	http.HandleFunc("/api/pets/", enableCORS(middleware.AuthMiddleware(handlers.PetHandler)))

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

	// Notifications
	notificationsHandler := &handlers.NotificationsHandler{DB: database.DB}
	http.HandleFunc("/api/notifications", enableCORS(middleware.AuthMiddleware(notificationsHandler.GetNotifications)))
	http.HandleFunc("/api/notifications/unread", enableCORS(middleware.AuthMiddleware(notificationsHandler.GetUnreadCount)))
	http.HandleFunc("/api/notifications/read-all", enableCORS(middleware.AuthMiddleware(notificationsHandler.MarkAllAsRead)))
	http.Handle("/api/notifications/", enableCORS(middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			notificationsHandler.MarkAsRead(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Organizations
	http.HandleFunc("/api/organizations/all", enableCORS(handlers.GetAllOrganizationsHandler))                          // Публичный endpoint
	http.HandleFunc("/api/organizations/my", enableCORS(middleware.AuthMiddleware(handlers.GetMyOrganizationsHandler))) // Мои организации для публикации
	http.HandleFunc("/api/organizations", enableCORS(middleware.AuthMiddleware(handlers.CreateOrganizationHandler)))
	http.HandleFunc("/api/organizations/", enableCORS(middleware.AuthMiddleware(handlers.OrganizationHandler))) // GET и PUT для конкретной организации
	http.HandleFunc("/api/organizations/user/", enableCORS(middleware.AuthMiddleware(handlers.GetUserOrganizationsHandler)))
	http.HandleFunc("/api/organizations/members/", enableCORS(middleware.AuthMiddleware(handlers.GetOrganizationMembersHandler)))
	http.HandleFunc("/api/organizations/members/add", enableCORS(middleware.AuthMiddleware(handlers.AddMemberHandler)))
	http.HandleFunc("/api/organizations/members/update", enableCORS(middleware.AuthMiddleware(handlers.UpdateMemberHandler)))
	http.HandleFunc("/api/organizations/members/remove", enableCORS(middleware.AuthMiddleware(handlers.RemoveMemberHandler)))

	// Messenger (личные чаты 1-1)
	http.HandleFunc("/api/chats", enableCORS(middleware.AuthMiddleware(handlers.GetChatsHandler(database.DB))))
	http.HandleFunc("/api/chats/", enableCORS(middleware.AuthMiddleware(handlers.GetChatMessagesHandler(database.DB))))
	http.HandleFunc("/api/messages/send", enableCORS(middleware.AuthMiddleware(handlers.SendMessageHandler(database.DB))))
	http.HandleFunc("/api/messages/send-media", enableCORS(middleware.AuthMiddleware(handlers.SendMediaMessageHandler(database.DB))))
	http.HandleFunc("/api/messages/unread", enableCORS(middleware.AuthMiddleware(handlers.GetUnreadCountHandler(database.DB))))

	// Favorites (избранные питомцы)
	http.HandleFunc("/api/favorites", enableCORS(middleware.AuthMiddleware(handlers.FavoritesHandler)))
	http.HandleFunc("/api/favorites/", enableCORS(middleware.AuthMiddleware(handlers.FavoriteDetailHandler)))

	// Roles (система ролей)
	http.HandleFunc("/api/roles/available", enableCORS(middleware.AuthMiddleware(handlers.GetAllRolesHandler(database.DB))))
	http.HandleFunc("/api/roles/user/", enableCORS(middleware.AuthMiddleware(handlers.GetUserRolesHandler(database.DB))))
	http.HandleFunc("/api/roles/grant", enableCORS(middleware.AuthMiddleware(handlers.GrantRoleHandler(database.DB))))
	http.HandleFunc("/api/roles/revoke", enableCORS(middleware.AuthMiddleware(handlers.RevokeRoleHandler(database.DB))))

	// Verification (верификация пользователей)
	http.HandleFunc("/api/verification/verify", enableCORS(middleware.AuthMiddleware(handlers.VerifyUserHandler(database.DB))))
	http.HandleFunc("/api/verification/unverify", enableCORS(middleware.AuthMiddleware(handlers.UnverifyUserHandler(database.DB))))
	http.HandleFunc("/api/verification/status/", enableCORS(handlers.GetUserVerificationStatusHandler(database.DB)))
	http.HandleFunc("/api/users/verified", enableCORS(handlers.GetVerifiedUsersHandler(database.DB)))

	// Admin Logs (логи действий администраторов)
	http.HandleFunc("/api/admin/logs", enableCORS(middleware.AuthMiddleware(handlers.AdminLogsHandler)))
	http.HandleFunc("/api/admin/logs/stats", enableCORS(middleware.AuthMiddleware(handlers.GetAdminLogStats)))

	// User Activity (отслеживание активности пользователей)
	http.HandleFunc("/api/activity/update", enableCORS(middleware.AuthMiddleware(handlers.UpdateUserActivityHandler(database.DB))))
	http.HandleFunc("/api/activity/online", enableCORS(handlers.GetOnlineUsersCountHandler(database.DB)))
	http.HandleFunc("/api/activity/stats", enableCORS(handlers.GetUserActivityStatsHandler(database.DB)))

	// User Logs (логи действий пользователей)
	http.HandleFunc("/api/users/logs/", enableCORS(middleware.AuthMiddleware(handlers.GetUserLogsHandler(database.DB))))
	http.HandleFunc("/api/users/storage/", enableCORS(middleware.AuthMiddleware(handlers.GetUserStorageStatsHandler(database.DB))))

	// Reports (система жалоб)
	http.HandleFunc("/api/reports", enableCORS(middleware.AuthMiddleware(handlers.CreateReportHandler)))

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
