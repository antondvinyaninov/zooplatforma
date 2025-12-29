package middleware

import (
	"database"
	"net/http"
	"time"
)

// responseWriter - обёртка для захвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware - middleware для логирования всех HTTP запросов
func LoggingMiddleware(logger *database.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Создаём обёртку для захвата статус кода
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Выполняем handler
			next(wrapped, r)

			// Логируем запрос
			duration := time.Since(start)
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			logger.LogRequest(r.Method, r.URL.Path, ip, wrapped.statusCode, duration)
		}
	}
}
