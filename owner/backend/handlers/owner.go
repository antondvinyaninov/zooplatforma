package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"owner/models"
	"path/filepath"
	"strconv"
	"strings"

	"time"

	pkgmiddleware "github.com/zooplatforma/pkg/middleware"
)

// GetMyPets возвращает список питомцев текущего пользователя
// Временно отключено - будет реализовано через интеграцию с PetID
/*
func GetMyPets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, user_id, name, species, breed, birth_date, gender, color,
			       chip_number, photo, status, created_at, updated_at
			FROM pets
			WHERE user_id = ?
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
*/

// GetPetEvents возвращает события питомца
// Временно отключено - будет реализовано через интеграцию с PetID
/*
func GetPetEvents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
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
		err := db.QueryRow("SELECT user_id FROM pets WHERE id = ?", petID).Scan(&ownerID)
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
*/

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

// GetMyPets возвращает список питомцев текущего пользователя через PetBase API
func GetMyPets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			log.Printf("❌ Unauthorized: no userID in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("🔍 Getting pets for user_id: %d from PetBase", userID)

		// Получаем токен из исходного запроса
		token := r.Header.Get("Authorization")
		if token == "" {
			// Если токена нет в header, пробуем cookie
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = "Bearer " + cookie.Value
			}
		}

		if token == "" {
			log.Printf("❌ No token found to forward to PetBase")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Запрос к PetBase API с токеном
		petbaseURL := fmt.Sprintf("http://localhost:8100/api/pets/user/%d", userID)
		req, err := http.NewRequest("GET", petbaseURL, nil)
		if err != nil {
			log.Printf("❌ Error creating request: %v", err)
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Передаем токен в PetBase
		req.Header.Set("Authorization", token)
		log.Printf("🔐 Forwarding token to PetBase: %s", token[:20]+"...")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ Error calling PetBase API: %v", err)
			http.Error(w, "Failed to fetch pets from PetBase", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("❌ PetBase returned status: %d", resp.StatusCode)
			http.Error(w, "Failed to fetch pets from PetBase", resp.StatusCode)
			return
		}

		// Читаем ответ от PetBase
		var petbaseResponse struct {
			Success bool                     `json:"success"`
			Data    []map[string]interface{} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&petbaseResponse); err != nil {
			log.Printf("❌ Error decoding PetBase response: %v", err)
			http.Error(w, "Failed to decode PetBase response", http.StatusInternalServerError)
			return
		}

		if !petbaseResponse.Success {
			log.Printf("❌ PetBase returned success=false")
			http.Error(w, "PetBase request failed", http.StatusInternalServerError)
			return
		}

		// Преобразуем данные: photo → photo_url, gender → sex
		for i := range petbaseResponse.Data {
			pet := petbaseResponse.Data[i]

			// Преобразуем photo в photo_url
			if photo, ok := pet["photo"].(string); ok && photo != "" {
				pet["photo_url"] = photo
			}

			// Преобразуем gender в sex
			if gender, ok := pet["gender"].(string); ok && gender != "" {
				pet["sex"] = gender
			}

			// Добавляем verification_status если его нет
			if _, ok := pet["verification_status"]; !ok {
				pet["verification_status"] = "verified"
			}
		}

		log.Printf("✅ Got %d pets from PetBase for user %d", len(petbaseResponse.Data), userID)

		// Возвращаем преобразованные данные
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pets": petbaseResponse.Data,
		})
	}
}

