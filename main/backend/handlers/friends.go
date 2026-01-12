package handlers

import (
	"backend/models"
	"database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SendFriendRequestHandler - отправить запрос в друзья
func SendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req models.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Нельзя добавить себя в друзья
	if req.FriendID == userID {
		sendErrorResponse(w, "Нельзя добавить себя в друзья", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли уже запрос
	var existingID int
	err := database.DB.QueryRow(`
		SELECT id FROM friendships 
		WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)
	`, userID, req.FriendID, req.FriendID, userID).Scan(&existingID)

	if err == nil {
		sendErrorResponse(w, "Запрос в друзья уже существует", http.StatusConflict)
		return
	}

	// Создаем запрос в друзья
	result, err := database.DB.Exec(`
		INSERT INTO friendships (user_id, friend_id, status, created_at, updated_at)
		VALUES (?, ?, 'pending', ?, ?)
	`, userID, req.FriendID, time.Now(), time.Now())

	if err != nil {
		sendErrorResponse(w, "Ошибка отправки запроса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Создаем уведомление для получателя запроса
	var senderName, senderLastName string
	err = database.DB.QueryRow(`
		SELECT name, COALESCE(last_name, '') FROM users WHERE id = ?
	`, userID).Scan(&senderName, &senderLastName)

	if err == nil {
		fullName := senderName
		if senderLastName != "" {
			fullName += " " + senderLastName
		}

		notifHandler := &NotificationsHandler{DB: database.DB}
		notifHandler.NotifyFriendRequest(req.FriendID, userID, int(id), fullName)
	}

	sendSuccessResponse(w, map[string]interface{}{
		"id":      id,
		"status":  "pending",
		"message": "Запрос в друзья отправлен",
	})
}

// AcceptFriendRequestHandler - принять запрос в друзья
func AcceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req models.FriendActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Обновляем статус запроса (только если текущий пользователь - получатель)
	result, err := database.DB.Exec(`
		UPDATE friendships 
		SET status = 'accepted', updated_at = ?
		WHERE id = ? AND friend_id = ? AND status = 'pending'
	`, time.Now(), req.FriendshipID, userID)

	if err != nil {
		sendErrorResponse(w, "Ошибка принятия запроса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		sendErrorResponse(w, "Запрос не найден или уже обработан", http.StatusNotFound)
		return
	}

	// Создаем уведомление для отправителя запроса
	var senderID int
	var acceptorName, acceptorLastName string
	err = database.DB.QueryRow(`
		SELECT f.user_id, u.name, COALESCE(u.last_name, '')
		FROM friendships f
		JOIN users u ON u.id = ?
		WHERE f.id = ?
	`, userID, req.FriendshipID).Scan(&senderID, &acceptorName, &acceptorLastName)

	if err == nil {
		fullName := acceptorName
		if acceptorLastName != "" {
			fullName += " " + acceptorLastName
		}

		notifHandler := &NotificationsHandler{DB: database.DB}
		notifHandler.NotifyFriendAccepted(senderID, userID, req.FriendshipID, fullName)
	}

	sendSuccessResponse(w, map[string]string{
		"message": "Запрос в друзья принят",
		"status":  "accepted",
	})
}

// RejectFriendRequestHandler - отклонить запрос в друзья
func RejectFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req models.FriendActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Удаляем запрос (только если текущий пользователь - получатель)
	result, err := database.DB.Exec(`
		DELETE FROM friendships 
		WHERE id = ? AND friend_id = ? AND status = 'pending'
	`, req.FriendshipID, userID)

	if err != nil {
		sendErrorResponse(w, "Ошибка отклонения запроса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		sendErrorResponse(w, "Запрос не найден или уже обработан", http.StatusNotFound)
		return
	}

	sendSuccessResponse(w, map[string]string{
		"message": "Запрос в друзья отклонен",
	})
}

// RemoveFriendHandler - удалить из друзей
func RemoveFriendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req models.FriendActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Удаляем дружбу (в любом направлении)
	result, err := database.DB.Exec(`
		DELETE FROM friendships 
		WHERE id = ? AND ((user_id = ? OR friend_id = ?)) AND status = 'accepted'
	`, req.FriendshipID, userID, userID)

	if err != nil {
		sendErrorResponse(w, "Ошибка удаления из друзей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		sendErrorResponse(w, "Дружба не найдена", http.StatusNotFound)
		return
	}

	sendSuccessResponse(w, map[string]string{
		"message": "Удалено из друзей",
	})
}

// GetFriendsHandler - получить список друзей
func GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Получаем всех друзей (где статус accepted)
	rows, err := database.DB.Query(`
		SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
		       u.id, u.name, u.last_name, u.email, u.avatar, u.location
		FROM friendships f
		JOIN users u ON (
			CASE 
				WHEN f.user_id = ? THEN u.id = f.friend_id
				ELSE u.id = f.user_id
			END
		)
		WHERE (f.user_id = ? OR f.friend_id = ?) AND f.status = 'accepted'
		ORDER BY f.created_at DESC
	`, userID, userID, userID)

	if err != nil {
		sendErrorResponse(w, "Ошибка получения друзей: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	friends := []models.FriendshipResponse{}
	for rows.Next() {
		var fr models.FriendshipResponse
		var friend models.UserResponse
		err := rows.Scan(
			&fr.ID, &fr.UserID, &fr.FriendID, &fr.Status, &fr.CreatedAt, &fr.UpdatedAt,
			&friend.ID, &friend.Name, &friend.LastName, &friend.Email, &friend.Avatar, &friend.Location,
		)
		if err != nil {
			continue
		}
		fr.Friend = friend
		friends = append(friends, fr)
	}

	sendSuccessResponse(w, friends)
}

// GetFriendRequestsHandler - получить входящие запросы в друзья
func GetFriendRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Получаем входящие запросы (где текущий пользователь - friend_id)
	rows, err := database.DB.Query(`
		SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
		       u.id, u.name, u.last_name, u.email, u.avatar, u.location
		FROM friendships f
		JOIN users u ON u.id = f.user_id
		WHERE f.friend_id = ? AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`, userID)

	if err != nil {
		sendErrorResponse(w, "Ошибка получения запросов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	requests := []models.FriendshipResponse{}
	for rows.Next() {
		var fr models.FriendshipResponse
		var friend models.UserResponse
		err := rows.Scan(
			&fr.ID, &fr.UserID, &fr.FriendID, &fr.Status, &fr.CreatedAt, &fr.UpdatedAt,
			&friend.ID, &friend.Name, &friend.LastName, &friend.Email, &friend.Avatar, &friend.Location,
		)
		if err != nil {
			continue
		}
		fr.Friend = friend
		requests = append(requests, fr)
	}

	sendSuccessResponse(w, requests)
}

// GetFriendshipStatusHandler - получить статус дружбы с пользователем
func GetFriendshipStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Получаем friend_id из query параметров
	friendIDStr := r.URL.Query().Get("friend_id")
	if friendIDStr == "" {
		sendErrorResponse(w, "friend_id обязателен", http.StatusBadRequest)
		return
	}

	var friendID int
	_, err := fmt.Sscanf(friendIDStr, "%d", &friendID)
	if err != nil {
		sendErrorResponse(w, "Неверный friend_id", http.StatusBadRequest)
		return
	}

	// Проверяем статус дружбы
	var friendship models.Friendship
	err = database.DB.QueryRow(`
		SELECT id, user_id, friend_id, status, created_at, updated_at
		FROM friendships
		WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)
	`, userID, friendID, friendID, userID).Scan(
		&friendship.ID, &friendship.UserID, &friendship.FriendID,
		&friendship.Status, &friendship.CreatedAt, &friendship.UpdatedAt,
	)

	if err != nil {
		// Нет дружбы
		sendSuccessResponse(w, map[string]interface{}{
			"status":      "none",
			"is_outgoing": false,
		})
		return
	}

	// Определяем, исходящий ли это запрос
	isOutgoing := friendship.UserID == userID

	sendSuccessResponse(w, map[string]interface{}{
		"id":          friendship.ID,
		"status":      friendship.Status,
		"is_outgoing": isOutgoing,
	})
}
