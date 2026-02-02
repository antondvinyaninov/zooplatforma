package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"volunteer/models"

	pkgmiddleware "github.com/zooplatforma/pkg/middleware"
)

// GetMyTasks возвращает список задач текущего волонтера
func GetMyTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := pkgmiddleware.GetUserID(r)
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
		userID, ok := pkgmiddleware.GetUserID(r)
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
		userID, ok := pkgmiddleware.GetUserID(r)
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

// TakeCustody - взять животное под опеку
func TakeCustody(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			PetID int `json:"pet_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if req.PetID == 0 {
			http.Error(w, `{"success":false,"error":"pet_id is required"}`, http.StatusBadRequest)
			return
		}

		// Проверяем что питомец существует
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pets WHERE id = ?)", req.PetID).Scan(&exists)
		if err != nil {
			log.Printf("❌ Error checking pet: %v", err)
			http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, `{"success":false,"error":"Pet not found"}`, http.StatusNotFound)
			return
		}

		// Проверяем что у питомца еще нет куратора
		var currentCuratorID sql.NullInt64
		err = db.QueryRow("SELECT curator_id FROM pets WHERE id = ?", req.PetID).Scan(&currentCuratorID)
		if err != nil {
			log.Printf("❌ Error checking curator: %v", err)
			http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		if currentCuratorID.Valid && currentCuratorID.Int64 != 0 {
			http.Error(w, `{"success":false,"error":"This pet already has a curator"}`, http.StatusConflict)
			return
		}

		// Назначаем куратора
		_, err = db.Exec("UPDATE pets SET curator_id = ? WHERE id = ?", userID, req.PetID)
		if err != nil {
			log.Printf("❌ Error updating pet curator: %v", err)
			http.Error(w, `{"success":false,"error":"Failed to take custody"}`, http.StatusInternalServerError)
			return
		}

		log.Printf("✅ User %d took custody of pet %d", userID, req.PetID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Successfully took custody of the pet",
		})
	}
}

// ReleaseCustody - снять опеку с животного
func ReleaseCustody(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			PetID int `json:"pet_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if req.PetID == 0 {
			http.Error(w, `{"success":false,"error":"pet_id is required"}`, http.StatusBadRequest)
			return
		}

		// Проверяем что пользователь является куратором этого питомца
		var curatorID sql.NullInt64
		err := db.QueryRow("SELECT curator_id FROM pets WHERE id = ?", req.PetID).Scan(&curatorID)
		if err == sql.ErrNoRows {
			http.Error(w, `{"success":false,"error":"Pet not found"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			log.Printf("❌ Error checking curator: %v", err)
			http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		if !curatorID.Valid || int(curatorID.Int64) != userID {
			http.Error(w, `{"success":false,"error":"You are not the curator of this pet"}`, http.StatusForbidden)
			return
		}

		// Снимаем опеку
		_, err = db.Exec("UPDATE pets SET curator_id = NULL WHERE id = ?", req.PetID)
		if err != nil {
			log.Printf("❌ Error releasing custody: %v", err)
			http.Error(w, `{"success":false,"error":"Failed to release custody"}`, http.StatusInternalServerError)
			return
		}

		log.Printf("✅ User %d released custody of pet %d", userID, req.PetID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Successfully released custody of the pet",
		})
	}
}

// CreatePet - создать животное (прокси к PetBase API)
// Волонтер создает животное с curator_id = его user_id
func CreatePet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🐾 CreatePet called: %s %s", r.Method, r.URL.Path)
		log.Printf("🔍 Headers: Authorization=%s, Cookie=%v", r.Header.Get("Authorization"), r.Header.Get("Cookie"))

		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			log.Printf("❌ GetUserID failed: userID not in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("✅ User authenticated: userID=%d", userID)

		// Читаем тело запроса
		var petData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&petData); err != nil {
			http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		// Устанавливаем curator_id = текущий пользователь
		petData["curator_id"] = userID

		// Устанавливаем статус по умолчанию для подопечных
		if _, ok := petData["status"]; !ok {
			petData["status"] = "looking_for_home"
		}

		// Получаем информацию о волонтере для заполнения curator_name и curator_phone
		var curatorName, curatorPhone string
		err := db.QueryRow("SELECT name, phone FROM users WHERE id = ?", userID).Scan(&curatorName, &curatorPhone)
		if err == nil {
			if curatorName != "" {
				petData["curator_name"] = curatorName
			}
			if curatorPhone != "" {
				petData["curator_phone"] = curatorPhone
			}
		}

		// Преобразуем обратно в JSON
		jsonData, err := json.Marshal(petData)
		if err != nil {
			log.Printf("❌ Error marshaling pet data: %v", err)
			http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Проксируем запрос к PetBase API
		petbaseURL := "http://localhost:8100/api/pets"
		req, err := http.NewRequest("POST", petbaseURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("❌ Error creating request to PetBase: %v", err)
			http.Error(w, `{"success":false,"error":"Failed to create pet"}`, http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		// КРИТИЧЕСКИ ВАЖНО: Передаем user_id через X-User-ID header
		// PetBase middleware читает этот заголовок и устанавливает user_id в контекст
		req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
		log.Printf("🔑 Setting X-User-ID header: %d", userID)

		// Также передаем токен авторизации (если есть)
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			req.Header.Set("Authorization", authHeader)
		}
		if cookie, err := r.Cookie("auth_token"); err == nil {
			req.AddCookie(cookie)
		}

		log.Printf("🔄 Proxying pet creation to PetBase: curator_id=%d, X-User-ID=%d", userID, userID)

		// Отправляем запрос к PetBase
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ Error calling PetBase API: %v", err)
			http.Error(w, `{"success":false,"error":"Failed to create pet"}`, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Читаем ответ от PetBase
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("❌ Error decoding PetBase response: %v", err)
			http.Error(w, `{"success":false,"error":"Failed to parse response"}`, http.StatusInternalServerError)
			return
		}

		// Возвращаем ответ от PetBase
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(result)

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			log.Printf("✅ Volunteer %d created pet via PetBase (curator_id=%d)", userID, userID)
		} else {
			log.Printf("❌ PetBase returned error: status=%d, response=%v", resp.StatusCode, result)
		}
	}
}
