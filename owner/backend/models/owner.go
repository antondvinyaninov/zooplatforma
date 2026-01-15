package models

import "time"

// Pet представляет питомца владельца (полная структура из PetID)
type Pet struct {
	ID                 int       `json:"id"`
	OwnerID            int       `json:"owner_id"`
	CuratorID          *int      `json:"curator_id,omitempty"`
	OrganizationID     *int      `json:"organization_id,omitempty"`
	Name               string    `json:"name"`
	Species            string    `json:"species"`
	Breed              *string   `json:"breed,omitempty"`
	Age                int       `json:"age"`
	Sex                string    `json:"sex"`
	Color              string    `json:"color"`
	Weight             *float64  `json:"weight,omitempty"`
	ChipNumber         *string   `json:"chip_number,omitempty"`
	TattooNumber       *string   `json:"tattoo_number,omitempty"`
	PassportNumber     *string   `json:"passport_number,omitempty"`
	PhotoURL           *string   `json:"photo_url,omitempty"`
	Status             string    `json:"status"`
	VerificationStatus string    `json:"verification_status"`
	Sterilized         bool      `json:"sterilized"`
	SterilizationDate  *string   `json:"sterilization_date,omitempty"`
	SpecialMarks       *string   `json:"special_marks,omitempty"`
	Character          *string   `json:"character,omitempty"`
	HealthNotes        *string   `json:"health_notes,omitempty"`
	Allergies          *string   `json:"allergies,omitempty"`
	ChronicConditions  *string   `json:"chronic_conditions,omitempty"`
	EmergencyContact   *string   `json:"emergency_contact,omitempty"`
	EmergencyPhone     *string   `json:"emergency_phone,omitempty"`
	InsuranceCompany   *string   `json:"insurance_company,omitempty"`
	InsurancePolicy    *string   `json:"insurance_policy,omitempty"`
	City               *string   `json:"city,omitempty"`
	Region             *string   `json:"region,omitempty"`
	Urgent             *bool     `json:"urgent,omitempty"`
	ContactPhone       *string   `json:"contact_phone,omitempty"`
	ContactName        *string   `json:"contact_name,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Treatment представляет обработку питомца (от блох, клещей, глистов)
type Treatment struct {
	ID         int        `json:"id"`
	PetID      int        `json:"pet_id"`
	UserID     int        `json:"user_id"`
	Date       time.Time  `json:"date"`
	Medication string     `json:"medication"`
	Dosage     string     `json:"dosage"`
	NextDate   *time.Time `json:"next_date,omitempty"`
	Notes      *string    `json:"notes,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// MedicalEvent представляет медицинское событие из PetID (от клиник)
type MedicalEvent struct {
	ID             int       `json:"id"`
	PetID          int       `json:"pet_id"`
	EventType      string    `json:"event_type"`
	UserID         int       `json:"user_id"`
	ClinicID       *int      `json:"clinic_id,omitempty"`
	OrganizationID *int      `json:"organization_id,omitempty"`
	EventDate      string    `json:"event_date"`
	Description    string    `json:"description"`
	Veterinarian   *string   `json:"veterinarian,omitempty"`
	NextVisitDate  *string   `json:"next_visit_date,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// StatusChangeRequest представляет запрос на изменение статуса питомца
type StatusChangeRequest struct {
	Status        string  `json:"status"`
	Date          string  `json:"date"`
	Location      *string `json:"location,omitempty"`
	Circumstances *string `json:"circumstances,omitempty"`
	ContactInfo   *string `json:"contact_info,omitempty"`
}

// TreatmentReminder представляет напоминание о предстоящей обработке
type TreatmentReminder struct {
	PetID       int    `json:"pet_id"`
	PetName     string `json:"pet_name"`
	TreatmentID int    `json:"treatment_id"`
	Medication  string `json:"medication"`
	NextDate    string `json:"next_date"`
	DaysUntil   int    `json:"days_until"`
	IsOverdue   bool   `json:"is_overdue"`
}

// User представляет пользователя
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
