package handlers

import (
	"clinic/middleware"
	"clinic/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	pkgmiddleware "github.com/zooplatforma/pkg/middleware"
)

// GetMyClinics возвращает список клиник, где пользователь является членом
func GetMyClinics(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			log.Printf("❌ GetMyClinics: user ID not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("🔍 GetMyClinics: userID=%d", userID)

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

		userID, ok := pkgmiddleware.GetUserID(r)
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
		_, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Получаем tenant ID из контекста
		_, ok = middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		// Пока возвращаем всех питомцев для демонстрации
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
		_, ok := pkgmiddleware.GetUserID(r)
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

// UpdateOrganization обновляет информацию о клинике
func UpdateOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		var data struct {
			Name         string `json:"name"`
			Description  string `json:"description"`
			Address      string `json:"address"`
			Phone        string `json:"phone"`
			Email        string `json:"email"`
			Website      string `json:"website"`
			WorkingHours string `json:"working_hours"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if data.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`
			UPDATE organizations
			SET name = ?, description = ?, address_city = ?, phone = ?, email = ?, website = ?, updated_at = datetime('now')
			WHERE id = ?
		`, data.Name, data.Description, data.Address, data.Phone, data.Email, data.Website, tenantID)

		if err != nil {
			log.Printf("❌ Error updating organization: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Updated clinic: ID=%d, name=%s", tenantID, data.Name)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Organization updated successfully",
		})
	}
}

// GetMembers возвращает список участников клиники
func GetMembers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			log.Printf("❌ GetMembers: Clinic not selected")
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		log.Printf("🔍 GetMembers: Loading members for clinic_id=%d", tenantID)

		rows, err := db.Query(`
			SELECT 
				om.id, om.user_id, om.role, om.position, om.joined_at,
				u.name, u.last_name, u.email, u.avatar
			FROM organization_members om
			JOIN users u ON om.user_id = u.id
			WHERE om.organization_id = ?
			ORDER BY 
				CASE om.role
					WHEN 'owner' THEN 1
					WHEN 'admin' THEN 2
					WHEN 'moderator' THEN 3
					ELSE 4
				END,
				om.joined_at ASC
		`, tenantID)

		if err != nil {
			log.Printf("❌ Error querying members: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		members := []map[string]interface{}{}
		for rows.Next() {
			var id, userID int
			var role, joinedAt string
			var position, name, lastName, email, avatar sql.NullString

			err := rows.Scan(&id, &userID, &role, &position, &joinedAt, &name, &lastName, &email, &avatar)
			if err != nil {
				log.Printf("❌ Error scanning member: %v", err)
				continue
			}

			member := map[string]interface{}{
				"id":        id,
				"user_id":   userID,
				"role":      role,
				"position":  "",
				"joined_at": joinedAt,
			}

			if name.Valid {
				member["name"] = name.String
			}
			if lastName.Valid {
				member["last_name"] = lastName.String
			}
			if email.Valid {
				member["email"] = email.String
			}
			if avatar.Valid {
				// Добавляем полный URL для аватарки
				if avatar.String != "" && !strings.HasPrefix(avatar.String, "http") {
					member["avatar"] = "http://localhost:8000" + avatar.String
				} else {
					member["avatar"] = avatar.String
				}
			}
			if position.Valid {
				member["position"] = position.String
			}

			members = append(members, member)
			log.Printf("✅ Found member: id=%d, user_id=%d, role=%s, name=%s", id, userID, role, name.String)
		}

		log.Printf("📋 Total members found: %d", len(members))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    members,
		})
	}
}

// AddMember добавляет участника в клинику
func AddMember(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		var data struct {
			UserID   int    `json:"user_id"`
			Role     string `json:"role"`
			Position string `json:"position"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверяем права доступа
		var canManage bool
		err := db.QueryRow(`
			SELECT can_manage_members FROM organization_members
			WHERE organization_id = ? AND user_id = ?
		`, tenantID, userID).Scan(&canManage)

		if err == sql.ErrNoRows || !canManage {
			http.Error(w, "You don't have permission to manage members", http.StatusForbidden)
			return
		}

		// Проверяем, что пользователь еще не является участником
		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM organization_members WHERE organization_id = ? AND user_id = ?)
		`, tenantID, data.UserID).Scan(&exists)

		if err != nil {
			log.Printf("❌ Error checking member existence: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if exists {
			http.Error(w, "User is already a member", http.StatusBadRequest)
			return
		}

		// Устанавливаем права в зависимости от роли
		canPost := data.Role == "owner" || data.Role == "admin" || data.Role == "moderator"
		canEdit := data.Role == "owner" || data.Role == "admin"
		canManageMembers := data.Role == "owner" || data.Role == "admin"

		// Добавляем участника
		_, err = db.Exec(`
			INSERT INTO organization_members (organization_id, user_id, role, position, can_post, can_edit, can_manage_members, joined_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))
		`, tenantID, data.UserID, data.Role, data.Position, canPost, canEdit, canManageMembers)

		if err != nil {
			log.Printf("❌ Error adding member: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Added member: user_id=%d, clinic_id=%d, role=%s", data.UserID, tenantID, data.Role)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Member added successfully",
		})
	}
}

// UpdateMember обновляет роль участника
func UpdateMember(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		var data struct {
			MemberID int    `json:"member_id"`
			Role     string `json:"role"`
			Position string `json:"position"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверяем права доступа
		var canManage bool
		err := db.QueryRow(`
			SELECT can_manage_members FROM organization_members
			WHERE organization_id = ? AND user_id = ?
		`, tenantID, userID).Scan(&canManage)

		if err == sql.ErrNoRows || !canManage {
			http.Error(w, "You don't have permission to manage members", http.StatusForbidden)
			return
		}

		// Устанавливаем права в зависимости от роли
		canPost := data.Role == "owner" || data.Role == "admin" || data.Role == "moderator"
		canEdit := data.Role == "owner" || data.Role == "admin"
		canManageMembers := data.Role == "owner" || data.Role == "admin"

		// Обновляем участника
		_, err = db.Exec(`
			UPDATE organization_members 
			SET role = ?, position = ?, can_post = ?, can_edit = ?, can_manage_members = ?
			WHERE id = ? AND organization_id = ?
		`, data.Role, data.Position, canPost, canEdit, canManageMembers, data.MemberID, tenantID)

		if err != nil {
			log.Printf("❌ Error updating member: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Updated member: id=%d, role=%s", data.MemberID, data.Role)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Member updated successfully",
		})
	}
}

