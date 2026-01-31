package handlers

import (
	"database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CatalogHandler - публичный endpoint для каталога питомцев (без авторизации)
// Возвращает только питомцев со статусами: looking_for_home, lost, found, needs_help
// Поддерживает фильтрацию: ?status=lost&city=Ижевск&region=Удмуртская Респ&urgent=true&species=dog&search=Бобик
func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("📋 CatalogHandler called: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		log.Printf("❌ Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем query параметры для фильтрации
	queryParams := r.URL.Query()
	statusFilter := queryParams.Get("status")
	cityFilter := queryParams.Get("city")
	regionFilter := queryParams.Get("region")
	urgentFilter := queryParams.Get("urgent")
	speciesFilter := queryParams.Get("species")
	searchFilter := queryParams.Get("search")
	organizationFilter := queryParams.Get("organization")          // Фильтр по типу организации (shelter, vet_clinic и т.д.)
	fromOrganizationFilter := queryParams.Get("from_organization") // Фильтр "только из организаций" (true/false)

	log.Printf("🔍 Filters: status=%s, city=%s, region=%s, urgent=%s, species=%s, search=%s, organization=%s, from_organization=%s",
		statusFilter, cityFilter, regionFilter, urgentFilter, speciesFilter, searchFilter, organizationFilter, fromOrganizationFilter)

	// Строим динамический SQL запрос с JOIN к organizations
	query := `
		SELECT 
			p.id, p.user_id, p.name, p.species, p.breed, p.gender, p.birth_date, p.color, p.size, p.weight,
			p.chip_number, p.tattoo_number, p.ear_tag_number, p.passport_number,
			p.is_sterilized, p.sterilization_date, p.is_vaccinated,
			p.health_notes, p.character_traits, p.special_needs,
			p.status, p.status_updated_at, p.photo, p.photos, p.story,
			p.created_at, p.updated_at,
			p.distinctive_marks, p.owner_name, p.owner_address, p.owner_phone, p.owner_email,
			p.blood_type, p.allergies, p.chronic_diseases, p.current_medications,
			p.pedigree_number, p.registration_org,
			p.curator_id, p.curator_name, p.curator_phone, p.location, p.foster_address, p.shelter_name,
			p.city, p.region, p.urgent, p.contact_name, p.contact_phone,
			p.organization_id, o.name as organization_name, o.type as organization_type
		FROM pets p
		LEFT JOIN organizations o ON p.organization_id = o.id
		WHERE p.status IN ('looking_for_home', 'lost', 'found', 'needs_help')
	`

	// Добавляем фильтры
	var args []interface{}

	if statusFilter != "" {
		query += " AND p.status = ?"
		args = append(args, statusFilter)
	}

	if cityFilter != "" {
		query += " AND p.city = ?"
		args = append(args, cityFilter)
	}

	if regionFilter != "" {
		query += " AND p.region = ?"
		args = append(args, regionFilter)
	}

	if urgentFilter == "true" {
		query += " AND p.urgent = 1"
	}

	if speciesFilter != "" {
		query += " AND p.species = ?"
		args = append(args, speciesFilter)
	}

	if searchFilter != "" {
		query += " AND (p.name LIKE ? OR p.breed LIKE ?)"
		searchPattern := "%" + searchFilter + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Фильтр по типу организации
	if organizationFilter != "" {
		query += " AND o.type = ?"
		args = append(args, organizationFilter)
	}

	// Фильтр "только из организаций"
	if fromOrganizationFilter == "true" {
		query += " AND p.organization_id IS NOT NULL"
	}

	// Сортировка: срочные первые, потом новые
	query += " ORDER BY p.urgent DESC, p.created_at DESC"

	log.Printf("🔍 Executing catalog query with %d filters...", len(args))
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ Error querying catalog pets: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		var sterilizationDate, statusUpdatedAt *string
		var organizationName, organizationType *string

		err := rows.Scan(
			&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Breed, &pet.Gender,
			&pet.BirthDate, &pet.Color, &pet.Size, &pet.Weight,
			&pet.ChipNumber, &pet.TattooNumber, &pet.EarTagNumber, &pet.PassportNumber,
			&pet.IsSterilized, &sterilizationDate, &pet.IsVaccinated,
			&pet.HealthNotes, &pet.CharacterTraits, &pet.SpecialNeeds,
			&pet.Status, &statusUpdatedAt, &pet.Photo, &pet.Photos, &pet.Story,
			&pet.CreatedAt, &pet.UpdatedAt,
			&pet.DistinctiveMarks, &pet.OwnerName, &pet.OwnerAddress, &pet.OwnerPhone, &pet.OwnerEmail,
			&pet.BloodType, &pet.Allergies, &pet.ChronicDiseases, &pet.CurrentMedications,
			&pet.PedigreeNumber, &pet.RegistrationOrg,
			&pet.CuratorID, &pet.CuratorName, &pet.CuratorPhone, &pet.Location, &pet.FosterAddress, &pet.ShelterName,
			&pet.City, &pet.Region, &pet.Urgent, &pet.ContactName, &pet.ContactPhone,
			&pet.OrganizationID, &organizationName, &organizationType,
		)
		if err != nil {
			log.Printf("❌ Error scanning catalog pet: %v", err)
			continue
		}

		pet.SterilizationDate = sterilizationDate
		pet.StatusUpdatedAt = statusUpdatedAt
		pet.OrganizationName = organizationName
		pet.OrganizationType = organizationType

		log.Printf("✅ Found pet: ID=%d, Name=%s, Status=%s, Org=%v", pet.ID, pet.Name, pet.Status, organizationName)
		pets = append(pets, pet)
	}

	if pets == nil {
		pets = []Pet{}
	}

	log.Printf("📦 Returning %d pets for catalog", len(pets))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    pets,
	})
}

