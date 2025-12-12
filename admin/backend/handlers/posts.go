package handlers

import (
	"database"
	"net/http"
)

// AdminPostsHandler - список постов
func AdminPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT p.id, p.content, p.post_type, p.created_at, u.name, u.email
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
		LIMIT 100
	`)
	if err != nil {
		sendError(w, "Ошибка получения постов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var content, postType, createdAt, userName, userEmail string
		if err := rows.Scan(&id, &content, &postType, &createdAt, &userName, &userEmail); err != nil {
			continue
		}
		posts = append(posts, map[string]interface{}{
			"id":         id,
			"content":    content,
			"post_type":  postType,
			"created_at": createdAt,
			"user_name":  userName,
			"user_email": userEmail,
		})
	}

	sendSuccess(w, posts)
}

// AdminPostHandler - действия с постом
func AdminPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать удаление/скрытие поста
	sendSuccess(w, map[string]string{"message": "Not implemented yet"})
}