// RemoveMember удаляет участника из клиники
func RemoveMember(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := pkgmiddleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		var data struct {
			MemberID int `json:"member_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверяем права доступа
		var canManage bool
		err := db.QueryRow(`
			SELECT can_manage_members FROM organization_members
			WHERE organization_id = ? AND user_id = ?
		`, tenantID, userID).Scan(&canManage)

		if err == sql.ErrNoRows || !canManage {
			http.Error(w, "You don't have permission to manage members", http.StatusForbidden)
			return
		}

		// Проверяем, что не удаляем owner
		var memberRole string
		var memberUserID int
		err = db.QueryRow(`
			SELECT user_id, role FROM organization_members WHERE id = ? AND organization_id = ?
		`, data.MemberID, tenantID).Scan(&memberUserID, &memberRole)

		if err == sql.ErrNoRows {
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}

		if memberRole == "owner" {
			http.Error(w, "Cannot remove owner", http.StatusForbidden)
			return
		}

		// Нельзя удалить самого себя
		if memberUserID == userID {
			http.Error(w, "Cannot remove yourself", http.StatusForbidden)
			return
		}

		// Удаляем участника
		_, err = db.Exec("DELETE FROM organization_members WHERE id = ? AND organization_id = ?", data.MemberID, tenantID)

		if err != nil {
			log.Printf("❌ Error removing member: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Removed member: id=%d", data.MemberID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Member removed successfully",
		})
	}
}

// SearchUsers ищет пользователей для добавления в клинику
func SearchUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		tenantID, ok := middleware.GetTenantID(r)
		if !ok {
			http.Error(w, "Clinic not selected", http.StatusBadRequest)
			return
		}

		// Ищем пользователей, которые еще не являются участниками
		rows, err := db.Query(`
			SELECT u.id, u.name, u.last_name, u.email, u.avatar
			FROM users u
			WHERE (u.name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ?)
			AND u.id NOT IN (
				SELECT user_id FROM organization_members WHERE organization_id = ?
			)
			LIMIT 10
		`, "%"+query+"%", "%"+query+"%", "%"+query+"%", tenantID)

		if err != nil {
			log.Printf("❌ Error searching users: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var name, lastName, email, avatar sql.NullString

			err := rows.Scan(&id, &name, &lastName, &email, &avatar)
			if err != nil {
				log.Printf("❌ Error scanning user: %v", err)
				continue
			}

			user := map[string]interface{}{
				"id": id,
			}

			if name.Valid {
				user["name"] = name.String
			}
			if lastName.Valid {
				user["last_name"] = lastName.String
			}
			if email.Valid {
				user["email"] = email.String
			}
			if avatar.Valid {
				user["avatar"] = avatar.String
			}

			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    users,
		})
	}
}
