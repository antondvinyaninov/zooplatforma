package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// ErrorLog структура для логирования ошибок
type ErrorLog struct {
	ID        int       `json:"id"`
	Service   string    `json:"service"`
	Endpoint  string    `json:"endpoint"`
	Method    string    `json:"method"`
	ErrorMsg  string    `json:"error_message"`
	UserID    *int      `json:"user_id,omitempty"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// SystemMetrics метрики системы
type SystemMetrics struct {
	TotalRequests    int     `json:"total_requests"`
	TotalErrors      int     `json:"total_errors"`
	ErrorRate        float64 `json:"error_rate"`
	AvgResponseTime  int64   `json:"avg_response_time_ms"`
	ActiveUsers      int     `json:"active_users"`
	DatabaseSize     int64   `json:"database_size_mb"`
	LastHourErrors   int     `json:"last_hour_errors"`
	Last24HourErrors int     `json:"last_24hour_errors"`
}

// GetRecentErrorsHandler возвращает последние ошибки
func GetRecentErrorsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем последние 50 ошибок
		query := `
			SELECT id, service, endpoint, method, error_message, 
			       user_id, ip_address, user_agent, created_at
			FROM error_logs
			ORDER BY created_at DESC
			LIMIT 50
		`

		rows, err := db.Query(query)
		if err != nil {
			log.Printf("❌ Error fetching error logs: %v", err)
			http.Error(w, "Failed to fetch error logs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var errors []ErrorLog
		for rows.Next() {
			var errLog ErrorLog
			err := rows.Scan(
				&errLog.ID, &errLog.Service, &errLog.Endpoint, &errLog.Method,
				&errLog.ErrorMsg, &errLog.UserID, &errLog.IPAddress,
				&errLog.UserAgent, &errLog.CreatedAt,
			)
			if err != nil {
				log.Printf("⚠️ Error scanning error log: %v", err)
				continue
			}
			errors = append(errors, errLog)
		}

		if errors == nil {
			errors = []ErrorLog{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    errors,
		})
	}
}

// GetSystemMetricsHandler возвращает метрики системы
func GetSystemMetricsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		metrics := SystemMetrics{}

		// Общее количество ошибок
		db.QueryRow("SELECT COUNT(*) FROM error_logs").Scan(&metrics.TotalErrors)

		// Ошибки за последний час
		db.QueryRow(`
			SELECT COUNT(*) FROM error_logs 
			WHERE created_at > datetime('now', '-1 hour')
		`).Scan(&metrics.LastHourErrors)

		// Ошибки за последние 24 часа
		db.QueryRow(`
			SELECT COUNT(*) FROM error_logs 
			WHERE created_at > datetime('now', '-24 hours')
		`).Scan(&metrics.Last24HourErrors)

		// Активные пользователи (онлайн в последние 5 минут)
		db.QueryRow(`
			SELECT COUNT(*) FROM user_activity 
			WHERE last_seen > datetime('now', '-5 minutes')
		`).Scan(&metrics.ActiveUsers)

		// Размер базы данных (примерно)
		var pageCount, pageSize int64
		db.QueryRow("PRAGMA page_count").Scan(&pageCount)
		db.QueryRow("PRAGMA page_size").Scan(&pageSize)
		metrics.DatabaseSize = (pageCount * pageSize) / (1024 * 1024) // MB

		// Процент ошибок (примерно, если есть таблица request_logs)
		var totalRequests int
		err := db.QueryRow("SELECT COUNT(*) FROM request_logs").Scan(&totalRequests)
		if err == nil && totalRequests > 0 {
			metrics.TotalRequests = totalRequests
			metrics.ErrorRate = float64(metrics.TotalErrors) / float64(totalRequests) * 100
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    metrics,
		})
	}
}

// GetErrorStatsByServiceHandler возвращает статистику ошибок по сервисам
func GetErrorStatsByServiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := `
			SELECT service, COUNT(*) as error_count
			FROM error_logs
			WHERE created_at > datetime('now', '-24 hours')
			GROUP BY service
			ORDER BY error_count DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			log.Printf("❌ Error fetching error stats: %v", err)
			http.Error(w, "Failed to fetch error stats", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		stats := make(map[string]int)
		for rows.Next() {
			var service string
			var count int
			if err := rows.Scan(&service, &count); err != nil {
				continue
			}
			stats[service] = count
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    stats,
		})
	}
}

// LogErrorMiddleware middleware для логирования ошибок
func LogErrorMiddleware(db *sql.DB, serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Создаем ResponseWriter wrapper для перехвата статус кода
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Засекаем время
			start := time.Now()

			// Выполняем запрос
			next.ServeHTTP(rw, r)

			// Логируем если ошибка (статус >= 400)
			if rw.statusCode >= 400 {
				userID := getUserIDFromContext(r)

				go func() {
					_, err := db.Exec(`
						INSERT INTO error_logs (
							service, endpoint, method, error_message,
							user_id, ip_address, user_agent, created_at
						) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
					`,
						serviceName,
						r.URL.Path,
						r.Method,
						http.StatusText(rw.statusCode),
						userID,
						r.RemoteAddr,
						r.Header.Get("User-Agent"),
						time.Now(),
					)
					if err != nil {
						log.Printf("⚠️ Failed to log error: %v", err)
					}
				}()

				log.Printf("❌ [%s] %s %s -> %d (took %v)",
					serviceName, r.Method, r.URL.Path, rw.statusCode, time.Since(start))
			}
		})
	}
}

// responseWriter wrapper для перехвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func getUserIDFromContext(r *http.Request) *int {
	if userID, ok := r.Context().Value("userID").(int); ok {
		return &userID
	}
	return nil
}
