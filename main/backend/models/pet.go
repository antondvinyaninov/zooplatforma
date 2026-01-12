package models

type Pet struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	Name           string `json:"name"`
	Species        string `json:"species"`
	Photo          string `json:"photo"`
	Status         string `json:"status"`
	City           string `json:"city,omitempty"`
	Region         string `json:"region,omitempty"`
	Urgent         bool   `json:"urgent"`
	ContactName    string `json:"contact_name,omitempty"`
	ContactPhone   string `json:"contact_phone,omitempty"`
	OrganizationID *int   `json:"organization_id,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type CreatePetRequest struct {
	Name    string `json:"name"`
	Species string `json:"species"`
	Photo   string `json:"photo"`
}

type UpdatePetRequest struct {
	Name    string `json:"name,omitempty"`
	Species string `json:"species,omitempty"`
	Photo   string `json:"photo,omitempty"`
}
