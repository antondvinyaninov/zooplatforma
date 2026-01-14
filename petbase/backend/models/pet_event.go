package models

import "time"

// PetEvent представляет событие в жизни питомца
type PetEvent struct {
	ID    int    `json:"id"`
	PetID int    `json:"pet_id"`
	Type  string `json:"event_type"` // registration, ownership_change, sterilization, vaccination, treatment, lost, found, death, shelter_intake, adoption
	Date  string `json:"event_date"`

	// Кто создал событие
	CreatedByUserID         *int `json:"created_by_user_id,omitempty"`
	CreatedByClinicID       *int `json:"created_by_clinic_id,omitempty"`
	CreatedByOrganizationID *int `json:"created_by_organization_id,omitempty"`

	// Детали события
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Details     *string `json:"details,omitempty"` // JSON

	// Медицинские данные
	VaccineName    *string `json:"vaccine_name,omitempty"`
	VaccineBatch   *string `json:"vaccine_batch,omitempty"`
	MedicationName *string `json:"medication_name,omitempty"`
	Dosage         *string `json:"dosage,omitempty"`
	NextDate       *string `json:"next_date,omitempty"`

	// Смена владельца
	PreviousOwnerID *int    `json:"previous_owner_id,omitempty"`
	NewOwnerID      *int    `json:"new_owner_id,omitempty"`
	TransferReason  *string `json:"transfer_reason,omitempty"`

	// Потеря/находка
	Location      *string `json:"location,omitempty"`
	Circumstances *string `json:"circumstances,omitempty"`
	ContactPhone  *string `json:"contact_phone,omitempty"`
	ContactName   *string `json:"contact_name,omitempty"`

	// Смерть
	DeathReason              *string `json:"death_reason,omitempty"` // natural, euthanasia, accident, disease
	DeathConfirmedByClinicID *int    `json:"death_confirmed_by_clinic_id,omitempty"`

	// Приют
	ShelterID        *int    `json:"shelter_id,omitempty"`
	AdoptionContract *string `json:"adoption_contract,omitempty"`

	// Верификация
	IsVerified       bool    `json:"is_verified"`
	VerifiedByUserID *int    `json:"verified_by_user_id,omitempty"`
	VerifiedAt       *string `json:"verified_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Дополнительные поля для отображения
	CreatedByUserName         *string `json:"created_by_user_name,omitempty"`
	CreatedByClinicName       *string `json:"created_by_clinic_name,omitempty"`
	CreatedByOrganizationName *string `json:"created_by_organization_name,omitempty"`
	PreviousOwnerName         *string `json:"previous_owner_name,omitempty"`
	NewOwnerName              *string `json:"new_owner_name,omitempty"`
	ShelterName               *string `json:"shelter_name,omitempty"`
}

// EventType константы для типов событий
const (
	EventTypeRegistration    = "registration"
	EventTypeOwnershipChange = "ownership_change"
	EventTypeSterilization   = "sterilization"
	EventTypeVaccination     = "vaccination"
	EventTypeTreatment       = "treatment"
	EventTypeLost            = "lost"
	EventTypeFound           = "found"
	EventTypeDeath           = "death"
	EventTypeShelterIntake   = "shelter_intake"
	EventTypeAdoption        = "adoption"
)

// DeathReason константы для причин смерти
const (
	DeathReasonNatural    = "natural"
	DeathReasonEuthanasia = "euthanasia"
	DeathReasonAccident   = "accident"
	DeathReasonDisease    = "disease"
)

// CreatePetEventRequest запрос на создание события
type CreatePetEventRequest struct {
	PetID       int    `json:"pet_id"`
	EventType   string `json:"event_type"`
	EventDate   string `json:"event_date,omitempty"` // Опционально, по умолчанию текущая дата
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details,omitempty"`

	// Медицинские данные
	VaccineName    string `json:"vaccine_name,omitempty"`
	VaccineBatch   string `json:"vaccine_batch,omitempty"`
	MedicationName string `json:"medication_name,omitempty"`
	Dosage         string `json:"dosage,omitempty"`
	NextDate       string `json:"next_date,omitempty"`

	// Смена владельца
	PreviousOwnerID int    `json:"previous_owner_id,omitempty"`
	NewOwnerID      int    `json:"new_owner_id,omitempty"`
	TransferReason  string `json:"transfer_reason,omitempty"`

	// Потеря/находка
	Location      string `json:"location,omitempty"`
	Circumstances string `json:"circumstances,omitempty"`
	ContactPhone  string `json:"contact_phone,omitempty"`
	ContactName   string `json:"contact_name,omitempty"`

	// Смерть
	DeathReason              string `json:"death_reason,omitempty"`
	DeathConfirmedByClinicID int    `json:"death_confirmed_by_clinic_id,omitempty"`

	// Приют
	ShelterID        int    `json:"shelter_id,omitempty"`
	AdoptionContract string `json:"adoption_contract,omitempty"`
}
