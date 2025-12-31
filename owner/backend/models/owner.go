package models

import "time"

// Pet представляет питомца владельца
type Pet struct {
	ID         int       `json:"id"`
	OwnerID    int       `json:"owner_id"`
	Name       string    `json:"name"`
	Species    string    `json:"species"`
	Breed      string    `json:"breed"`
	BirthDate  string    `json:"birth_date"`
	Gender     string    `json:"gender"`
	Color      string    `json:"color"`
	ChipNumber string    `json:"chip_number"`
	Photo      string    `json:"photo"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PetEvent представляет событие в жизни питомца
type PetEvent struct {
	ID          int       `json:"id"`
	PetID       int       `json:"pet_id"`
	EventType   string    `json:"event_type"`
	EventDate   string    `json:"event_date"`
	Description string    `json:"description"`
	ClinicID    *int      `json:"clinic_id"`
	CreatedAt   time.Time `json:"created_at"`
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
