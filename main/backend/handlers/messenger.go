package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetChatsHandler возвращает список диалогов пользователя
func GetChatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		// Получаем все чаты пользователя
		query := `
			SELECT 
				c.id, c.user1_id, c.user2_id, c.last_message_id, c.last_message_at, c.created_at
			FROM chats c
			WHERE c.user1_id = ? OR c.user2_id = ?
			ORDER BY c.last_message_at DESC NULLS LAST, c.created_at DESC
		`

		rows, err := db.Query(query, userID, userID)
		if err != nil {
			log.Printf("❌ Error fetching chats: %v", err)
			http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var chats []models.Chat
		for rows.Next() {
			var chat models.Chat
			err := rows.Scan(
				&chat.ID, &chat.User1ID, &chat.User2ID,
				&chat.LastMessageID, &chat.LastMessageAt, &chat.CreatedAt,
			)
			if err != nil {
				log.Printf("❌ Error scanning chat: %v", err)
				continue
			}

			// Определяем ID собеседника
			otherUserID := chat.User1ID
			if chat.User1ID == userID {
				otherUserID = chat.User2ID
			}

			// Получаем информацию о собеседнике
			otherUser, err := getUserByID(db, otherUserID)
			if err != nil {
				log.Printf("❌ Error fetching other user: %v", err)
				continue
			}
			chat.OtherUser = otherUser

			// Получаем последнее сообщение
			if chat.LastMessageID != nil {
				lastMessage, err := getMessageByID(db, *chat.LastMessageID)
				if err == nil {
					chat.LastMessage = lastMessage
				}
			}

			// Подсчитываем непрочитанные сообщения
			unreadCount, err := getUnreadCount(db, chat.ID, userID)
			if err == nil {
				chat.UnreadCount = unreadCount
			}

			chats = append(chats, chat)
		}

		if chats == nil {
			chats = []models.Chat{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chats)
	}
}

