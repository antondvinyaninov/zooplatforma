package handlers

import (
	"backend/middleware"
	"backend/models"
	"database"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// logSystemEvent - логирует событие в системе
func logSystemEvent(level, category, action, message string, userID *int, ipAddress string) {
	query := `
		INSERT INTO system_logs (level, category, action, message, user_id, ip_address, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	database.DB.Exec(query, level, category, action, message, userID, ipAddress, time.Now())
}

// getUserRoles получает роли пользователя из таблицы admins
func getUserRoles(userID int) []string {
	roles := []string{"user"} // По умолчанию все пользователи имеют роль "user"

	// Проверяем, есть ли у пользователя роль админа
	var adminRole string
	err := database.DB.QueryRow("SELECT role FROM admins WHERE user_id = ?", userID).Scan(&adminRole)
	if err == nil {
		roles = append(roles, adminRole)
	}

	return roles
}

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

	// Логируем регистрацию
	ipAddress := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	CreateUserLog(database.DB, int(id), "register", "Пользователь зарегистрировался", ipAddress, userAgent)

	// Получаем роли пользователя (по умолчанию только "user")
	roles := getUserRoles(int(id))

	// Generate token
	token, err := middleware.GenerateToken(int(id), req.Email, roles)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set httpOnly cookie (доступен всем поддоменам для SSO)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Domain:   "", // Пустой для localhost, ".zooplatforma.ru" для production
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

	// Get user from database - используем sql.NullString для nullable полей
	var user models.User
	var lastName, bio, phone, location, avatar, coverPhoto sql.NullString
	var profileVisibility, showPhone, showEmail, allowMessages, showOnline sql.NullString

	err = database.DB.QueryRow(`
		SELECT id, name, last_name, email, bio, phone, location, avatar, cover_photo, 
		       profile_visibility, show_phone, show_email, allow_messages, show_online, created_at 
		FROM users WHERE id = ?`, token.UserID).
		Scan(&user.ID, &user.Name, &lastName, &user.Email, &bio, &phone, &location, &avatar, &coverPhoto,
			&profileVisibility, &showPhone, &showEmail, &allowMessages, &showOnline, &user.CreatedAt)

	if err != nil {
		sendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Преобразуем NullString в обычные строки
	user.LastName = lastName.String
	user.Bio = bio.String
	user.Phone = phone.String
	user.Location = location.String
	user.Avatar = avatar.String
	user.CoverPhoto = coverPhoto.String
	user.ProfileVisibility = profileVisibility.String
	user.ShowPhone = showPhone.String
	user.ShowEmail = showEmail.String
	user.AllowMessages = allowMessages.String
	user.ShowOnline = showOnline.String

	// Устанавливаем значения по умолчанию, если пусто
	if user.ProfileVisibility == "" {
		user.ProfileVisibility = "public"
	}
	if user.ShowPhone == "" {
		user.ShowPhone = "nobody"
	}
	if user.ShowEmail == "" {
		user.ShowEmail = "nobody"
	}
	if user.AllowMessages == "" {
		user.AllowMessages = "everyone"
	}
	if user.ShowOnline == "" {
		user.ShowOnline = "yes"
	}

	sendSuccess(w, models.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		LastName:          user.LastName,
		Email:             user.Email,
		Bio:               user.Bio,
		Phone:             user.Phone,
		Location:          user.Location,
		Avatar:            user.Avatar,
		CoverPhoto:        user.CoverPhoto,
		ProfileVisibility: user.ProfileVisibility,
		ShowPhone:         user.ShowPhone,
		ShowEmail:         user.ShowEmail,
		AllowMessages:     user.AllowMessages,
		ShowOnline:        user.ShowOnline,
		CreatedAt:         user.CreatedAt,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем user_id из токена перед удалением cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		token, err := middleware.ParseToken(cookie.Value)
		if err == nil {
			// Логируем выход
			ipAddress := r.RemoteAddr
			userAgent := r.Header.Get("User-Agent")
			CreateUserLog(database.DB, token.UserID, "logout", "Выход из системы", ipAddress, userAgent)
		}
	}

	// Clear cookie (для всех поддоменов)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   "", // Пустой для localhost, ".zooplatforma.ru" для production
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete cookie
	})

	sendSuccess(w, map[string]string{"message": "Logged out successfully"})
}

func VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		sendError(w, "Токен не найден", http.StatusUnauthorized)
		return
	}

	// Parse token
	token, err := middleware.ParseToken(cookie.Value)
	if err != nil {
		sendError(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	// Verify user exists
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM users WHERE id = ?", token.UserID).Scan(&exists)
	if err != nil {
		sendError(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	sendSuccess(w, map[string]interface{}{
		"user_id": token.UserID,
		"email":   token.Email,
		"roles":   token.Roles,
		"valid":   true,
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

	// Получаем роли пользователя
	roles := getUserRoles(user.ID)

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Email, roles)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set httpOnly cookie (доступен всем поддоменам для SSO)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Domain:   "", // Пустой для localhost, ".zooplatforma.ru" для production
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	})

	// Логируем успешный вход
	ipAddress := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	logSystemEvent("info", "auth", "login", "Пользователь вошел в систему", &user.ID, ipAddress)
	CreateUserLog(database.DB, user.ID, "login", "Вход в систему", ipAddress, userAgent)

	sendSuccess(w, models.AuthResponse{
		User: models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
