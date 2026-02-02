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
var jwtSecret []byte

// InitJWTSecret инициализирует JWT secret из переменной окружения
// ДОЛЖНА быть вызвана ПОСЛЕ godotenv.Load()
func InitJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// В разработке используем дефолтный ключ
		// В продакшене ОБЯЗАТЕЛЬНО установить JWT_SECRET!
		secret = "dev-secret-key-change-in-production"
		fmt.Printf("⚠️ JWT_SECRET not set, using default\n")
	} else {
		fmt.Printf("✅ JWT_SECRET loaded: %s...\n", secret[:10])
	}
	jwtSecret = []byte(secret)
}

func getJWTSecret() string {
	return string(jwtSecret)
}

// Claims структура для JWT токена (совместимая с Main backend)
type Claims struct {
	UserID int      `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
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
	fmt.Printf("🔍 Validating token: %s...\n", tokenString[:min(20, len(tokenString))])
	fmt.Printf("🔑 Using JWT secret: %s...\n", string(jwtSecret)[:min(10, len(jwtSecret))])

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
		return 0, fmt.Errorf("parse error: %w", err)
	}

	if !token.Valid {
		fmt.Printf("❌ Token not valid\n")
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		fmt.Printf("❌ Invalid claims type\n")
		return 0, fmt.Errorf("invalid claims")
	}

	fmt.Printf("✅ Token valid, user_id=%d, email=%s\n", claims.UserID, claims.Email)
	return claims.UserID, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetUserIDFromRequest получает user_id из JWT токена или временных методов
func GetUserIDFromRequest(r *http.Request) (int, bool) {
	// 1. Пробуем получить из JWT токена (приоритет)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			fmt.Printf("🔍 Validating JWT token: %s...\n", parts[1][:20])
			userID, err := ValidateToken(parts[1])
			if err != nil {
				fmt.Printf("❌ JWT validation failed: %v\n", err)
			} else if userID > 0 {
				fmt.Printf("✅ JWT valid, user_id=%d\n", userID)
				return userID, true
			}
		}
	}

	// 2. Временно: получаем из заголовка X-User-ID (только для разработки!)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		var userID int
		if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err == nil {
			fmt.Printf("⚠️ Using X-User-ID header: %d\n", userID)
			return userID, true
		}
	}

	// 3. Пробуем получить из cookie
	cookie, err := r.Cookie("user_id")
	if err == nil {
		var userID int
		if _, err := fmt.Sscanf(cookie.Value, "%d", &userID); err == nil {
			fmt.Printf("⚠️ Using cookie user_id: %d\n", userID)
			return userID, true
		}
	}

	// 4. Пробуем получить из контекста (если уже установлен)
	if userID, ok := r.Context().Value("user_id").(int); ok {
		fmt.Printf("⚠️ Using context user_id: %d\n", userID)
		return userID, ok
	}

	fmt.Printf("❌ No user_id found in request\n")
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
		authHeader := r.Header.Get("Authorization")
		fmt.Printf("🔐 RequireAuth: Authorization header: %s\n", authHeader)

		userID, ok := GetUserIDFromRequest(r)
		if !ok || userID == 0 {
			fmt.Printf("❌ RequireAuth: No valid user_id found\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success": false, "error": "Unauthorized: authentication required"}`))
			return
		}

		fmt.Printf("✅ RequireAuth: user_id=%d\n", userID)
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
