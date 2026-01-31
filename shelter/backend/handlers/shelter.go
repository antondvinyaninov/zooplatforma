package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"shelter/middleware"
	"shelter/models"

	pkgmiddleware "github.com/zooplatforma/pkg/middleware"
)

// GetMyShelters возвращает список приютов, где пользователь является членом
func GetMyShelters(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем userID из контекста (установлен AuthMiddleware)
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT o.id, o.name, o.type, o.description, o.address_city, o.logo
			FROM organizations o
			JOIN organization_members om ON o.id = om.organization_id
			WHERE om.user_id = ? AND o.type = 'shelter'
			ORDER BY om.joined_at ASC
		`, userID)

		if err != nil {
			log.Printf("❌ Error querying user shelters: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		shelters := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var name, orgType string
			var description, address, logo sql.NullString

			err := rows.Scan(&id, &name, &orgType, &description, &address, &logo)
			if err != nil {
				log.Printf("❌ Error scanning shelter: %v", err)
				continue
			}

			shelter := map[string]interface{}{
				"id":   id,
				"name": name,
				"type": orgType,
			}

			if description.Valid {
				shelter["description"] = description.String
			}
			if address.Valid {
				shelter["address"] = address.String
			}
			if logo.Valid {
				shelter["logo"] = logo.String
			}

			shelters = append(shelters, shelter)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    shelters,
		})
	}
}

// CreateShelter создает новый приют
func CreateShelter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получаем userID из контекста
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Парсим данные
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

		// Валидация
		if data.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		if data.Type == "" {
			data.Type = "shelter"
		}

		// Создаем организацию
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

		// Добавляем создателя как администратора
		_, err = db.Exec(`
			INSERT INTO organization_members (organization_id, user_id, role, joined_at)
			VALUES (?, ?, 'admin', datetime('now'))
		`, orgID, userID)

		if err != nil {
			log.Printf("❌ Error adding organization member: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Created shelter: ID=%d, name=%s, owner=%d", orgID, data.Name, userID)

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

// GetAnimals возвращает список животных текущего приюта
func GetAnimals(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем tenant ID из контекста
		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Tenant not found", http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(`
			SELECT id, organization_id, name, species, breed, age_years, age_months, 
			       gender, color, size, weight, chip_number, status, arrival_date, 
			       arrival_reason, health_status, vaccinated, sterilized, 
			       character_traits, special_needs, photo, photos, story, 
			       created_at, updated_at
			FROM shelter_animals
			WHERE organization_id = ?
			ORDER BY created_at DESC
		`, tenantID)

		if err != nil {
			log.Printf("❌ Error querying animals: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		animals := []models.ShelterAnimal{}
		for rows.Next() {
			var animal models.ShelterAnimal
			err := rows.Scan(
				&animal.ID, &animal.OrganizationID, &animal.Name, &animal.Species,
				&animal.Breed, &animal.AgeYears, &animal.AgeMonths, &animal.Gender,
				&animal.Color, &animal.Size, &animal.Weight, &animal.ChipNumber,
				&animal.Status, &animal.ArrivalDate, &animal.ArrivalReason,
				&animal.HealthStatus, &animal.Vaccinated, &animal.Sterilized,
				&animal.CharacterTraits, &animal.SpecialNeeds, &animal.Photo,
				&animal.Photos, &animal.Story, &animal.CreatedAt, &animal.UpdatedAt,
			)
			if err != nil {
				log.Printf("❌ Error scanning animal: %v", err)
				continue
			}
			animals = append(animals, animal)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    animals,
		})
	}
}

// GetStats возвращает статистику приюта
func GetStats(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Tenant not found", http.StatusInternalServerError)
			return
		}

		var totalAnimals, adoptedThisMonth, activeVolunteers, pendingRequests int

		// Всего животных в приюте
		db.QueryRow(`
			SELECT COUNT(*) FROM shelter_animals 
			WHERE organization_id = ? AND status = 'in_shelter'
		`, tenantID).Scan(&totalAnimals)

		// Пристроено за месяц
		db.QueryRow(`
			SELECT COUNT(*) FROM shelter_animals 
			WHERE organization_id = ? AND status = 'adopted' 
			AND updated_at >= date('now', '-1 month')
		`, tenantID).Scan(&adoptedThisMonth)

		// Активных волонтеров
		db.QueryRow(`
			SELECT COUNT(*) FROM shelter_volunteers 
			WHERE organization_id = ? AND status = 'active'
		`, tenantID).Scan(&activeVolunteers)

		// Заявок на пристройство
		db.QueryRow(`
			SELECT COUNT(*) FROM adoption_requests 
			WHERE organization_id = ? AND status = 'pending'
		`, tenantID).Scan(&pendingRequests)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]int{
				"total_animals":      totalAnimals,
				"adopted_this_month": adoptedThisMonth,
				"active_volunteers":  activeVolunteers,
				"pending_requests":   pendingRequests,
			},
		})
	}
}

// GetOrganization возвращает информацию о текущем приюте
func GetOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Tenant not found", http.StatusInternalServerError)
			return
		}

		var org models.Organization
		err := db.QueryRow(`
			SELECT 
				id, 
				name, 
				COALESCE(type, 'shelter') as type, 
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