// GetChatMessagesHandler возвращает сообщения конкретного чата
func GetChatMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		// Получаем ID чата из URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		chatID, err := strconv.Atoi(pathParts[3])
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		// Проверяем, что пользователь является участником чата
		if !isUserInChat(db, chatID, userID) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Получаем сообщения
		query := `
			SELECT 
				m.id, m.chat_id, m.sender_id, m.receiver_id, 
				m.content, m.is_read, m.read_at, m.created_at
			FROM messages m
			WHERE m.chat_id = ?
			ORDER BY m.created_at ASC
		`

		rows, err := db.Query(query, chatID)
		if err != nil {
			log.Printf("❌ Error fetching messages: %v", err)
			http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var messages []models.Message
		for rows.Next() {
			var msg models.Message
			err := rows.Scan(
				&msg.ID, &msg.ChatID, &msg.SenderID, &msg.ReceiverID,
				&msg.Content, &msg.IsRead, &msg.ReadAt, &msg.CreatedAt,
			)
			if err != nil {
				log.Printf("❌ Error scanning message: %v", err)
				continue
			}

			// Получаем информацию об отправителе
			sender, err := getUserByID(db, msg.SenderID)
			if err == nil {
				msg.Sender = sender
			}

			messages = append(messages, msg)
		}

		// Отмечаем все сообщения как прочитанные
		go markMessagesAsRead(db, chatID, userID)

		if messages == nil {
			messages = []models.Message{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

// SendMessageHandler отправляет текстовое сообщение
func SendMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var req struct {
			ReceiverID int    `json:"receiver_id"`
			Content    string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.ReceiverID == 0 || req.Content == "" {
			http.Error(w, "Receiver ID and content are required", http.StatusBadRequest)
			return
		}

		if req.ReceiverID == userID {
			http.Error(w, "Cannot send message to yourself", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли получатель
		receiverExists, err := userExists(db, req.ReceiverID)
		if err != nil || !receiverExists {
			http.Error(w, "Receiver not found", http.StatusNotFound)
			return
		}

		// Ищем или создаем чат
		chatID, err := getOrCreateChat(db, userID, req.ReceiverID)
		if err != nil {
			log.Printf("❌ Error getting/creating chat: %v", err)
			http.Error(w, "Failed to create chat", http.StatusInternalServerError)
			return
		}

		// Создаем сообщение
		result, err := db.Exec(`
			INSERT INTO messages (chat_id, sender_id, receiver_id, content, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, chatID, userID, req.ReceiverID, req.Content, time.Now())

		if err != nil {
			log.Printf("❌ Error creating message: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		messageID, _ := result.LastInsertId()

		// Обновляем last_message в чате
		_, err = db.Exec(`
			UPDATE chats 
			SET last_message_id = ?, last_message_at = ?
			WHERE id = ?
		`, messageID, time.Now(), chatID)

		if err != nil {
			log.Printf("⚠️ Warning: Failed to update chat last_message: %v", err)
		}

		// Получаем созданное сообщение
		message, err := getMessageByID(db, int(messageID))
		if err != nil {
			log.Printf("❌ Error fetching created message: %v", err)
			http.Error(w, "Message sent but failed to fetch", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Message sent: user %d -> user %d in chat %d", userID, req.ReceiverID, chatID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
	}
}

// GetUnreadCountHandler возвращает количество непрочитанных сообщений
func GetUnreadCountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var count int
		err := db.QueryRow(`
			SELECT COUNT(*) 
			FROM messages 
			WHERE receiver_id = ? AND is_read = 0
		`, userID).Scan(&count)

		if err != nil {
			log.Printf("❌ Error counting unread messages: %v", err)
			http.Error(w, "Failed to count unread messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": count})
	}
}

// Вспомогательные функции

func getOrCreateChat(db *sql.DB, user1ID, user2ID int) (int, error) {
	// Нормализуем порядок пользователей (меньший ID всегда первый)
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	// Ищем существующий чат
	var chatID int
	err := db.QueryRow(`
		SELECT id FROM chats 
		WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)
	`, user1ID, user2ID, user2ID, user1ID).Scan(&chatID)

	if err == nil {
		return chatID, nil
	}

	// Создаем новый чат
	result, err := db.Exec(`
		INSERT INTO chats (user1_id, user2_id, created_at)
		VALUES (?, ?, ?)
	`, user1ID, user2ID, time.Now())

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func isUserInChat(db *sql.DB, chatID, userID int) bool {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM chats 
		WHERE id = ? AND (user1_id = ? OR user2_id = ?)
	`, chatID, userID, userID).Scan(&count)

	return err == nil && count > 0
}

func getMessageByID(db *sql.DB, messageID int) (*models.Message, error) {
	var msg models.Message
	err := db.QueryRow(`
		SELECT id, chat_id, sender_id, receiver_id, content, is_read, read_at, created_at
		FROM messages WHERE id = ?
	`, messageID).Scan(
		&msg.ID, &msg.ChatID, &msg.SenderID, &msg.ReceiverID,
		&msg.Content, &msg.IsRead, &msg.ReadAt, &msg.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Получаем отправителя
	sender, err := getUserByID(db, msg.SenderID)
	if err == nil {
		msg.Sender = sender
	}

	return &msg, nil
}

func getUnreadCount(db *sql.DB, chatID, userID int) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM messages 
		WHERE chat_id = ? AND receiver_id = ? AND is_read = 0
	`, chatID, userID).Scan(&count)

	return count, err
}

func markMessagesAsRead(db *sql.DB, chatID, userID int) {
	_, err := db.Exec(`
		UPDATE messages 
		SET is_read = 1, read_at = ?
		WHERE chat_id = ? AND receiver_id = ? AND is_read = 0
	`, time.Now(), chatID, userID)

	if err != nil {
		log.Printf("⚠️ Warning: Failed to mark messages as read: %v", err)
	}
}

func userExists(db *sql.DB, userID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	return count > 0, err
}

func getUserByID(db *sql.DB, userID int) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, email, name, last_name, avatar, cover_photo, bio, 
		       location, phone, created_at
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Email, &user.Name, &user.LastName, &user.Avatar,
		&user.CoverPhoto, &user.Bio, &user.Location, &user.Phone,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
