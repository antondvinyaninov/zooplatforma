package models

import "time"

type Pet struct {
	ID         int       `json:"id"`
	OwnerID    int       `json:"owner_id"`
	Name       string    `json:"name"`
	Species    string    `json:"species"`
	Breed      *string   `json:"breed,omitempty"`
	BirthDate  *string   `json:"birth_date,omitempty"`
	Gender     *string   `json:"gender,omitempty"`
	Color      *string   `json:"color,omitempty"`
	ChipNumber *string   `json:"chip_number,omitempty"`
	Photo      *string   `json:"photo,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Appointment struct {
	ID        int       `json:"id"`
	PetID     int       `json:"pet_id"`
	ClinicID  int       `json:"clinic_id"`
	Date      time.Time `json:"date"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	LastName  *string   `json:"last_name,omitempty"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

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
