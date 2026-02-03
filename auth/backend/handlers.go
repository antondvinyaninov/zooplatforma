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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte
var isPostgres bool

// sqlPlaceholder возвращает правильный placeholder для SQL запроса
func sqlPlaceholder(index int) string {
	if isPostgres {
		return fmt.Sprintf("$%d", index)
	}
	return "?"
}

// getCookieDomain возвращает правильный домен для cookie
func getCookieDomain() string {
	if os.Getenv("ENVIRONMENT") == "production" {
		// В production используем домен из переменной окружения или пустую строку
		// Пустая строка означает что cookie будет работать для текущего домена
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
		return http.SameSiteNoneMode // None для cross-origin в production
	}
	return http.SameSiteLaxMode // Lax для localhost
}

// InitJWTSecret - инициализировать JWT secret (вызывается после загрузки .env)
func InitJWTSecret() {
	// Определяем тип БД
	isPostgres = os.Getenv("ENVIRONMENT") == "production"
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("❌ JWT_SECRET not set in environment")
	}
	jwtSecret = []byte(secret)
	log.Printf("✅ JWT Secret initialized")
}

// Claims - JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// User - структура пользователя
type User struct {
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

// registerHandler - регистрация нового пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
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
	var checkQuery string
	if isPostgres {
		checkQuery = "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	} else {
		checkQuery = "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	}
	err := db.QueryRow(checkQuery, req.Email).Scan(&exists)
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
	var insertQuery string
	if isPostgres {
		insertQuery = `
			INSERT INTO users (email, password, name, last_name)
			VALUES ($1, $2, $3, $4)
		`
	} else {
		insertQuery = `
			INSERT INTO users (email, password, name, last_name)
			VALUES (?, ?, ?, ?)
		`
	}
	result, err := db.Exec(insertQuery, req.Email, string(hashedPassword), req.Name, req.LastName)

	if err != nil {
		log.Printf("❌ Failed to create user: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()

	// Создать JWT токен
	token, err := createJWT(int(userID), req.Email, "user")
	if err != nil {
		log.Printf("❌ Failed to create token: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to create token"}`, http.StatusInternalServerError)
		return
	}

	// Установить cookie с токеном (для работы между портами localhost)
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
	user := User{
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

// loginHandler - вход в систему
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Найти пользователя
	var user User
	var passwordHash string
	var avatar sql.NullString

	var selectQuery string
	if isPostgres {
		selectQuery = `
			SELECT id, email, password, name, last_name, avatar
			FROM users WHERE email = $1
		`
	} else {
		selectQuery = `
			SELECT id, email, password, name, last_name, avatar
			FROM users WHERE email = ?
		`
	}

	err := db.QueryRow(selectQuery, req.Email).Scan(&user.ID, &user.Email, &passwordHash, &user.Name, &user.LastName, &avatar)

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
	rows, err := db.Query(sqlQuery("SELECT role FROM user_roles WHERE user_id = ?"), user.ID)
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

	// Определить основную роль (приоритет: superadmin > admin > moderator > user)
	user.Role = "user" // По умолчанию
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

	// Установить cookie с токеном (для работы между портами localhost)
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

	log.Printf("✅ User logged in: %s (role=%s, all_roles=%v)", req.Email, user.Role, roles)
}

// getMeHandler - получить текущего пользователя
func getMeHandler(w http.ResponseWriter, r *http.Request) {
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
	var user User
	var avatar, lastName, bio, phone, location, coverPhoto sql.NullString
	var createdAt sql.NullString

	err = db.QueryRow(sqlQuery(`
		SELECT id, email, name, last_name, avatar, bio, phone, location, cover_photo, created_at
		FROM users WHERE id = ?
	`), claims.UserID).Scan(&user.ID, &user.Email, &user.Name, &lastName, &avatar, &bio, &phone, &location, &coverPhoto, &createdAt)

	// СНАЧАЛА проверяем ошибку!
	if err != nil {
		log.Printf("❌ Failed to get user: %v", err)
		http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)
		return
	}

	// ПОТОМ копируем данные из NullString
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

	// Логируем полученные данные для отладки
	log.Printf("🔍 User data from DB: id=%d, name=%s, last_name=%s, phone=%s, location=%s, bio=%s",
		user.ID, user.Name, user.LastName, user.Phone, user.Location, user.Bio)

	// Загрузить роли из user_roles
	roles := []string{}
	rows, err := db.Query(sqlQuery("SELECT role FROM user_roles WHERE user_id = ?"), user.ID)
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

	// Определить основную роль (приоритет: superadmin > admin > moderator > user)
	user.Role = "user" // По умолчанию
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

	log.Printf("✅ User loaded: id=%d, email=%s, role=%s, all_roles=%v", user.ID, user.Email, user.Role, roles)

	// Логируем данные перед отправкой
	log.Printf("📤 Sending user to client: id=%d, name=%s, last_name=%s, phone=%s, location=%s, bio=%s",
		user.ID, user.Name, user.LastName, user.Phone, user.Location, user.Bio)

	// Ответ
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user":  user,
			"token": token,
			"roles": roles, // Все роли пользователя
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// verifyTokenHandler - проверить токен (для других сервисов)
func verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, `{"success":false,"error":"No token provided"}`, http.StatusUnauthorized)
		return
	}

	// Убрать "Bearer "
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Парсить JWT
	claims, err := parseJWT(token)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"data": map[string]interface{}{
				"valid": false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Получить пользователя
	var user User
	var avatar sql.NullString

	err = db.QueryRow(sqlQuery(`
		SELECT id, email, name, last_name, avatar
		FROM users WHERE id = ?
	`), claims.UserID).Scan(&user.ID, &user.Email, &user.Name, &user.LastName, &avatar)

	if avatar.Valid {
		user.Avatar = avatar.String
	}

	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"data": map[string]interface{}{
				"valid": false,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Загрузить роли из user_roles
	roles := []string{}
	rows, err := db.Query(sqlQuery("SELECT role FROM user_roles WHERE user_id = ?"), user.ID)
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
			"valid": true,
			"user":  user,
			"roles": roles,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Заглушки для остальных handlers
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Удалить cookie (для всех портов localhost)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   getCookieDomain(),
		HttpOnly: false,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(), // Lax для localhost
		MaxAge:   -1,                  // Удалить cookie
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"message":"Logged out successfully"}`))

	log.Printf("✅ User logged out")
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":false,"error":"Not implemented yet"}`))
}

func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":false,"error":"Not implemented yet"}`))
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":false,"error":"Not implemented yet"}`))
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":false,"error":"Not implemented yet"}`))
}

func sendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":false,"error":"Not implemented yet"}`))
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":false,"error":"Not implemented yet"}`))
}

// Helper функции
func createJWT(userID int, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func parseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
