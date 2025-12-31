package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

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

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

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

func AuthMiddleware(db interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				log.Printf("❌ Cookie not found: %v\n", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"success": false, "error": "Unauthorized: no auth token"}`))
				return
			}

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

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			log.Printf("✅ Auth successful, userID=%d\n", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
