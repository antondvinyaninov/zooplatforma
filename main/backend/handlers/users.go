package handlers

import (
	"backend/models"
	"database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		handleGetUsers(w, r)
	case http.MethodPost:
		handleCreateUser(w, r)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 UserHandler: Received %s request to %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	id := extractID(r.URL.Path)
	log.Printf("🔍 UserHandler: Extracted ID: %d", id)

	if id == 0 {
		sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Printf("🔍 UserHandler: Calling handleGetUser for id=%d", id)
		// GET публичный - можно просматривать профили без авторизации
		handleGetUser(w, r, id)
	case http.MethodPut:
		// PUT требует авторизации - проверяем токен
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Можно редактировать только свой профиль
		if userID != id {
			sendError(w, "Forbidden", http.StatusForbidden)
			return
		}
		handleUpdateUser(w, r, id)
	case http.MethodDelete:
		// DELETE требует авторизации
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Можно удалить только свой профиль
		if userID != id {
			sendError(w, "Forbidden", http.StatusForbidden)
			return
		}
		handleDeleteUser(w, r, id)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetUsers(w http.ResponseWriter, _ *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, last_name, email, avatar, created_at, verified, verified_at FROM users")
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Avatar, &user.CreatedAt, &user.Verified, &user.VerifiedAt); err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	sendSuccess(w, users)
}

func handleGetUser(w http.ResponseWriter, _ *http.Request, id int) {
	// Получаем данные пользователя из локальной БД Main Backend
	var user models.User
	query := `SELECT u.id, u.name, u.last_name, u.email, u.bio, u.phone, u.location, u.avatar, u.cover_photo,
	          u.profile_visibility, u.show_phone, u.show_email, u.allow_messages, u.show_online, 
	          u.verified, u.verified_at, u.created_at,
	          ua.last_seen
	          FROM users u
	          LEFT JOIN user_activity ua ON u.id = ua.user_id
	          WHERE u.id = ?`

	var lastSeenTime *string
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.LastName, &user.Email, &user.Bio, &user.Phone,
		&user.Location, &user.Avatar, &user.CoverPhoto,
		&user.ProfileVisibility, &user.ShowPhone, &user.ShowEmail, &user.AllowMessages, &user.ShowOnline,
		&user.Verified, &user.VerifiedAt, &user.CreatedAt,
		&lastSeenTime,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			sendError(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("❌ Database error: %v", err)
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем LastSeen если есть
	if lastSeenTime != nil {
		user.LastSeen = parseTime(*lastSeenTime)
	}

	// Проверяем онлайн статус (активен в последние 5 минут)
	if user.LastSeen != nil {
		fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
		user.IsOnline = user.LastSeen.After(fiveMinutesAgo)
	}

	// Возвращаем данные пользователя
	sendSuccess(w, user)
	log.Printf("✅ User profile loaded from Main Backend: id=%d, name=%s, last_name=%s, is_online=%v", id, user.Name, user.LastName, user.IsOnline)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		sendError(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?)", req.Name, req.Email)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	sendSuccess(w, map[string]interface{}{"id": id, "name": req.Name, "email": req.Email})
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", req.Name, req.Email, id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, map[string]interface{}{"id": id, "name": req.Name, "email": req.Email})
}

func handleDeleteUser(w http.ResponseWriter, _ *http.Request, id int) {
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, map[string]string{"message": "User deleted"})
}

func extractID(path string) int {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return 0
	}
	id, _ := strconv.Atoi(parts[3])
	return id
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(models.Response{Success: true, Data: data})
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.Response{Success: false, Error: message})
}
