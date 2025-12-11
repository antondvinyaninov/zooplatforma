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
		sendError(w, "Name, email and password are required", http.StatusBadRequest)
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
		sendError(w, "Email already exists", http.StatusConflict)
		return
	}

	id, _ := result.LastInsertId()

	// Generate token
	token, err := middleware.GenerateToken(int(id), req.Email)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:    int(id),
			Name:  req.Name,
			Email: req.Email,
		},
	})
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
		sendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get user
	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &hashedPassword)

	if err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
