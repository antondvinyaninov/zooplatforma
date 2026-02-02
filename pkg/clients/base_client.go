package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BaseClient - базовый HTTP клиент для всех сервисов
type BaseClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

// APIResponse - стандартный формат ответа API
type APIResponse struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewBaseClient создает новый базовый клиент
func NewBaseClient(baseURL string, timeout time.Duration) *BaseClient {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &BaseClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		Timeout: timeout,
	}
}

// Get выполняет GET запрос
func (c *BaseClient) Get(endpoint string, token string) (*APIResponse, error) {
	return c.request("GET", endpoint, nil, token)
}

// Post выполняет POST запрос
func (c *BaseClient) Post(endpoint string, body interface{}, token string) (*APIResponse, error) {
	return c.request("POST", endpoint, body, token)
}

// Put выполняет PUT запрос
func (c *BaseClient) Put(endpoint string, body interface{}, token string) (*APIResponse, error) {
	return c.request("PUT", endpoint, body, token)
}

// Delete выполняет DELETE запрос
func (c *BaseClient) Delete(endpoint string, token string) (*APIResponse, error) {
	return c.request("DELETE", endpoint, nil, token)
}

// request выполняет HTTP запрос
func (c *BaseClient) request(method, endpoint string, body interface{}, token string) (*APIResponse, error) {
	url := c.BaseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Заголовки
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Выполнить запрос
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Прочитать ответ
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Распарсить JSON
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Проверить статус
	if resp.StatusCode >= 400 {
		return &apiResp, fmt.Errorf("API error (status %d): %s", resp.StatusCode, apiResp.Error)
	}

	return &apiResp, nil
}

// GetWithRetry выполняет GET запрос с повторными попытками
func (c *BaseClient) GetWithRetry(endpoint string, token string, maxRetries int) (*APIResponse, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		resp, err := c.Get(endpoint, token)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
