package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ServiceStatus struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Status    string `json:"status"`  // "online", "offline", "slow"
	Latency   int64  `json:"latency"` // в миллисекундах
	CheckedAt string `json:"checked_at"`
}

type HealthResponse struct {
	Services []ServiceStatus `json:"services"`
	Overall  string          `json:"overall"` // "healthy", "degraded", "down"
}

// HealthCheckHandler проверяет работоспособность всех backend и frontend сервисов
// Backend проверяются через GET /api/health
// Frontend проверяются через GET /favicon.ico (легковесный запрос, не вызывает рендеринг)
// Примечание: Admin Backend не проверяет сам себя, чтобы избежать циклической зависимости
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем backend и frontend сервисы
	// Backend проверяем через /api/health
	// Frontend проверяем через /favicon.ico (легковесный запрос)
	// Примечание: Admin Backend не проверяет сам себя, чтобы избежать циклической зависимости
	services := []struct {
		Name      string
		URL       string
		IsBackend bool
	}{
		{"Auth Service", "http://localhost:7100/api/health", true},
		{"Main Backend", "http://localhost:8000/api/health", true},
		{"Main Frontend", "http://localhost:3000", false},
		{"Admin Frontend", "http://localhost:4000", false},
		{"PetBase Backend", "http://localhost:8100/api/health", true},
		{"PetBase Frontend", "http://localhost:4100", false},
		{"Shelter Backend", "http://localhost:8200/api/health", true},
		{"Shelter Frontend", "http://localhost:5100", false},
		{"Owner Backend", "http://localhost:8400/api/health", true},
		{"Owner Frontend", "http://localhost:6100", false},
		{"Volunteer Backend", "http://localhost:8500/api/health", true},
		{"Volunteer Frontend", "http://localhost:6200", false},
		{"Clinic Backend", "http://localhost:8600/api/health", true},
		{"Clinic Frontend", "http://localhost:6300", false},
	}

	// Параллельная проверка всех сервисов
	var wg sync.WaitGroup
	statusChan := make(chan ServiceStatus, len(services))

	for _, service := range services {
		wg.Add(1)
		go func(name, url string, isBackend bool) {
			defer wg.Done()
			statusChan <- checkService(name, url, isBackend)
		}(service.Name, service.URL, service.IsBackend)
	}

	// Ждем завершения всех проверок
	go func() {
		wg.Wait()
		close(statusChan)
	}()

	// Собираем результаты
	var statuses []ServiceStatus
	onlineCount := 0
	slowCount := 0

	for status := range statusChan {
		statuses = append(statuses, status)

		switch status.Status {
		case "online":
			onlineCount++
		case "slow":
			slowCount++
		}
	}

	// Определяем общий статус
	overall := "healthy"
	if onlineCount == 0 {
		overall = "down"
	} else if slowCount > 0 || onlineCount < len(services) {
		overall = "degraded"
	}

	response := HealthResponse{
		Services: statuses,
		Overall:  overall,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkService(name, url string, isBackend bool) ServiceStatus {
	start := time.Now()

	// Уменьшаем таймаут до 1 секунды для быстрой проверки
	client := &http.Client{
		Timeout: 1 * time.Second,
		// Не следуем редиректам для быстрой проверки
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var resp *http.Response
	var err error

	// Для frontend используем GET к favicon.ico (легковесный запрос)
	// Для backend используем GET к /api/health
	if isBackend {
		resp, err = client.Get(url)
	} else {
		// Для frontend просто проверяем, что сервер отвечает
		resp, err = client.Get(url + "/favicon.ico")
	}

	latency := time.Since(start).Milliseconds()

	status := ServiceStatus{
		Name:      name,
		URL:       url,
		Latency:   latency,
		CheckedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err != nil {
		status.Status = "offline"
		return status
	}
	defer resp.Body.Close()

	// Для frontend принимаем 200 и 404 (favicon может не существовать, но сервер работает)
	// Для backend только 200
	if isBackend {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if latency > 1000 {
				status.Status = "slow"
			} else {
				status.Status = "online"
			}
		} else {
			status.Status = "offline"
		}
	} else {
		// Frontend: если сервер ответил (даже 404), значит он работает
		if resp.StatusCode >= 200 && resp.StatusCode < 500 {
			if latency > 500 {
				status.Status = "slow"
			} else {
				status.Status = "online"
			}
		} else {
			status.Status = "offline"
		}
	}

	return status
}
