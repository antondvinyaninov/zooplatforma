package handlers

import (
	"backend/models"
	"database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	// Проверяем, это запрос на лайк?
	if strings.HasSuffix(path, "/like") {
		postIDStr := strings.TrimSuffix(path, "/like")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			sendErrorResponse(w, "Неверный ID поста", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodPost:
			toggleLike(w, r, postID, userID)
		case http.MethodGet:
			getLikeStatus(w, r, postID, userID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Проверяем, это запрос на список лайкнувших?
	if strings.HasSuffix(path, "/likers") {
		postIDStr := strings.TrimSuffix(path, "/likers")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			sendErrorResponse(w, "Неверный ID поста", http.StatusBadRequest)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		getLikers(w, r, postID)
		return
	}

	// Обычная обработка поста
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

func PetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID питомца из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/pet/")
	petID, err := strconv.Atoi(path)
	if err != nil {
		sendErrorResponse(w, "Неверный ID питомца", http.StatusBadRequest)
		return
	}

	getPetPosts(w, r, petID)
}

// getAllPosts получает все посты для Feed
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста (может быть 0 для неавторизованных)
	userID, _ := r.Context().Value("userID").(int)

	query := `
		SELECT p.id, p.author_id, p.author_type, p.content, p.attached_pets, 
		       p.attachments, p.tags, p.status, p.scheduled_at, p.created_at, p.updated_at,
		       u.name, u.email, u.avatar,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		WHERE p.is_deleted = 0 AND p.status = 'published'
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
		post, err := scanPost(rows)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	// Загружаем опросы для всех постов
	posts = loadPollsForPosts(posts, userID)

	sendSuccessResponse(w, posts)
}

// getDrafts получает черновики пользователя
func getDrafts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT p.id, p.author_id, p.author_type, p.content, p.attached_pets, 
		       p.attachments, p.tags, p.status, p.scheduled_at, p.created_at, p.updated_at,
		       u.name, u.email, u.avatar,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		WHERE p.author_id = ? AND p.author_type = 'user' AND p.status = 'draft' AND p.is_deleted = 0
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения черновиков: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var drafts []models.Post
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		drafts = append(drafts, post)
	}

	if drafts == nil {
		drafts = []models.Post{}
	}

	// Загружаем опросы для всех черновиков
	drafts = loadPollsForPosts(drafts, userID)

	sendSuccessResponse(w, drafts)
}

// DraftsHandler - отдельный handler для черновиков
func DraftsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	getDrafts(w, r)
}

// getUserPosts получает посты конкретного пользователя (Wall)
func getUserPosts(w http.ResponseWriter, r *http.Request, userID int) {
	// Получаем текущего пользователя из контекста
	currentUserID, _ := r.Context().Value("userID").(int)

	query := `
		SELECT p.id, p.author_id, p.author_type, p.content, p.attached_pets, 
		       p.attachments, p.tags, p.status, p.scheduled_at, p.created_at, p.updated_at,
		       u.name, u.email, u.avatar,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		WHERE p.author_id = ? AND p.author_type = 'user' AND p.is_deleted = 0
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
		post, err := scanPost(rows)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	// Загружаем опросы для всех постов
	posts = loadPollsForPosts(posts, currentUserID)

	sendSuccessResponse(w, posts)
}

// getPetPosts получает посты, в которых упоминается питомец
func getPetPosts(w http.ResponseWriter, r *http.Request, petID int) {
	// Получаем текущего пользователя из контекста
	currentUserID, _ := r.Context().Value("userID").(int)

	query := `
		SELECT p.id, p.author_id, p.author_type, p.content, p.attached_pets, 
		       p.attachments, p.tags, p.status, p.scheduled_at, p.created_at, p.updated_at,
		       u.name, u.email, u.avatar,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		INNER JOIN post_pets pp ON p.id = pp.post_id
		WHERE pp.pet_id = ? AND p.is_deleted = 0 AND p.status = 'published'
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(query, petID)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения постов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	// Загружаем опросы для всех постов
	posts = loadPollsForPosts(posts, currentUserID)

	sendSuccessResponse(w, posts)
}

// createPost создаёт новый пост
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

	// Валидация: хотя бы одно поле должно быть заполнено
	if req.Content == "" && len(req.AttachedPets) == 0 && len(req.Attachments) == 0 && req.Poll == nil {
		sendErrorResponse(w, "Пост должен содержать текст, фото, прикреплённых питомцев или опрос", http.StatusBadRequest)
		return
	}

	// Сериализуем массивы в JSON
	attachedPetsJSON, _ := json.Marshal(req.AttachedPets)
	attachmentsJSON, _ := json.Marshal(req.Attachments)
	tagsJSON, _ := json.Marshal(req.Tags)

	// Определяем статус поста
	status := "published"
	if req.Status != "" {
		status = req.Status
	}

	var scheduledAt *string
	if req.ScheduledAt != nil {
		scheduledAt = req.ScheduledAt
	}

	query := `INSERT INTO posts (author_id, author_type, content, attached_pets, attachments, tags, status, scheduled_at) 
	          VALUES (?, 'user', ?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, userID, req.Content, string(attachedPetsJSON), string(attachmentsJSON), string(tagsJSON), status, scheduledAt)
	if err != nil {
		sendErrorResponse(w, "Ошибка создания поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	postID, _ := result.LastInsertId()

	// Добавляем связи в post_pets для быстрых запросов
	for _, petID := range req.AttachedPets {
		_, err := database.DB.Exec("INSERT INTO post_pets (post_id, pet_id) VALUES (?, ?)", postID, petID)
		if err != nil {
			// Логируем ошибку, но не прерываем создание поста
			continue
		}
	}

	// Создаем опрос, если он есть
	if req.Poll != nil {
		err := createPollForPost(int(postID), req.Poll)
		if err != nil {
			// Логируем ошибку, но не прерываем создание поста
		}
	}

	// Получаем созданный пост
	post, err := getPostByID(int(postID), userID)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения поста", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, post)
}

// updatePost обновляет существующий пост
func updatePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Проверяем, что пост принадлежит пользователю
	var authorID int
	var authorType string
	err := database.DB.QueryRow("SELECT author_id, author_type FROM posts WHERE id = ? AND is_deleted = 0", postID).Scan(&authorID, &authorType)
	if err != nil {
		sendErrorResponse(w, "Пост не найден", http.StatusNotFound)
		return
	}

	if authorID != userID || authorType != "user" {
		sendErrorResponse(w, "Нет прав на редактирование этого поста", http.StatusForbidden)
		return
	}

	var req models.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Сериализуем массивы в JSON
	attachedPetsJSON, _ := json.Marshal(req.AttachedPets)
	attachmentsJSON, _ := json.Marshal(req.Attachments)
	tagsJSON, _ := json.Marshal(req.Tags)

	query := `UPDATE posts SET content = ?, attached_pets = ?, attachments = ?, tags = ?, updated_at = ? WHERE id = ?`
	_, err = database.DB.Exec(query, req.Content, string(attachedPetsJSON), string(attachmentsJSON), string(tagsJSON), time.Now().Format("2006-01-02 15:04:05"), postID)
	if err != nil {
		sendErrorResponse(w, "Ошибка обновления поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновляем связи в post_pets
	database.DB.Exec("DELETE FROM post_pets WHERE post_id = ?", postID)
	for _, petID := range req.AttachedPets {
		database.DB.Exec("INSERT INTO post_pets (post_id, pet_id) VALUES (?, ?)", postID, petID)
	}

	// Получаем обновлённый пост
	post, err := getPostByID(postID, userID)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения поста", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, post)
}

// deletePost удаляет пост (мягкое удаление)
func deletePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		sendErrorResponse(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Проверяем, что пост принадлежит пользователю
	var authorID int
	var authorType string
	err := database.DB.QueryRow("SELECT author_id, author_type FROM posts WHERE id = ? AND is_deleted = 0", postID).Scan(&authorID, &authorType)
	if err != nil {
		sendErrorResponse(w, "Пост не найден", http.StatusNotFound)
		return
	}

	if authorID != userID || authorType != "user" {
		sendErrorResponse(w, "Нет прав на удаление этого поста", http.StatusForbidden)
		return
	}

	// Мягкое удаление
	_, err = database.DB.Exec("UPDATE posts SET is_deleted = 1 WHERE id = ?", postID)
	if err != nil {
		sendErrorResponse(w, "Ошибка удаления поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]string{"message": "Пост удален"})
}

// scanPost сканирует строку БД в структуру Post
func scanPost(rows interface {
	Scan(dest ...interface{}) error
}) (models.Post, error) {
	var post models.Post
	var user models.User
	var attachedPetsJSON, attachmentsJSON, tagsJSON string
	var userName, userEmail, userAvatar *string

	err := rows.Scan(
		&post.ID, &post.AuthorID, &post.AuthorType, &post.Content,
		&attachedPetsJSON, &attachmentsJSON, &tagsJSON,
		&post.Status, &post.ScheduledAt,
		&post.CreatedAt, &post.UpdatedAt,
		&userName, &userEmail, &userAvatar,
		&post.CommentsCount,
	)
	if err != nil {
		return post, err
	}

	// Десериализуем JSON массивы
	json.Unmarshal([]byte(attachedPetsJSON), &post.AttachedPets)
	json.Unmarshal([]byte(attachmentsJSON), &post.Attachments)
	json.Unmarshal([]byte(tagsJSON), &post.Tags)

	// Инициализируем пустые массивы если nil
	if post.AttachedPets == nil {
		post.AttachedPets = []int{}
	}
	if post.Attachments == nil {
		post.Attachments = []models.Attachment{}
	}
	if post.Tags == nil {
		post.Tags = []string{}
	}

	// Добавляем данные автора если это user
	if post.AuthorType == "user" && userName != nil {
		user.ID = post.AuthorID
		user.Name = *userName
		if userEmail != nil {
			user.Email = *userEmail
		}
		if userAvatar != nil {
			user.Avatar = *userAvatar
		}
		post.User = &user
	}

	// Загружаем данные прикреплённых питомцев
	if len(post.AttachedPets) > 0 {
		post.Pets = loadPetsForPost(post.AttachedPets)
	}

	return post, nil
}

// getPostByID получает пост по ID
func getPostByID(postID int, userID int) (models.Post, error) {
	query := `
		SELECT p.id, p.author_id, p.author_type, p.content, p.attached_pets, 
		       p.attachments, p.tags, p.status, p.scheduled_at, p.created_at, p.updated_at,
		       u.name, u.email, u.avatar,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		WHERE p.id = ? AND p.is_deleted = 0
	`

	row := database.DB.QueryRow(query, postID)
	post, err := scanPost(row)
	if err != nil {
		return post, err
	}

	// Загружаем опрос, если есть
	poll, err := loadPollForPost(postID, userID)
	if err == nil {
		post.Poll = poll
	}

	return post, nil
}

// loadPetsForPost загружает данные питомцев для поста
func loadPetsForPost(petIDs []int) []models.Pet {
	if len(petIDs) == 0 {
		return []models.Pet{}
	}

	// Создаём плейсхолдеры для IN запроса
	placeholders := strings.Repeat("?,", len(petIDs)-1) + "?"
	query := `SELECT id, user_id, name, species, photo, created_at FROM pets WHERE id IN (` + placeholders + `)`

	// Конвертируем []int в []interface{} для Exec
	args := make([]interface{}, len(petIDs))
	for i, id := range petIDs {
		args[i] = id
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return []models.Pet{}
	}
	defer rows.Close()

	var pets []models.Pet
	for rows.Next() {
		var pet models.Pet
		rows.Scan(&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Photo, &pet.CreatedAt)
		pets = append(pets, pet)
	}

	return pets
}

// loadPollsForPosts загружает опросы для списка постов
func loadPollsForPosts(posts []models.Post, userID int) []models.Post {
	for i := range posts {
		poll, err := loadPollForPost(posts[i].ID, userID)
		if err == nil {
			posts[i].Poll = poll
		}
	}
	return posts
}
