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
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := []string{
			"http://localhost:3000",                                  // Main Frontend (dev)
			"http://localhost:4000",                                  // Admin Frontend (dev)
			"http://localhost:4100",                                  // PetBase Frontend (dev)
			"http://localhost:5100",                                  // Shelter Frontend (dev)
			"http://localhost:6100",                                  // Owner Frontend (dev)
			"http://localhost:6200",                                  // Volunteer Frontend (dev)
			"http://localhost:6300",                                  // Clinic Frontend (dev)
			"https://my-projects-zooplatforma.crv1ic.easypanel.host", // Main Frontend (prod)
			"https://my-projects-admin.crv1ic.easypanel.host",        // Admin Frontend (prod)
			"https://my-projects-petbase.crv1ic.easypanel.host",      // PetBase Frontend (prod)
			"https://my-projects-shelter.crv1ic.easypanel.host",      // Shelter Frontend (prod)
			"https://my-projects-owner.crv1ic.easypanel.host",        // Owner Frontend (prod)
			"https://my-projects-volunteer.crv1ic.easypanel.host",    // Volunteer Frontend (prod)
			"https://my-projects-clinic.crv1ic.easypanel.host",       // Clinic Frontend (prod)
		}

		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Email, X-User-Role")

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
