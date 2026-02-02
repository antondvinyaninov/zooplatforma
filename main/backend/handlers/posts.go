package handlers

import (
	"backend/models"
	"database"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET публичный - можно просматривать ленту без авторизации
		getAllPosts(w, r)
	case http.MethodPost:
		// POST требует авторизации - проверка внутри createPost
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
	case http.MethodGet:
		getPost(w, r, id)
	case http.MethodPut:
		updatePost(w, r, id)
	case http.MethodDelete:
		deletePost(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UserPostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("📥 UserPostsHandler: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID пользователя из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/user/")
	log.Printf("🔍 UserPostsHandler: Extracted path: %s", path)

	userID, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("❌ UserPostsHandler: Invalid user ID: %s, error: %v", path, err)
		sendErrorResponse(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	log.Printf("✅ UserPostsHandler: Calling getUserPosts for userID=%d", userID)
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

func OrganizationPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID организации из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/organization/")
	orgID, err := strconv.Atoi(path)
	if err != nil {
		sendErrorResponse(w, "Неверный ID организации", http.StatusBadRequest)
		return
	}

	getOrganizationPosts(w, r, orgID)
}

// checkCanEditPost проверяет может ли пользователь редактировать пост
func checkCanEditPost(userID int, post *models.Post) bool {
	if userID == 0 {
		log.Printf("🔒 checkCanEditPost: userID=0, can_edit=false")
		return false
	}

	// Если пост от пользователя - проверяем ID
	if post.AuthorType == "user" && post.AuthorID == userID {
		log.Printf("✅ checkCanEditPost: post %d by user %d, can_edit=true", post.ID, userID)
		return true
	}

	// Если пост от организации - проверяем членство с правами
	if post.AuthorType == "organization" {
		var role string
		err := database.DB.QueryRow(`
			SELECT role FROM organization_members 
			WHERE organization_id = ? AND user_id = ?
		`, post.AuthorID, userID).Scan(&role)

		if err == nil && (role == "owner" || role == "admin" || role == "moderator") {
			log.Printf("✅ checkCanEditPost: post %d by org %d, user %d has role %s, can_edit=true", post.ID, post.AuthorID, userID, role)
			return true
		}
		log.Printf("🔒 checkCanEditPost: post %d by org %d, user %d has no rights, can_edit=false", post.ID, post.AuthorID, userID)
	}

	log.Printf("🔒 checkCanEditPost: post %d, userID %d, can_edit=false", post.ID, userID)
	return false
}

// getAllPosts получает все посты для Feed
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста (может быть 0 для неавторизованных)
	userID, _ := r.Context().Value("userID").(int)

	// Получаем параметр фильтра
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "for-you"
	}

	log.Printf("🔍 getAllPosts: userID=%d, filter=%s", userID, filter)

	// Получаем параметры пагинации
	limitStr := r.URL.Query().Get("limit")
	limit := 20 // По умолчанию 20 постов
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	// Получаем город пользователя для фильтра "city"
	var userCity string
	if filter == "city" && userID > 0 {
		database.DB.QueryRow(ConvertPlaceholders("SELECT location FROM users WHERE id = ?"), userID).Scan(&userCity)
		log.Printf("🏙️ User city: %s", userCity)
	}

	// Базовый запрос
	query := `
		SELECT p.id, p.author_id, p.author_type, p.content, p.attached_pets, 
		       p.attachments, p.tags, p.status, p.scheduled_at, p.created_at, p.updated_at,
		       o.name as org_name, o.short_name as org_short_name, o.logo as org_logo,
		       u.name as user_name, u.last_name as user_last_name, u.avatar as user_avatar,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count,
		       CASE 
		           WHEN p.author_type = 'user' AND EXISTS (
		               SELECT 1 FROM friendships f 
		               WHERE ((f.user_id = ? AND f.friend_id = p.author_id) 
		                   OR (f.friend_id = ? AND f.user_id = p.author_id))
		                   AND f.status = 'accepted'
		           ) THEN 1
		           ELSE 0
		       END as is_friend
		FROM posts p
		LEFT JOIN organizations o ON p.author_id = o.id AND p.author_type = 'organization'
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		WHERE p.is_deleted = FALSE AND p.status = 'published'
	`

	// Добавляем фильтры в зависимости от типа
	args := []interface{}{userID, userID}

	switch filter {
	case "following":
		// Только посты от друзей (не свои)
		if userID > 0 {
			query += ` AND p.author_type = 'user' AND p.author_id != ? AND EXISTS (
				SELECT 1 FROM friendships f 
				WHERE ((f.user_id = ? AND f.friend_id = p.author_id) 
					OR (f.friend_id = ? AND f.user_id = p.author_id))
					AND f.status = 'accepted'
			)`
			args = append(args, userID, userID, userID)
			log.Printf("🔍 Following filter: excluding userID=%d, checking friendships", userID)
		}
	case "city":
		// Только посты из города пользователя
		if userCity != "" {
			query += ` AND (
				(p.author_type = 'user' AND u.location = ?) OR
				(p.author_type = 'organization' AND o.address_city = ?)
			)`
			args = append(args, userCity, userCity)
			log.Printf("🏙️ City filter: filtering by city=%s", userCity)
		}
	}

	query += ` ORDER BY is_friend DESC, p.created_at DESC LIMIT ?`
	args = append(args, limit)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		sendErrorResponse(w, "Ошибка получения постов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var attachedPetsJSON, attachmentsJSON, tagsJSON string
		var orgName, orgShortName, orgLogo *string
		var userName, userLastName, userAvatar *string
		var isFriend int

		err := rows.Scan(
			&post.ID, &post.AuthorID, &post.AuthorType, &post.Content,
			&attachedPetsJSON, &attachmentsJSON, &tagsJSON,
			&post.Status, &post.ScheduledAt,
			&post.CreatedAt, &post.UpdatedAt,
			&orgName, &orgShortName, &orgLogo,
			&userName, &userLastName, &userAvatar,
			&post.CommentsCount,
			&isFriend,
		)
		if err != nil {
			sendErrorResponse(w, "Ошибка чтения данных: "+err.Error(), http.StatusInternalServerError)
			return
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

		// Добавляем данные организации если это organization
		if post.AuthorType == "organization" && orgName != nil {
			org := models.Organization{
				ID:        post.AuthorID,
				Name:      *orgName,
				ShortName: orgShortName,
				Logo:      orgLogo,
			}
			post.Organization = &org
		}

		// Добавляем данные пользователя если это user
		if post.AuthorType == "user" && userName != nil {
			user := models.User{
				ID:   post.AuthorID,
				Name: *userName,
			}
			if userLastName != nil {
				user.LastName = *userLastName // Если LastName это *string в модели
			}
			if userAvatar != nil {
				user.Avatar = *userAvatar // Если Avatar это *string в модели
			}
			post.User = &user
		}

		posts = append(posts, post)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	// ✅ ОПТИМИЗАЦИЯ: Загружаем питомцев одним запросом для всех постов
	posts = loadPetsForPostsBatch(posts)

	// ✅ ОПТИМИЗАЦИЯ: Загружаем опросы одним запросом для всех постов
	includePolls := r.URL.Query().Get("include_polls")
	if includePolls == "true" {
		posts = loadPollsForPostsBatch(posts, userID)
	}

	// ✅ Проверяем права на редактирование для каждого поста
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(userID, &posts[i])
	}

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
		       o.name as org_name, o.short_name as org_short_name, o.logo as org_logo,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		LEFT JOIN organizations o ON p.author_id = o.id AND p.author_type = 'organization'
		WHERE p.author_id = ? AND p.author_type = 'user' AND p.status = 'draft' AND p.is_deleted = FALSE
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

	// ✅ Проверяем права на редактирование для каждого черновика
	for i := range drafts {
		drafts[i].CanEdit = checkCanEditPost(userID, &drafts[i])
	}

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
	log.Printf("🔍 getUserPosts: Starting for userID=%d", userID)

	// Получаем текущего пользователя из контекста
	currentUserID, _ := r.Context().Value("userID").(int)
	log.Printf("🔍 getUserPosts: currentUserID=%d", currentUserID)

	// Простой запрос только ID постов (без JOIN)
	log.Printf("🔍 getUserPosts: Fetching post IDs...")

	// Получаем параметры пагинации из query
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // По умолчанию 20 постов
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	offset := 0 // По умолчанию с начала
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	log.Printf("🔍 getUserPosts: Pagination - limit=%d, offset=%d", limit, offset)

	simpleQuery := `SELECT id FROM posts WHERE author_id = ? AND author_type = 'user' AND is_deleted = FALSE ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := database.DB.Query(simpleQuery, userID, limit, offset)
	if err != nil {
		log.Printf("❌ getUserPosts: Query error: %v", err)
		sendErrorResponse(w, "Ошибка получения постов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var postIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("❌ getUserPosts: Scan error: %v", err)
			continue
		}
		postIDs = append(postIDs, id)
	}
	log.Printf("✅ getUserPosts: Found %d post IDs", len(postIDs))

	// Загружаем полные данные для каждого поста отдельным запросом
	var posts []models.Post
	for _, postID := range postIDs {
		// Простой запрос без JOIN
		query := `
			SELECT id, author_id, author_type, content, attached_pets, 
			       attachments, tags, status, scheduled_at, created_at, updated_at
			FROM posts
			WHERE id = ?
		`
		var post models.Post
		var attachedPetsJSON, attachmentsJSON, tagsJSON sql.NullString
		var scheduledAt sql.NullTime

		err := database.DB.QueryRow(query, postID).Scan(
			&post.ID, &post.AuthorID, &post.AuthorType, &post.Content,
			&attachedPetsJSON, &attachmentsJSON, &tagsJSON,
			&post.Status, &scheduledAt, &post.CreatedAt, &post.UpdatedAt,
		)

		if err != nil {
			log.Printf("⚠️ getUserPosts: Failed to load post %d: %v", postID, err)
			continue
		}

		// Парсим JSON поля
		if attachedPetsJSON.Valid && attachedPetsJSON.String != "" && attachedPetsJSON.String != "null" {
			json.Unmarshal([]byte(attachedPetsJSON.String), &post.AttachedPets)
		}
		if post.AttachedPets == nil {
			post.AttachedPets = []int{}
		}

		if attachmentsJSON.Valid && attachmentsJSON.String != "" && attachmentsJSON.String != "null" {
			json.Unmarshal([]byte(attachmentsJSON.String), &post.Attachments)
		}
		if post.Attachments == nil {
			post.Attachments = []models.Attachment{}
		}

		if tagsJSON.Valid && tagsJSON.String != "" && tagsJSON.String != "null" {
			json.Unmarshal([]byte(tagsJSON.String), &post.Tags)
		}
		if post.Tags == nil {
			post.Tags = []string{}
		}

		if scheduledAt.Valid {
			timeStr := scheduledAt.Time.Format(time.RFC3339)
			post.ScheduledAt = &timeStr
		}

		// User и Organization будут nil (можно загрузить позже если нужно)
		post.User = nil
		post.Organization = nil
		post.Pets = []models.Pet{}
		post.CommentsCount = 0

		posts = append(posts, post)
	}
	log.Printf("✅ getUserPosts: Loaded %d posts", len(posts))

	if posts == nil {
		posts = []models.Post{}
	}

	// Загружаем опросы для всех постов
	log.Printf("🔍 getUserPosts: Loading polls...")
	posts = loadPollsForPosts(posts, currentUserID)
	log.Printf("✅ getUserPosts: Polls loaded")

	// Проверяем права на редактирование для каждого поста
	log.Printf("🔍 getUserPosts: Checking edit permissions...")
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(currentUserID, &posts[i])
	}
	log.Printf("✅ getUserPosts: Edit permissions checked")

	log.Printf("✅ getUserPosts: Sending response with %d posts", len(posts))
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
		       o.name as org_name, o.short_name as org_short_name, o.logo as org_logo,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		LEFT JOIN organizations o ON p.author_id = o.id AND p.author_type = 'organization'
		INNER JOIN post_pets pp ON p.id = pp.post_id
		WHERE pp.pet_id = ? AND p.is_deleted = FALSE AND p.status = 'published'
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

	// ✅ Проверяем права на редактирование для каждого поста
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(currentUserID, &posts[i])
	}

	sendSuccessResponse(w, posts)
}

// getOrganizationPosts получает посты организации
func getOrganizationPosts(w http.ResponseWriter, r *http.Request, orgID int) {
	// Получаем текущего пользователя из контекста
	currentUserID, _ := r.Context().Value("userID").(int)

	query := `
		SELECT p.id, p.author_id, p.author_type, p.content, p.attached_pets, 
		       p.attachments, p.tags, p.status, p.scheduled_at, p.created_at, p.updated_at,
		       u.name, u.email, u.avatar,
		       o.name as org_name, o.short_name as org_short_name, o.logo as org_logo,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id AND p.author_type = 'user'
		LEFT JOIN organizations o ON p.author_id = o.id AND p.author_type = 'organization'
		WHERE p.author_id = ? AND p.author_type = 'organization' AND p.is_deleted = FALSE AND p.status = 'published'
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(query, orgID)
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

	// ✅ Проверяем права на редактирование для каждого поста
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(currentUserID, &posts[i])
	}

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

	// Определяем автора поста
	authorType := "user"
	authorID := userID
	if req.AuthorType == "organization" && req.OrganizationID != nil {
		// Проверяем права пользователя на публикацию от имени организации
		var canPost bool
		err := database.DB.QueryRow(`
			SELECT can_post FROM organization_members
			WHERE organization_id = ? AND user_id = ?
		`, *req.OrganizationID, userID).Scan(&canPost)

		if err != nil || !canPost {
			sendErrorResponse(w, "Нет прав на публикацию от имени этой организации", http.StatusForbidden)
			return
		}

		authorType = "organization"
		authorID = *req.OrganizationID
	}

	query := `INSERT INTO posts (author_id, author_type, content, attached_pets, attachments, tags, status, scheduled_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, authorID, authorType, req.Content, string(attachedPetsJSON), string(attachmentsJSON), string(tagsJSON), status, scheduledAt)
	if err != nil {
		sendErrorResponse(w, "Ошибка создания поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	postID, _ := result.LastInsertId()

	// Добавляем связи в post_pets для быстрых запросов
	for _, petID := range req.AttachedPets {
		_, err := database.DB.Exec(ConvertPlaceholders("INSERT INTO post_pets (post_id, pet_id) VALUES (?, ?)"), postID, petID)
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

	// Логируем создание поста
	ipAddress := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	details := "Создан пост"
	if authorType == "organization" {
		details = "Создан пост от имени организации"
	}
	CreateUserLog(database.DB, userID, "post_create", details, ipAddress, userAgent)

	sendSuccessResponse(w, post)
}

// getPost возвращает один пост по ID
func getPost(w http.ResponseWriter, r *http.Request, postID int) {
	// Получаем userID если пользователь авторизован (опционально)
	userID := 0
	if uid, ok := r.Context().Value("userID").(int); ok {
		userID = uid
	}

	post, err := getPostByID(postID, userID)
	if err != nil {
		sendErrorResponse(w, "Пост не найден", http.StatusNotFound)
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

	// Получаем пост для проверки прав
	post, err := getPostByID(postID, userID)
	if err != nil {
		sendErrorResponse(w, "Пост не найден", http.StatusNotFound)
		return
	}

	// Проверяем права на редактирование
	if !checkCanEditPost(userID, &post) {
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
	database.DB.Exec(ConvertPlaceholders("DELETE FROM post_pets WHERE post_id = ?"), postID)
	for _, petID := range req.AttachedPets {
		database.DB.Exec(ConvertPlaceholders("INSERT INTO post_pets (post_id, pet_id) VALUES (?, ?)"), postID, petID)
	}

	// Получаем обновлённый пост
	post, err = getPostByID(postID, userID)
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

	// Получаем пост для проверки прав
	post, err := getPostByID(postID, userID)
	if err != nil {
		sendErrorResponse(w, "Пост не найден", http.StatusNotFound)
		return
	}

	// Проверяем права на удаление
	if !checkCanEditPost(userID, &post) {
		sendErrorResponse(w, "Нет прав на удаление этого поста", http.StatusForbidden)
		return
	}

	// Мягкое удаление
	_, err = database.DB.Exec(ConvertPlaceholders("UPDATE posts SET is_deleted = TRUE WHERE id = ?"), postID)
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
	var orgName, orgShortName, orgLogo *string

	err := rows.Scan(
		&post.ID, &post.AuthorID, &post.AuthorType, &post.Content,
		&attachedPetsJSON, &attachmentsJSON, &tagsJSON,
		&post.Status, &post.ScheduledAt,
		&post.CreatedAt, &post.UpdatedAt,
		&userName, &userEmail, &userAvatar,
		&orgName, &orgShortName, &orgLogo,
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

	// Добавляем данные организации если это organization
	if post.AuthorType == "organization" && orgName != nil {
		org := models.Organization{
			ID:        post.AuthorID,
			Name:      *orgName,
			ShortName: orgShortName,
			Logo:      orgLogo,
		}
		post.Organization = &org
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
		       o.name as org_name, o.short_name as org_short_name, o.logo as org_logo,
		       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comments_count
		FROM posts p
		LEFT JOIN organizations o ON p.author_id = o.id AND p.author_type = 'organization'
		WHERE p.id = ? AND p.is_deleted = FALSE
	`

	var post models.Post
	var orgName, orgShortName, orgLogo sql.NullString
	var attachedPetsJSON, attachmentsJSON, tagsJSON, scheduledAt sql.NullString

	err := database.DB.QueryRow(query, postID).Scan(
		&post.ID, &post.AuthorID, &post.AuthorType, &post.Content,
		&attachedPetsJSON, &attachmentsJSON, &tagsJSON,
		&post.Status, &scheduledAt, &post.CreatedAt, &post.UpdatedAt,
		&orgName, &orgShortName, &orgLogo,
		&post.CommentsCount,
	)

	if err != nil {
		return post, err
	}

	// Парсим JSON поля
	if attachedPetsJSON.Valid && attachedPetsJSON.String != "" {
		json.Unmarshal([]byte(attachedPetsJSON.String), &post.AttachedPets)
	}
	if attachmentsJSON.Valid && attachmentsJSON.String != "" {
		json.Unmarshal([]byte(attachmentsJSON.String), &post.Attachments)
	}
	if tagsJSON.Valid && tagsJSON.String != "" {
		json.Unmarshal([]byte(tagsJSON.String), &post.Tags)
	}
	if scheduledAt.Valid {
		post.ScheduledAt = &scheduledAt.String
	}

	// Загружаем данные автора
	switch post.AuthorType {
	case "user":
		// 🔥 Загружаем данные пользователя через Auth Service
		authServiceURL := os.Getenv("AUTH_SERVICE_URL")
		if authServiceURL == "" {
			authServiceURL = "http://localhost:7100"
		}

		resp, err := http.Get(fmt.Sprintf("%s/api/users/%d", authServiceURL, post.AuthorID))
		if err != nil {
			log.Printf("❌ Failed to fetch user from Auth Service: %v", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var authResp struct {
					Success bool        `json:"success"`
					Data    models.User `json:"data"`
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("❌ Failed to read Auth Service response: %v", err)
				} else {
					if err := json.Unmarshal(body, &authResp); err != nil {
						log.Printf("❌ Failed to unmarshal Auth Service response: %v", err)
					} else if authResp.Success {
						post.User = &authResp.Data
						log.Printf("✅ Loaded user data from Auth Service: id=%d, name=%s, avatar=%s", authResp.Data.ID, authResp.Data.Name, authResp.Data.Avatar)
					} else {
						log.Printf("⚠️ Auth Service returned success=false")
					}
				}
			} else {
				log.Printf("⚠️ Auth Service returned status %d", resp.StatusCode)
			}
		}
	case "organization":
		// Данные организации уже загружены из JOIN
		if orgName.Valid {
			shortName := orgShortName.String
			logo := orgLogo.String
			post.Organization = &models.Organization{
				ID:        post.AuthorID,
				Name:      orgName.String,
				ShortName: &shortName,
				Logo:      &logo,
			}
		}
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
	query := `
		SELECT 
			p.id, p.user_id, p.name, p.species, p.breed, p.gender, p.birth_date, 
			p.color, p.size, p.photo, p.status, p.city, p.region, p.urgent, p.story,
			p.organization_id, o.name as organization_name, o.type as organization_type,
			p.created_at
		FROM pets p
		LEFT JOIN organizations o ON p.organization_id = o.id
		WHERE p.id IN (` + placeholders + `)
	`

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
		var organizationName, organizationType sql.NullString

		err := rows.Scan(
			&pet.ID, &pet.UserID, &pet.Name, &pet.Species, &pet.Breed, &pet.Gender, &pet.BirthDate,
			&pet.Color, &pet.Size, &pet.Photo, &pet.Status, &pet.City, &pet.Region, &pet.Urgent, &pet.Story,
			&pet.OrganizationID, &organizationName, &organizationType,
			&pet.CreatedAt,
		)

		if err != nil {
			continue
		}

		if organizationName.Valid {
			pet.OrganizationName = organizationName.String
		}
		if organizationType.Valid {
			pet.OrganizationType = organizationType.String
		}

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

// loadUsersForPostsBatch загружает данные пользователей для списка постов через Auth Service
func loadUsersForPostsBatch(posts []models.Post) []models.Post {
	if len(posts) == 0 {
		return posts
	}

	// Собираем уникальные ID пользователей
	userIDs := make(map[int]bool)
	for _, post := range posts {
		if post.AuthorType == "user" {
			userIDs[post.AuthorID] = true
		}
	}

	if len(userIDs) == 0 {
		return posts
	}

	// Загружаем данных пользователей через Auth Service
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:7100"
	}

	// Создаем map для быстрого доступа к данным пользователей
	usersMap := make(map[int]*models.User)

	for userID := range userIDs {
		resp, err := http.Get(fmt.Sprintf("%s/api/users/%d", authServiceURL, userID))
		if err != nil {
			log.Printf("❌ Failed to fetch user %d from Auth Service: %v", userID, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var authResp struct {
				Success bool        `json:"success"`
				Data    models.User `json:"data"`
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("❌ Failed to read Auth Service response for user %d: %v", userID, err)
				continue
			}

			if err := json.Unmarshal(body, &authResp); err != nil {
				log.Printf("❌ Failed to unmarshal Auth Service response for user %d: %v", userID, err)
				continue
			}

			if authResp.Success {
				usersMap[userID] = &authResp.Data
			}
		}
	}

	log.Printf("✅ Loaded %d users from Auth Service for %d posts", len(usersMap), len(posts))

	// Присваиваем данные пользователей постам
	for i := range posts {
		if posts[i].AuthorType == "user" {
			if user, ok := usersMap[posts[i].AuthorID]; ok {
				posts[i].User = user
			}
		}
	}

	return posts
}