// GetPet возвращает информацию о конкретном питомце через PetBase API
func GetPet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Извлекаем pet_id из URL
		petIDStr := strings.TrimPrefix(r.URL.Path, "/api/pets/")
		petID, err := strconv.Atoi(petIDStr)
		if err != nil {
			http.Error(w, "Invalid pet ID", http.StatusBadRequest)
			return
		}

		log.Printf("🔍 Getting pet %d for user %d from PetBase", petID, userID)

		// Получаем токен из исходного запроса
		token := r.Header.Get("Authorization")
		if token == "" {
			// Если токена нет в header, пробуем cookie
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = "Bearer " + cookie.Value
			}
		}

		if token == "" {
			log.Printf("❌ No token found to forward to PetBase")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Запрос к PetBase API с токеном
		petbaseURL := fmt.Sprintf("http://localhost:8100/api/pets/%d", petID)
		req, err := http.NewRequest("GET", petbaseURL, nil)
		if err != nil {
			log.Printf("❌ Error creating request: %v", err)
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Передаем токен в PetBase
		req.Header.Set("Authorization", token)
		log.Printf("🔐 Forwarding token to PetBase: %s", token[:20]+"...")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ Error calling PetBase API: %v", err)
			http.Error(w, "Failed to fetch pet from PetBase", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("❌ PetBase returned status: %d", resp.StatusCode)
			http.Error(w, "Failed to fetch pet from PetBase", resp.StatusCode)
			return
		}

		// Читаем ответ от PetBase
		var petbaseResponse struct {
			Success bool                   `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&petbaseResponse); err != nil {
			log.Printf("❌ Error decoding PetBase response: %v", err)
			http.Error(w, "Failed to decode PetBase response", http.StatusInternalServerError)
			return
		}

		if !petbaseResponse.Success {
			log.Printf("❌ PetBase returned success=false")
			http.Error(w, "PetBase request failed", http.StatusInternalServerError)
			return
		}

		pet := petbaseResponse.Data

		// Проверяем, что питомец принадлежит пользователю
		if ownerID, ok := pet["user_id"].(float64); ok && int(ownerID) != userID {
			log.Printf("❌ Access denied: pet owner %d != user %d", int(ownerID), userID)
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Преобразуем данные: photo → photo_url, gender → sex
		if photo, ok := pet["photo"].(string); ok && photo != "" {
			pet["photo_url"] = photo
		}

		if gender, ok := pet["gender"].(string); ok && gender != "" {
			pet["sex"] = gender
		}

		if _, ok := pet["verification_status"]; !ok {
			pet["verification_status"] = "verified"
		}

		log.Printf("✅ Found pet %d from PetBase", petID)

		// Возвращаем преобразованные данные
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pet": pet,
		})
	}
}

// GetBreeds возвращает список пород для указанного вида животного
func GetBreeds(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		species := r.URL.Query().Get("species")
		if species == "" {
			http.Error(w, "species parameter is required", http.StatusBadRequest)
			return
		}

		// Преобразуем species в species_id
		var speciesID int
		switch species {
		case "dog":
			speciesID = 1
		case "cat":
			speciesID = 2
		default:
			http.Error(w, "Invalid species. Use 'dog' or 'cat'", http.StatusBadRequest)
			return
		}

		log.Printf("🔍 Getting breeds for species: %s (id: %d)", species, speciesID)

		rows, err := db.Query(`
			SELECT id, name, name_en
			FROM breeds
			WHERE species_id = ?
			ORDER BY 
				CASE 
					WHEN name IN ('Метис', 'Беспородная') THEN 0
					ELSE 1
				END,
				name
		`, speciesID)

		if err != nil {
			log.Printf("❌ Error querying breeds: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Breed struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			NameEn string `json:"name_en"`
		}

		breeds := []Breed{}
		for rows.Next() {
			var breed Breed
			err := rows.Scan(&breed.ID, &breed.Name, &breed.NameEn)
			if err != nil {
				log.Printf("❌ Error scanning breed: %v", err)
				continue
			}
			breeds = append(breeds, breed)
		}

		log.Printf("✅ Found %d breeds for %s", len(breeds), species)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"breeds": breeds,
		})
	}
}

// CreatePet создаёт нового питомца через PetBase API
func CreatePet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Читаем данные из запроса
		var petData struct {
			Name       string   `json:"name"`
			Species    string   `json:"species"`
			Breed      *string  `json:"breed"`
			Age        int      `json:"age"`
			Sex        string   `json:"sex"`
			Color      string   `json:"color"`
			Weight     *float64 `json:"weight"`
			ChipNumber *string  `json:"chip_number"`
			Sterilized bool     `json:"sterilized"`
		}

		if err := json.NewDecoder(r.Body).Decode(&petData); err != nil {
			log.Printf("❌ Error decoding request: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Printf("🐾 Creating pet for user %d: %s (%s) via PetBase", userID, petData.Name, petData.Species)

		// Вычисляем дату рождения из возраста
		birthDate := ""
		if petData.Age > 0 {
			birthDate = fmt.Sprintf("%d-01-01", time.Now().Year()-petData.Age)
		}

		// Подготавливаем данные для PetBase API
		petbaseData := map[string]interface{}{
			"user_id":       userID,
			"name":          petData.Name,
			"species":       petData.Species,
			"breed":         petData.Breed,
			"birth_date":    birthDate,
			"gender":        petData.Sex,
			"color":         petData.Color,
			"weight":        petData.Weight,
			"chip_number":   petData.ChipNumber,
			"is_sterilized": petData.Sterilized,
			"status":        "home",
		}

		// Отправляем запрос к PetBase API
		jsonData, err := json.Marshal(petbaseData)
		if err != nil {
			log.Printf("❌ Error marshaling data: %v", err)
			http.Error(w, "Failed to prepare request", http.StatusInternalServerError)
			return
		}

		petbaseURL := "http://localhost:8100/api/pets"
		req, err := http.NewRequest("POST", petbaseURL, strings.NewReader(string(jsonData)))
		if err != nil {
			log.Printf("❌ Error creating request: %v", err)
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		// Получаем токен из исходного запроса
		token := r.Header.Get("Authorization")
		if token == "" {
			// Если токена нет в header, пробуем cookie
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = "Bearer " + cookie.Value
			}
		}

		if token == "" {
			log.Printf("❌ No token found to forward to PetBase")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Передаем токен в PetBase
		req.Header.Set("Authorization", token)
		log.Printf("🔐 Forwarding token to PetBase for pet creation")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("❌ Error calling PetBase API: %v", err)
			http.Error(w, "Failed to create pet in PetBase", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("❌ PetBase returned status %d: %s", resp.StatusCode, string(body))
			http.Error(w, "Failed to create pet in PetBase", resp.StatusCode)
			return
		}

		// Читаем ответ от PetBase
		var petbaseResponse struct {
			Success bool                   `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&petbaseResponse); err != nil {
			log.Printf("❌ Error decoding PetBase response: %v", err)
			http.Error(w, "Failed to decode PetBase response", http.StatusInternalServerError)
			return
		}

		if !petbaseResponse.Success {
			log.Printf("❌ PetBase returned success=false")
			http.Error(w, "PetBase request failed", http.StatusInternalServerError)
			return
		}

		pet := petbaseResponse.Data

		// Преобразуем данные: photo → photo_url, gender → sex
		if photo, ok := pet["photo"].(string); ok && photo != "" {
			pet["photo_url"] = photo
		}

		if gender, ok := pet["gender"].(string); ok && gender != "" {
			pet["sex"] = gender
		}

		if _, ok := pet["verification_status"]; !ok {
			pet["verification_status"] = "verified"
		}

		log.Printf("✅ Created pet %d via PetBase: %s", pet["id"], petData.Name)

		// Возвращаем преобразованные данные
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pet": pet,
		})
	}
}