// Pet представляет полную карточку питомца
type Pet struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	Name              string    `json:"name"`
	Species           string    `json:"species"`
	Breed             string    `json:"breed,omitempty"`
	Gender            string    `json:"gender,omitempty"`
	BirthDate         string    `json:"birth_date,omitempty"`
	Color             string    `json:"color,omitempty"`
	Size              string    `json:"size,omitempty"`
	Weight            float64   `json:"weight,omitempty"`
	ChipNumber        string    `json:"chip_number,omitempty"`
	TattooNumber      string    `json:"tattoo_number,omitempty"`
	EarTagNumber      string    `json:"ear_tag_number,omitempty"`
	PassportNumber    string    `json:"passport_number,omitempty"`
	IsSterilized      bool      `json:"is_sterilized"`
	SterilizationDate *string   `json:"sterilization_date,omitempty"`
	IsVaccinated      bool      `json:"is_vaccinated"`
	HealthNotes       string    `json:"health_notes,omitempty"`
	CharacterTraits   string    `json:"character_traits,omitempty"`
	SpecialNeeds      string    `json:"special_needs,omitempty"`
	Status            string    `json:"status"`
	StatusUpdatedAt   *string   `json:"status_updated_at,omitempty"`
	Photo             string    `json:"photo,omitempty"`
	Photos            string    `json:"photos,omitempty"`
	Story             string    `json:"story,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	// Поля паспорта (миграция 012)
	DistinctiveMarks   string `json:"distinctive_marks,omitempty"`
	OwnerName          string `json:"owner_name,omitempty"`
	OwnerAddress       string `json:"owner_address,omitempty"`
	OwnerPhone         string `json:"owner_phone,omitempty"`
	OwnerEmail         string `json:"owner_email,omitempty"`
	BloodType          string `json:"blood_type,omitempty"`
	Allergies          string `json:"allergies,omitempty"`
	ChronicDiseases    string `json:"chronic_diseases,omitempty"`
	CurrentMedications string `json:"current_medications,omitempty"`
	PedigreeNumber     string `json:"pedigree_number,omitempty"`
	RegistrationOrg    string `json:"registration_org,omitempty"`
	// Куратор и локация (миграция 015)
	CuratorID     *int   `json:"curator_id,omitempty"`
	CuratorName   string `json:"curator_name,omitempty"`
	CuratorPhone  string `json:"curator_phone,omitempty"`
	Location      string `json:"location,omitempty"`
	FosterAddress string `json:"foster_address,omitempty"`
	ShelterName   string `json:"shelter_name,omitempty"`
	// Экстренные контакты и страховка (миграция 016)
	EmergencyContactName     string  `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone    string  `json:"emergency_contact_phone,omitempty"`
	EmergencyContactRelation string  `json:"emergency_contact_relation,omitempty"`
	VetClinicName            string  `json:"vet_clinic_name,omitempty"`
	VetClinicPhone           string  `json:"vet_clinic_phone,omitempty"`
	VetClinicAddress         string  `json:"vet_clinic_address,omitempty"`
	InsuranceCompany         string  `json:"insurance_company,omitempty"`
	InsurancePolicyNumber    string  `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate      *string `json:"insurance_expiry_date,omitempty"`
	// Карточки породы (миграция 017-018)
	CardTitle           string `json:"card_title,omitempty"`
	CardDescription     string `json:"card_description,omitempty"`
	CardCharacteristics string `json:"card_characteristics,omitempty"`
	CardCareTips        string `json:"card_care_tips,omitempty"`
	CardHealthInfo      string `json:"card_health_info,omitempty"`
	CardNutrition       string `json:"card_nutrition,omitempty"`
	CardPhotos          string `json:"card_photos,omitempty"`
	CardIsPublished     bool   `json:"card_is_published"`
	// Поля для каталога (миграция 020, 023)
	City             string  `json:"city,omitempty"`
	Region           string  `json:"region,omitempty"`
	Urgent           bool    `json:"urgent"`
	ContactName      string  `json:"contact_name,omitempty"`
	ContactPhone     string  `json:"contact_phone,omitempty"`
	OrganizationID   *int    `json:"organization_id,omitempty"`
	OrganizationName *string `json:"organization_name,omitempty"`
	OrganizationType *string `json:"organization_type,omitempty"`
}

// PetSummary представляет краткую информацию о питомце (для постов)
type PetSummary struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Photo   string `json:"photo,omitempty"`
	Status  string `json:"status"`
}

// CreatePetRequest запрос на создание питомца
type CreatePetRequest struct {
	UserID             int     `json:"user_id"`
	Name               string  `json:"name"`
	Species            string  `json:"species"`
	Breed              string  `json:"breed,omitempty"`
	Gender             string  `json:"gender,omitempty"`
	BirthDate          string  `json:"birth_date,omitempty"`
	Color              string  `json:"color,omitempty"`
	Size               string  `json:"size,omitempty"`
	Weight             float64 `json:"weight,omitempty"`
	ChipNumber         string  `json:"chip_number,omitempty"`
	TattooNumber       string  `json:"tattoo_number,omitempty"`
	EarTagNumber       string  `json:"ear_tag_number,omitempty"`
	PassportNumber     string  `json:"passport_number,omitempty"`
	IsSterilized       bool    `json:"is_sterilized"`
	SterilizationDate  *string `json:"sterilization_date,omitempty"`
	IsVaccinated       bool    `json:"is_vaccinated"`
	HealthNotes        string  `json:"health_notes,omitempty"`
	CharacterTraits    string  `json:"character_traits,omitempty"`
	SpecialNeeds       string  `json:"special_needs,omitempty"`
	Status             string  `json:"status,omitempty"`
	Photo              string  `json:"photo,omitempty"`
	Photos             string  `json:"photos,omitempty"`
	Story              string  `json:"story,omitempty"`
	DistinctiveMarks   string  `json:"distinctive_marks,omitempty"`
	OwnerName          string  `json:"owner_name,omitempty"`
	OwnerAddress       string  `json:"owner_address,omitempty"`
	OwnerPhone         string  `json:"owner_phone,omitempty"`
	OwnerEmail         string  `json:"owner_email,omitempty"`
	BloodType          string  `json:"blood_type,omitempty"`
	Allergies          string  `json:"allergies,omitempty"`
	ChronicDiseases    string  `json:"chronic_diseases,omitempty"`
	CurrentMedications string  `json:"current_medications,omitempty"`
	PedigreeNumber     string  `json:"pedigree_number,omitempty"`
	RegistrationOrg    string  `json:"registration_org,omitempty"`
	CuratorID          *int    `json:"curator_id,omitempty"`
	CuratorName        string  `json:"curator_name,omitempty"`
	CuratorPhone       string  `json:"curator_phone,omitempty"`
	Location           string  `json:"location,omitempty"`
	FosterAddress      string  `json:"foster_address,omitempty"`
	ShelterName        string  `json:"shelter_name,omitempty"`
	// Поля для каталога
	City           string `json:"city,omitempty"`
	Region         string `json:"region,omitempty"`
	Urgent         bool   `json:"urgent"`
	ContactName    string `json:"contact_name,omitempty"`
	ContactPhone   string `json:"contact_phone,omitempty"`
	OrganizationID *int   `json:"organization_id,omitempty"`
}

// PetsHandler обрабатывает /api/pets
func PetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllPets(w, r)
	case http.MethodPost:
		createPet(w, r)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// PetDetailHandler обрабатывает /api/pets/:id
func PetDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/pets/")

	// Проверяем специальные endpoints
	if strings.HasPrefix(path, "user/") {
		userID := strings.TrimPrefix(path, "user/")
		getPetsByUser(w, r, userID)
		return
	}

	if strings.HasPrefix(path, "status/") {
		status := strings.TrimPrefix(path, "status/")
		getPetsByStatus(w, r, status)
		return
	}

	// Обычный ID питомца
	parts := strings.Split(path, "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		sendError(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	// Проверяем /summary endpoint
	if len(parts) > 1 && parts[1] == "summary" {
		getPetSummary(w, r, id)
		return
	}

	// Проверяем /organization endpoint
	if len(parts) > 1 && parts[1] == "organization" {
		SetPetOrganizationHandler(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getPet(w, r, id)
	case http.MethodPut:
		updatePet(w, r, id)
	case http.MethodDelete:
		deletePet(w, r, id)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllPets получает список всех питомцев
func getAllPets(w http.ResponseWriter, _ *http.Request) {
	query := `
		SELECT id, user_id, name, species, breed, gender, birth_date, 
		       color, size, weight, chip_number, tattoo_number, ear_tag_number, passport_number,
		       is_sterilized, sterilization_date, is_vaccinated, health_notes, character_traits,
		       special_needs, status, status_updated_at, photo, photos, story,
		       created_at, updated_at,
		       distinctive_marks, owner_name, owner_address, owner_phone, owner_email,
		       blood_type, allergies, chronic_diseases, current_medications,
		       pedigree_number, registration_org,
		       curator_id, curator_name, curator_phone, location, foster_address, shelter_name,
		       city, region, urgent, contact_name, contact_phone, organization_id
		FROM pets
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		sendError(w, "Failed to fetch pets: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		err := rows.Scan(
			&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Breed,
			&pet.Gender, &pet.BirthDate, &pet.Color, &pet.Size, &pet.Weight,
			&pet.ChipNumber, &pet.TattooNumber, &pet.EarTagNumber, &pet.PassportNumber,
			&pet.IsSterilized, &pet.SterilizationDate, &pet.IsVaccinated,
			&pet.HealthNotes, &pet.CharacterTraits, &pet.SpecialNeeds,
			&pet.Status, &pet.StatusUpdatedAt, &pet.Photo, &pet.Photos, &pet.Story,
			&pet.CreatedAt, &pet.UpdatedAt,
			&pet.DistinctiveMarks, &pet.OwnerName, &pet.OwnerAddress, &pet.OwnerPhone, &pet.OwnerEmail,
			&pet.BloodType, &pet.Allergies, &pet.ChronicDiseases, &pet.CurrentMedications,
			&pet.PedigreeNumber, &pet.RegistrationOrg,
			&pet.CuratorID, &pet.CuratorName, &pet.CuratorPhone, &pet.Location, &pet.FosterAddress, &pet.ShelterName,
			&pet.City, &pet.Region, &pet.Urgent, &pet.ContactName, &pet.ContactPhone, &pet.OrganizationID,
		)
		if err != nil {
			sendError(w, "Failed to scan pet: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pets = append(pets, pet)
	}

	if pets == nil {
		pets = []Pet{}
	}

	sendSuccess(w, pets)
}

