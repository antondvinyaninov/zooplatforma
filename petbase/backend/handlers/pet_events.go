package handlers

import (
	"database"
	"encoding/json"
	"log"
	"net/http"
	"petbase/models"
	"strconv"
	"strings"
	"time"
)

// GetPetEventsHandler получает историю событий питомца
// GET /api/petid/:id/events
func GetPetEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID питомца из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	petID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		sendError(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	// Получаем фильтр по типу события (опционально)
	eventType := r.URL.Query().Get("type")

	// Строим запрос
	query := `
		SELECT 
			e.id, e.pet_id, e.event_type, e.event_date,
			e.created_by_user_id, e.created_by_clinic_id, e.created_by_organization_id,
			e.title, e.description, e.details,
			e.vaccine_name, e.vaccine_batch, e.medication_name, e.dosage, e.next_date,
			e.previous_owner_id, e.new_owner_id, e.transfer_reason,
			e.location, e.circumstances, e.contact_phone, e.contact_name,
			e.death_reason, e.death_confirmed_by_clinic_id,
			e.shelter_id, e.adoption_contract,
			e.is_verified, e.verified_by_user_id, e.verified_at,
			e.created_at, e.updated_at,
			u1.name as created_by_user_name,
			o1.name as created_by_clinic_name,
			o2.name as created_by_organization_name,
			u2.name as previous_owner_name,
			u3.name as new_owner_name,
			o3.name as shelter_name
		FROM pet_events e
		LEFT JOIN users u1 ON e.created_by_user_id = u1.id
		LEFT JOIN organizations o1 ON e.created_by_clinic_id = o1.id
		LEFT JOIN organizations o2 ON e.created_by_organization_id = o2.id
		LEFT JOIN users u2 ON e.previous_owner_id = u2.id
		LEFT JOIN users u3 ON e.new_owner_id = u3.id
		LEFT JOIN organizations o3 ON e.shelter_id = o3.id
		WHERE e.pet_id = ?
	`

	args := []interface{}{petID}

	if eventType != "" {
		query += " AND e.event_type = ?"
		args = append(args, eventType)
	}

	query += " ORDER BY e.event_date DESC, e.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying pet events: %v", err)
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	events := []models.PetEvent{}
	for rows.Next() {
		var event models.PetEvent
		err := rows.Scan(
			&event.ID, &event.PetID, &event.Type, &event.Date,
			&event.CreatedByUserID, &event.CreatedByClinicID, &event.CreatedByOrganizationID,
			&event.Title, &event.Description, &event.Details,
			&event.VaccineName, &event.VaccineBatch, &event.MedicationName, &event.Dosage, &event.NextDate,
			&event.PreviousOwnerID, &event.NewOwnerID, &event.TransferReason,
			&event.Location, &event.Circumstances, &event.ContactPhone, &event.ContactName,
			&event.DeathReason, &event.DeathConfirmedByClinicID,
			&event.ShelterID, &event.AdoptionContract,
			&event.IsVerified, &event.VerifiedByUserID, &event.VerifiedAt,
			&event.CreatedAt, &event.UpdatedAt,
			&event.CreatedByUserName, &event.CreatedByClinicName, &event.CreatedByOrganizationName,
			&event.PreviousOwnerName, &event.NewOwnerName, &event.ShelterName,
		)
		if err != nil {
			log.Printf("Error scanning pet event: %v", err)
			continue
		}
		events = append(events, event)
	}

	sendSuccess(w, events)
}

