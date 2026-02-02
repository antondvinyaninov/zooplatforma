package handlers

import (
	"database"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

// LikesHandler обрабатывает лайки поста
func LikesHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID поста из URL: /api/posts/{id}/like
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	path = strings.TrimSuffix(path, "/like")
	postID, err := strconv.Atoi(path)
	if err != nil {
		sendErrorResponse(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	// Для GET запросов userID опционален (OptionalAuthMiddleware)
	// Для POST запросов userID обязателен (AuthMiddleware)
	userID, _ := r.Context().Value("userID").(int)

	switch r.Method {
	case http.MethodPost:
		// POST требует авторизации
		if userID == 0 {
			sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
			return
		}
		toggleLike(w, r, postID, userID)
	case http.MethodGet:
		// GET работает и без авторизации (userID может быть 0)
		getLikeStatus(w, r, postID, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// toggleLike добавляет или удаляет лайк
func toggleLike(w http.ResponseWriter, r *http.Request, postID int, userID int) {
	// Проверяем, есть ли уже лайк
	var exists bool
	err := database.DB.QueryRow(ConvertPlaceholders("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)"), userID, postID).Scan(&exists)
	if err != nil {
		sendErrorResponse(w, "Ошибка проверки лайка: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ipAddress := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")

	if exists {
		// Удаляем лайк
		_, err = database.DB.Exec(ConvertPlaceholders("DELETE FROM likes WHERE user_id = ? AND post_id = ?"), userID, postID)
		if err != nil {
			sendErrorResponse(w, "Ошибка удаления лайка: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Логируем удаление лайка
		CreateUserLog(database.DB, userID, "like_remove", "Удалён лайк с поста", ipAddress, userAgent)
	} else {
		// Добавляем лайк
		_, err = database.DB.Exec(ConvertPlaceholders("INSERT INTO likes (user_id, post_id) VALUES (?, ?)"), userID, postID)
		if err != nil {
			sendErrorResponse(w, "Ошибка добавления лайка: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Логируем добавление лайка
		CreateUserLog(database.DB, userID, "like_add", "Добавлен лайк на пост", ipAddress, userAgent)

		// Создаем уведомление для автора поста
		var postAuthorID int
		var likerName, likerLastName sql.NullString
		err = database.DB.QueryRow(`
			SELECT p.author_id, u.name, u.last_name 
			FROM posts p 
			JOIN users u ON u.id = ? 
			WHERE p.id = ?
		`, userID, postID).Scan(&postAuthorID, &likerName, &likerLastName)

		if err == nil && postAuthorID != userID {
			// Формируем имя лайкнувшего
			fullName := likerName.String
			if likerLastName.Valid && likerLastName.String != "" {
				fullName += " " + likerLastName.String
			}

			// Создаем уведомление
			notifHandler := &NotificationsHandler{DB: database.DB}
			notifHandler.NotifyLike(postAuthorID, userID, postID, fullName)
		}
	}

	// Получаем обновленное количество лайков
	var likesCount int
	err = database.DB.QueryRow(ConvertPlaceholders("SELECT COUNT(*) FROM likes WHERE post_id = ?"), postID).Scan(&likesCount)
	if err != nil {
		sendErrorResponse(w, "Ошибка подсчета лайков: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"liked":       !exists,
		"likes_count": likesCount,
	})
}

// getLikeStatus получает статус лайка и количество
func getLikeStatus(w http.ResponseWriter, _ *http.Request, postID int, userID int) {
	var liked bool

	// Если пользователь авторизован - проверяем его лайк
	if userID > 0 {
		err := database.DB.QueryRow(ConvertPlaceholders("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)"), userID, postID).Scan(&liked)
		if err != nil {
			sendErrorResponse(w, "Ошибка проверки лайка: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Если не авторизован - liked = false

	var likesCount int
	err := database.DB.QueryRow(ConvertPlaceholders("SELECT COUNT(*) FROM likes WHERE post_id = ?"), postID).Scan(&likesCount)
	if err != nil {
		sendErrorResponse(w, "Ошибка подсчета лайков: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"liked":       liked,
		"likes_count": likesCount,
	})
}

// getLikers получает список пользователей, которые лайкнули пост
func getLikers(w http.ResponseWriter, _ *http.Request, postID int) {
	query := `
		SELECT u.id, u.name, u.last_name, u.avatar
		FROM likes l
		JOIN users u ON l.user_id = u.id
		WHERE l.post_id = ?
		ORDER BY l.created_at DESC
	`

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения списка лайков: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Liker struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		LastName *string `json:"last_name"`
		Avatar   *string `json:"avatar"`
	}

	var likers []Liker
	for rows.Next() {
		var liker Liker
		err := rows.Scan(&liker.ID, &liker.Name, &liker.LastName, &liker.Avatar)
		if err != nil {
			continue
		}
		likers = append(likers, liker)
	}

	if likers == nil {
		likers = []Liker{}
	}

	sendSuccessResponse(w, likers)
}
