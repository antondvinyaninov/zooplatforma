package models

import "time"

type Species struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	NameEn      string    `json:"name_en"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
}
