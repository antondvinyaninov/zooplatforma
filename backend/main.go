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
	"github.com/zooplatforma/pkg/clients"
	"github.com/zooplatforma/pkg/middleware"
)

// Global AuthClient
var authClient *clients.AuthClient

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

	// ✅ Auth Service URL будет автоматически прочитан из AUTH_SERVICE_URL в .env
	// pkg/middleware использует os.Getenv("AUTH_SERVICE_URL") внутри
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:7100"
		log.Printf("⚠️ AUTH_SERVICE_URL not set, using default: %s\n", authServiceURL)
	} else {
		log.Printf("🔐 Auth Service URL: %s\n", authServiceURL)
	}

	// Initialize AuthClient
	authClient = clients.NewAuthClient(authServiceURL)
	log.Printf("✅ AuthClient initialized: %s\n", authServiceURL)

	// ✅ КРИТИЧНО: Инициализировать AuthMiddleware с URL Auth Service
	middleware.InitAuthMiddleware(authServiceURL)
	log.Printf("✅ AuthMiddleware initialized with Auth Service: %s\n", authServiceURL)

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
	http.Handle("/api/users", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.UsersHandler))))
	http.Handle("/api/profile", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateProfileHandler))))
	http.Handle("/api/profile/avatar", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadAvatarHandler))))
	http.Handle("/api/profile/avatar/delete", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteAvatarHandler))))
	http.Handle("/api/profile/cover", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadCoverPhotoHandler))))
	http.Handle("/api/profile/cover/delete", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteCoverPhotoHandler))))
	http.Handle("/api/posts/drafts", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.DraftsHandler))))

	// /api/posts - GET с опциональной авторизацией, POST требует авторизации
	http.Handle("/api/posts", enableCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.AuthMiddleware(http.HandlerFunc(handlers.PostsHandler)).ServeHTTP(w, r)
		} else {
			middleware.OptionalAuthMiddleware(http.HandlerFunc(handlers.PostsHandler)).ServeHTTP(w, r)
		}
	})))

	// /api/posts/ - универсальный обработчик для всех подпутей
	http.Handle("/api/posts/", enableCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// /like endpoint
		if strings.HasSuffix(path, "/like") {
			if r.Method == http.MethodGet {
				middleware.OptionalAuthMiddleware(http.HandlerFunc(handlers.LikesHandler)).ServeHTTP(w, r)
			} else {
				middleware.AuthMiddleware(http.HandlerFunc(handlers.LikesHandler)).ServeHTTP(w, r)
			}
			return
		}

		// Обычные посты /api/posts/{id}
		if r.Method == http.MethodGet {
			middleware.OptionalAuthMiddleware(http.HandlerFunc(handlers.PostHandler)).ServeHTTP(w, r)
		} else {
			middleware.AuthMiddleware(http.HandlerFunc(handlers.PostHandler)).ServeHTTP(w, r)
		}
	})))

	// Comments
	http.Handle("/api/comments/post/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.CommentsHandler))))
	http.Handle("/api/comments/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteCommentHandler))))

	// Polls
	http.Handle("/api/polls/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.VoteHandler))))

	// Pets
	http.Handle("/api/pets", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.PetsHandler))))
	http.HandleFunc("/api/pets/user/", enableCORS(handlers.UserPetsHandler))       // Публичный endpoint для просмотра питомцев
	http.HandleFunc("/api/pets/curated/", enableCORS(handlers.CuratedPetsHandler)) // Публичный endpoint для просмотра курируемых питомцев
	// /api/pets/:id - GET публичный, DELETE требует авторизации
	http.Handle("/api/pets/", enableCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			middleware.AuthMiddleware(http.HandlerFunc(handlers.PetHandler)).ServeHTTP(w, r)
		} else {
			handlers.PetHandler(w, r)
		}
	})))

	// Pet Announcements
	http.Handle("/api/announcements", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.AnnouncementsHandler))))
	http.Handle("/api/announcements/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.AnnouncementHandler))))
	http.Handle("/api/announcements/posts/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.AnnouncementPostsHandler))))
	http.Handle("/api/announcements/donations/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.AnnouncementDonationsHandler))))

	// Friends
	http.Handle("/api/friends", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFriendsHandler))))
	http.Handle("/api/friends/requests", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFriendRequestsHandler))))
	http.Handle("/api/friends/send", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.SendFriendRequestHandler))))
	http.Handle("/api/friends/accept", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.AcceptFriendRequestHandler))))
	http.Handle("/api/friends/reject", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.RejectFriendRequestHandler))))
	http.Handle("/api/friends/remove", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.RemoveFriendHandler))))
	http.Handle("/api/friends/status", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFriendshipStatusHandler))))

	// Notifications
	notificationsHandler := &handlers.NotificationsHandler{DB: database.DB}
	http.Handle("/api/notifications", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(notificationsHandler.GetNotifications))))
	http.Handle("/api/notifications/unread", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(notificationsHandler.GetUnreadCount))))
	http.Handle("/api/notifications/read-all", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(notificationsHandler.MarkAllAsRead))))
	http.Handle("/api/notifications/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			notificationsHandler.MarkAsRead(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Organizations
	http.HandleFunc("/api/organizations/all", enableCORS(handlers.GetAllOrganizationsHandler))                                               // Публичный endpoint
	http.Handle("/api/organizations/my", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetMyOrganizationsHandler)))) // Мои организации для публикации
	http.Handle("/api/organizations", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateOrganizationHandler))))
	http.Handle("/api/organizations/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.OrganizationHandler)))) // GET и PUT для конкретной организации
	http.Handle("/api/organizations/user/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserOrganizationsHandler))))
	http.Handle("/api/organizations/members/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetOrganizationMembersHandler))))
	http.Handle("/api/organizations/members/add", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.AddMemberHandler))))
	http.Handle("/api/organizations/members/update", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateMemberHandler))))
	http.Handle("/api/organizations/members/remove", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.RemoveMemberHandler))))

	// Messenger (личные чаты 1-1)
	http.Handle("/api/chats", enableCORSHandler(middleware.AuthMiddleware(handlers.GetChatsHandler(database.DB))))
	http.Handle("/api/chats/", enableCORSHandler(middleware.AuthMiddleware(handlers.GetChatMessagesHandler(database.DB))))
	http.Handle("/api/messages/send", enableCORSHandler(middleware.AuthMiddleware(handlers.SendMessageHandler(database.DB))))
	http.Handle("/api/messages/send-media", enableCORSHandler(middleware.AuthMiddleware(handlers.SendMediaMessageHandler(database.DB))))
	http.Handle("/api/messages/unread", enableCORSHandler(middleware.AuthMiddleware(handlers.GetUnreadCountHandler(database.DB))))

	// Favorites (избранные питомцы)
	http.Handle("/api/favorites", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.FavoritesHandler))))
	http.Handle("/api/favorites/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.FavoriteDetailHandler))))

	// Roles (система ролей)
	http.Handle("/api/roles/available", enableCORSHandler(middleware.AuthMiddleware(handlers.GetAllRolesHandler(database.DB))))
	http.Handle("/api/roles/user/", enableCORSHandler(middleware.AuthMiddleware(handlers.GetUserRolesHandler(database.DB))))
	http.Handle("/api/roles/grant", enableCORSHandler(middleware.AuthMiddleware(handlers.GrantRoleHandler(database.DB))))
	http.Handle("/api/roles/revoke", enableCORSHandler(middleware.AuthMiddleware(handlers.RevokeRoleHandler(database.DB))))

	// Verification (верификация пользователей)
	http.Handle("/api/verification/verify", enableCORSHandler(middleware.AuthMiddleware(handlers.VerifyUserHandler(database.DB))))
	http.Handle("/api/verification/unverify", enableCORSHandler(middleware.AuthMiddleware(handlers.UnverifyUserHandler(database.DB))))
	http.HandleFunc("/api/verification/status/", enableCORS(handlers.GetUserVerificationStatusHandler(database.DB)))
	http.HandleFunc("/api/users/verified", enableCORS(handlers.GetVerifiedUsersHandler(database.DB)))

	// Admin Logs (логи действий администраторов)
	http.Handle("/api/admin/logs", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.AdminLogsHandler))))
	http.Handle("/api/admin/logs/stats", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetAdminLogStats))))

	// User Activity (отслеживание активности пользователей)
	http.Handle("/api/activity/update", enableCORSHandler(middleware.AuthMiddleware(handlers.UpdateUserActivityHandler(database.DB))))
	http.HandleFunc("/api/activity/online", enableCORS(handlers.GetOnlineUsersCountHandler(database.DB)))
	http.HandleFunc("/api/activity/stats", enableCORS(handlers.GetUserActivityStatsHandler(database.DB)))

	// User Logs (логи действий пользователей)
	http.Handle("/api/users/logs/", enableCORSHandler(middleware.AuthMiddleware(handlers.GetUserLogsHandler(database.DB))))
	http.Handle("/api/users/storage/", enableCORSHandler(middleware.AuthMiddleware(handlers.GetUserStorageStatsHandler(database.DB))))

	// Reports (система жалоб)
	http.Handle("/api/reports", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateReportHandler))))

	// Media - более специфичные роуты должны быть первыми
	mediaHandler := handlers.NewMediaHandler(database.DB)
	http.Handle("/api/media/upload", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(mediaHandler.UploadMedia))))
	http.Handle("/api/media/stats", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(mediaHandler.GetMediaStats))))
	http.Handle("/api/media/user/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(mediaHandler.GetUserMedia))))
	http.HandleFunc("/api/media/file/", enableCORS(mediaHandler.GetMediaFile)) // Public для отображения
	http.Handle("/api/media/delete/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(mediaHandler.DeleteMedia))))

	// Chunked Upload
	chunkedHandler := handlers.NewChunkedUploadHandler(database.DB)
	http.Handle("/api/media/chunked/initiate", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(chunkedHandler.InitiateUpload))))
	http.Handle("/api/media/chunked/upload", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(chunkedHandler.UploadChunk))))
	http.Handle("/api/media/chunked/complete", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(chunkedHandler.CompleteUpload))))

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
