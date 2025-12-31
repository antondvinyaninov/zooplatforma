package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"volunteer/middleware"
	"volunteer/models"
)

// GetMyTasks возвращает список задач текущего волонтера
func GetMyTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Пока возвращаем пустой список, таблица tasks будет создана позже
		tasks := []models.Task{}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    tasks,
		})
	}
}

// GetMyPets возвращает список питомцев, за которыми ухаживает волонтер
func GetMyPets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Получаем питомцев, где пользователь является куратором
		rows, err := db.Query(`
			SELECT id, user_id, name, species, breed, birth_date, gender, color, 
			       chip_number, photo, status, created_at, updated_at
			FROM pets
			WHERE curator_id = ?
			ORDER BY created_at DESC
		`, userID)

		if err != nil {
			log.Printf("❌ Error querying pets: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		pets := []models.Pet{}
		for rows.Next() {
			var pet models.Pet
			err := rows.Scan(
				&pet.ID, &pet.OwnerID, &pet.Name, &pet.Species, &pet.Breed,
				&pet.BirthDate, &pet.Gender, &pet.Color, &pet.ChipNumber,
				&pet.Photo, &pet.Status, &pet.CreatedAt, &pet.UpdatedAt,
			)
			if err != nil {
				log.Printf("❌ Error scanning pet: %v", err)
				continue
			}
			pets = append(pets, pet)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    pets,
		})
	}
}

// GetProfile возвращает профиль текущего пользователя
func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user models.User
		err := db.QueryRow(`
			SELECT id, name, last_name, email, avatar, created_at
			FROM users
			WHERE id = ?
		`, userID).Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Avatar, &user.CreatedAt)

		if err != nil {
			log.Printf("❌ Error getting user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    user,
		})
	}
}