// getPet получает полную карточку питомца
func getPet(w http.ResponseWriter, _ *http.Request, id int) {
	query := `
		SELECT id, user_id, name, species, breed, gender, birth_date, 
		       color, size, weight, chip_number, tattoo_number, ear_tag_number, passport_number,
		       is_sterilized, sterilization_date, is_vaccinated, health_notes, character_traits,
		       special_needs, status, status_updated_at, photo, photos, story,
		       created_at, updated_at,
		       distinctive_marks, owner_name, owner_address, owner_phone, owner_email,
		       blood_type, allergies, chronic_diseases, current_medications,
		       pedigree_number, registration_org,
		       curator_id, curator_name, curator_phone, location, foster_address, shelter_name,
		       city, region, urgent, contact_name, contact_phone, organization_id
		FROM pets
		WHERE id = ?
	`

	var pet Pet
	err := database.DB.QueryRow(query, id).Scan(
		&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Breed,
		&pet.Gender, &pet.BirthDate, &pet.Color, &pet.Size, &pet.Weight,
		&pet.ChipNumber, &pet.TattooNumber, &pet.EarTagNumber, &pet.PassportNumber,
		&pet.IsSterilized, &pet.SterilizationDate, &pet.IsVaccinated,
		&pet.HealthNotes, &pet.CharacterTraits, &pet.SpecialNeeds,
		&pet.Status, &pet.StatusUpdatedAt, &pet.Photo, &pet.Photos, &pet.Story,
		&pet.CreatedAt, &pet.UpdatedAt,
		&pet.DistinctiveMarks, &pet.OwnerName, &pet.OwnerAddress, &pet.OwnerPhone, &pet.OwnerEmail,
		&pet.BloodType, &pet.Allergies, &pet.ChronicDiseases, &pet.CurrentMedications,
		&pet.PedigreeNumber, &pet.RegistrationOrg,
		&pet.CuratorID, &pet.CuratorName, &pet.CuratorPhone, &pet.Location, &pet.FosterAddress, &pet.ShelterName,
		&pet.City, &pet.Region, &pet.Urgent, &pet.ContactName, &pet.ContactPhone, &pet.OrganizationID,
	)

	if err != nil {
		sendError(w, "Pet not found", http.StatusNotFound)
		return
	}

	sendSuccess(w, pet)
}

