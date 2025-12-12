package middleware

import (
	"context"
	"database"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int      `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// SuperAdminMiddleware проверяет, что пользователь является суперадмином
func SuperAdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем общий токен из cookie (SSO)
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			sendError(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		// Парсим токен
		token, err := ParseToken(cookie.Value)
		if err != nil {
			sendError(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Проверяем, что у пользователя есть роль superadmin
		if !contains(token.Roles, "superadmin") {
			sendError(w, "Доступ запрещен. Требуются права суперадмина", http.StatusForbidden)
			return
		}

		// Получаем adminID из таблицы admins
		var adminID int
		err = database.DB.QueryRow(`
			SELECT id FROM admins WHERE user_id = ? AND role = 'superadmin'
		`, token.UserID).Scan(&adminID)

		if err != nil {
			sendError(w, "Доступ запрещен. Требуются права суперадмина", http.StatusForbidden)
			return
		}

		// Добавляем данные в контекст
		ctx := context.WithValue(r.Context(), "userID", token.UserID)
		ctx = context.WithValue(ctx, "adminID", adminID)
		ctx = context.WithValue(ctx, "roles", token.Roles)

		// Логируем запрос
		logAdminAction(adminID, "api_request", r.Method+" "+r.URL.Path, r.RemoteAddr)

		next(w, r.WithContext(ctx))
	}
}

// contains проверяет наличие строки в массиве
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ParseToken парсит JWT токен (с поддержкой roles для SSO)
func ParseToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	log.Printf("🔍 Parsing token with secret: %s... (len: %d)\n", secret[:10], len(secret))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		log.Printf("🔑 Using secret bytes (len: %d)\n", len([]byte(secret)))
		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("❌ ParseToken error: %v\n", err)
		log.Printf("   Secret used: %s... (full len: %d)\n", secret[:10], len(secret))
		return nil, err
	}

	log.Println("✅ Token parsed successfully")

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email")
	}

	// Парсим roles
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		for _, r := range rolesInterface {
			if roleStr, ok := r.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return &Claims{
		UserID: int(userID),
		Email:  email,
		Roles:  roles,
	}, nil
}

// GenerateToken больше не нужен - используем общий токен из Main Backend (SSO)

// logAdminAction логирует действие администратора
func logAdminAction(adminID int, action, details, ipAddress string) {
	// Извлекаем IP без порта
	ip := strings.Split(ipAddress, ":")[0]

	database.DB.Exec(`
		INSERT INTO admin_logs (admin_id, action, details, ip_address)
		VALUES (?, ?, ?, ?)
	`, adminID, action, details, ip)
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"success": false, "error": "%s"}`, message)
}
