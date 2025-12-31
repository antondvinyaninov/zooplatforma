package handlers

import (
"clinic/middleware"
"clinic/models"
"database/sql"
"encoding/json"
"log"
"net/http"
)

// GetMyClinics возвращает список клиник, где пользователь является членом
func GetMyClinics(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT o.id, o.name, o.type, o.description, o.address_city, o.logo
			FROM organizations o
			JOIN organization_members om ON o.id = om.organization_id
			WHERE om.user_id = ? AND o.type = 'clinic'
			ORDER BY om.joined_at ASC
		`, userID)

		if err != nil {
			log.Printf("❌ Error querying user clinics: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		clinics := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var name, orgType string
			var description, address, logo sql.NullString

			err := rows.Scan(&id, &name, &orgType, &description, &address, &logo)
			if err != nil {
				log.Printf("❌ Error scanning clinic: %v", err)
				continue
			}

			clinic := map[string]interface{}{
				"id":   id,
				"name": name,
				"type": orgType,
			}

			if description.Valid {
				clinic["description"] = description.String
			}
			if address.Valid {
				clinic["address"] = address.String
			}
			if logo.Valid {
				clinic["logo"] = logo.String
			}

			clinics = append(clinics, clinic)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
"success": true,
"data":    clinics,
})
	}
}

// CreateClinic создает новую клинику
func CreateClinic(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var data struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Description string `json:"description"`
			AddressCity string `json:"address_city"`
			Phone       string `json:"phone"`
			Email       string `json:"email"`
			Website     string `json:"website"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if data.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		if data.Type == "" {
			data.Type = "clinic"
		}

		result, err := db.Exec(`
			INSERT INTO organizations (name, type, description, address_city, phone, email, website, owner_user_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		`, data.Name, data.Type, data.Description, data.AddressCity, data.Phone, data.Email, data.Website, userID)

		if err != nil {
			log.Printf("❌ Error creating organization: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		orgID, err := result.LastInsertId()
		if err != nil {
			log.Printf("❌ Error getting organization ID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO organization_members (organization_id, user_id, role, joined_at)
			VALUES (?, ?, 'admin', datetime('now'))
		`, orgID, userID)

		if err != nil {
			log.Printf("❌ Error adding organization member: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Created clinic: ID=%d, name=%s, owner=%d", orgID, data.Name, userID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
"success": true,
"data": map[string]interface{}{
"id":   orgID,
"name": data.Name,
},
})
	}
}

// GetMyPatients возвращает список пациентов клиники
func GetMyPatients(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, ok = middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT id, user_id, name, species, breed, birth_date, gender, color, 
			       chip_number, photo, status, created_at, updated_at
			FROM pets
			ORDER BY created_at DESC
			LIMIT 50
		`)

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

// GetAppointments возвращает список записей на приём
func GetAppointments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		appointments := []models.Appointment{}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
"success": true,
"data":    appointments,
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

// GetOrganization возвращает информацию о текущей клинике
func GetOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		var org models.Organization
		err := db.QueryRow(`
			SELECT 
				id, 
				name, 
				COALESCE(type, 'clinic') as type, 
				COALESCE(description, '') as description, 
				COALESCE(address_city, '') as address_city, 
				COALESCE(phone, '') as phone, 
				COALESCE(email, '') as email, 
				COALESCE(website, '') as website, 
				COALESCE(logo, '') as logo, 
				created_at
			FROM organizations
			WHERE id = ?
		`, tenantID).Scan(
&org.ID, &org.Name, &org.Type, &org.Description,
			&org.Address, &org.Phone, &org.Email, &org.Website,
			&org.Logo, &org.CreatedAt,
		)

		if err != nil {
			log.Printf("❌ Error getting organization: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
"success": true,
"data":    org,
})
	}
}
