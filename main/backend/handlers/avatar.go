package handlers

import (
	"database"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// UploadAvatarHandler - загрузка аватара пользователя
func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID пользователя из контекста
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	ipAddress := r.RemoteAddr

	// Парсим multipart form (максимум 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logSystemEvent("error", "profile", "upload_avatar", fmt.Sprintf("Ошибка парсинга формы: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		logSystemEvent("error", "profile", "upload_avatar", fmt.Sprintf("Файл не найден: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Файл не найден", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверяем тип файла
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		logSystemEvent("warning", "profile", "upload_avatar", fmt.Sprintf("Попытка загрузить не изображение: %s", contentType), &userID, ipAddress)
		sendErrorResponse(w, "Разрешены только изображения", http.StatusBadRequest)
		return
	}

	// Проверяем размер (максимум 10MB)
	if header.Size > 10<<20 {
		logSystemEvent("warning", "profile", "upload_avatar", fmt.Sprintf("Файл слишком большой: %d bytes", header.Size), &userID, ipAddress)
		sendErrorResponse(w, "Файл слишком большой (максимум 10MB)", http.StatusBadRequest)
		return
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Создаем папку для аватаров пользователя
	uploadDir := fmt.Sprintf("uploads/users/%d/avatars", userID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logSystemEvent("error", "profile", "upload_avatar", fmt.Sprintf("Ошибка создания директории: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка создания директории", http.StatusInternalServerError)
		return
	}

	// Сохраняем файл
	filePath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		logSystemEvent("error", "profile", "upload_avatar", fmt.Sprintf("Ошибка создания файла: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		logSystemEvent("error", "profile", "upload_avatar", fmt.Sprintf("Ошибка сохранения файла: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	// Формируем URL для доступа к файлу
	avatarURL := fmt.Sprintf("/uploads/users/%d/avatars/%s", userID, fileName)

	// Обновляем аватар в базе данных
	query := `UPDATE users SET avatar = ? WHERE id = ?`
	_, err = database.DB.Exec(query, avatarURL, userID)
	if err != nil {
		logSystemEvent("error", "profile", "upload_avatar", fmt.Sprintf("Ошибка обновления БД: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка обновления базы данных", http.StatusInternalServerError)
		return
	}

	logSystemEvent("info", "profile", "upload_avatar", fmt.Sprintf("Аватар успешно загружен: %s (размер: %d bytes)", fileName, header.Size), &userID, ipAddress)

	sendSuccessResponse(w, map[string]interface{}{
		"avatar_url": avatarURL,
		"message":    "Аватар успешно загружен",
	})
}

// UploadCoverPhotoHandler - загрузка обложки профиля
func UploadCoverPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID пользователя из контекста
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	ipAddress := r.RemoteAddr

	// Парсим multipart form (максимум 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logSystemEvent("error", "profile", "upload_cover", fmt.Sprintf("Ошибка парсинга формы: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("cover")
	if err != nil {
		logSystemEvent("error", "profile", "upload_cover", fmt.Sprintf("Файл не найден: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Файл не найден", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверяем тип файла
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		logSystemEvent("warning", "profile", "upload_cover", fmt.Sprintf("Попытка загрузить не изображение: %s", contentType), &userID, ipAddress)
		sendErrorResponse(w, "Разрешены только изображения", http.StatusBadRequest)
		return
	}

	// Проверяем размер (максимум 10MB)
	if header.Size > 10<<20 {
		logSystemEvent("warning", "profile", "upload_cover", fmt.Sprintf("Файл слишком большой: %d bytes", header.Size), &userID, ipAddress)
		sendErrorResponse(w, "Файл слишком большой (максимум 10MB)", http.StatusBadRequest)
		return
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Создаем папку для обложек пользователя
	uploadDir := fmt.Sprintf("uploads/users/%d/covers", userID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logSystemEvent("error", "profile", "upload_cover", fmt.Sprintf("Ошибка создания директории: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка создания директории", http.StatusInternalServerError)
		return
	}

	// Сохраняем файл
	filePath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		logSystemEvent("error", "profile", "upload_cover", fmt.Sprintf("Ошибка создания файла: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		logSystemEvent("error", "profile", "upload_cover", fmt.Sprintf("Ошибка сохранения файла: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	// Формируем URL для доступа к файлу
	coverURL := fmt.Sprintf("/uploads/users/%d/covers/%s", userID, fileName)

	// Обновляем обложку в базе данных
	query := `UPDATE users SET cover_photo = ? WHERE id = ?`
	_, err = database.DB.Exec(query, coverURL, userID)
	if err != nil {
		logSystemEvent("error", "profile", "upload_cover", fmt.Sprintf("Ошибка обновления БД: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка обновления базы данных", http.StatusInternalServerError)
		return
	}

	logSystemEvent("info", "profile", "upload_cover", fmt.Sprintf("Обложка успешно загружена: %s (размер: %d bytes)", fileName, header.Size), &userID, ipAddress)

	sendSuccessResponse(w, map[string]interface{}{
		"cover_url": coverURL,
		"message":   "Обложка успешно загружена",
	})
}

// DeleteAvatarHandler - удаление аватара пользователя
func DeleteAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID пользователя из контекста
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	ipAddress := r.RemoteAddr

	// Обновляем аватар в базе данных (устанавливаем NULL)
	query := `UPDATE users SET avatar = NULL WHERE id = ?`
	_, err := database.DB.Exec(query, userID)
	if err != nil {
		logSystemEvent("error", "profile", "delete_avatar", fmt.Sprintf("Ошибка обновления БД: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка удаления аватара", http.StatusInternalServerError)
		return
	}

	logSystemEvent("info", "profile", "delete_avatar", "Аватар успешно удален", &userID, ipAddress)

	sendSuccessResponse(w, map[string]interface{}{
		"message": "Аватар успешно удален",
	})
}

// DeleteCoverPhotoHandler - удаление обложки профиля
func DeleteCoverPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID пользователя из контекста
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	ipAddress := r.RemoteAddr

	// Обновляем обложку в базе данных (устанавливаем NULL)
	query := `UPDATE users SET cover_photo = NULL WHERE id = ?`
	_, err := database.DB.Exec(query, userID)
	if err != nil {
		logSystemEvent("error", "profile", "delete_cover", fmt.Sprintf("Ошибка обновления БД: %v", err), &userID, ipAddress)
		sendErrorResponse(w, "Ошибка удаления обложки", http.StatusInternalServerError)
		return
	}

	logSystemEvent("info", "profile", "delete_cover", "Обложка успешно удалена", &userID, ipAddress)

	sendSuccessResponse(w, map[string]interface{}{
		"message": "Обложка успешно удалена",
	})
}
