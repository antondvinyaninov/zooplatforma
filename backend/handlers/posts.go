package handlers

import (
	"backend/models"
	"database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllPosts(w, r)
	case http.MethodPost:
		createPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendErrorResponse(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updatePost(w, r, id)
	case http.MethodDelete:
		deletePost(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UserPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID пользователя из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/user/")
	userID, err := strconv.Atoi(path)
	if err != nil {
		sendErrorResponse(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	getUserPosts(w, r, userID)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT p.id, p.user_id, p.content, p.post_type, p.created_at,
		       u.name, u.email, u.avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения постов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var user models.User
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Content, &post.PostType, &post.CreatedAt,
			&user.Name, &user.Email, &user.Avatar,
		)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		user.ID = post.UserID
		post.User = &user
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	sendSuccessResponse(w, posts)
}

func getUserPosts(w http.ResponseWriter, r *http.Request, userID int) {
	query := `
		SELECT p.id, p.user_id, p.content, p.post_type, p.created_at,
		       u.name, u.email, u.avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения постов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var user models.User
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Content, &post.PostType, &post.CreatedAt,
			&user.Name, &user.Email, &user.Avatar,
		)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		user.ID = post.UserID
		post.User = &user
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	sendSuccessResponse(w, posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		sendErrorResponse(w, "Содержимое поста не может быть пустым", http.StatusBadRequest)
		return
	}

	if req.PostType == "" {
		req.PostType = "post"
	}

	query := `INSERT INTO posts (user_id, content, post_type) VALUES (?, ?, ?)`
	result, err := database.DB.Exec(query, userID, req.Content, req.PostType)
	if err != nil {
		sendErrorResponse(w, "Ошибка создания поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Получаем созданный пост с данными пользователя
	var post models.Post
	var user models.User
	query = `
		SELECT p.id, p.user_id, p.content, p.post_type, p.created_at,
		       u.name, u.email, u.avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`
	err = database.DB.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &post.Content, &post.PostType, &post.CreatedAt,
		&user.Name, &user.Email, &user.Avatar,
	)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения поста", http.StatusInternalServerError)
		return
	}

	user.ID = post.UserID
	post.User = &user

	sendSuccessResponse(w, post)
}

func updatePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Проверяем, что пост принадлежит пользователю
	var ownerID int
	err := database.DB.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&ownerID)
	if err != nil {
		sendErrorResponse(w, "Пост не найден", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		sendErrorResponse(w, "Нет прав на редактирование этого поста", http.StatusForbidden)
		return
	}

	var req models.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	query := `UPDATE posts SET content = ?, post_type = ? WHERE id = ?`
	_, err = database.DB.Exec(query, req.Content, req.PostType, postID)
	if err != nil {
		sendErrorResponse(w, "Ошибка обновления поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]string{"message": "Пост обновлен"})
}

func deletePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Проверяем, что пост принадлежит пользователю
	var ownerID int
	err := database.DB.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&ownerID)
	if err != nil {
		sendErrorResponse(w, "Пост не найден", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		sendErrorResponse(w, "Нет прав на удаление этого поста", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		sendErrorResponse(w, "Ошибка удаления поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]string{"message": "Пост удален"})
}
