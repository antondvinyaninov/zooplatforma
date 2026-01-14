package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

			// Получаем attachments
			attachments, err := getMessageAttachments(db, msg.ID)
			if err == nil {
				msg.Attachments = attachments
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

// SendMediaMessageHandler отправляет сообщение с медиа (фото/видео)
func SendMediaMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		// Парсим multipart form
		err := r.ParseMultipartForm(50 << 20) // 50 MB max
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Получаем receiver_id
		receiverIDStr := r.FormValue("receiver_id")
		if receiverIDStr == "" {
			http.Error(w, "Receiver ID is required", http.StatusBadRequest)
			return
		}

		receiverID, err := strconv.Atoi(receiverIDStr)
		if err != nil {
			http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
			return
		}

		if receiverID == userID {
			http.Error(w, "Cannot send message to yourself", http.StatusBadRequest)
			return
		}

		// Получаем текст сообщения (опционально)
		content := r.FormValue("content")

		// Получаем файлы
		files := r.MultipartForm.File["media"]
		if len(files) == 0 {
			http.Error(w, "At least one media file is required", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли получатель
		receiverExists, err := userExists(db, receiverID)
		if err != nil || !receiverExists {
			http.Error(w, "Receiver not found", http.StatusNotFound)
			return
		}

		// Ищем или создаем чат
		chatID, err := getOrCreateChat(db, userID, receiverID)
		if err != nil {
			log.Printf("❌ Error getting/creating chat: %v", err)
			http.Error(w, "Failed to create chat", http.StatusInternalServerError)
			return
		}

		// Создаем сообщение
		result, err := db.Exec(`
			INSERT INTO messages (chat_id, sender_id, receiver_id, content, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, chatID, userID, receiverID, content, time.Now())

		if err != nil {
			log.Printf("❌ Error creating message: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		messageID, _ := result.LastInsertId()

		// Сохраняем файлы и создаем attachments
		var attachments []models.MessageAttachment
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				log.Printf("❌ Error opening file: %v", err)
				continue
			}
			defer file.Close()

			// Определяем тип файла
			fileType := "image"
			contentType := fileHeader.Header.Get("Content-Type")
			if strings.HasPrefix(contentType, "video/") {
				fileType = "video"
			}

			// Сохраняем файл
			filePath, err := saveUploadedFile(file, fileHeader.Filename)
			if err != nil {
				log.Printf("❌ Error saving file: %v", err)
				continue
			}

			// Создаем запись в БД
			attachResult, err := db.Exec(`
				INSERT INTO message_attachments (message_id, file_path, file_type, file_size, created_at)
				VALUES (?, ?, ?, ?, ?)
			`, messageID, filePath, fileType, fileHeader.Size, time.Now())

			if err != nil {
				log.Printf("❌ Error creating attachment: %v", err)
				continue
			}

			attachID, _ := attachResult.LastInsertId()
			attachments = append(attachments, models.MessageAttachment{
				ID:        int(attachID),
				MessageID: int(messageID),
				FilePath:  filePath,
				FileType:  fileType,
				FileSize:  int(fileHeader.Size),
				CreatedAt: time.Now(),
			})
		}

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

		// Добавляем attachments к сообщению
		message.Attachments = attachments

		log.Printf("✅ Media message sent: user %d -> user %d in chat %d (%d attachments)", userID, receiverID, chatID, len(attachments))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
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

	// Получаем attachments
	attachments, err := getMessageAttachments(db, messageID)
	if err == nil {
		msg.Attachments = attachments
	}

	return &msg, nil
}

func getMessageAttachments(db *sql.DB, messageID int) ([]models.MessageAttachment, error) {
	rows, err := db.Query(`
		SELECT id, message_id, file_path, file_type, file_size, created_at
		FROM message_attachments
		WHERE message_id = ?
		ORDER BY created_at ASC
	`, messageID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.MessageAttachment
	for rows.Next() {
		var attachment models.MessageAttachment
		err := rows.Scan(
			&attachment.ID, &attachment.MessageID, &attachment.FilePath,
			&attachment.FileType, &attachment.FileSize, &attachment.CreatedAt,
		)
		if err != nil {
			log.Printf("⚠️ Error scanning attachment: %v", err)
			continue
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

func saveUploadedFile(file multipart.File, filename string) (string, error) {
	// Генерируем уникальное имя файла
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Создаем папку для медиа сообщений
	uploadDir := "../../uploads/messages"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Сохраняем файл
	filePath := filepath.Join(uploadDir, newFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Возвращаем относительный путь для URL
	return fmt.Sprintf("/uploads/messages/%s", newFilename), nil
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
