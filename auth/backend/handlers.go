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
	Bio           string  `json:"bio"`
	Phone         string  `json:"phone"`
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
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
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
	result, err := db.Exec(`
		INSERT INTO users (email, password, name, last_name)
		VALUES ($1, $2, $3, $4)
	`, req.Email, string(hashedPassword), req.Name, req.LastName)

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
		Domain:   "localhost",          // Явно указываем localhost
		HttpOnly: false,                // false чтобы JS мог читать (для отладки)
		Secure:   false,                // false для localhost
		SameSite: http.SameSiteLaxMode, // Lax для localhost (None требует Secure=true)
		MaxAge:   86400 * 7,            // 7 дней
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

	err := db.QueryRow(`
		SELECT id, email, password, name, last_name, avatar
		FROM users WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Email, &passwordHash, &user.Name, &user.LastName, &avatar)

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
	rows, err := db.Query("SELECT role FROM user_roles WHERE user_id = $1", user.ID)
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
		Domain:   "localhost",          // Явно указываем localhost
		HttpOnly: false,                // false чтобы JS мог читать (для отладки)
		Secure:   false,                // false для localhost
		SameSite: http.SameSiteLaxMode, // Lax для localhost (None требует Secure=true)
		MaxAge:   86400 * 7,            // 7 дней
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
	err = db.QueryRow(`
		SELECT id, email, name, last_name, avatar, email_verified, created_at
		FROM users WHERE id = $1
	`, claims.UserID).Scan(&user.ID, &user.Email, &user.Name, &user.LastName, &user.Avatar, &user.EmailVerified, &user.CreatedAt)

	if err != nil {
		log.Printf("❌ Failed to get user: %v", err)
		http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)
		return
	}

	// Загрузить роли из user_roles
	roles := []string{}
	rows, err := db.Query("SELECT role FROM user_roles WHERE user_id = $1", user.ID)
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
	err = db.QueryRow(`
		SELECT id, email, name, last_name, avatar
		FROM users WHERE id = $1
	`, claims.UserID).Scan(&user.ID, &user.Email, &user.Name, &user.LastName, &user.Avatar)

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
	rows, err := db.Query("SELECT role FROM user_roles WHERE user_id = $1", user.ID)
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
		Domain:   "localhost",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode, // Lax для localhost
		MaxAge:   -1,                   // Удалить cookie
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
