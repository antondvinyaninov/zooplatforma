package main

import (
	"os"
)

// Service представляет backend сервис
type Service struct {
	Name    string
	URL     string
	Timeout int // секунды
}

// Services содержит все backend сервисы
type Services struct {
	Auth      *Service
	Main      *Service
	PetBase   *Service
	Clinic    *Service
	Owner     *Service
	Shelter   *Service
	Volunteer *Service
	Admin     *Service
}

// InitServices инициализирует все сервисы
func InitServices() *Services {
	return &Services{
		Auth: &Service{
			Name:    "Auth Service",
			URL:     getEnv("AUTH_SERVICE_URL", "http://localhost:7100"),
			Timeout: 5,
		},
		Main: &Service{
			Name:    "Main Backend",
			URL:     getEnv("MAIN_SERVICE_URL", "http://localhost:8000"),
			Timeout: 10,
		},
		PetBase: &Service{
			Name:    "PetBase Backend",
			URL:     getEnv("PETBASE_SERVICE_URL", "http://localhost:8100"),
			Timeout: 5,
		},
		Clinic: &Service{
			Name:    "Clinic Backend",
			URL:     getEnv("CLINIC_SERVICE_URL", "http://localhost:8600"),
			Timeout: 10,
		},
		Owner: &Service{
			Name:    "Owner Backend",
			URL:     getEnv("OWNER_SERVICE_URL", "http://localhost:8400"),
			Timeout: 10,
		},
		Shelter: &Service{
			Name:    "Shelter Backend",
			URL:     getEnv("SHELTER_SERVICE_URL", "http://localhost:8200"),
			Timeout: 10,
		},
		Volunteer: &Service{
			Name:    "Volunteer Backend",
			URL:     getEnv("VOLUNTEER_SERVICE_URL", "http://localhost:8500"),
			Timeout: 10,
		},
		Admin: &Service{
			Name:    "Admin Backend",
			URL:     getEnv("ADMIN_SERVICE_URL", "http://localhost:9000"),
			Timeout: 10,
		},
	}
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
