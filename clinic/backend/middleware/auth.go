package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

var jwtSecret []byte

func InitJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("⚠️ JWT_SECRET not set, using default (NOT FOR PRODUCTION)")
		secret = "your-secret-key-here-change-in-production"
	}
	jwtSecret = []byte(secret)
	log.Printf("✅ JWT Secret initialized (length: %d)", len(jwtSecret))
}

func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🔐 Auth check for: %s %s", r.Method, r.URL.Path)

			// Получаем токен из cookie
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				log.Printf("❌ No auth_token cookie: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			log.Printf("🍪 Token from cookie: %s...", tokenString[:min(len(tokenString), 20)])

			// Парсим и валидируем токен
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Printf("❌ Unexpected signing method: %v", token.Header["alg"])
					return nil, http.ErrAbortHandler
				}
				return jwtSecret, nil
			})

			if err != nil {
				log.Printf("❌ Token parse error: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				log.Printf("❌ Token is invalid")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Извлекаем claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Printf("❌ Failed to extract claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				log.Printf("❌ user_id not found in claims or wrong type")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			log.Printf("✅ User authenticated: ID=%d", int(userID))

			// Добавляем user_id в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func extractTokenFromHeader(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
