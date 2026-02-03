package main

import (
	"backend/handlers"
	"database"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Разрешённые origins
		allowedOrigins := map[string]bool{
			"http://localhost:3000":                                  true, // Main frontend (dev)
			"http://localhost:4000":                                  true, // Admin frontend (dev)
			"http://localhost:4100":                                  true, // PetBase frontend (dev)
			"http://localhost:5100":                                  true, // Shelter frontend (dev)
			"http://localhost:6100":                                  true, // Owner frontend (dev)
			"http://localhost:6200":                                  true, // Volunteer frontend (dev)
			"http://localhost:6300":                                  true, // Clinic frontend (dev)
			"https://my-projects-zooplatforma.crv1ic.easypanel.host": true, // Main frontend (prod)
			"https://my-projects-admin.crv1ic.easypanel.host":        true, // Admin frontend (prod)
			"https://my-projects-petbase.crv1ic.easypanel.host":      true, // PetBase frontend (prod)
			"https://my-projects-shelter.crv1ic.easypanel.host":      true, // Shelter frontend (prod)
			"https://my-projects-owner.crv1ic.easypanel.host":        true, // Owner frontend (prod)
			"https://my-projects-volunteer.crv1ic.easypanel.host":    true, // Volunteer frontend (prod)
			"https://my-projects-clinic.crv1ic.easypanel.host":       true, // Clinic frontend (prod)
		}

		// Если origin в списке разрешённых, используем его
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			// Если origin не указан, используем дефолтный
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		} else {
			// Origin не разрешён - не устанавливаем заголовок
			log.Printf("⚠️ Blocked request from unauthorized origin: %s", origin)
		}

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

