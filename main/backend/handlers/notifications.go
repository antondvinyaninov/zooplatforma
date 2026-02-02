package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Notification struct {
	ID         int          `json:"id"`
	UserID     int          `json:"user_id"`
	Type       string       `json:"type"`
	ActorID    int          `json:"actor_id"`
	EntityType string       `json:"entity_type,omitempty"`
	EntityID   int          `json:"entity_id,omitempty"`
	Message    string       `json:"message"`
	IsRead     bool         `json:"is_read"`
	CreatedAt  time.Time    `json:"created_at"`
	Actor      *models.User `json:"actor,omitempty"`
}

type NotificationsHandler struct {
	DB *sql.DB
}

// GetNotifications - получить список уведомлений текущего пользователя
func (h *NotificationsHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok || userID == 0 {
		http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	query := `
		SELECT n.id, n.user_id, n.type, n.actor_id, n.entity_type, n.entity_id, 
		       n.message, n.is_read, n.created_at,
		       u.id, u.name, u.last_name, u.email, u.avatar
		FROM notifications n
		LEFT JOIN users u ON n.actor_id = u.id
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC
		LIMIT 50
	`

	rows, err := h.DB.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notifications := []Notification{}
	for rows.Next() {
		var n Notification
		var actor models.User
		var entityType, entityID sql.NullString
		var actorLastName, actorAvatar sql.NullString

		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.ActorID, &entityType, &entityID,
			&n.Message, &n.IsRead, &n.CreatedAt,
			&actor.ID, &actor.Name, &actorLastName, &actor.Email, &actorAvatar,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if entityType.Valid {
			n.EntityType = entityType.String
		}
		if entityID.Valid {
			id, _ := strconv.Atoi(entityID.String)
			n.EntityID = id
		}
		if actorLastName.Valid {
			actor.LastName = actorLastName.String
		}
		if actorAvatar.Valid {
			actor.Avatar = actorAvatar.String
		}

		n.Actor = &actor
		notifications = append(notifications, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    notifications,
	})
}

// GetUnreadCount - получить количество непрочитанных уведомлений
func (h *NotificationsHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok || userID == 0 {
		http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var count int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0", userID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]int{
			"count": count,
		},
	})
}

// MarkAsRead - отметить уведомление как прочитанное
func (h *NotificationsHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok || userID == 0 {
		http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Извлекаем ID из URL: /api/notifications/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	notificationID := strings.TrimSpace(path)

	if notificationID == "" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	// Проверяем, что уведомление принадлежит пользователю
	var ownerID int
	err := h.DB.QueryRow("SELECT user_id FROM notifications WHERE id = ?", notificationID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Отмечаем как прочитанное
	_, err = h.DB.Exec("UPDATE notifications SET is_read = 1 WHERE id = ?", notificationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"message": "Notification marked as read",
		},
	})
}

// MarkAllAsRead - отметить все уведомления как прочитанные
func (h *NotificationsHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok || userID == 0 {
		http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	_, err := h.DB.Exec("UPDATE notifications SET is_read = 1 WHERE user_id = ? AND is_read = 0", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"message": "All notifications marked as read",
		},
	})
}

// CreateNotification - создать уведомление (вспомогательная функция)
func (h *NotificationsHandler) CreateNotification(userID, actorID int, notifType, entityType string, entityID int, message string) error {
	// Не создаем уведомление, если пользователь сам совершил действие
	if userID == actorID {
		return nil
	}

	query := `
		INSERT INTO notifications (user_id, type, actor_id, entity_type, entity_id, message)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := h.DB.Exec(query, userID, notifType, actorID, entityType, entityID, message)
	return err
}

// Вспомогательные функции для создания уведомлений

func (h *NotificationsHandler) NotifyComment(postAuthorID, commenterID, postID int, commenterName string) error {
	message := fmt.Sprintf("%s прокомментировал ваш пост", commenterName)
	return h.CreateNotification(postAuthorID, commenterID, "comment", "post", postID, message)
}

func (h *NotificationsHandler) NotifyLike(postAuthorID, likerID, postID int, likerName string) error {
	message := fmt.Sprintf("%s лайкнул ваш пост", likerName)
	return h.CreateNotification(postAuthorID, likerID, "like", "post", postID, message)
}

func (h *NotificationsHandler) NotifyFriendRequest(recipientID, senderID, friendshipID int, senderName string) error {
	message := fmt.Sprintf("%s отправил вам запрос в друзья", senderName)
	return h.CreateNotification(recipientID, senderID, "friend_request", "friendship", friendshipID, message)
}

func (h *NotificationsHandler) NotifyFriendAccepted(recipientID, acceptorID, friendshipID int, acceptorName string) error {
	message := fmt.Sprintf("%s принял ваш запрос в друзья", acceptorName)
	return h.CreateNotification(recipientID, acceptorID, "friend_accepted", "friendship", friendshipID, message)
}
