package handlers

import (
	"admin/middleware"
	"database"
	"encoding/json"
	"net/http"
)

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// AdminLoginHandler - редирект на основной сайт для входа (SSO)
// Этот endpoint больше не нужен, так как используется общая авторизация
// Оставляем для обратной совместимости, но возвращаем сообщение о редиректе
func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sendError(w, "Используйте основной сайт для входа (SSO)", http.StatusBadRequest)
}

// AdminLogoutHandler - выход из всех сервисов (SSO)
func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Удаляем общий cookie (выход из всех сервисов)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   "", // Пустой для localhost, ".zooplatforma.ru" для production
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	sendSuccess(w, map[string]string{"message": "Logged out successfully"})
}

// AdminMeHandler - получить текущего админа (SSO)
func AdminMeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем общий токен (SSO)
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		sendError(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Парсим токен
	token, err := middleware.ParseToken(cookie.Value)
	if err != nil {
		sendError(w, "Неверный токен: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Проверяем, что у пользователя есть роль superadmin
	hasSuperadmin := false
	for _, role := range token.Roles {
		if role == "superadmin" {
			hasSuperadmin = true
			break
		}
	}

	if !hasSuperadmin {
		sendError(w, "Доступ запрещен. Требуются права суперадмина", http.StatusForbidden)
		return
	}

	// Получаем данные пользователя
	var name, email string
	err = database.DB.QueryRow(`
		SELECT name, email FROM users WHERE id = ?
	`, token.UserID).Scan(&name, &email)

	if err != nil {
		sendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Получаем роль из таблицы admins
	var adminRole string
	database.DB.QueryRow(`
		SELECT role FROM admins WHERE user_id = ?
	`, token.UserID).Scan(&adminRole)

	sendSuccess(w, AdminResponse{
		ID:    token.UserID,
		Name:  name,
		Email: email,
		Role:  adminRole,
	})
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
