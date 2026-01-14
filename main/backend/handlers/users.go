package handlers

import (
	"backend/models"
	"database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	w.Header().Set("Content-Type", "application/json")

	id := extractID(r.URL.Path)
	if id == 0 {
		sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
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
	var user models.User
	err := database.DB.QueryRow("SELECT id, name, last_name, email, avatar, bio, location, phone, created_at, verified, verified_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Avatar, &user.Bio, &user.Location, &user.Phone, &user.CreatedAt, &user.Verified, &user.VerifiedAt)

	if err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	sendSuccess(w, user)
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