// getPetSummary получает краткую информацию о питомце (для постов)
func getPetSummary(w http.ResponseWriter, _ *http.Request, id int) {
	query := `
		SELECT id, name, species, photo, status
		FROM pets
		WHERE id = ?
	`

	var summary PetSummary
	err := database.DB.QueryRow(query, id).Scan(
		&summary.ID, &summary.Name, &summary.Species, &summary.Photo, &summary.Status,
	)

	if err != nil {
		sendError(w, "Pet not found", http.StatusNotFound)
		return
	}

	sendSuccess(w, summary)
}

// Owner представляет информацию о владельце питомца
type Owner struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// PetsWithOwnerResponse ответ с владельцем и питомцами
type PetsWithOwnerResponse struct {
	Owner Owner `json:"owner"`
	Pets  []Pet `json:"pets"`
}

// getPetsByUser получает питомцев пользователя
// Поддерживает параметр ?include_owner=true для включения информации о владельце
func getPetsByUser(w http.ResponseWriter, r *http.Request, userIDStr string) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Проверяем параметр include_owner
	includeOwner := r.URL.Query().Get("include_owner") == "true"

	query := `
		SELECT id, user_id, name, species, breed, gender, birth_date, 
		       color, size, weight, chip_number, tattoo_number, ear_tag_number, passport_number,
		       is_sterilized, sterilization_date, is_vaccinated, health_notes, character_traits,
		       special_needs, status, status_updated_at, photo, photos, story,
		       created_at, updated_at,
		       distinctive_marks, owner_name, owner_address, owner_phone, owner_email,
		       blood_type, allergies, chronic_diseases, current_medications,
		       pedigree_number, registration_org,
		       curator_id, curator_name, curator_phone, location, foster_address, shelter_name,
		       city, region, urgent, contact_name, contact_phone, organization_id
		FROM pets
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		sendError(w, "Failed to fetch pets: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		err := rows.Scan(
			&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Breed,
			&pet.Gender, &pet.BirthDate, &pet.Color, &pet.Size, &pet.Weight,
			&pet.ChipNumber, &pet.TattooNumber, &pet.EarTagNumber, &pet.PassportNumber,
			&pet.IsSterilized, &pet.SterilizationDate, &pet.IsVaccinated,
			&pet.HealthNotes, &pet.CharacterTraits, &pet.SpecialNeeds,
			&pet.Status, &pet.StatusUpdatedAt, &pet.Photo, &pet.Photos, &pet.Story,
			&pet.CreatedAt, &pet.UpdatedAt,
			&pet.DistinctiveMarks, &pet.OwnerName, &pet.OwnerAddress, &pet.OwnerPhone, &pet.OwnerEmail,
			&pet.BloodType, &pet.Allergies, &pet.ChronicDiseases, &pet.CurrentMedications,
			&pet.PedigreeNumber, &pet.RegistrationOrg,
			&pet.CuratorID, &pet.CuratorName, &pet.CuratorPhone, &pet.Location, &pet.FosterAddress, &pet.ShelterName,
			&pet.City, &pet.Region, &pet.Urgent, &pet.ContactName, &pet.ContactPhone, &pet.OrganizationID,
		)
		if err != nil {
			sendError(w, "Failed to scan pet: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pets = append(pets, pet)
	}

	if pets == nil {
		pets = []Pet{}
	}

	// Если запрошена информация о владельце, загружаем её
	if includeOwner {
		var owner Owner
		ownerQuery := `
			SELECT id, name, last_name, email, phone, avatar
			FROM users
			WHERE id = ?
		`
		err := database.DB.QueryRow(ownerQuery, userID).Scan(
			&owner.ID, &owner.Name, &owner.LastName, &owner.Email, &owner.Phone, &owner.Avatar,
		)
		if err != nil {
			log.Printf("⚠️ Failed to fetch owner info: %v", err)
			// Если владелец не найден, возвращаем только питомцев
			sendSuccess(w, pets)
			return
		}

		// Возвращаем владельца + питомцев
		response := PetsWithOwnerResponse{
			Owner: owner,
			Pets:  pets,
		}
		sendSuccess(w, response)
		return
	}

	// Обычный ответ - только питомцы
	sendSuccess(w, pets)
}

// getPetsByStatus получает питомцев по статусу
func getPetsByStatus(w http.ResponseWriter, _ *http.Request, status string) {
	query := `
		SELECT id, user_id, name, species, breed, gender, birth_date, 
		       color, size, weight, chip_number, tattoo_number, ear_tag_number, passport_number,
		       is_sterilized, sterilization_date, is_vaccinated, health_notes, character_traits,
		       special_needs, status, status_updated_at, photo, photos, story,
		       created_at, updated_at,
		       distinctive_marks, owner_name, owner_address, owner_phone, owner_email,
		       blood_type, allergies, chronic_diseases, current_medications,
		       pedigree_number, registration_org,
		       curator_id, curator_name, curator_phone, location, foster_address, shelter_name,
		       city, region, urgent, contact_name, contact_phone, organization_id
		FROM pets
		WHERE status = ?
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, status)
	if err != nil {
		sendError(w, "Failed to fetch pets: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		err := rows.Scan(
			&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Breed,
			&pet.Gender, &pet.BirthDate, &pet.Color, &pet.Size, &pet.Weight,
			&pet.ChipNumber, &pet.TattooNumber, &pet.EarTagNumber, &pet.PassportNumber,
			&pet.IsSterilized, &pet.SterilizationDate, &pet.IsVaccinated,
			&pet.HealthNotes, &pet.CharacterTraits, &pet.SpecialNeeds,
			&pet.Status, &pet.StatusUpdatedAt, &pet.Photo, &pet.Photos, &pet.Story,
			&pet.CreatedAt, &pet.UpdatedAt,
			&pet.DistinctiveMarks, &pet.OwnerName, &pet.OwnerAddress, &pet.OwnerPhone, &pet.OwnerEmail,
			&pet.BloodType, &pet.Allergies, &pet.ChronicDiseases, &pet.CurrentMedications,
			&pet.PedigreeNumber, &pet.RegistrationOrg,
			&pet.CuratorID, &pet.CuratorName, &pet.CuratorPhone, &pet.Location, &pet.FosterAddress, &pet.ShelterName,
			&pet.City, &pet.Region, &pet.Urgent, &pet.ContactName, &pet.ContactPhone, &pet.OrganizationID,
		)
		if err != nil {
			sendError(w, "Failed to scan pet: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pets = append(pets, pet)
	}

	if pets == nil {
		pets = []Pet{}
	}

	sendSuccess(w, pets)
}

// createPet создаёт новую карточку питомца
func createPet(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из контекста (установлен middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		sendError(w, "Unauthorized: user_id not found", http.StatusUnauthorized)
		return
	}

	var req CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Name == "" {
		sendError(w, "Name is required", http.StatusBadRequest)
		return
	}

	// ВАЖНО: Логика установки владельца/куратора
	// Если передан curator_id - это волонтёр создаёт бездомное животное (НЕ устанавливаем user_id)
	// Если НЕ передан curator_id - это обычный пользователь создаёт своего питомца (устанавливаем user_id из токена)
	if req.CuratorID == nil || *req.CuratorID == 0 {
		// Обычный пользователь создаёт своего питомца
		req.UserID = userID
	} else {
		// Волонтёр создаёт бездомное животное - НЕ устанавливаем user_id
		req.UserID = 0
	}

	// Устанавливаем статус по умолчанию
	if req.Status == "" {
		req.Status = "home"
	}

	// Устанавливаем локацию по умолчанию
	if req.Location == "" {
		req.Location = "home"
	}

	query := `
		INSERT INTO pets (
			user_id, name, species, breed, gender, birth_date,
			color, size, weight, chip_number, tattoo_number, ear_tag_number, passport_number,
			is_sterilized, sterilization_date, is_vaccinated, health_notes, character_traits,
			special_needs, status, photo, photos, story,
			distinctive_marks, owner_name, owner_address, owner_phone, owner_email,
			blood_type, allergies, chronic_diseases, current_medications,
			pedigree_number, registration_org,
			curator_id, curator_name, curator_phone, location, foster_address, shelter_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := database.DB.Exec(
		query,
		req.UserID, req.Name, req.Species, req.Breed, req.Gender, req.BirthDate,
		req.Color, req.Size, req.Weight, req.ChipNumber, req.TattooNumber, req.EarTagNumber, req.PassportNumber,
		req.IsSterilized, req.SterilizationDate, req.IsVaccinated, req.HealthNotes, req.CharacterTraits,
		req.SpecialNeeds, req.Status, req.Photo, req.Photos, req.Story,
		req.DistinctiveMarks, req.OwnerName, req.OwnerAddress, req.OwnerPhone, req.OwnerEmail,
		req.BloodType, req.Allergies, req.ChronicDiseases, req.CurrentMedications,
		req.PedigreeNumber, req.RegistrationOrg,
		req.CuratorID, req.CuratorName, req.CuratorPhone, req.Location, req.FosterAddress, req.ShelterName,
	)

	if err != nil {
		// Логируем детали для отладки
		log.Printf("Failed to create pet. Error: %v", err)
		log.Printf("Query: %s", query)
		log.Printf("Values count: UserID=%v, Name=%v, Species=%v, Breed=%v, Gender=%v, BirthDate=%v",
			req.UserID, req.Name, req.Species, req.Breed, req.Gender, req.BirthDate)
		sendError(w, "Failed to create pet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Получаем созданного питомца
	getPet(w, nil, int(id))
}

// updatePet обновляет карточку питомца
func updatePet(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("🐾 updatePet called: pet_id=%d", id)

	// Получаем user_id из контекста (установлен middleware)
	userID, ok := r.Context().Value("user_id").(int)
	log.Printf("🔍 Context user_id: %d, ok=%v", userID, ok)

	if !ok || userID == 0 {
		log.Printf("❌ No user_id in context")
		sendError(w, "Unauthorized: user_id not found", http.StatusUnauthorized)
		return
	}

	log.Printf("✅ User authenticated: user_id=%d", userID)

	// Проверяем, что питомец принадлежит пользователю
	var ownerID int
	err := database.DB.QueryRow("SELECT user_id FROM pets WHERE id = ?", id).Scan(&ownerID)
	if err != nil {
		log.Printf("❌ Pet not found in DB: pet_id=%d, error=%v", id, err)
		sendError(w, "Pet not found", http.StatusNotFound)
		return
	}

	log.Printf("🔍 Pet owner check: pet_id=%d, owner_id=%d, requester_id=%d", id, ownerID, userID)

	if ownerID != userID {
		log.Printf("❌ Forbidden: user %d trying to edit pet owned by %d", userID, ownerID)
		sendError(w, "Forbidden: you can only edit your own pets", http.StatusForbidden)
		return
	}

	log.Printf("✅ Ownership verified, proceeding with update")

	var req CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE pets SET
			name = ?, species = ?, breed = ?, gender = ?, birth_date = ?,
			color = ?, size = ?, weight = ?, chip_number = ?, tattoo_number = ?, ear_tag_number = ?, passport_number = ?,
			is_sterilized = ?, sterilization_date = ?, is_vaccinated = ?, health_notes = ?, character_traits = ?,
			special_needs = ?, status = ?, photo = ?, photos = ?, story = ?,
			distinctive_marks = ?, owner_name = ?, owner_address = ?, owner_phone = ?, owner_email = ?,
			blood_type = ?, allergies = ?, chronic_diseases = ?, current_medications = ?,
			pedigree_number = ?, registration_org = ?,
			curator_id = ?, curator_name = ?, curator_phone = ?, location = ?, foster_address = ?, shelter_name = ?,
			city = ?, region = ?, urgent = ?, contact_name = ?, contact_phone = ?, organization_id = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err = database.DB.Exec(
		query,
		req.Name, req.Species, req.Breed, req.Gender, req.BirthDate,
		req.Color, req.Size, req.Weight, req.ChipNumber, req.TattooNumber, req.EarTagNumber, req.PassportNumber,
		req.IsSterilized, req.SterilizationDate, req.IsVaccinated, req.HealthNotes, req.CharacterTraits,
		req.SpecialNeeds, req.Status, req.Photo, req.Photos, req.Story,
		req.DistinctiveMarks, req.OwnerName, req.OwnerAddress, req.OwnerPhone, req.OwnerEmail,
		req.BloodType, req.Allergies, req.ChronicDiseases, req.CurrentMedications,
		req.PedigreeNumber, req.RegistrationOrg,
		req.CuratorID, req.CuratorName, req.CuratorPhone, req.Location, req.FosterAddress, req.ShelterName,
		req.City, req.Region, req.Urgent, req.ContactName, req.ContactPhone, req.OrganizationID,
		id,
	)

	if err != nil {
		sendError(w, "Failed to update pet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем обновлённого питомца
	getPet(w, nil, id)
}

// deletePet удаляет карточку питомца
func deletePet(w http.ResponseWriter, r *http.Request, id int) {
	// Получаем user_id из контекста (установлен middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		sendError(w, "Unauthorized: user_id not found", http.StatusUnauthorized)
		return
	}

	// Проверяем, что питомец принадлежит пользователю
	var ownerID int
	err := database.DB.QueryRow("SELECT user_id FROM pets WHERE id = ?", id).Scan(&ownerID)
	if err != nil {
		sendError(w, "Pet not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		sendError(w, "Forbidden: you can only delete your own pets", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec("DELETE FROM pets WHERE id = ?", id)
	if err != nil {
		sendError(w, "Failed to delete pet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, map[string]string{"message": "Pet deleted successfully"})
}

// SearchPetsHandler - поиск питомцев для клиники
// GET /api/pets/search?owner_phone=...&chip_number=...&name=...&status=...
func SearchPetsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 SearchPetsHandler called: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры поиска
	ownerPhone := r.URL.Query().Get("owner_phone")
	chipNumber := r.URL.Query().Get("chip_number")
	petName := r.URL.Query().Get("name")
	status := r.URL.Query().Get("status")

	if ownerPhone == "" && chipNumber == "" && petName == "" {
		sendError(w, "At least one search parameter required", http.StatusBadRequest)
		return
	}

	// Строим динамический запрос
	query := `
		SELECT 
			p.id, p.user_id, p.name, p.species, p.breed, p.gender, p.birth_date, p.color,
			p.chip_number, p.is_sterilized, p.sterilization_date, p.status, p.photo, p.created_at,
			u.name as owner_name, u.phone as owner_phone, u.email as owner_email
		FROM pets p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE 1=1
	`

	var args []interface{}

	if ownerPhone != "" {
		query += " AND u.phone LIKE ?"
		args = append(args, "%"+ownerPhone+"%")
	}

	if chipNumber != "" {
		query += " AND p.chip_number = ?"
		args = append(args, chipNumber)
	}

	if petName != "" {
		query += " AND p.name LIKE ?"
		args = append(args, "%"+petName+"%")
	}

	if status != "" {
		query += " AND p.status = ?"
		args = append(args, status)
	}

	query += " ORDER BY p.created_at DESC LIMIT 50"

	log.Printf("🔍 Search query: %s, args: %v", query, args)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ Database error: %v", err)
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pets []map[string]interface{}
	for rows.Next() {
		var (
			id, userID                                         int
			name, species, status                              string
			breed, gender, birthDate, color, chipNumber, photo *string
			isSterilized                                       bool
			sterilizationDate                                  *string
			createdAt                                          time.Time
			ownerName, ownerPhone, ownerEmail                  *string
		)

		err := rows.Scan(
			&id, &userID, &name, &species, &breed, &gender, &birthDate, &color,
			&chipNumber, &isSterilized, &sterilizationDate, &status, &photo, &createdAt,
			&ownerName, &ownerPhone, &ownerEmail,
		)

		if err != nil {
			log.Printf("❌ Scan error: %v", err)
			continue
		}

		pet := map[string]interface{}{
			"id":                 id,
			"user_id":            userID,
			"name":               name,
			"species":            species,
			"breed":              breed,
			"gender":             gender,
			"birth_date":         birthDate,
			"color":              color,
			"chip_number":        chipNumber,
			"is_sterilized":      isSterilized,
			"sterilization_date": sterilizationDate,
			"status":             status,
			"photo":              photo,
			"created_at":         createdAt,
			"owner_name":         ownerName,
			"owner_phone":        ownerPhone,
			"owner_email":        ownerEmail,
		}

		pets = append(pets, pet)
	}

	log.Printf("✅ Found %d pets", len(pets))

	sendSuccess(w, pets)
}
