package middleware

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "your-secret-key-change-in-production"))

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from cookie first
		cookie, err := r.Cookie("auth_token")
		var tokenString string

		if err == nil {
			tokenString = cookie.Value
		} else {
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

		if err != nil || !token.Valid {
			sendError(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func GenerateToken(userID int, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
	})

	return token.SignedString(jwtSecret)
}

type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
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

	return &TokenClaims{
		UserID: int(userID),
		Email:  email,
	}, nil
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.Response{Success: false, Error: message})
}
