package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// sqlPlaceholder возвращает правильный placeholder для SQL запроса
func sqlPlaceholder(index int) string {
	if isPostgres {
		return fmt.Sprintf("$%d", index)
	}
	return "?"
}

// sqlQuery конвертирует ? в $1, $2 для PostgreSQL
func sqlQuery(query string) string {
	if !isPostgres {
		return query
	}

	result := ""
	paramIndex := 1
	for _, char := range query {
		if char == '?' {
			result += fmt.Sprintf("$%d", paramIndex)
			paramIndex++
		} else {
			result += string(char)
		}
	}
	return result
}

// getCookieDomain возвращает правильный домен для cookie
func getCookieDomain() string {
	if os.Getenv("ENVIRONMENT") == "production" {
		return ""
	}
	return "localhost"
}

// getCookieSecure возвращает правильное значение Secure для cookie
func getCookieSecure() bool {
	return os.Getenv("ENVIRONMENT") == "production"
}

// getCookieSameSite возвращает правильное значение SameSite для cookie
func getCookieSameSite() http.SameSite {
	if os.Getenv("ENVIRONMENT") == "production" {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

// User - структура пользователя
type AuthUser struct {
	ID            int     `json:"id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	LastName      string  `json:"last_name"`
	Avatar        string  `json:"avatar"`
	CoverPhoto    string  `json:"cover_photo"`
	Bio           string  `json:"bio"`
	Phone         string  `json:"phone"`
	Location      string  `json:"location"`
	DateOfBirth   *string `json:"date_of_birth,omitempty"`
	Gender        string  `json:"gender"`
	City          string  `json:"city"`
	Country       string  `json:"country"`
	Role          string  `json:"role"`
	EmailVerified bool    `json:"email_verified"`
	CreatedAt     string  `json:"created_at"`
}

// RegisterHandler - регистрация нового пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		LastName string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, `{"success":false,"error":"Email, password and name are required"}`, http.StatusBadRequest)
		return
	}

	// Проверить существование email
	var exists bool
	err := authDB.QueryRow(sqlQuery("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"), req.Email).Scan(&exists)
	if err != nil {
		log.Printf("❌ Database error: %v", err)
		http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, `{"success":false,"error":"Email already exists"}`, http.StatusConflict)
		return
	}

	// Хешировать пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ Failed to hash password: %v", err)
		http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Создать пользователя
	var userID int64
	err = authDB.QueryRow(sqlQuery(`
		INSERT INTO users (email, password, name, last_name)
		VALUES (?, ?, ?, ?)
		RETURNING id
	`), req.Email, string(hashedPassword), req.Name, req.LastName).Scan(&userID)

	if err != nil {
		log.Printf("❌ Failed to create user: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	// Создать JWT токен
	token, err := createJWT(int(userID), req.Email, "user")
	if err != nil {
		log.Printf("❌ Failed to create token: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to create token"}`, http.StatusInternalServerError)
		return
	}

	// Установить cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Domain:   getCookieDomain(),
		HttpOnly: false,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
		MaxAge:   86400 * 7, // 7 дней
	})

	// Ответ
	user := AuthUser{
		ID:       int(userID),
		Email:    req.Email,
		Name:     req.Name,
		LastName: req.LastName,
		Role:     "user",
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"token": token,
			"user":  user,
		},
		"message": "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ User registered: %s", req.Email)
}

