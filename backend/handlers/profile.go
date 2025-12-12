package handlers

import (
	"backend/models"
	"database"
	"encoding/json"
	"net/http"
)

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID пользователя из контекста (установлен middleware)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name       string `json:"name"`
		Bio        string `json:"bio"`
		Phone      string `json:"phone"`
		Location   string `json:"location"`
		Avatar     string `json:"avatar"`
		CoverPhoto string `json:"cover_photo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Обновляем профиль
	query := `UPDATE users SET name = ?, bio = ?, phone = ?, location = ?, avatar = ?, cover_photo = ? WHERE id = ?`
	_, err := database.DB.Exec(query, req.Name, req.Bio, req.Phone, req.Location, req.Avatar, req.CoverPhoto, userID)
	if err != nil {
		sendErrorResponse(w, "Ошибка обновления профиля: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем обновленные данные пользователя
	var user models.User
	query = `SELECT id, name, email, bio, phone, location, avatar, cover_photo, created_at FROM users WHERE id = ?`
	err = database.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Bio, &user.Phone,
		&user.Location, &user.Avatar, &user.CoverPhoto, &user.CreatedAt,
	)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения данных пользователя", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, models.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Bio:        user.Bio,
		Phone:      user.Phone,
		Location:   user.Location,
		Avatar:     user.Avatar,
		CoverPhoto: user.CoverPhoto,
		CreatedAt:  user.CreatedAt,
	})
}
