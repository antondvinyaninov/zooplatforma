package models

type Pet struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Species   string `json:"species"`
	Photo     string `json:"photo"`
	CreatedAt string `json:"created_at"`
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