// enableCORSHandler - версия для http.Handler (используется с middleware)
func enableCORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:3000":                                  true,
			"http://localhost:4000":                                  true,
			"http://localhost:4100":                                  true,
			"http://localhost:5100":                                  true,
			"http://localhost:6100":                                  true,
			"http://localhost:6200":                                  true,
			"http://localhost:6300":                                  true,
			"https://my-projects-zooplatforma.crv1ic.easypanel.host": true,
			"https://my-projects-admin.crv1ic.easypanel.host":        true,
			"https://my-projects-petbase.crv1ic.easypanel.host":      true,
			"https://my-projects-shelter.crv1ic.easypanel.host":      true,
			"https://my-projects-owner.crv1ic.easypanel.host":        true,
			"https://my-projects-volunteer.crv1ic.easypanel.host":    true,
			"https://my-projects-clinic.crv1ic.easypanel.host":       true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		} else {
			log.Printf("⚠️ Blocked request from unauthorized origin: %s", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:7100"
		log.Printf("⚠️ AUTH_SERVICE_URL not set, using default: %s\n", authServiceURL)
	} else {
		log.Printf("🔐 Auth Service URL: %s\n", authServiceURL)
	}

	// ✅ Gateway теперь обрабатывает авторизацию
	log.Printf("🚀 Running behind API Gateway - auth handled by Gateway")

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

	// Protected routes (Gateway уже проверил авторизацию)
	http.HandleFunc("/api/users", enableCORS(handlers.UsersHandler))
	http.HandleFunc("/api/profile", enableCORS(handlers.UpdateProfileHandler))
	http.HandleFunc("/api/profile/avatar", enableCORS(handlers.UploadAvatarHandler))
	http.HandleFunc("/api/profile/avatar/delete", enableCORS(handlers.DeleteAvatarHandler))
	http.HandleFunc("/api/profile/cover", enableCORS(handlers.UploadCoverPhotoHandler))
	http.HandleFunc("/api/profile/cover/delete", enableCORS(handlers.DeleteCoverPhotoHandler))
	http.HandleFunc("/api/posts/drafts", enableCORS(handlers.DraftsHandler))

	// /api/posts - GET публичный, POST требует авторизации (Gateway проверяет)
	http.HandleFunc("/api/posts", enableCORS(handlers.PostsHandler))

	// /api/posts/ - универсальный обработчик для всех подпутей
	http.HandleFunc("/api/posts/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Специфичные роуты - проверяем первыми
		if strings.HasPrefix(path, "/api/posts/user/") {
			handlers.UserPostsHandler(w, r)
			return
		}
		if strings.HasPrefix(path, "/api/posts/pet/") {
			handlers.PetPostsHandler(w, r)
			return
		}
		if strings.HasPrefix(path, "/api/posts/organization/") {
			handlers.OrganizationPostsHandler(w, r)
			return
		}

		// /like endpoint (Gateway проверяет авторизацию для POST)
		if strings.HasSuffix(path, "/like") {
			handlers.LikesHandler(w, r)
			return
		}

		// Обычные посты /api/posts/{id}
		handlers.PostHandler(w, r)
	}))

	// Comments (Gateway проверяет авторизацию)
	http.HandleFunc("/api/comments/post/", enableCORS(handlers.CommentsHandler))
	http.HandleFunc("/api/comments/", enableCORS(handlers.DeleteCommentHandler))

	// Polls (Gateway проверяет авторизацию)
	http.HandleFunc("/api/polls/", enableCORS(handlers.VoteHandler))

	// Pets (Gateway проверяет авторизацию для защищенных endpoints)
	http.HandleFunc("/api/pets", enableCORS(handlers.PetsHandler))
	http.HandleFunc("/api/pets/user/", enableCORS(handlers.UserPetsHandler))       // Публичный endpoint
	http.HandleFunc("/api/pets/curated/", enableCORS(handlers.CuratedPetsHandler)) // Публичный endpoint
	http.HandleFunc("/api/pets/", enableCORS(handlers.PetHandler))                 // Gateway проверяет для DELETE

	// Pet Announcements (Gateway проверяет авторизацию)
	http.HandleFunc("/api/announcements", enableCORS(handlers.AnnouncementsHandler))
	http.HandleFunc("/api/announcements/", enableCORS(handlers.AnnouncementHandler))
	http.HandleFunc("/api/announcements/posts/", enableCORS(handlers.AnnouncementPostsHandler))
	http.HandleFunc("/api/announcements/donations/", enableCORS(handlers.AnnouncementDonationsHandler))

	// Friends (Gateway проверяет авторизацию)
	http.HandleFunc("/api/friends", enableCORS(handlers.GetFriendsHandler))
	http.HandleFunc("/api/friends/requests", enableCORS(handlers.GetFriendRequestsHandler))
	http.HandleFunc("/api/friends/send", enableCORS(handlers.SendFriendRequestHandler))
	http.HandleFunc("/api/friends/accept", enableCORS(handlers.AcceptFriendRequestHandler))
	http.HandleFunc("/api/friends/reject", enableCORS(handlers.RejectFriendRequestHandler))
	http.HandleFunc("/api/friends/remove", enableCORS(handlers.RemoveFriendHandler))
	http.HandleFunc("/api/friends/status", enableCORS(handlers.GetFriendshipStatusHandler))

	// Notifications (Gateway проверяет авторизацию)
	notificationsHandler := &handlers.NotificationsHandler{DB: database.DB}
	http.HandleFunc("/api/notifications", enableCORS(notificationsHandler.GetNotifications))
	http.HandleFunc("/api/notifications/unread", enableCORS(notificationsHandler.GetUnreadCount))
	http.HandleFunc("/api/notifications/read-all", enableCORS(notificationsHandler.MarkAllAsRead))
	http.HandleFunc("/api/notifications/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			notificationsHandler.MarkAsRead(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Organizations (Gateway проверяет авторизацию для защищенных endpoints)
	http.HandleFunc("/api/organizations/all", enableCORS(handlers.GetAllOrganizationsHandler))         // Публичный endpoint
	http.HandleFunc("/api/organizations/my", enableCORS(handlers.GetMyOrganizationsHandler))           // Gateway проверяет
	http.HandleFunc("/api/organizations", enableCORS(handlers.CreateOrganizationHandler))              // Gateway проверяет
	http.HandleFunc("/api/organizations/", enableCORS(handlers.OrganizationHandler))                   // Gateway проверяет
	http.HandleFunc("/api/organizations/user/", enableCORS(handlers.GetUserOrganizationsHandler))      // Gateway проверяет
	http.HandleFunc("/api/organizations/members/", enableCORS(handlers.GetOrganizationMembersHandler)) // Gateway проверяет
	http.HandleFunc("/api/organizations/members/add", enableCORS(handlers.AddMemberHandler))           // Gateway проверяет
	http.HandleFunc("/api/organizations/members/update", enableCORS(handlers.UpdateMemberHandler))     // Gateway проверяет
	http.HandleFunc("/api/organizations/members/remove", enableCORS(handlers.RemoveMemberHandler))     // Gateway проверяет

	// Messenger (личные чаты 1-1) (Gateway проверяет авторизацию)
	http.HandleFunc("/api/chats", enableCORS(handlers.GetChatsHandler(database.DB)))
	http.HandleFunc("/api/chats/", enableCORS(handlers.GetChatMessagesHandler(database.DB)))
	http.HandleFunc("/api/messages/send", enableCORS(handlers.SendMessageHandler(database.DB)))
	http.HandleFunc("/api/messages/send-media", enableCORS(handlers.SendMediaMessageHandler(database.DB)))
	http.HandleFunc("/api/messages/unread", enableCORS(handlers.GetUnreadCountHandler(database.DB)))

	// Favorites (избранные питомцы) (Gateway проверяет авторизацию)
	http.HandleFunc("/api/favorites", enableCORS(handlers.FavoritesHandler))
	http.HandleFunc("/api/favorites/", enableCORS(handlers.FavoriteDetailHandler))

	// Roles (система ролей) (Gateway проверяет авторизацию)
	http.HandleFunc("/api/roles/available", enableCORS(handlers.GetAllRolesHandler(database.DB)))
	http.HandleFunc("/api/roles/user/", enableCORS(handlers.GetUserRolesHandler(database.DB)))
	http.HandleFunc("/api/roles/grant", enableCORS(handlers.GrantRoleHandler(database.DB)))
	http.HandleFunc("/api/roles/revoke", enableCORS(handlers.RevokeRoleHandler(database.DB)))

	// Verification (верификация пользователей) (Gateway проверяет авторизацию для защищенных endpoints)
	http.HandleFunc("/api/verification/verify", enableCORS(handlers.VerifyUserHandler(database.DB)))
	http.HandleFunc("/api/verification/unverify", enableCORS(handlers.UnverifyUserHandler(database.DB)))
	http.HandleFunc("/api/verification/status/", enableCORS(handlers.GetUserVerificationStatusHandler(database.DB)))
	http.HandleFunc("/api/users/verified", enableCORS(handlers.GetVerifiedUsersHandler(database.DB)))

	// Admin Logs (логи действий администраторов) (Gateway проверяет авторизацию)
	http.HandleFunc("/api/admin/logs", enableCORS(handlers.AdminLogsHandler))
	http.HandleFunc("/api/admin/logs/stats", enableCORS(handlers.GetAdminLogStats))

	// User Activity (отслеживание активности пользователей) (Gateway проверяет авторизацию для защищенных endpoints)
	http.HandleFunc("/api/activity/update", enableCORS(handlers.UpdateUserActivityHandler(database.DB)))
	http.HandleFunc("/api/activity/online", enableCORS(handlers.GetOnlineUsersCountHandler(database.DB)))
	http.HandleFunc("/api/activity/stats", enableCORS(handlers.GetUserActivityStatsHandler(database.DB)))

	// User Logs (логи действий пользователей) (Gateway проверяет авторизацию)
	http.HandleFunc("/api/users/logs/", enableCORS(handlers.GetUserLogsHandler(database.DB)))
	http.HandleFunc("/api/users/storage/", enableCORS(handlers.GetUserStorageStatsHandler(database.DB)))

	// Reports (система жалоб) (Gateway проверяет авторизацию)
	http.HandleFunc("/api/reports", enableCORS(handlers.CreateReportHandler))

	// Media - более специфичные роуты должны быть первыми (Gateway проверяет авторизацию)
	mediaHandler := handlers.NewMediaHandler(database.DB)
	http.HandleFunc("/api/media/upload", enableCORS(mediaHandler.UploadMedia))
	http.HandleFunc("/api/media/stats", enableCORS(mediaHandler.GetMediaStats))
	http.HandleFunc("/api/media/user/", enableCORS(mediaHandler.GetUserMedia))
	http.HandleFunc("/api/media/file/", enableCORS(mediaHandler.GetMediaFile)) // Public для отображения
	http.HandleFunc("/api/media/delete/", enableCORS(mediaHandler.DeleteMedia))

	// Chunked Upload (Gateway проверяет авторизацию)
	chunkedHandler := handlers.NewChunkedUploadHandler(database.DB)
	http.HandleFunc("/api/media/chunked/initiate", enableCORS(chunkedHandler.InitiateUpload))
	http.HandleFunc("/api/media/chunked/upload", enableCORS(chunkedHandler.UploadChunk))
	http.HandleFunc("/api/media/chunked/complete", enableCORS(chunkedHandler.CompleteUpload))

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
