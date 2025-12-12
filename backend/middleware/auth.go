package middleware

import (
	"backend/models"
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
