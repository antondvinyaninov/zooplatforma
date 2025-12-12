package handlers

import (
	"database"
	"database/sql"
	"log"
	"net/http"
)

// AdminLogsHandler - системные логи
func AdminLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT l.id, l.level, l.category, l.action, l.user_id, l.target_type, l.target_id, 
		       l.message, l.details, l.ip_address, l.user_agent, l.created_at, u.name as user_name
		FROM system_logs l
		LEFT JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC
		LIMIT 200
	`

	rows, err := database.DB.Query(query)

	if err != nil {
		// Логируем ошибку
		sendError(w, "Ошибка получения логов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int
		var level, category, action string
		var message, details, ipAddress, userAgent, createdAt string
		var userID, targetID *int
		var targetType, userName *string

		// Используем sql.NullString для полей которые могут быть NULL
		var messageNull, detailsNull, ipNull, userAgentNull sql.NullString
		var targetTypeNull, userNameNull sql.NullString

		err := rows.Scan(
			&id, &level, &category, &action, &userID, &targetTypeNull, &targetID,
			&messageNull, &detailsNull, &ipNull, &userAgentNull, &createdAt, &userNameNull,
		)
		if err != nil {
			// Логируем ошибку сканирования
			log.Printf("Error scanning log row: %v", err)
			continue
		}

		// Преобразуем NullString в обычные строки
		if messageNull.Valid {
			message = messageNull.String
		}
		if detailsNull.Valid {
			details = detailsNull.String
		}
		if ipNull.Valid {
			ipAddress = ipNull.String
		}
		if userAgentNull.Valid {
			userAgent = userAgentNull.String
		}
		if targetTypeNull.Valid {
			targetType = &targetTypeNull.String
		}
		if userNameNull.Valid {
			userName = &userNameNull.String
		}

		log := map[string]interface{}{
			"id":         id,
			"level":      level,
			"category":   category,
			"action":     action,
			"message":    message,
			"details":    details,
			"ip_address": ipAddress,
			"user_agent": userAgent,
			"created_at": createdAt,
		}

		if userID != nil {
			log["user_id"] = *userID
		}
		if userName != nil {
			log["admin_name"] = *userName // Для совместимости с frontend
		}
		if targetType != nil {
			log["target_type"] = *targetType
		}
		if targetID != nil {
			log["target_id"] = *targetID
		}

		logs = append(logs, log)
	}

	// Если логов нет - возвращаем пустой массив
	if logs == nil {
		logs = []map[string]interface{}{}
	}

	sendSuccess(w, logs)
}
