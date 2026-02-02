package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// UpdateProfileRequest - запрос на обновление профиля
type UpdateProfileRequest struct {
	Name        *string `json:"name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     *string `json:"country,omitempty"`
}

// getUserProfileHandler - получить профиль пользователя
func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"success":false,"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	// Получить пользователя из БД
	var user User
	err = db.QueryRow(`
		SELECT id, email, name, last_name, avatar, bio, phone, date_of_birth, 
		       gender, city, country, role, email_verified, created_at
		FROM users 
		WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Email, &user.Name, &user.LastName, &user.Avatar,
		&user.Bio, &user.Phone, &user.DateOfBirth,
		&user.Gender, &user.City, &user.Country,
		&user.Role, &user.EmailVerified, &user.CreatedAt,
	)

	if err != nil {
		log.Printf("❌ User not found: %v", err)
		http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)
		return
	}

	// Убрать пароль из ответа (его нет в структуре User, но на всякий случай)
	response := map[string]interface{}{
		"success": true,
		"data":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// updateUserProfileHandler - обновить профиль пользователя
func updateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"success":false,"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	// TODO: Проверить что пользователь обновляет свой профиль (из JWT)
	// Пока пропускаем для простоты

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Построить динамический UPDATE запрос
	updates := []string{}
	args := []interface{}{}

	if req.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.LastName != nil {
		updates = append(updates, "last_name = ?")
		args = append(args, *req.LastName)
	}
	if req.Bio != nil {
		updates = append(updates, "bio = ?")
		args = append(args, *req.Bio)
	}
	if req.Phone != nil {
		updates = append(updates, "phone = ?")
		args = append(args, *req.Phone)
	}
	if req.DateOfBirth != nil {
		updates = append(updates, "date_of_birth = ?")
		args = append(args, *req.DateOfBirth)
	}
	if req.Gender != nil {
		updates = append(updates, "gender = ?")
		args = append(args, *req.Gender)
	}
	if req.City != nil {
		updates = append(updates, "city = ?")
		args = append(args, *req.City)
	}
	if req.Country != nil {
		updates = append(updates, "country = ?")
		args = append(args, *req.Country)
	}

	if len(updates) == 0 {
		http.Error(w, `{"success":false,"error":"No fields to update"}`, http.StatusBadRequest)
		return
	}

	// Добавить updated_at
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	// Добавить userID в конец args
	args = append(args, userID)

	query := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = ?"

	_, err = db.Exec(query, args...)
	if err != nil {
		log.Printf("❌ Failed to update profile: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to update profile"}`, http.StatusInternalServerError)
		return
	}

	// Получить обновленный профиль
	var user User
	err = db.QueryRow(`
		SELECT id, email, name, last_name, avatar, bio, phone, date_of_birth,
		       gender, city, country, role, email_verified, created_at
		FROM users 
		WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Email, &user.Name, &user.LastName, &user.Avatar,
		&user.Bio, &user.Phone, &user.DateOfBirth,
		&user.Gender, &user.City, &user.Country,
		&user.Role, &user.EmailVerified, &user.CreatedAt,
	)

	if err != nil {
		log.Printf("❌ Failed to get updated profile: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to get updated profile"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data":    user,
		"message": "Profile updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// uploadAvatarHandler - загрузить аватар пользователя
func uploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"success":false,"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	// TODO: Проверить что пользователь загружает свой аватар (из JWT)

	// Парсинг multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, `{"success":false,"error":"File too large"}`, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, `{"success":false,"error":"No file uploaded"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверить тип файла
	ext := filepath.Ext(handler.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[strings.ToLower(ext)] {
		http.Error(w, `{"success":false,"error":"Invalid file type. Allowed: jpg, jpeg, png, gif, webp"}`, http.StatusBadRequest)
		return
	}

	// Создать папку для аватаров если не существует
	uploadsDir := "../../uploads/users"
	userDir := filepath.Join(uploadsDir, strconv.Itoa(userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		log.Printf("❌ Failed to create directory: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to create directory"}`, http.StatusInternalServerError)
		return
	}

	// Сгенерировать имя файла
	filename := "avatar" + ext
	filepath := filepath.Join(userDir, filename)

	// Создать файл
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("❌ Failed to create file: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Скопировать содержимое
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("❌ Failed to copy file: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to save file"}`, http.StatusInternalServerError)
		return
	}

	// Обновить путь к аватару в БД
	avatarURL := "/uploads/users/" + strconv.Itoa(userID) + "/" + filename
	_, err = db.Exec("UPDATE users SET avatar = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", avatarURL, userID)
	if err != nil {
		log.Printf("❌ Failed to update avatar in DB: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to update avatar"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"avatar_url": avatarURL,
		},
		"message": "Avatar uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// deleteAvatarHandler - удалить аватар пользователя
func deleteAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"success":false,"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	// TODO: Проверить что пользователь удаляет свой аватар (из JWT)

	// Обновить БД
	_, err = db.Exec("UPDATE users SET avatar = '', updated_at = CURRENT_TIMESTAMP WHERE id = ?", userID)
	if err != nil {
		log.Printf("❌ Failed to delete avatar: %v", err)
		http.Error(w, `{"success":false,"error":"Failed to delete avatar"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Avatar deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
