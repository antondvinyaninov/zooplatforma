package handlers

import (
	"database"
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

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		toggleLike(w, r, postID, userID)
	case http.MethodGet:
		getLikeStatus(w, r, postID, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// toggleLike добавляет или удаляет лайк
func toggleLike(w http.ResponseWriter, r *http.Request, postID int, userID int) {
	// Проверяем, есть ли уже лайк
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)", userID, postID).Scan(&exists)
	if err != nil {
		sendErrorResponse(w, "Ошибка проверки лайка: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		// Удаляем лайк
		_, err = database.DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			sendErrorResponse(w, "Ошибка удаления лайка: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Добавляем лайк
		_, err = database.DB.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
		if err != nil {
			sendErrorResponse(w, "Ошибка добавления лайка: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Получаем обновленное количество лайков
	var likesCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", postID).Scan(&likesCount)
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
func getLikeStatus(w http.ResponseWriter, r *http.Request, postID int, userID int) {
	var liked bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)", userID, postID).Scan(&liked)
	if err != nil {
		sendErrorResponse(w, "Ошибка проверки лайка: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var likesCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", postID).Scan(&likesCount)
	if err != nil {
		sendErrorResponse(w, "Ошибка подсчета лайков: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"liked":       liked,
		"likes_count": likesCount,
	})
}
