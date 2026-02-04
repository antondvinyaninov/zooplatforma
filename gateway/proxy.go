package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ProxyHandler создает handler для проксирования запросов к сервису
func ProxyHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Парсим URL сервиса
		targetURL, err := url.Parse(service.URL)
		if err != nil {
			log.Printf("❌ Invalid service URL: %v", err)
			sendError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Создаем новый запрос к backend сервису
		proxyReq, err := http.NewRequest(r.Method, targetURL.String()+r.URL.Path, r.Body)
		if err != nil {
			log.Printf("❌ Failed to create proxy request: %v", err)
			sendError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Копируем query parameters
		proxyReq.URL.RawQuery = r.URL.RawQuery

		// Копируем заголовки
		for key, values := range r.Header {
			// Пропускаем некоторые заголовки
			if key == "Connection" || key == "Keep-Alive" || key == "Proxy-Authenticate" || key == "Proxy-Authorization" || key == "Te" || key == "Trailers" || key == "Transfer-Encoding" || key == "Upgrade" {
				continue
			}
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Добавляем X-Forwarded-* заголовки
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
		proxyReq.Header.Set("X-Forwarded-Proto", getScheme(r))
		proxyReq.Header.Set("X-Forwarded-Host", r.Host)

		// Отправляем запрос к backend сервису
		client := &http.Client{
			Timeout: time.Duration(service.Timeout) * time.Second,
		}

		start := time.Now()
		resp, err := client.Do(proxyReq)
		duration := time.Since(start)

		if err != nil {
			log.Printf("❌ Proxy error to %s: %v (took %v)", service.Name, err, duration)
			sendError(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Копируем заголовки ответа от backend
		// ВАЖНО: Пропускаем CORS заголовки (Access-Control-*) потому что:
		// 1. CORS управляется только Gateway (CORSMiddleware)
		// 2. Backend сервисы НЕ должны устанавливать CORS
		// 3. Это предотвращает дублирование заголовков
		for key, values := range resp.Header {
			// Пропускаем CORS заголовки - они уже установлены в CORSMiddleware
			if strings.HasPrefix(key, "Access-Control-") {
				continue
			}
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Копируем status code
		w.WriteHeader(resp.StatusCode)

		// Копируем body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("⚠️ Error copying response body: %v", err)
		}

		// Логируем успешный проксированный запрос
		log.Printf("✅ Proxied to %s: %s %s → %d (took %v)",
			service.Name,
			r.Method,
			r.URL.Path,
			resp.StatusCode,
			duration,
		)
	}
}

// getScheme определяет схему (http/https)
func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}

// HealthCheckHandler проверяет здоровье всех сервисов
func HealthCheckHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allServices := []*Service{
			services.Main,
			services.PetBase,
			services.Clinic,
			services.Owner,
			services.Shelter,
			services.Volunteer,
			services.Admin,
		}

		results := make(map[string]interface{})
		allHealthy := true

		for _, service := range allServices {
			healthy := checkServiceHealth(service)
			results[strings.ToLower(strings.ReplaceAll(service.Name, " ", "_"))] = map[string]interface{}{
				"url":     service.URL,
				"healthy": healthy,
			}
			if !healthy {
				allHealthy = false
			}
		}

		status := "healthy"
		statusCode := http.StatusOK
		if !allHealthy {
			status = "degraded"
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		response := map[string]interface{}{
			"success":  allHealthy,
			"status":   status,
			"gateway":  "API Gateway",
			"version":  "1.0.0",
			"services": results,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("❌ Failed to encode health check response: %v", err)
		}
	}
}

// checkServiceHealth проверяет здоровье одного сервиса
func checkServiceHealth(service *Service) bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(service.URL + "/api/health")
	if err != nil {
		log.Printf("⚠️ Health check failed for %s: %v", service.Name, err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
