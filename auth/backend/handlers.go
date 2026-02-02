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
	err = db.QueryRow(sqlQuery(`SELECT id, email, name, last_name, avatar, email_verified, created_at FROM users WHERE id = ?`), claims.UserID).Scan(&user.ID, &user.Email, &user.Name, &user.LastName, &user.Avatar, &user.EmailVerified, &user.CreatedAt)

	if err != nil {
		log.Printf("❌ Failed to get user: %v", err)
		http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)
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
	err = db.QueryRow(sqlQuery(`SELECT id, email, name, last_name, avatar FROM users WHERE id = ?`), claims.UserID).Scan(&user.ID, &user.Email, &user.Name, &user.LastName, &user.Avatar)

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
