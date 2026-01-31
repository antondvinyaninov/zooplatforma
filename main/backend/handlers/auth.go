package handlers

import (
	"backend/middleware"
	"backend/models"
	"bytes"
	"database"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

	// 🔥 НОВОЕ: Используем Auth Service
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:7100"
	}

	// Отправляем запрос к Auth Service
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(authServiceURL+"/api/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Auth Service error: %v", err)
		sendError(w, "Auth service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ от Auth Service
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Передаем ошибку от Auth Service
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// Парсим ответ
	var authResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
			User  struct {
				ID    int    `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		sendError(w, "Invalid auth response", http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookie с токеном от Auth Service
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authResp.Data.Token,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	})

	// Логируем регистрацию
	ipAddress := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	userID := authResp.Data.User.ID
	CreateUserLog(database.DB, userID, "register", "Пользователь зарегистрировался через Auth Service", ipAddress, userAgent)

	// Возвращаем ответ клиенту
	w.Write(body)

	log.Printf("✅ User registered via Auth Service: %s", authResp.Data.User.Email)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from Authorization header (priority) or cookie
	var token string

	// 1. Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Remove "Bearer " prefix if present
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			token = authHeader
		}
	}

	// 2. If no header, try cookie
	if token == "" {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			sendError(w, "Не авторизован", http.StatusUnauthorized)
			return
		}
		token = cookie.Value
	}

	// 3. If still no token, return 401
	if token == "" {
		sendError(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// 🔥 НОВОЕ: Используем Auth Service для получения данных пользователя
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:7100"
	}

	// Создаем запрос к Auth Service
	req, err := http.NewRequest("GET", authServiceURL+"/api/auth/me", nil)
	if err != nil {
		log.Printf("❌ Failed to create request: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Добавляем токен в заголовок
	req.Header.Set("Authorization", "Bearer "+token)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Auth Service error: %v", err)
		sendError(w, "Auth service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Передаем ошибку от Auth Service
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// Парсим ответ от Auth Service
	var authResp struct {
		Success bool `json:"success"`
		Data    struct {
			User struct {
				ID        int       `json:"id"`
				Email     string    `json:"email"`
				Name      string    `json:"name"`
				Bio       string    `json:"bio"`
				Phone     string    `json:"phone"`
				Avatar    string    `json:"avatar"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		log.Printf("❌ Failed to parse auth response: %v", err)
		sendError(w, "Invalid auth response", http.StatusInternalServerError)
		return
	}

	// Формируем ответ в формате Main Backend (для обратной совместимости)
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":         authResp.Data.User.ID,
				"name":       authResp.Data.User.Name,
				"email":      authResp.Data.User.Email,
				"bio":        authResp.Data.User.Bio,
				"phone":      authResp.Data.User.Phone,
				"avatar":     authResp.Data.User.Avatar,
				"created_at": authResp.Data.User.CreatedAt,
				// Поля для обратной совместимости (пока не в Auth Service)
				"last_name":          "",
				"location":           "",
				"cover_photo":        "",
				"profile_visibility": "public",
				"show_phone":         "nobody",
				"show_email":         "nobody",
				"allow_messages":     "everyone",
				"show_online":        "yes",
			},
			"token": token,
		},
		"token": token,
	}

	json.NewEncoder(w).Encode(response)
	log.Printf("✅ User profile loaded via Auth Service: %s", authResp.Data.User.Email)
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
		Domain:   "localhost", // Пустой для localhost, ".zooplatforma.ru" для production
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

	// 🔥 НОВОЕ: Используем Auth Service
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:7100"
	}

	// Отправляем запрос к Auth Service
	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(authServiceURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Auth Service error: %v", err)
		sendError(w, "Auth service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ от Auth Service
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Передаем ошибку от Auth Service
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// Парсим ответ
	var authResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
			User  struct {
				ID    int    `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		sendError(w, "Invalid auth response", http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookie с токеном от Auth Service
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authResp.Data.Token,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	})

	// Логируем успешный вход
	ipAddress := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	userID := authResp.Data.User.ID
	logSystemEvent("info", "auth", "login", "Пользователь вошел в систему (Auth Service)", &userID, ipAddress)
	CreateUserLog(database.DB, userID, "login", "Вход в систему через Auth Service", ipAddress, userAgent)

	// Возвращаем ответ клиенту
	w.Write(body)

	log.Printf("✅ User logged in via Auth Service: %s", authResp.Data.User.Email)
}
