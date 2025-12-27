package handlers

import (
	"database/sql"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"backend/models"

	"github.com/google/uuid"
)

const (
	MaxUploadSize = 10 * 1024 * 1024 // 10MB
	UploadDir     = "uploads"
)

type MediaHandler struct {
	DB *sql.DB
}

func NewMediaHandler(db *sql.DB) *MediaHandler {
	return &MediaHandler{DB: db}
}

// UploadMedia загружает медиа-файл
func (h *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	log.Printf("📥 UploadMedia вызван: метод=%s, путь=%s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("❌ Неверный метод: %s", r.Method)
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем user_id из контекста
	userID, ok := r.Context().Value("userID").(int)
	log.Printf("🔑 userID из контекста: %v, ok=%v", userID, ok)
	if !ok {
		log.Printf("❌ Не удалось получить userID из контекста")
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("✅ Пользователь авторизован: userID=%d", userID)

	// Ограничение размера
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		sendErrorResponse(w, "File too large. Max size: 10MB", http.StatusBadRequest)
		return
	}

	// Получаем файл
	file, header, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем тип медиа
	mediaType := r.FormValue("media_type")
	if mediaType == "" {
		mediaType = "photo" // По умолчанию
	}

	// Проверяем MIME type
	mimeType := header.Header.Get("Content-Type")
	if !isAllowedMimeType(mimeType, mediaType) {
		sendErrorResponse(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(header.Filename)
	fileName := uuid.New().String() + ext

	// Создаем путь к файлу
	now := time.Now()
	relativePath := filepath.Join("users", strconv.Itoa(userID), mediaType+"s",
		strconv.Itoa(now.Year()), fmt.Sprintf("%02d", now.Month()), fileName)
	fullPath := filepath.Join(UploadDir, relativePath)

	// Создаем директории
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		sendErrorResponse(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// Сохраняем файл
	dst, err := os.Create(fullPath)
	if err != nil {
		sendErrorResponse(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	fileSize, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(fullPath) // Удаляем файл при ошибке
		sendErrorResponse(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Получаем размеры изображения (если это фото)
	var width, height *int
	if mediaType == "photo" {
		file.Seek(0, 0) // Возвращаемся в начало файла
		img, _, err := image.DecodeConfig(file)
		if err == nil {
			w := img.Width
			h := img.Height
			width = &w
			height = &h
		}
	}

	// Сохраняем в БД
	query := `
		INSERT INTO user_media (user_id, file_name, original_name, file_path, file_size, mime_type, media_type, width, height)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := h.DB.Exec(query, userID, fileName, header.Filename, relativePath, fileSize, mimeType, mediaType, width, height)
	if err != nil {
		os.Remove(fullPath) // Удаляем файл при ошибке БД
		sendErrorResponse(w, "Failed to save to database", http.StatusInternalServerError)
		return
	}

	mediaID, _ := result.LastInsertId()

	// Формируем ответ
	media := models.UserMedia{
		ID:           int(mediaID),
		UserID:       userID,
		FileName:     fileName,
		OriginalName: header.Filename,
		FilePath:     relativePath,
		FileSize:     fileSize,
		MimeType:     mimeType,
		MediaType:    mediaType,
		Width:        width,
		Height:       height,
		UploadedAt:   now,
		URL:          "/api/media/file/" + strconv.Itoa(int(mediaID)),
	}

	sendSuccessResponse(w, media)
}

// GetUserMedia получает все медиа пользователя
func (h *MediaHandler) GetUserMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем user_id из URL
	userIDStr := strings.TrimPrefix(r.URL.Path, "/api/media/user/")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем тип медиа из query параметров
	mediaType := r.URL.Query().Get("type")

	// Формируем запрос
	query := `
		SELECT id, user_id, file_name, original_name, file_path, file_size, mime_type, media_type, width, height, duration, uploaded_at
		FROM user_media
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	if mediaType != "" {
		query += " AND media_type = ?"
		args = append(args, mediaType)
	}

	query += " ORDER BY uploaded_at DESC"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		sendErrorResponse(w, "Failed to fetch media", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var mediaList []models.UserMedia
	for rows.Next() {
		var media models.UserMedia
		err := rows.Scan(
			&media.ID, &media.UserID, &media.FileName, &media.OriginalName,
			&media.FilePath, &media.FileSize, &media.MimeType, &media.MediaType,
			&media.Width, &media.Height, &media.Duration, &media.UploadedAt,
		)
		if err != nil {
			continue
		}
		media.URL = "/api/media/file/" + strconv.Itoa(media.ID)
		mediaList = append(mediaList, media)
	}

	if mediaList == nil {
		mediaList = []models.UserMedia{}
	}

	sendSuccessResponse(w, mediaList)
}

// GetMediaFile отдает файл по ID
func (h *MediaHandler) GetMediaFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем media_id из URL
	mediaIDStr := strings.TrimPrefix(r.URL.Path, "/api/media/file/")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		sendErrorResponse(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле из БД
	var media models.UserMedia
	query := `SELECT file_path, mime_type FROM user_media WHERE id = ?`
	err = h.DB.QueryRow(query, mediaID).Scan(&media.FilePath, &media.MimeType)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		sendErrorResponse(w, "Failed to fetch media", http.StatusInternalServerError)
		return
	}

	// Открываем файл
	fullPath := filepath.Join(UploadDir, media.FilePath)
	file, err := os.Open(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	// Отдаем файл
	w.Header().Set("Content-Type", media.MimeType)
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Кеш на год
	io.Copy(w, file)
}

// DeleteMedia удаляет медиа-файл
func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем user_id из контекста
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем media_id из URL
	mediaIDStr := strings.TrimPrefix(r.URL.Path, "/api/media/delete/")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		sendErrorResponse(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	// Проверяем, что файл принадлежит пользователю
	var filePath string
	var ownerID int
	query := `SELECT user_id, file_path FROM user_media WHERE id = ?`
	err = h.DB.QueryRow(query, mediaID).Scan(&ownerID, &filePath)
	if err == sql.ErrNoRows {
		sendErrorResponse(w, "Media not found", http.StatusNotFound)
		return
	}
	if err != nil {
		sendErrorResponse(w, "Failed to fetch media", http.StatusInternalServerError)
		return
	}

	if ownerID != userID {
		sendErrorResponse(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Удаляем файл
	fullPath := filepath.Join(UploadDir, filePath)
	os.Remove(fullPath)

	// Удаляем из БД
	_, err = h.DB.Exec("DELETE FROM user_media WHERE id = ?", mediaID)
	if err != nil {
		sendErrorResponse(w, "Failed to delete from database", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]string{"message": "Media deleted successfully"})
}

// GetMediaStats получает статистику использования медиа
func (h *MediaHandler) GetMediaStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем user_id из контекста
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT 
			COUNT(*) as total_files,
			COALESCE(SUM(file_size), 0) as total_size,
			COALESCE(SUM(CASE WHEN media_type = 'photo' THEN 1 ELSE 0 END), 0) as photos_count,
			COALESCE(SUM(CASE WHEN media_type = 'video' THEN 1 ELSE 0 END), 0) as videos_count,
			COALESCE(SUM(CASE WHEN media_type = 'document' THEN 1 ELSE 0 END), 0) as docs_count
		FROM user_media
		WHERE user_id = ?
	`

	var stats models.MediaStats
	err := h.DB.QueryRow(query, userID).Scan(
		&stats.TotalFiles, &stats.TotalSize, &stats.PhotosCount, &stats.VideosCount, &stats.DocsCount,
	)
	if err != nil {
		sendErrorResponse(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, stats)
}

// Вспомогательные функции

func isAllowedMimeType(mimeType, mediaType string) bool {
	allowedTypes := map[string][]string{
		"photo": {
			"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
		},
		"video": {
			"video/mp4", "video/mpeg", "video/quicktime", "video/webm",
		},
		"document": {
			"application/pdf", "application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
	}

	allowed, ok := allowedTypes[mediaType]
	if !ok {
		return false
	}

	for _, t := range allowed {
		if t == mimeType {
			return true
		}
	}
	return false
}
