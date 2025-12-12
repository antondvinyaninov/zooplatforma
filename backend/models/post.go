package models

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	PostType  string `json:"post_type"`
	CreatedAt string `json:"created_at"`
	User      *User  `json:"user,omitempty"`
}

type CreatePostRequest struct {
	Content  string `json:"content"`
	PostType string `json:"post_type"`
}

type UpdatePostRequest struct {
	Content  string `json:"content,omitempty"`
	PostType string `json:"post_type,omitempty"`
}
