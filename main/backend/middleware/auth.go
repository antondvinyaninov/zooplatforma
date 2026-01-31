package middleware

import (
	"backend/models"
	"context"
	"database"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
		log.Println("⚠️  JWT_SECRET not set, using default")
	}
	jwtSecret = []byte(secret)
	log.Printf("🔑 JWT Secret initialized: length=%d, first 10=%s\n", len(jwtSecret), string(jwtSecret[:10]))
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from cookie first
		cookie, err := r.Cookie("auth_token")
		var tokenString string

		if err == nil {
			tokenString = cookie.Value
			log.Printf("✅ Cookie found: %s...\n", tokenString[:20])
		} else {
			log.Printf("❌ Cookie not found: %v\n", err)
			// Fallback to Authorization header for backward compatibility
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendError(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}

			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				sendError(w, "Неверный формат авторизации", http.StatusUnauthorized)
				return
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("❌ Token parse error: %v\n", err)
			sendError(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("❌ Token invalid\n")
			sendError(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Извлекаем данные из токена и добавляем в контекст
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("❌ Failed to parse claims\n")
			sendError(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("❌ user_id not found in token\n")
			sendError(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", int(userID))

		// Обновляем активность пользователя
		if database.DB != nil {
			go updateUserActivity(int(userID), r.RemoteAddr, r.Header.Get("User-Agent"))
		}

		log.Printf("✅ Token valid, userID=%d, calling next handler\n", int(userID))
		next(w, r.WithContext(ctx))
	}
}

// updateUserActivity обновляет время последней активности пользователя
func updateUserActivity(userID int, ipAddress, userAgent string) {
	_, err := database.DB.Exec(`
		INSERT INTO user_activity (user_id, last_seen, ip_address, user_agent)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			last_seen = excluded.last_seen,
			ip_address = excluded.ip_address,
			user_agent = excluded.user_agent
	`, userID, time.Now(), ipAddress, userAgent)

	if err != nil {
		log.Printf("⚠️  Failed to update user activity: %v", err)
	}
}

func GenerateToken(userID int, email string, roles []string) (string, error) {
	log.Printf("🔑 GenerateToken: secret length=%d, first 10 chars=%s\n", len(jwtSecret), string(jwtSecret[:10]))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"roles":   roles,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 дней
		"iat":     time.Now().Unix(),
	})

	return token.SignedString(jwtSecret)
}

type TokenClaims struct {
	UserID int      `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

func ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// Roles могут отсутствовать в старых токенах
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		for _, r := range rolesInterface {
			if roleStr, ok := r.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return &TokenClaims{
		UserID: int(userID),
		Email:  email,
		Roles:  roles,
	}, nil
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.Response{Success: false, Error: message})
}

// OptionalAuthMiddleware извлекает userID если пользователь авторизован, но не требует авторизации
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from cookie
		cookie, err := r.Cookie("auth_token")
		var tokenString string

		if err == nil {
			tokenString = cookie.Value
		} else {
			// Fallback to Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Если токена нет - продолжаем без авторизации (userID = 0)
		if tokenString == "" {
			log.Printf("🔓 OptionalAuth: no token, continuing without auth")
			next(w, r)
			return
		}

		// Парсим токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Printf("⚠️ OptionalAuth: invalid token, continuing without auth")
			next(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("⚠️ OptionalAuth: invalid claims, continuing without auth")
			next(w, r)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("⚠️ OptionalAuth: no user_id in claims, continuing without auth")
			next(w, r)
			return
		}

		// Добавляем userID в контекст
		ctx := context.WithValue(r.Context(), "userID", int(userID))
		log.Printf("✅ OptionalAuth: userID=%d extracted", int(userID))
		next(w, r.WithContext(ctx))
	}
}