// CreatePetEventHandler создаёт новое событие для питомца
// POST /api/petid/:id/events
func CreatePetEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID питомца из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	petID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		sendError(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из контекста (если авторизован)
	var userID *int
	if uid, ok := r.Context().Value("userID").(int); ok {
		userID = &uid
	}

	// Парсим тело запроса
	var req models.CreatePetEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.EventType == "" {
		sendError(w, "Event type is required", http.StatusBadRequest)
		return
	}

	// Проверяем, что тип события валидный
	validTypes := []string{
		models.EventTypeRegistration,
		models.EventTypeOwnershipChange,
		models.EventTypeSterilization,
		models.EventTypeVaccination,
		models.EventTypeTreatment,
		models.EventTypeLost,
		models.EventTypeFound,
		models.EventTypeDeath,
		models.EventTypeShelterIntake,
		models.EventTypeAdoption,
	}
	isValidType := false
	for _, t := range validTypes {
		if req.EventType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		sendError(w, "Invalid event type", http.StatusBadRequest)
		return
	}

	// Дата события (по умолчанию текущая)
	eventDate := time.Now().Format("2006-01-02 15:04:05")
	if req.EventDate != "" {
		eventDate = req.EventDate
	}

	// Создаём событие
	query := `
		INSERT INTO pet_events (
			pet_id, event_type, event_date,
			created_by_user_id,
			title, description, details,
			vaccine_name, vaccine_batch, medication_name, dosage, next_date,
			previous_owner_id, new_owner_id, transfer_reason,
			location, circumstances, contact_phone, contact_name,
			death_reason, death_confirmed_by_clinic_id,
			shelter_id, adoption_contract
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := database.DB.Exec(
		query,
		petID, req.EventType, eventDate,
		userID,
		nullString(req.Title), nullString(req.Description), nullString(req.Details),
		nullString(req.VaccineName), nullString(req.VaccineBatch), nullString(req.MedicationName),
		nullString(req.Dosage), nullString(req.NextDate),
		nullInt(req.PreviousOwnerID), nullInt(req.NewOwnerID), nullString(req.TransferReason),
		nullString(req.Location), nullString(req.Circumstances), nullString(req.ContactPhone), nullString(req.ContactName),
		nullString(req.DeathReason), nullInt(req.DeathConfirmedByClinicID),
		nullInt(req.ShelterID), nullString(req.AdoptionContract),
	)

	if err != nil {
		log.Printf("Error creating pet event: %v", err)
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}

	eventID, _ := result.LastInsertId()

	log.Printf("✅ Pet event created: ID=%d, PetID=%d, Type=%s", eventID, petID, req.EventType)

	sendSuccess(w, map[string]interface{}{
		"id":         eventID,
		"pet_id":     petID,
		"event_type": req.EventType,
		"event_date": eventDate,
	})
}

// GetPetMedicalHistoryHandler получает медицинскую историю питомца
// GET /api/petid/:id/medical
func GetPetMedicalHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID питомца из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	petID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		sendError(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	// Получаем только медицинские события
	query := `
		SELECT 
			e.id, e.pet_id, e.event_type, e.event_date,
			e.created_by_user_id, e.created_by_clinic_id,
			e.title, e.description,
			e.vaccine_name, e.vaccine_batch, e.medication_name, e.dosage, e.next_date,
			e.is_verified, e.created_at,
			u.name as created_by_user_name,
			o.name as created_by_clinic_name
		FROM pet_events e
		LEFT JOIN users u ON e.created_by_user_id = u.id
		LEFT JOIN organizations o ON e.created_by_clinic_id = o.id
		WHERE e.pet_id = ? 
		  AND e.event_type IN ('vaccination', 'treatment', 'sterilization')
		ORDER BY e.event_date DESC
	`

	rows, err := database.DB.Query(query, petID)
	if err != nil {
		log.Printf("Error querying medical history: %v", err)
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	events := []map[string]interface{}{}
	for rows.Next() {
		var (
			id, petID                                                   int
			eventType, eventDate                                        string
			createdByUserID, createdByClinicID                          *int
			title, description                                          *string
			vaccineName, vaccineBatch, medicationName, dosage, nextDate *string
			isVerified                                                  bool
			createdAt                                                   time.Time
			createdByUserName, createdByClinicName                      *string
		)

		err := rows.Scan(
			&id, &petID, &eventType, &eventDate,
			&createdByUserID, &createdByClinicID,
			&title, &description,
			&vaccineName, &vaccineBatch, &medicationName, &dosage, &nextDate,
			&isVerified, &createdAt,
			&createdByUserName, &createdByClinicName,
		)
		if err != nil {
			log.Printf("Error scanning medical event: %v", err)
			continue
		}

		event := map[string]interface{}{
			"id":          id,
			"pet_id":      petID,
			"event_type":  eventType,
			"event_date":  eventDate,
			"is_verified": isVerified,
			"created_at":  createdAt,
		}

		if title != nil {
			event["title"] = *title
		}
		if description != nil {
			event["description"] = *description
		}
		if vaccineName != nil {
			event["vaccine_name"] = *vaccineName
		}
		if vaccineBatch != nil {
			event["vaccine_batch"] = *vaccineBatch
		}
		if medicationName != nil {
			event["medication_name"] = *medicationName
		}
		if dosage != nil {
			event["dosage"] = *dosage
		}
		if nextDate != nil {
			event["next_date"] = *nextDate
		}
		if createdByUserName != nil {
			event["created_by_user_name"] = *createdByUserName
		}
		if createdByClinicName != nil {
			event["created_by_clinic_name"] = *createdByClinicName
		}

		events = append(events, event)
	}

	sendSuccess(w, events)
}

// Вспомогательные функции
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullInt(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}
