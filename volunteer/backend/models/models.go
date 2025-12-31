package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type Pet struct {
	ID         int       `json:"id"`
	OwnerID    int       `json:"user_id"`
	Name       string    `json:"name"`
	Species    string    `json:"species"`
	Breed      *string   `json:"breed"`
	BirthDate  *string   `json:"birth_date"`
	Gender     *string   `json:"gender"`
	Color      *string   `json:"color"`
	ChipNumber *string   `json:"chip_number"`
	Photo      *string   `json:"photo"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Task struct {
	ID          int       `json:"id"`
	VolunteerID int       `json:"volunteer_id"`
	PetID       *int      `json:"pet_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     *string   `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
