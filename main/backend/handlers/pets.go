package handlers

import (
	"backend/models"
	"database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func UserPetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID пользователя из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/pets/user/")
	userID, err := strconv.Atoi(path)
	if err != nil {
		sendErrorResponse(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	getUserPets(w, r, userID)
}

func PetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createPet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PetHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/pets/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendErrorResponse(w, "Неверный ID питомца", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		deletePet(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUserPets(w http.ResponseWriter, _ *http.Request, userID int) {
	query := `SELECT id, user_id, name, species, photo, created_at FROM pets WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения питомцев: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pets []models.Pet
	for rows.Next() {
		var pet models.Pet
		err := rows.Scan(&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Photo, &pet.CreatedAt)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pets = append(pets, pet)
	}

	if pets == nil {
		pets = []models.Pet{}
	}

	sendSuccessResponse(w, pets)
}

func createPet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req models.CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		sendErrorResponse(w, "Имя питомца не может быть пустым", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO pets (user_id, name, species, photo) VALUES (?, ?, ?, ?)`
	result, err := database.DB.Exec(query, userID, req.Name, req.Species, req.Photo)
	if err != nil {
		sendErrorResponse(w, "Ошибка добавления питомца: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Получаем созданного питомца
	var pet models.Pet
	query = `SELECT id, user_id, name, species, photo, created_at FROM pets WHERE id = ?`
	err = database.DB.QueryRow(query, id).Scan(&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Photo, &pet.CreatedAt)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения питомца", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, pet)
}

func deletePet(w http.ResponseWriter, r *http.Request, petID int) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Проверяем, что питомец принадлежит пользователю
	var ownerID int
	err := database.DB.QueryRow("SELECT user_id FROM pets WHERE id = ?", petID).Scan(&ownerID)
	if err != nil {
		sendErrorResponse(w, "Питомец не найден", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		sendErrorResponse(w, "Нет прав на удаление этого питомца", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec("DELETE FROM pets WHERE id = ?", petID)
	if err != nil {
		sendErrorResponse(w, "Ошибка удаления питомца: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]string{"message": "Питомец удален"})
}
