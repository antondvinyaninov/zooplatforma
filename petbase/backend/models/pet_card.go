package models

import "time"

type PetCard struct {
	ID              int       `json:"id"`
	BreedID         int       `json:"breed_id"`
	BreedName       string    `json:"breed_name,omitempty"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Characteristics string    `json:"characteristics"` // JSON
	CareTips        string    `json:"care_tips"`
	HealthInfo      string    `json:"health_info"`
	Nutrition       string    `json:"nutrition"`
	Photos          string    `json:"photos"` // JSON array
	IsPublished     bool      `json:"is_published"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
