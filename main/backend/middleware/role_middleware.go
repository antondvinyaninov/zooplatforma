package middleware

import (
	"backend/models"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

// RoleMiddleware проверяет, что у пользователя есть требуемая роль
func RoleMiddleware(db *sql.DB, requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Проверяем наличие роли
			if !hasRole(db, userID, requiredRole) {
				log.Printf("⚠️ Access denied: user %d doesn't have role '%s'", userID, requiredRole)
				http.Error(w, "Access denied: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AnyRoleMiddleware проверяет, что у пользователя есть хотя бы одна из требуемых ролей
func AnyRoleMiddleware(db *sql.DB, requiredRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Проверяем наличие хотя бы одной роли
			hasAnyRole := false
			for _, role := range requiredRoles {
				if hasRole(db, userID, role) {
					hasAnyRole = true
					break
				}
			}

			if !hasAnyRole {
				log.Printf("⚠️ Access denied: user %d doesn't have any of required roles %v", userID, requiredRoles)
				http.Error(w, "Access denied: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// EnrichWithRoles добавляет роли пользователя в контекст запроса
func EnrichWithRoles(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			// Получаем роли пользователя
			roles, err := getUserRoles(db, userID)
			if err != nil {
				log.Printf("⚠️ Error fetching user roles: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			// Добавляем роли в контекст
			ctx := context.WithValue(r.Context(), "userRoles", roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Вспомогательные функции

func hasRole(db *sql.DB, userID int, role string) bool {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM user_roles 
		WHERE user_id = ? AND role = ? AND is_active = 1
		AND (expires_at IS NULL OR expires_at > ?)
	`, userID, role, time.Now()).Scan(&count)

	return err == nil && count > 0
}

func getUserRoles(db *sql.DB, userID int) ([]string, error) {
	rows, err := db.Query(`
		SELECT role FROM user_roles 
		WHERE user_id = ? AND is_active = 1
		AND (expires_at IS NULL OR expires_at > ?)
		ORDER BY granted_at DESC
	`, userID, time.Now())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// HasPermission проверяет, есть ли у пользователя определенное право
func HasPermission(db *sql.DB, userID int, permission string) bool {
	roles, err := getUserRoles(db, userID)
	if err != nil {
		return false
	}

	for _, role := range roles {
		if models.HasPermission(role, permission) {
			return true
		}
	}

	return false
}
