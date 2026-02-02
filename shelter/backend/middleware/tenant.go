package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	pkgmiddleware "github.com/zooplatforma/pkg/middleware"
)

type tenantContextKey string

const TenantIDKey tenantContextKey = "tenantID"
const OrganizationKey tenantContextKey = "organization"

// TenantMiddleware определяет текущий приют (tenant) для пользователя
func TenantMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем userID из контекста (установлен централизованным AuthMiddleware)
			userID, ok := pkgmiddleware.GetUserID(r)
			if !ok {
				log.Printf("❌ User ID not found in context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Проверяем, передан ли organization_id в query параметрах (для переключения)
			orgIDStr := r.URL.Query().Get("organization_id")
			var organizationID int

			if orgIDStr != "" {
				// Пользователь хочет переключиться на другой приют
				var err error
				organizationID, err = strconv.Atoi(orgIDStr)
				if err != nil {
					http.Error(w, "Invalid organization_id", http.StatusBadRequest)
					return
				}
			} else {
				// Получаем первый приют пользователя (по умолчанию)
				err := db.QueryRow(`
					SELECT om.organization_id 
					FROM organization_members om
					JOIN organizations o ON o.id = om.organization_id
					WHERE om.user_id = ? AND o.type = 'shelter'
					ORDER BY om.joined_at ASC
					LIMIT 1
				`, userID).Scan(&organizationID)

				if err == sql.ErrNoRows {
					log.Printf("❌ User %d is not a member of any shelter", userID)
					http.Error(w, "You are not a member of any shelter", http.StatusForbidden)
					return
				} else if err != nil {
					log.Printf("❌ Error getting user's shelter: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}

			// Проверяем, что пользователь действительно член этого приюта
			var exists bool
			err := db.QueryRow(`
				SELECT EXISTS(
					SELECT 1 FROM organization_members om
					JOIN organizations o ON o.id = om.organization_id
					WHERE om.user_id = ? AND om.organization_id = ? AND o.type = 'shelter'
				)
			`, userID, organizationID).Scan(&exists)

			if err != nil || !exists {
				log.Printf("❌ User %d is not a member of shelter %d", userID, organizationID)
				http.Error(w, "Access denied to this shelter", http.StatusForbidden)
				return
			}

			// Получаем информацию о приюте
			var orgName, orgType string
			err = db.QueryRow(`
				SELECT name, type FROM organizations WHERE id = ?
			`, organizationID).Scan(&orgName, &orgType)

			if err != nil {
				log.Printf("❌ Error getting organization info: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			log.Printf("✅ Tenant: User %d accessing shelter %d (%s)", userID, organizationID, orgName)

			// Добавляем tenant ID в контекст
			ctx := context.WithValue(r.Context(), TenantIDKey, organizationID)
			ctx = context.WithValue(ctx, OrganizationKey, map[string]interface{}{
				"id":   organizationID,
				"name": orgName,
				"type": orgType,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetTenantID извлекает tenant ID из контекста
func GetTenantID(r *http.Request) (int, bool) {
	tenantID, ok := r.Context().Value(TenantIDKey).(int)
	return tenantID, ok
}

// GetOrganization извлекает информацию об организации из контекста
func GetOrganization(r *http.Request) (map[string]interface{}, bool) {
	org, ok := r.Context().Value(OrganizationKey).(map[string]interface{})
	return org, ok
}
