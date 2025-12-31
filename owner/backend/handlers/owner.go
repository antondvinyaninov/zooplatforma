package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"owner/middleware"
	"owner/models"
)

// GetMyPets возвращает список питомцев текущего пользователя
func GetMyPets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, owner_id, name, species, breed, birth_date, gender, color, 
			       chip_number, photo, status, created_at, updated_at
			FROM pets
			WHERE owner_id = ?
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

// GetPetEvents возвращает события питомца
func GetPetEvents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		petID := r.URL.Query().Get("pet_id")
		if petID == "" {
			http.Error(w, "pet_id is required", http.StatusBadRequest)
			return
		}

		// Проверяем, что питомец принадлежит пользователю
		var ownerID int
		err := db.QueryRow("SELECT owner_id FROM pets WHERE id = ?", petID).Scan(&ownerID)
		if err != nil {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}

		if ownerID != userID {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT id, pet_id, event_type, event_date, description, clinic_id, created_at
			FROM pet_events
			WHERE pet_id = ?
			ORDER BY event_date DESC
		`, petID)

		if err != nil {
			log.Printf("❌ Error querying pet events: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		events := []models.PetEvent{}
		for rows.Next() {
			var event models.PetEvent
			err := rows.Scan(
				&event.ID, &event.PetID, &event.EventType, &event.EventDate,
				&event.Description, &event.ClinicID, &event.CreatedAt,
			)
			if err != nil {
				log.Printf("❌ Error scanning event: %v", err)
				continue
			}
			events = append(events, event)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    events,
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
