package handlers

import (
	"database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Organization структура организации
type Organization struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	ShortName     *string   `json:"short_name"`
	Type          string    `json:"type"`
	INN           *string   `json:"inn"`
	Email         *string   `json:"email"`
	Phone         *string   `json:"phone"`
	AddressCity   *string   `json:"address_city"`
	AddressRegion *string   `json:"address_region"`
	IsVerified    bool      `json:"is_verified"`
	IsActive      bool      `json:"is_active"`
	Status        string    `json:"status"`
	OwnerUserID   int       `json:"owner_user_id"`
	OwnerName     *string   `json:"owner_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AdminOrganizationsHandler - список всех организаций для модерации
func AdminOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Параметры фильтрации
	status := r.URL.Query().Get("status") // all, verified, pending, rejected
	orgType := r.URL.Query().Get("type")  // shelter, vet_clinic, pet_shop, foundation, kennel, other
	search := r.URL.Query().Get("search") // поиск по названию, ИНН
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Значения по умолчанию
	pageNum := 1
	limitNum := 20

	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}

	offset := (pageNum - 1) * limitNum

	// Построение SQL запроса
	query := `
		SELECT 
			o.id, o.name, o.short_name, o.type, o.inn, o.email, o.phone,
			o.address_city, o.address_region, o.is_verified, o.is_active, o.status,
			o.owner_user_id, o.created_at, o.updated_at,
			u.name as owner_name
		FROM organizations o
		LEFT JOIN users u ON o.owner_user_id = u.id
		WHERE 1=1
	`

	args := []interface{}{}

	// Фильтр по статусу верификации
	if status == "verified" {
		query += " AND o.is_verified = 1"
	} else if status == "pending" {
		query += " AND o.is_verified = 0 AND o.status = 'active'"
	} else if status == "rejected" {
		query += " AND o.status = 'blocked'"
	}

	// Фильтр по типу организации
	if orgType != "" {
		query += " AND o.type = ?"
		args = append(args, orgType)
	}

	// Поиск по названию или ИНН
	if search != "" {
		query += " AND (o.name LIKE ? OR o.inn LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Сортировка: сначала непроверенные, потом по дате создания
	query += " ORDER BY o.is_verified ASC, o.created_at DESC"

	// Пагинация
	query += " LIMIT ? OFFSET ?"
	args = append(args, limitNum, offset)

	// Выполнение запроса
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying organizations: %v", err)
		http.Error(w, "Failed to fetch organizations", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	organizations := []Organization{}

	for rows.Next() {
		var org Organization
		err := rows.Scan(
			&org.ID, &org.Name, &org.ShortName, &org.Type, &org.INN, &org.Email, &org.Phone,
			&org.AddressCity, &org.AddressRegion, &org.IsVerified, &org.IsActive, &org.Status,
			&org.OwnerUserID, &org.CreatedAt, &org.UpdatedAt, &org.OwnerName,
		)
		if err != nil {
			log.Printf("Error scanning organization: %v", err)
			continue
		}
		organizations = append(organizations, org)
	}

	// Подсчет общего количества
	countQuery := `
		SELECT COUNT(*) FROM organizations o WHERE 1=1
	`
	countArgs := []interface{}{}

	if status == "verified" {
		countQuery += " AND o.is_verified = 1"
	} else if status == "pending" {
		countQuery += " AND o.is_verified = 0 AND o.status = 'active'"
	} else if status == "rejected" {
		countQuery += " AND o.status = 'blocked'"
	}

	if orgType != "" {
		countQuery += " AND o.type = ?"
		countArgs = append(countArgs, orgType)
	}

	if search != "" {
		countQuery += " AND (o.name LIKE ? OR o.inn LIKE ?)"
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern)
	}

	var total int
	err = database.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		log.Printf("Error counting organizations: %v", err)
		total = len(organizations)
	}

	// Ответ
	response := map[string]interface{}{
		"organizations": organizations,
		"pagination": map[string]interface{}{
			"page":       pageNum,
			"limit":      limitNum,
			"total":      total,
			"totalPages": (total + limitNum - 1) / limitNum,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdminVerifyOrganizationHandler - верификация/отклонение организации
func AdminVerifyOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID организации из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/organizations/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "verify" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	orgID, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	// Парсим тело запроса
	var req struct {
		Action string `json:"action"` // verify, reject
		Reason string `json:"reason"` // причина отклонения (опционально)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем ID админа из контекста
	adminID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Выполняем действие
	var query string
	var args []interface{}

	if req.Action == "verify" {
		query = `UPDATE organizations SET is_verified = 1, status = 'active', updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		args = []interface{}{orgID}
	} else if req.Action == "reject" {
		query = `UPDATE organizations SET is_verified = 0, status = 'blocked', updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		args = []interface{}{orgID}
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error updating organization: %v", err)
		http.Error(w, "Failed to update organization", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}

	// Логируем действие админа
	logAction := "verify_organization"
	if req.Action == "reject" {
		logAction = "reject_organization"
	}

	logQuery := `
		INSERT INTO admin_logs (admin_id, action, target_type, target_id, details, created_at)
		VALUES (?, ?, 'organization', ?, ?, CURRENT_TIMESTAMP)
	`
	details := fmt.Sprintf("Action: %s", req.Action)
	if req.Reason != "" {
		details += fmt.Sprintf(", Reason: %s", req.Reason)
	}

	_, err = database.DB.Exec(logQuery, adminID, logAction, orgID, details)
	if err != nil {
		log.Printf("Error logging admin action: %v", err)
	}

	// Ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Organization %s successfully", req.Action+"ed"),
	})
}

// AdminOrganizationStatsHandler - статистика организаций
func AdminOrganizationStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Общая статистика
	stats := map[string]interface{}{}

	// Всего организаций
	var total int
	database.DB.QueryRow("SELECT COUNT(*) FROM organizations").Scan(&total)
	stats["total"] = total

	// Верифицированных
	var verified int
	database.DB.QueryRow("SELECT COUNT(*) FROM organizations WHERE is_verified = 1").Scan(&verified)
	stats["verified"] = verified

	// Ожидают проверки
	var pending int
	database.DB.QueryRow("SELECT COUNT(*) FROM organizations WHERE is_verified = 0 AND status = 'active'").Scan(&pending)
	stats["pending"] = pending

	// Отклоненных
	var rejected int
	database.DB.QueryRow("SELECT COUNT(*) FROM organizations WHERE status = 'blocked'").Scan(&rejected)
	stats["rejected"] = rejected

	// По типам
	typeStats := []map[string]interface{}{}
	rows, err := database.DB.Query(`
		SELECT type, COUNT(*) as count 
		FROM organizations 
		GROUP BY type 
		ORDER BY count DESC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var orgType string
			var count int
			if err := rows.Scan(&orgType, &count); err == nil {
				typeStats = append(typeStats, map[string]interface{}{
					"type":  orgType,
					"count": count,
				})
			}
		}
	}
	stats["by_type"] = typeStats

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
