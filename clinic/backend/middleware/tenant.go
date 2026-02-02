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

// TenantMiddleware проверяет наличие выбранной клиники в заголовке X-Clinic-ID
func TenantMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🏥 Tenant check for: %s %s", r.Method, r.URL.Path)

			// Получаем user_id из контекста (установлен централизованным AuthMiddleware)
			userID, ok := pkgmiddleware.GetUserID(r)
			if !ok {
				log.Printf("❌ User ID not found in context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Получаем clinic ID из заголовка
			clinicIDStr := r.Header.Get("X-Clinic-ID")
			if clinicIDStr == "" {
				log.Printf("❌ X-Clinic-ID header not found")
				http.Error(w, "Clinic not selected", http.StatusBadRequest)
				return
			}

			clinicID, err := strconv.Atoi(clinicIDStr)
			if err != nil {
				log.Printf("❌ Invalid clinic ID: %s", clinicIDStr)
				http.Error(w, "Invalid clinic ID", http.StatusBadRequest)
				return
			}

			// Проверяем, что пользователь является членом этой клиники
			var count int
			err = db.QueryRow(`
				SELECT COUNT(*) 
				FROM organization_members 
				WHERE organization_id = ? AND user_id = ?
			`, clinicID, userID).Scan(&count)

			if err != nil {
				log.Printf("❌ Error checking clinic membership: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if count == 0 {
				log.Printf("❌ User %d is not a member of clinic %d", userID, clinicID)
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			log.Printf("✅ Tenant verified: clinic_id=%d, user_id=%d", clinicID, userID)

			// Добавляем clinic ID в контекст
			ctx := context.WithValue(r.Context(), TenantIDKey, clinicID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetTenantID извлекает tenant ID из контекста
func GetTenantID(r *http.Request) (int, bool) {
	tenantID, ok := r.Context().Value(TenantIDKey).(int)
	return tenantID, ok
}
