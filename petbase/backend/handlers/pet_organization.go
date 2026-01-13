package handlers

import (
	"database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SetPetOrganizationHandler - привязка питомца к организации
// PUT /api/pets/:id/organization
func SetPetOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID питомца из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/pets/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "organization" {
		sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	petID, err := strconv.Atoi(parts[0])
	if err != nil {
		sendError(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста (установлен middleware)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсим тело запроса
	var req struct {
		OrganizationID *int `json:"organization_id"` // null для отвязки
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("🏢 Setting organization for pet %d: org_id=%v, user_id=%d", petID, req.OrganizationID, userID)

	// Проверяем, что питомец существует и принадлежит пользователю
	var currentUserID int
	var currentOrgID *int
	err = database.DB.QueryRow(`
		SELECT user_id, organization_id FROM pets WHERE id = ?
	`, petID).Scan(&currentUserID, &currentOrgID)

	if err != nil {
		log.Printf("❌ Pet not found: %v", err)
		sendError(w, "Pet not found", http.StatusNotFound)
		return
	}

	// Проверяем права: либо владелец питомца, либо член организации
	hasPermission := false

	if currentUserID == userID {
		// Владелец питомца может привязать к любой организации
		hasPermission = true
		log.Printf("✅ User is pet owner")
	} else if req.OrganizationID != nil {
		// Проверяем, является ли пользователь членом организации
		var memberCount int
		err = database.DB.QueryRow(`
			SELECT COUNT(*) FROM organization_members 
			WHERE organization_id = ? AND user_id = ? AND role IN ('owner', 'admin', 'moderator')
		`, *req.OrganizationID, userID).Scan(&memberCount)

		if err == nil && memberCount > 0 {
			hasPermission = true
			log.Printf("✅ User is organization member")
		}
	}

	if !hasPermission {
		log.Printf("❌ User has no permission")
		sendError(w, "You don't have permission to modify this pet", http.StatusForbidden)
		return
	}

	// Если привязываем к организации, проверяем что она существует и верифицирована
	if req.OrganizationID != nil {
		var isVerified bool
		var orgStatus string
		err = database.DB.QueryRow(`
			SELECT is_verified, status FROM organizations WHERE id = ?
		`, *req.OrganizationID).Scan(&isVerified, &orgStatus)

		if err != nil {
			log.Printf("❌ Organization not found: %v", err)
			sendError(w, "Organization not found", http.StatusNotFound)
			return
		}

		if !isVerified || orgStatus != "active" {
			log.Printf("❌ Organization not verified or not active")
			sendError(w, "Organization must be verified and active", http.StatusBadRequest)
			return
		}
	}

	// Обновляем organization_id у питомца
	_, err = database.DB.Exec(`
		UPDATE pets 
		SET organization_id = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, req.OrganizationID, petID)

	if err != nil {
		log.Printf("❌ Failed to update pet organization: %v", err)
		sendError(w, "Failed to update pet organization", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Pet %d organization updated to %v", petID, req.OrganizationID)

	// Возвращаем успех
	sendSuccess(w, map[string]interface{}{
		"message":         "Pet organization updated successfully",
		"pet_id":          petID,
		"organization_id": req.OrganizationID,
	})
}
