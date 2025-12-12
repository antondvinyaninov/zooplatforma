package handlers

import (
	"database"
	"net/http"
	"strconv"
	"strings"
)

// AdminUsersHandler - список пользователей
func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Добавить пагинацию и фильтры
	rows, err := database.DB.Query(`
		SELECT id, name, email, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT 100
	`)
	if err != nil {
		sendError(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email, createdAt string
		if err := rows.Scan(&id, &name, &email, &createdAt); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{
			"id":         id,
			"name":       name,
			"email":      email,
			"created_at": createdAt,
		})
	}

	sendSuccess(w, users)
}

// AdminUserHandler - действия с конкретным пользователем
func AdminUserHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		sendError(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		sendError(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getUserDetails(w, r, userID)
	case http.MethodPut:
		// Проверяем действие (block/unblock)
		if len(parts) > 1 {
			action := parts[1]
			switch action {
			case "block":
				blockUser(w, r, userID)
			case "unblock":
				unblockUser(w, r, userID)
			default:
				sendError(w, "Неизвестное действие", http.StatusBadRequest)
			}
		} else {
			sendError(w, "Действие не указано", http.StatusBadRequest)
		}
	case http.MethodDelete:
		deleteUser(w, r, userID)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUserDetails(w http.ResponseWriter, r *http.Request, userID int) {
	var name, email, bio, phone, location, createdAt string
	err := database.DB.QueryRow(`
		SELECT name, email, bio, phone, location, created_at
		FROM users WHERE id = ?
	`, userID).Scan(&name, &email, &bio, &phone, &location, &createdAt)

	if err != nil {
		sendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Получаем статистику
	var postsCount, petsCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&postsCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM pets WHERE user_id = ?", userID).Scan(&petsCount)

	sendSuccess(w, map[string]interface{}{
		"id":          userID,
		"name":        name,
		"email":       email,
		"bio":         bio,
		"phone":       phone,
		"location":    location,
		"created_at":  createdAt,
		"posts_count": postsCount,
		"pets_count":  petsCount,
	})
}

func blockUser(w http.ResponseWriter, r *http.Request, userID int) {
	// TODO: Добавить поле blocked в таблицу users
	adminID := r.Context().Value("adminID").(int)

	// Логируем действие
	database.DB.Exec(`
		INSERT INTO admin_logs (admin_id, action, target_type, target_id, details, ip_address)
		VALUES (?, 'user_block', 'user', ?, 'User blocked', ?)
	`, adminID, userID, r.RemoteAddr)

	sendSuccess(w, map[string]string{"message": "Пользователь заблокирован"})
}

func unblockUser(w http.ResponseWriter, r *http.Request, userID int) {
	// TODO: Добавить поле blocked в таблицу users
	adminID := r.Context().Value("adminID").(int)

	// Логируем действие
	database.DB.Exec(`
		INSERT INTO admin_logs (admin_id, action, target_type, target_id, details, ip_address)
		VALUES (?, 'user_unblock', 'user', ?, 'User unblocked', ?)
	`, adminID, userID, r.RemoteAddr)

	sendSuccess(w, map[string]string{"message": "Пользователь разблокирован"})
}

func deleteUser(w http.ResponseWriter, r *http.Request, userID int) {
	adminID := r.Context().Value("adminID").(int)

	// Soft delete - помечаем как удаленного
	// TODO: Добавить поле deleted_at в таблицу users

	// Логируем действие
	database.DB.Exec(`
		INSERT INTO admin_logs (admin_id, action, target_type, target_id, details, ip_address)
		VALUES (?, 'user_delete', 'user', ?, 'User deleted', ?)
	`, adminID, userID, r.RemoteAddr)

	sendSuccess(w, map[string]string{"message": "Пользователь удален"})
}
