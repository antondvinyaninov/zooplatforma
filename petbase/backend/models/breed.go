package models

import "time"

type Breed struct {
	ID          int       `json:"id"`
	SpeciesID   int       `json:"species_id"`
	SpeciesName string    `json:"species_name,omitempty"`
	Name        string    `json:"name"`
	NameEn      string    `json:"name_en"`
	Description string    `json:"description"`
	Origin      string    `json:"origin"`
	Size        string    `json:"size"`
	WeightMin   float64   `json:"weight_min"`
	WeightMax   float64   `json:"weight_max"`
	LifespanMin int       `json:"lifespan_min"`
	LifespanMax int       `json:"lifespan_max"`
	Temperament string    `json:"temperament"`
	CareLevel   string    `json:"care_level"`
	Photo       string    `json:"photo"`
	CreatedAt   time.Time `json:"created_at"`
}
