package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWT инициализирует JWT secret
func InitJWT() {
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

// UserContext - ключ для контекста пользователя
type contextKey string

const UserContextKey contextKey = "user"

// UserInfo - информация о пользователе из JWT
type UserInfo struct {
	ID    int
	Email string
	Role  string
}

// AuthMiddleware - проверяет JWT токен
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получить токен
		token := extractToken(r)
		if token == "" {
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Парсить JWT
		claims, err := parseJWT(token)
		if err != nil {
			log.Printf("⚠️ Invalid token: %v", err)
			sendError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Проверить что токен не истек
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			sendError(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Добавить информацию о пользователе в контекст
		userInfo := &UserInfo{
			ID:    claims.UserID,
			Email: claims.Email,
			Role:  claims.Role,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userInfo)

		// Добавить заголовки для backend сервисов
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", userInfo.ID))
		r.Header.Set("X-User-Email", userInfo.Email)
		r.Header.Set("X-User-Role", userInfo.Role)

		log.Printf("✅ Authenticated: user_id=%d, email=%s, role=%s", userInfo.ID, userInfo.Email, userInfo.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnlyMiddleware - проверяет что пользователь админ
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(*UserInfo)
		if !ok {
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверить роль
		if userInfo.Role != "admin" && userInfo.Role != "superadmin" {
			log.Printf("⚠️ Access denied: user_id=%d, role=%s (required: admin)", userInfo.ID, userInfo.Role)
			sendError(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		log.Printf("✅ Admin access granted: user_id=%d, role=%s", userInfo.ID, userInfo.Role)
		next.ServeHTTP(w, r)
	})
}

// extractToken извлекает токен из запроса
func extractToken(r *http.Request) string {
	// 1. Попробовать Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Убрать "Bearer " prefix
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:]
		}
		return authHeader
	}

	// 2. Попробовать cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		return cookie.Value
	}

	return ""
}

// parseJWT парсит JWT токен
func parseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверить метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// VerifyTokenWithAuthService проверяет токен через Auth Service (альтернативный метод)
func VerifyTokenWithAuthService(token string, authServiceURL string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", authServiceURL+"/api/auth/verify", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			Valid bool `json:"valid"`
			User  struct {
				ID    int    `json:"id"`
				Email string `json:"email"`
				Role  string `json:"role"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.Success || !result.Data.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &UserInfo{
		ID:    result.Data.User.ID,
		Email: result.Data.User.Email,
		Role:  result.Data.User.Role,
	}, nil
}

// sendError отправляет ошибку в формате JSON
func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// GetUserFromContext получает информацию о пользователе из контекста
func GetUserFromContext(r *http.Request) (*UserInfo, bool) {
	userInfo, ok := r.Context().Value(UserContextKey).(*UserInfo)
	return userInfo, ok
}
