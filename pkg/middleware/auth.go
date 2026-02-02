package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Используем строковые ключи для совместимости с существующими handlers
const (
	UserIDKey    = "userID"
	UserEmailKey = "userEmail"
	UserRoleKey  = "userRole"
)

var authServiceURL string

// InitAuthMiddleware инициализирует middleware с URL Auth Service
func InitAuthMiddleware(url string) {
	authServiceURL = url
}

// getAuthServiceURL возвращает URL Auth Service
func getAuthServiceURL() string {
	if authServiceURL != "" {
		return authServiceURL
	}
	// Fallback на переменную окружения
	url := os.Getenv("AUTH_SERVICE_URL")
	if url == "" {
		url = "http://localhost:7100"
	}
	return url
}

// User - структура пользователя из Auth Service
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// AuthMiddleware - универсальный middleware для проверки авторизации через Auth Service
// Работает для всех микросервисов
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Получить токен из Authorization header (приоритет)
		token := r.Header.Get("Authorization")

		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// 2. Если нет в header, попробовать cookie
		if token == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		// 3. Если токена нет - 401
		if token == "" {
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 4. Проверить токен через Auth Service (SSO)
		authURL := getAuthServiceURL()

		// Создаем запрос к Auth Service /api/auth/me
		req, err := http.NewRequest("GET", authURL+"/api/auth/me", nil)
		if err != nil {
			sendError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Передаем токен в заголовке
		req.Header.Set("Authorization", "Bearer "+token)

		// Отправляем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ AuthMiddleware: Auth Service unavailable: %v", err)
			sendError(w, "Auth service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Читаем ответ
		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			log.Printf("❌ AuthMiddleware: Auth Service returned %d", resp.StatusCode)
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Парсим ответ от Auth Service
		var authResp struct {
			Success bool `json:"success"`
			Data    struct {
				User *User `json:"user"`
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &authResp); err != nil {
			log.Printf("❌ AuthMiddleware: Failed to parse auth response: %v", err)
			sendError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !authResp.Success || authResp.Data.User == nil {
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 5. Добавить данные пользователя в контекст
		user := authResp.Data.User
		ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
		ctx = context.WithValue(ctx, UserEmailKey, user.Email)
		ctx = context.WithValue(ctx, UserRoleKey, user.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware - опциональная авторизация (не требует токен)
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пробуем получить токен
		token := r.Header.Get("Authorization")

		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		// Если токена нет - просто пропускаем дальше
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Если токен есть - проверяем через Auth Service (SSO)
		authURL := getAuthServiceURL()

		req, err := http.NewRequest("GET", authURL+"/api/auth/me", nil)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			next.ServeHTTP(w, r)
			return
		}

		var authResp struct {
			Success bool `json:"success"`
			Data    struct {
				User *User `json:"user"`
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &authResp); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if authResp.Success && authResp.Data.User != nil {
			user := authResp.Data.User
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			ctx = context.WithValue(ctx, UserEmailKey, user.Email)
			ctx = context.WithValue(ctx, UserRoleKey, user.Role)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID получает user_id из контекста
func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}

// GetUserEmail получает email из контекста
func GetUserEmail(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(UserEmailKey).(string)
	return email, ok
}

// GetUserRole получает роль из контекста
func GetUserRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(UserRoleKey).(string)
	return role, ok
}

// RequireSuperAdmin - middleware для проверки роли superadmin
func RequireSuperAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetUserRole(r)
		if !ok || role != "superadmin" {
			sendError(w, "Forbidden: superadmin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// VerifyTokenViaAuthService - проверяет токен через Auth Service (SSO) и возвращает пользователя
func VerifyTokenViaAuthService(token string) (*User, error) {
	authURL := getAuthServiceURL()

	req, err := http.NewRequest("GET", authURL+"/api/auth/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unauthorized: status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var authResp struct {
		Success bool `json:"success"`
		Data    struct {
			User *User `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !authResp.Success || authResp.Data.User == nil {
		return nil, fmt.Errorf("invalid token")
	}

	return authResp.Data.User, nil
}

// LoginViaAuthService - выполняет вход через Auth Service (SSO)
func LoginViaAuthService(email, password string) (string, *User, error) {
	authURL := getAuthServiceURL()

	reqBody := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(authURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("login failed: %s", string(body))
	}

	var authResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
			User  *User  `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !authResp.Success {
		return "", nil, fmt.Errorf("login failed")
	}

	return authResp.Data.Token, authResp.Data.User, nil
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   message,
	})
}
