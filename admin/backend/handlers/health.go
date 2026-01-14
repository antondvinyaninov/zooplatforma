package handlers

import (
	"encoding/json"
	"net/http"
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

// HealthCheckHandler проверяет работоспособность всех сервисов
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	services := []struct {
		Name string
		URL  string
	}{
		{"Main Backend", "http://localhost:8000/api/health"},
		{"Main Frontend", "http://localhost:3000"},
		{"Admin Backend", "http://localhost:9000/api/health"},
		{"Admin Frontend", "http://localhost:4000"},
		{"PetBase Backend", "http://localhost:8100/api/health"},
		{"PetBase Frontend", "http://localhost:4100"},
		{"Shelter Backend", "http://localhost:8200/api/health"},
		{"Shelter Frontend", "http://localhost:5100"},
		{"Owner Backend", "http://localhost:8400/api/health"},
		{"Owner Frontend", "http://localhost:6100"},
		{"Volunteer Backend", "http://localhost:8500/api/health"},
		{"Volunteer Frontend", "http://localhost:6200"},
		{"Clinic Backend", "http://localhost:8600/api/health"},
		{"Clinic Frontend", "http://localhost:6300"},
	}

	var statuses []ServiceStatus
	onlineCount := 0
	slowCount := 0

	for _, service := range services {
		status := checkService(service.Name, service.URL)
		statuses = append(statuses, status)

		if status.Status == "online" {
			onlineCount++
		} else if status.Status == "slow" {
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

func checkService(name, url string) ServiceStatus {
	start := time.Now()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
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

	// Проверяем статус код и задержку
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if latency > 1000 {
			status.Status = "slow"
		} else {
			status.Status = "online"
		}
	} else {
		status.Status = "offline"
	}

	return status
}