// LoginHandler - вход в систему
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Найти пользователя
	var user AuthUser
	var passwordHash string
	var avatar sql.NullString

	err := authDB.QueryRow(sqlQuery(`
		SELECT id, email, password, name, last_name, avatar
		FROM users WHERE email = ?
	`), req.Email).Scan(&user.ID, &user.Email, &passwordHash, &user.Name, &user.LastName, &avatar)

	if avatar.Valid {
		user.Avatar = avatar.String
	}

	if err == sql.ErrNoRows {
		http.Error(w, `{"success":false,"error":"Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("❌ Database error: %v", err)
		http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Проверить пароль
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		http.Error(w, `{"success":false,"error":"Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Загрузить роли из user_roles
	roles := []string{}
	rows, err := authDB.Query(sqlQuery("SELECT role FROM user_roles WHERE user_id = ?"), user.ID)
	if err != nil {
		log.Printf("❌ Failed to load roles: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var role string
			if err := rows.Scan(&role); err == nil {
				roles = append(roles, role)
			}
		}
	}

	// Определить основную роль
	user.Role = "user"
	for _, role := range roles {
		if role == "superadmin" {
			user.Role = "superadmin"
			break
		} else if role == "admin" && user.Role != "superadmin" {
			user.Role = "admin"
		} else if role == "moderator" && user.Role != "superadmin" && user.Role != "admin" {
			user.Role = "moderator"
		}
	}

	// Создать JWT токен
	token, err := createJWT(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("❌ Failed to create token: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to create token"}`, http.StatusInternalServerError)
		return
	}

	// Установить cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Domain:   getCookieDomain(),
		HttpOnly: false,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
		MaxAge:   86400 * 7, // 7 дней
	})

	// Ответ
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"token": token,
			"user":  user,
			"roles": roles,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ User logged in: %s (role=%s)", req.Email, user.Role)
}

// GetMeHandler - получить текущего пользователя
func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		// Попробовать cookie
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token = "Bearer " + cookie.Value
		}
	}

	if token == "" {
		http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Убрать "Bearer "
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Парсить JWT
	claims, err := parseJWT(token)
	if err != nil {
		http.Error(w, `{"success":false,"error":"Invalid token"}`, http.StatusUnauthorized)
		return
	}

	// Получить пользователя из БД
	var user AuthUser
	var avatar, lastName, bio, phone, location, coverPhoto sql.NullString
	var createdAt sql.NullString

	err = authDB.QueryRow(sqlQuery(`
		SELECT id, email, name, last_name, avatar, bio, phone, location, cover_photo, created_at
		FROM users WHERE id = ?
	`), claims.UserID).Scan(&user.ID, &user.Email, &user.Name, &lastName, &avatar, &bio, &phone, &location, &coverPhoto, &createdAt)

	if err != nil {
		log.Printf("❌ Failed to get user: %v", err)
		http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)
		return
	}

	// Копируем данные из NullString
	if lastName.Valid {
		user.LastName = lastName.String
	}
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if bio.Valid {
		user.Bio = bio.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	if location.Valid {
		user.Location = location.String
	}
	if coverPhoto.Valid {
		user.CoverPhoto = coverPhoto.String
	}
	if createdAt.Valid {
		user.CreatedAt = createdAt.String
	}

	// Загрузить роли
	roles := []string{}
	rows, err := authDB.Query(sqlQuery("SELECT role FROM user_roles WHERE user_id = ?"), user.ID)
	if err != nil {
		log.Printf("❌ Failed to load roles: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var role string
			if err := rows.Scan(&role); err == nil {
				roles = append(roles, role)
			}
		}
	}

	// Определить основную роль
	user.Role = "user"
	for _, role := range roles {
		if role == "superadmin" {
			user.Role = "superadmin"
			break
		} else if role == "admin" && user.Role != "superadmin" {
			user.Role = "admin"
		} else if role == "moderator" && user.Role != "superadmin" && user.Role != "admin" {
			user.Role = "moderator"
		}
	}

	// Ответ
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user":  user,
			"token": token,
			"roles": roles,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ User loaded: id=%d, email=%s, role=%s", user.ID, user.Email, user.Role)
}

// LogoutHandler - выход из системы
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   getCookieDomain(),
		HttpOnly: false,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"message":"Logged out successfully"}`))

	log.Printf("✅ User logged out")
}

// Helper функции для работы с JWT уже определены в auth.go:
// - createJWT(userID int, email, role string) (string, error)
// - parseJWT(tokenString string) (*Claims, error)

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
