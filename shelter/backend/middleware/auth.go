package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT secret из переменной окружения
var jwtSecret []byte

// Context keys
type contextKey string

const UserIDKey contextKey = "user_id"

func InitJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-key-change-in-production"
		log.Println("⚠️  JWT_SECRET not set, using default")
	}
	jwtSecret = []byte(secret)
	log.Printf("🔑 JWT Secret initialized: length=%d\n", len(jwtSecret))
}

// Claims структура для JWT токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken создаёт JWT токен для пользователя
func GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken проверяет JWT токен и возвращает user_id
func ValidateToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	return claims.UserID, nil
}

// GetUserIDFromRequest получает user_id из JWT токена или временных методов
func GetUserIDFromRequest(r *http.Request) (int, bool) {
	// 1. Пробуем получить из JWT токена (приоритет)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			userID, err := ValidateToken(parts[1])
			if err == nil && userID > 0 {
				return userID, true
			}
		}
	}

	// 2. Временно: получаем из заголовка X-User-ID (только для разработки!)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		var userID int
		if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err == nil {
			return userID, true
		}
	}

	// 3. Пробуем получить из cookie
	cookie, err := r.Cookie("user_id")
	if err == nil {
		var userID int
		if _, err := fmt.Sscanf(cookie.Value, "%d", &userID); err == nil {
			return userID, true
		}
	}

	// 4. Пробуем получить из контекста (если уже установлен)
	if userID, ok := r.Context().Value(UserIDKey).(int); ok {
		return userID, ok
	}

	return 0, false
}

// SetUserID устанавливает user_id в контекст запроса
func SetUserID(r *http.Request, userID int) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, userID)
	return r.WithContext(ctx)
}

// RequireAuth middleware - требует аутентификации
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromRequest(r)
		if !ok || userID == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success": false, "error": "Unauthorized: authentication required"}`))
			return
		}

		// Добавляем user_id в контекст
		r = SetUserID(r, userID)
		next(w, r)
	}
}

// OptionalAuth middleware - опциональная аутентификация
func OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromRequest(r)
		if ok && userID > 0 {
			r = SetUserID(r, userID)
		}
		next(w, r)
	}
}

// AuthMiddleware - middleware для проверки JWT токена из cookie
// Возвращает middleware для использования с http.Handler
func AuthMiddleware(db interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из cookie
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				log.Printf("❌ Cookie not found: %v\n", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"success": false, "error": "Unauthorized: no auth token"}`))
				return
			}

			// Валидируем токен
			userID, err := ValidateToken(cookie.Value)
			if err != nil {
				log.Printf("❌ Token validation failed: %v\n", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"success": false, "error": "Unauthorized: invalid token"}`))
				return
			}

			if userID == 0 {
				log.Printf("❌ Invalid userID in token\n")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"success": false, "error": "Unauthorized: invalid user"}`))
				return
			}

			// Добавляем userID в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			log.Printf("✅ Auth successful, userID=%d\n", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
