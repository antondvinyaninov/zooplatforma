package models

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Bio        string `json:"bio"`
	Phone      string `json:"phone"`
	Location   string `json:"location"`
	Avatar     string `json:"avatar"`
	CoverPhoto string `json:"cover_photo"`
	CreatedAt  string `json:"created_at"`
}

type UserResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Bio        string `json:"bio"`
	Phone      string `json:"phone"`
	Location   string `json:"location"`
	Avatar     string `json:"avatar"`
	CoverPhoto string `json:"cover_photo"`
	CreatedAt  string `json:"created_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
