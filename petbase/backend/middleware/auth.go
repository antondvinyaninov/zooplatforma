package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT secret из переменной окружения
var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// В разработке используем дефолтный ключ
		// В продакшене ОБЯЗАТЕЛЬНО установить JWT_SECRET!
		return "dev-secret-key-change-in-production"
	}
	return secret
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
	if userID, ok := r.Context().Value("user_id").(int); ok {
		return userID, ok
	}

	return 0, false
}

// SetUserID устанавливает user_id в контекст запроса
func SetUserID(r *http.Request, userID int) *http.Request {
	ctx := context.WithValue(r.Context(), "user_id", userID)
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
