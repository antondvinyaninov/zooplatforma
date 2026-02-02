package clients

import (
	"fmt"
)

// AuthClient - клиент для Auth Service
type AuthClient struct {
	*BaseClient
}

// NewAuthClient создает новый клиент для Auth Service
func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		BaseClient: NewBaseClient(baseURL, 0),
	}
}

// User - структура пользователя
type User struct {
	ID            int     `json:"id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	LastName      string  `json:"last_name"`
	Avatar        string  `json:"avatar"`
	Bio           string  `json:"bio"`
	Phone         string  `json:"phone"`
	DateOfBirth   *string `json:"date_of_birth,omitempty"`
	Gender        string  `json:"gender"`
	City          string  `json:"city"`
	Country       string  `json:"country"`
	Role          string  `json:"role"`
	EmailVerified bool    `json:"email_verified"`
	CreatedAt     string  `json:"created_at"`
}

// VerifyTokenResponse - ответ проверки токена
type VerifyTokenResponse struct {
	Valid bool  `json:"valid"`
	User  *User `json:"user,omitempty"`
}

// VerifyToken проверяет JWT токен
func (c *AuthClient) VerifyToken(token string) (*VerifyTokenResponse, error) {
	resp, err := c.Get("/api/auth/verify", token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if !resp.Success {
		return &VerifyTokenResponse{Valid: false}, nil
	}

	// Парсим данные пользователя
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	userMap, ok := dataMap["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data format")
	}

	user := parseUser(userMap)

	return &VerifyTokenResponse{
		Valid: true,
		User:  user,
	}, nil
}

// GetMe получает информацию о текущем пользователе
func (c *AuthClient) GetMe(token string) (*User, error) {
	resp, err := c.Get("/api/auth/me", token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("auth failed: %s", resp.Error)
	}

	// Парсим данные
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	userMap, ok := dataMap["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data format")
	}

	user := parseUser(userMap)
	return user, nil
}

// LoginRequest - запрос на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse - ответ на вход
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// Login выполняет вход в систему
func (c *AuthClient) Login(email, password string) (*LoginResponse, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.Post("/api/auth/login", req, "")
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("login failed: %s", resp.Error)
	}

	// Парсим данные
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	token, _ := dataMap["token"].(string)
	userMap, _ := dataMap["user"].(map[string]interface{})

	user := parseUser(userMap)

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Helper функция для безопасного получения строки
func getStringOrEmpty(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetUserProfile получает профиль пользователя по ID
func (c *AuthClient) GetUserProfile(userID int, token string) (*User, error) {
	resp, err := c.Get(fmt.Sprintf("/api/users/%d", userID), token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get profile: %s", resp.Error)
	}

	// Парсим данные
	userMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	user := parseUser(userMap)
	return user, nil
}

// UpdateProfileRequest - запрос на обновление профиля
type UpdateProfileRequest struct {
	Name        *string `json:"name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     *string `json:"country,omitempty"`
}

// UpdateUserProfile обновляет профиль пользователя
func (c *AuthClient) UpdateUserProfile(userID int, req UpdateProfileRequest, token string) (*User, error) {
	resp, err := c.Put(fmt.Sprintf("/api/users/%d", userID), req, token)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to update profile: %s", resp.Error)
	}

	// Парсим данные
	userMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	user := parseUser(userMap)
	return user, nil
}

// RegisterRequest - запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	LastName string `json:"last_name,omitempty"`
}

// Register регистрирует нового пользователя
func (c *AuthClient) Register(req RegisterRequest) (*LoginResponse, error) {
	resp, err := c.Post("/api/auth/register", req, "")
	if err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("registration failed: %s", resp.Error)
	}

	// Парсим данные
	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	token, _ := dataMap["token"].(string)
	userMap, _ := dataMap["user"].(map[string]interface{})

	user := parseUser(userMap)

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// parseUser парсит map в структуру User
func parseUser(userMap map[string]interface{}) *User {
	user := &User{
		ID:            int(userMap["id"].(float64)),
		Email:         getStringOrEmpty(userMap, "email"),
		Name:          getStringOrEmpty(userMap, "name"),
		LastName:      getStringOrEmpty(userMap, "last_name"),
		Avatar:        getStringOrEmpty(userMap, "avatar"),
		Bio:           getStringOrEmpty(userMap, "bio"),
		Phone:         getStringOrEmpty(userMap, "phone"),
		Gender:        getStringOrEmpty(userMap, "gender"),
		City:          getStringOrEmpty(userMap, "city"),
		Country:       getStringOrEmpty(userMap, "country"),
		Role:          getStringOrEmpty(userMap, "role"),
		EmailVerified: getBoolOrFalse(userMap, "email_verified"),
		CreatedAt:     getStringOrEmpty(userMap, "created_at"),
	}

	// DateOfBirth может быть nil
	if dob := getStringOrEmpty(userMap, "date_of_birth"); dob != "" {
		user.DateOfBirth = &dob
	}

	return user
}