// UploadPetPhoto загружает фото питомца
func UploadPetPhoto(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📸 === UPLOAD PHOTO REQUEST START ===")
		log.Printf("📸 Method: %s, URL: %s", r.Method, r.URL.Path)
		log.Printf("📸 Content-Type: %s", r.Header.Get("Content-Type"))
		log.Printf("📸 Content-Length: %s", r.Header.Get("Content-Length"))

		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			log.Printf("❌ No userID in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("📸 User ID from context: %d", userID)

		// Извлекаем pet_id из URL: /api/pets/{id}/photo
		parts := strings.Split(r.URL.Path, "/")
		log.Printf("📸 URL parts: %v (length: %d)", parts, len(parts))

		// parts = ["", "api", "pets", "{id}", "photo"]
		if len(parts) < 5 {
			log.Printf("❌ Invalid URL structure: %s (parts: %v)", r.URL.Path, parts)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		petIDStr := parts[3] // Берём 4-й элемент (индекс 3)
		log.Printf("📸 Pet ID string: %s", petIDStr)

		petID, err := strconv.Atoi(petIDStr)
		if err != nil {
			log.Printf("❌ Invalid pet ID: %s, error: %v", petIDStr, err)
			http.Error(w, "Invalid pet ID", http.StatusBadRequest)
			return
		}
		log.Printf("📸 Pet ID parsed: %d", petID)

		// Проверяем, что питомец принадлежит пользователю
		log.Printf("📸 Checking pet ownership...")
		var ownerID int
		err = db.QueryRow("SELECT user_id FROM pets WHERE id = ?", petID).Scan(&ownerID)
		if err == sql.ErrNoRows {
			log.Printf("❌ Pet %d not found", petID)
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}
		if err != nil {
			log.Printf("❌ Error checking pet owner: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Printf("📸 Pet %d owner: %d, current user: %d", petID, ownerID, userID)

		if ownerID != userID {
			log.Printf("❌ Access denied: pet owner %d != user %d", ownerID, userID)
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		log.Printf("✅ Ownership verified")

		// Парсим multipart form
		log.Printf("📸 Parsing multipart form (max 10MB)...")
		err = r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			log.Printf("❌ Error parsing form: %v", err)
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		log.Printf("✅ Form parsed successfully")

		log.Printf("📸 Getting file from form field 'photo'...")
		file, header, err := r.FormFile("photo")
		if err != nil {
			log.Printf("❌ Error getting file: %v", err)
			http.Error(w, "No file provided", http.StatusBadRequest)
			return
		}
		defer file.Close()
		log.Printf("📸 File received: %s, size: %d bytes", header.Filename, header.Size)

		// Проверяем тип файла
		contentType := header.Header.Get("Content-Type")
		log.Printf("📸 File content type: %s", contentType)
		if !strings.HasPrefix(contentType, "image/") {
			log.Printf("❌ Invalid content type: %s", contentType)
			http.Error(w, "File must be an image", http.StatusBadRequest)
			return
		}
		log.Printf("✅ Content type valid")

		// Создаём директорию для загрузок в общей папке uploads
		uploadDir := "../../uploads/pets"
		log.Printf("📸 Creating upload directory: %s", uploadDir)
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("❌ Error creating upload directory: %v", err)
			http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Upload directory ready")

		// Генерируем уникальное имя файла
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d_%d%s", petID, time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)
		log.Printf("📸 Saving file to: %s", filePath)

		// Сохраняем файл
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("❌ Error creating file: %v", err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		bytesWritten, err := io.Copy(dst, file)
		if err != nil {
			log.Printf("❌ Error copying file: %v", err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		log.Printf("✅ File saved: %d bytes written", bytesWritten)

		// Обновляем путь к фото в базе данных (только /uploads/pets/filename)
		photoURL := "/uploads/pets/" + filename
		log.Printf("📸 Updating database with photo URL: %s", photoURL)
		_, err = db.Exec("UPDATE pets SET photo = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", photoURL, petID)
		if err != nil {
			log.Printf("❌ Error updating pet photo: %v", err)
			http.Error(w, "Failed to update pet photo", http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Database updated")

		// Обновляем фото в PetID (ЗооБаза)
		log.Printf("📸 Updating photo in PetID...")
		petIDURL := fmt.Sprintf("http://localhost:8100/api/pets/%d", petID)

		// Получаем токен из исходного запроса
		token := r.Header.Get("Authorization")
		if token == "" {
			// Если токена нет в header, пробуем cookie
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = "Bearer " + cookie.Value
			}
		}

		// Получаем текущие данные питомца из PetID
		petIDReq, err := http.NewRequest("GET", petIDURL, nil)
		if err != nil {
			log.Printf("⚠️ Failed to create PetID request: %v", err)
		} else {
			// Передаем токен вместо X-User-ID
			if token != "" {
				petIDReq.Header.Set("Authorization", token)
			}
			petIDResp, err := http.DefaultClient.Do(petIDReq)
			if err != nil {
				log.Printf("⚠️ Failed to get pet from PetID: %v", err)
			} else if petIDResp.StatusCode == 200 {
				var petIDData map[string]interface{}
				if err := json.NewDecoder(petIDResp.Body).Decode(&petIDData); err == nil {
					if petData, ok := petIDData["data"].(map[string]interface{}); ok {
						// Обновляем поле photo
						petData["photo"] = photoURL

						// Отправляем обновлённые данные обратно в PetID
						petDataJSON, _ := json.Marshal(petData)
						updateReq, err := http.NewRequest("PUT", petIDURL, strings.NewReader(string(petDataJSON)))
						if err == nil {
							updateReq.Header.Set("Content-Type", "application/json")
							// Передаем токен вместо X-User-ID
							if token != "" {
								updateReq.Header.Set("Authorization", token)
							}
							updateResp, err := http.DefaultClient.Do(updateReq)
							if err != nil {
								log.Printf("⚠️ Failed to update PetID: %v", err)
							} else {
								log.Printf("✅ PetID updated: %d", updateResp.StatusCode)
								updateResp.Body.Close()
							}
						}
					}
				}
				petIDResp.Body.Close()
			}
		}

		log.Printf("🎉 Photo upload complete for pet %d: %s", petID, photoURL)
		log.Printf("📸 === UPLOAD PHOTO REQUEST END ===")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"photo_url": photoURL,
		})
	}
}
