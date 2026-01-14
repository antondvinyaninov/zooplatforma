package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const (
	BaseURL = "http://localhost:8000"
)

// TestUser для тестов
var testUser = struct {
	Email    string
	Password string
	Token    string
}{
	Email:    "test@example.com",
	Password: "testpassword123",
}

// Response структура для API ответов
type Response struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// TestHealthCheck проверяет работоспособность API
func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(BaseURL + "/api/health")
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Ожидался status=ok, получен %s", result["status"])
	}

	t.Log("✅ Health check пройден")
}

// TestAuthRegister проверяет регистрацию пользователя
func TestAuthRegister(t *testing.T) {
	// Генерируем уникальный email для каждого теста
	testEmail := "test_" + time.Now().Format("20060102150405") + "@example.com"

	payload := map[string]string{
		"name":     "Test User",
		"email":    testEmail,
		"password": testUser.Password,
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(
		BaseURL+"/api/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус 200/201, получен %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if !result.Success {
		t.Errorf("Регистрация не удалась: %s", result.Error)
	}

	t.Log("✅ Регистрация пройдена")
}

// TestAuthLogin проверяет вход в систему
func TestAuthLogin(t *testing.T) {
	// Сначала регистрируем пользователя
	testEmail := "login_test_" + time.Now().Format("20060102150405") + "@example.com"

	registerPayload := map[string]string{
		"name":     "Login Test User",
		"email":    testEmail,
		"password": testUser.Password,
	}
	body, _ := json.Marshal(registerPayload)
	http.Post(BaseURL+"/api/auth/register", "application/json", bytes.NewBuffer(body))

	// Теперь пытаемся войти
	loginPayload := map[string]string{
		"email":    testEmail,
		"password": testUser.Password,
	}

	body, _ = json.Marshal(loginPayload)
	resp, err := http.Post(
		BaseURL+"/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if !result.Success {
		t.Errorf("Вход не удался: %s", result.Error)
	}

	// Проверяем наличие cookie с токеном
	cookies := resp.Cookies()
	hasAuthToken := false
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			hasAuthToken = true
			testUser.Token = cookie.Value
			break
		}
	}

	if !hasAuthToken {
		t.Error("Cookie auth_token не установлен")
	}

	t.Log("✅ Вход пройден, токен получен")
}

// TestGetPosts проверяет получение ленты постов
func TestGetPosts(t *testing.T) {
	resp, err := http.Get(BaseURL + "/api/posts")
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if !result.Success {
		t.Errorf("Получение постов не удалось: %s", result.Error)
	}

	t.Log("✅ Получение постов пройдено")
}

// TestCreatePost проверяет создание поста (требует авторизации)
func TestCreatePost(t *testing.T) {
	// Сначала логинимся
	testEmail := "post_test_" + time.Now().Format("20060102150405") + "@example.com"

	// Регистрация
	registerPayload := map[string]string{
		"name":     "Post Test User",
		"email":    testEmail,
		"password": testUser.Password,
	}
	body, _ := json.Marshal(registerPayload)
	http.Post(BaseURL+"/api/auth/register", "application/json", bytes.NewBuffer(body))

	// Логин
	loginPayload := map[string]string{
		"email":    testEmail,
		"password": testUser.Password,
	}
	body, _ = json.Marshal(loginPayload)
	loginResp, _ := http.Post(BaseURL+"/api/auth/login", "application/json", bytes.NewBuffer(body))
	defer loginResp.Body.Close()

	// Получаем cookie
	var authCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}

	if authCookie == nil {
		t.Fatal("Не удалось получить auth_token")
	}

	// Создаём пост
	postPayload := map[string]interface{}{
		"content":       "Тестовый пост из автотеста",
		"author_type":   "user",
		"attached_pets": []int{},
		"tags":          []string{"тест"},
	}

	body, _ = json.Marshal(postPayload)
	req, _ := http.NewRequest("POST", BaseURL+"/api/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(authCookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус 200/201, получен %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if !result.Success {
		t.Errorf("Создание поста не удалось: %s", result.Error)
	}

	t.Log("✅ Создание поста пройдено")
}

// TestUnauthorizedAccess проверяет защиту endpoints
func TestUnauthorizedAccess(t *testing.T) {
	// Пытаемся создать пост без авторизации
	postPayload := map[string]interface{}{
		"content":     "Тест без авторизации",
		"author_type": "user",
	}

	body, _ := json.Marshal(postPayload)
	resp, err := http.Post(
		BaseURL+"/api/posts",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Ожидался статус 401, получен %d", resp.StatusCode)
	}

	t.Log("✅ Защита от неавторизованного доступа работает")
}

// TestReportCreation проверяет создание жалобы
func TestReportCreation(t *testing.T) {
	// Логинимся
	testEmail := "report_test_" + time.Now().Format("20060102150405") + "@example.com"

	registerPayload := map[string]string{
		"name":     "Report Test User",
		"email":    testEmail,
		"password": testUser.Password,
	}
	body, _ := json.Marshal(registerPayload)
	http.Post(BaseURL+"/api/auth/register", "application/json", bytes.NewBuffer(body))

	loginPayload := map[string]string{
		"email":    testEmail,
		"password": testUser.Password,
	}
	body, _ = json.Marshal(loginPayload)
	loginResp, _ := http.Post(BaseURL+"/api/auth/login", "application/json", bytes.NewBuffer(body))
	defer loginResp.Body.Close()

	var authCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}

	if authCookie == nil {
		t.Fatal("Не удалось получить auth_token")
	}

	// Создаём жалобу
	reportPayload := map[string]interface{}{
		"target_type": "post",
		"target_id":   1,
		"reason":      "spam",
		"description": "Тестовая жалоба из автотеста",
	}

	body, _ = json.Marshal(reportPayload)
	req, _ := http.NewRequest("POST", BaseURL+"/api/reports", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(authCookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус 200/201, получен %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if !result.Success {
		t.Errorf("Создание жалобы не удалось: %s", result.Error)
	}

	t.Log("✅ Создание жалобы пройдено")
}
