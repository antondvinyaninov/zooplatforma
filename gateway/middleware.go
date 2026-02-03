package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// LoggingMiddleware логирует все запросы
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем ResponseWriter который записывает status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Логируем запрос
		log.Printf("📋 %s %s %d %v %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// responseWriter оборачивает http.ResponseWriter для записи status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORSMiddleware добавляет CORS заголовки
// ВАЖНО: Это единственное место где устанавливаются CORS заголовки!
// Backend сервисы НЕ должны устанавливать CORS - они фильтруются в ProxyHandler
func CORSMiddleware(next http.Handler) http.Handler {
	// Список разрешенных origins (все frontend приложения)
	allowedOrigins := map[string]bool{
		"http://localhost:3000":                                  true, // Main Frontend (dev)
		"http://localhost:4000":                                  true, // Admin Frontend (dev)
		"http://localhost:4100":                                  true, // PetBase Frontend (dev)
		"http://localhost:5100":                                  true, // Shelter Frontend (dev)
		"http://localhost:6100":                                  true, // Owner Frontend (dev)
		"http://localhost:6200":                                  true, // Volunteer Frontend (dev)
		"http://localhost:6300":                                  true, // Clinic Frontend (dev)
		"http://localhost:8000":                                  true, // Main Backend (dev) - если нужен
		"https://my-projects-zooplatforma.crv1ic.easypanel.host": true, // Main Frontend (prod)
		"https://my-projects-admin.crv1ic.easypanel.host":        true, // Admin Frontend (prod)
		"https://my-projects-petbase.crv1ic.easypanel.host":      true, // PetBase Frontend (prod)
		"https://my-projects-shelter.crv1ic.easypanel.host":      true, // Shelter Frontend (prod)
		"https://my-projects-owner.crv1ic.easypanel.host":        true, // Owner Frontend (prod)
		"https://my-projects-volunteer.crv1ic.easypanel.host":    true, // Volunteer Frontend (prod)
		"https://my-projects-clinic.crv1ic.easypanel.host":       true, // Clinic Frontend (prod)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Логируем для отладки
		log.Printf("🌐 CORS: Origin=%s, Method=%s, Path=%s", origin, r.Method, r.URL.Path)

		// Проверяем что origin разрешен
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role")
			w.Header().Set("Access-Control-Max-Age", "3600") // Кеш preflight на 1 час

			log.Printf("✅ CORS: Allowed origin %s", origin)
		} else if origin != "" {
			// Origin не разрешен
			log.Printf("⚠️ CORS: Blocked origin %s", origin)
		}

		// Обработать preflight запрос
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimiter хранит rate limiters для каждого IP
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

var rateLimiter = &RateLimiter{
	limiters: make(map[string]*rate.Limiter),
	rate:     rate.Limit(100), // 100 запросов в секунду
	burst:    200,             // burst до 200 запросов
}

// getLimiter получает или создает rate limiter для IP
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware ограничивает количество запросов с одного IP
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получить IP адрес
		ip := r.RemoteAddr

		// Получить rate limiter для этого IP
		limiter := rateLimiter.getLimiter(ip)

		// Проверить лимит
		if !limiter.Allow() {
			log.Printf("⚠️ Rate limit exceeded: %s %s from %s", r.Method, r.URL.Path, ip)
			sendError(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CleanupRateLimiters периодически очищает неиспользуемые rate limiters
func CleanupRateLimiters() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rateLimiter.mu.Lock()
		// Очистить все limiters (простая реализация)
		// В production можно отслеживать последнее использование
		rateLimiter.limiters = make(map[string]*rate.Limiter)
		rateLimiter.mu.Unlock()
		log.Println("🧹 Rate limiters cleaned up")
	}
}
