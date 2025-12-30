package models

import "time"

// ShelterAnimal представляет животное в приюте
type ShelterAnimal struct {
	ID              int       `json:"id"`
	OrganizationID  int       `json:"organization_id"`
	Name            string    `json:"name"`
	Species         string    `json:"species"`
	Breed           string    `json:"breed"`
	AgeYears        int       `json:"age_years"`
	AgeMonths       int       `json:"age_months"`
	Gender          string    `json:"gender"`
	Color           string    `json:"color"`
	Size            string    `json:"size"`
	Weight          float64   `json:"weight"`
	ChipNumber      string    `json:"chip_number"`
	Status          string    `json:"status"`
	ArrivalDate     string    `json:"arrival_date"`
	ArrivalReason   string    `json:"arrival_reason"`
	HealthStatus    string    `json:"health_status"`
	Vaccinated      bool      `json:"vaccinated"`
	Sterilized      bool      `json:"sterilized"`
	CharacterTraits string    `json:"character_traits"`
	SpecialNeeds    string    `json:"special_needs"`
	Photo           string    `json:"photo"`
	Photos          string    `json:"photos"`
	Story           string    `json:"story"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ShelterVolunteer представляет волонтера приюта
type ShelterVolunteer struct {
	ID               int       `json:"id"`
	OrganizationID   int       `json:"organization_id"`
	UserID           int       `json:"user_id"`
	Role             string    `json:"role"`
	Status           string    `json:"status"`
	JoinedDate       string    `json:"joined_date"`
	HoursContributed int       `json:"hours_contributed"`
	Specialization   string    `json:"specialization"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Joined data
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	UserPhone string `json:"user_phone,omitempty"`
}

// AdoptionRequest представляет заявку на пристройство
type AdoptionRequest struct {
	ID               int       `json:"id"`
	OrganizationID   int       `json:"organization_id"`
	AnimalID         int       `json:"animal_id"`
	UserID           int       `json:"user_id"`
	Status           string    `json:"status"`
	RequestDate      string    `json:"request_date"`
	DecisionDate     *string   `json:"decision_date"`
	DecisionBy       *int      `json:"decision_by"`
	ApplicantName    string    `json:"applicant_name"`
	ApplicantPhone   string    `json:"applicant_phone"`
	ApplicantEmail   string    `json:"applicant_email"`
	ApplicantAddress string    `json:"applicant_address"`
	LivingConditions string    `json:"living_conditions"`
	HasExperience    bool      `json:"has_experience"`
	HasOtherPets     bool      `json:"has_other_pets"`
	OtherPetsInfo    string    `json:"other_pets_info"`
	Motivation       string    `json:"motivation"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Joined data
	AnimalName    string `json:"animal_name,omitempty"`
	AnimalSpecies string `json:"animal_species,omitempty"`
	AnimalPhoto   string `json:"animal_photo,omitempty"`
}

// ShelterEvent представляет событие в приюте
type ShelterEvent struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	EventType      string    `json:"event_type"`
	EntityType     string    `json:"entity_type"`
	EntityID       int       `json:"entity_id"`
	UserID         *int      `json:"user_id"`
	Description    string    `json:"description"`
	Metadata       string    `json:"metadata"`
	CreatedAt      time.Time `json:"created_at"`
}

// Organization представляет организацию (приют)
type Organization struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	Logo        string    `json:"logo"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrganizationMember представляет членство пользователя в организации
type OrganizationMember struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	UserID         int       `json:"user_id"`
	Role           string    `json:"role"`
	JoinedAt       time.Time `json:"joined_at"`
}
