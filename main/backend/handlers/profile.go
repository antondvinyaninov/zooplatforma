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
		Name              string `json:"name"`
		LastName          string `json:"last_name"`
		Bio               string `json:"bio"`
		Phone             string `json:"phone"`
		Location          string `json:"location"`
		ProfileVisibility string `json:"profile_visibility"`
		ShowPhone         string `json:"show_phone"`
		ShowEmail         string `json:"show_email"`
		AllowMessages     string `json:"allow_messages"`
		ShowOnline        string `json:"show_online"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Обновляем профиль (БЕЗ avatar и cover_photo - для них отдельные endpoints)
	query := `UPDATE users SET name = ?, last_name = ?, bio = ?, phone = ?, location = ?,
	          profile_visibility = ?, show_phone = ?, show_email = ?, allow_messages = ?, show_online = ?
	          WHERE id = ?`
	_, err := database.DB.Exec(query, req.Name, req.LastName, req.Bio, req.Phone, req.Location,
		req.ProfileVisibility, req.ShowPhone, req.ShowEmail, req.AllowMessages, req.ShowOnline, userID)
	if err != nil {
		sendErrorResponse(w, "Ошибка обновления профиля: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем обновленные данные пользователя
	var user models.User
	query = `SELECT id, name, last_name, email, bio, phone, location, avatar, cover_photo,
	         profile_visibility, show_phone, show_email, allow_messages, show_online, created_at 
	         FROM users WHERE id = ?`
	err = database.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Name, &user.LastName, &user.Email, &user.Bio, &user.Phone,
		&user.Location, &user.Avatar, &user.CoverPhoto,
		&user.ProfileVisibility, &user.ShowPhone, &user.ShowEmail, &user.AllowMessages, &user.ShowOnline,
		&user.CreatedAt,
	)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения данных пользователя", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, models.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		LastName:          user.LastName,
		Email:             user.Email,
		Bio:               user.Bio,
		Phone:             user.Phone,
		Location:          user.Location,
		Avatar:            user.Avatar,
		CoverPhoto:        user.CoverPhoto,
		ProfileVisibility: user.ProfileVisibility,
		ShowPhone:         user.ShowPhone,
		ShowEmail:         user.ShowEmail,
		AllowMessages:     user.AllowMessages,
		ShowOnline:        user.ShowOnline,
		CreatedAt:         user.CreatedAt,
	})
}
