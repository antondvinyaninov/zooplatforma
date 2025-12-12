package handlers

import (
	"backend/middleware"
	"backend/models"
	"database"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		sendError(w, "Имя, email и пароль обязательны", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user
	result, err := database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		req.Name, req.Email, string(hashedPassword))
	if err != nil {
		// Log the actual error for debugging
		println("Database error:", err.Error())
		sendError(w, "Ошибка создания пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Generate token
	token, err := middleware.GenerateToken(int(id), req.Email)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	})

	sendSuccess(w, models.AuthResponse{
		User: models.UserResponse{
			ID:    int(id),
			Name:  req.Name,
			Email: req.Email,
		},
	})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		sendError(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Parse token
	token, err := middleware.ParseToken(cookie.Value)
	if err != nil {
		sendError(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	// Get user from database
	var user models.User
	err = database.DB.QueryRow("SELECT id, name, email, bio, phone, location, avatar, cover_photo, created_at FROM users WHERE id = ?", token.UserID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Bio, &user.Phone, &user.Location, &user.Avatar, &user.CoverPhoto, &user.CreatedAt)

	if err != nil {
		sendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	sendSuccess(w, models.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Bio:        user.Bio,
		Phone:      user.Phone,
		Location:   user.Location,
		Avatar:     user.Avatar,
		CoverPhoto: user.CoverPhoto,
		CreatedAt:  user.CreatedAt,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete cookie
	})

	sendSuccess(w, map[string]string{"message": "Logged out successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		sendError(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	// Get user
	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &hashedPassword)

	if err != nil {
		sendError(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		sendError(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	})

	sendSuccess(w, models.AuthResponse{
		User: models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
