package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Report struct {
	ID               int     `json:"id"`
	ReporterID       int     `json:"reporter_id"`
	TargetType       string  `json:"target_type"`
	TargetID         int     `json:"target_id"`
	Reason           string  `json:"reason"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
	ModeratorID      *int    `json:"moderator_id"`
	ModeratorAction  *string `json:"moderator_action"`
	ModeratorComment *string `json:"moderator_comment"`
	ReviewedAt       *string `json:"reviewed_at"`
	CreatedAt        string  `json:"created_at"`
	ReporterName     string  `json:"reporter_name"`
	ReporterEmail    string  `json:"reporter_email"`
	ModeratorName    *string `json:"moderator_name"`
}

// GetReportsHandler - получить список жалоб
func GetReportsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		if status == "" {
			status = "pending"
		}

		query := `
			SELECT 
				r.id, r.reporter_id, r.target_type, r.target_id, r.reason, 
				r.description, r.status, r.moderator_id, r.moderator_action,
				r.moderator_comment, r.reviewed_at, r.created_at,
				u.name as reporter_name, u.email as reporter_email,
				m.name as moderator_name
			FROM reports r
			LEFT JOIN users u ON r.reporter_id = u.id
			LEFT JOIN users m ON r.moderator_id = m.id
			WHERE r.status = ?
			ORDER BY r.created_at DESC
			LIMIT 100
		`

		rows, err := db.Query(query, status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		reports := []Report{}
		for rows.Next() {
			var r Report
			err := rows.Scan(
				&r.ID, &r.ReporterID, &r.TargetType, &r.TargetID, &r.Reason,
				&r.Description, &r.Status, &r.ModeratorID, &r.ModeratorAction,
				&r.ModeratorComment, &r.ReviewedAt, &r.CreatedAt,
				&r.ReporterName, &r.ReporterEmail, &r.ModeratorName,
			)
			if err != nil {
				continue
			}
			reports = append(reports, r)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    reports,
		})
	}
}

// ReviewReportHandler - рассмотреть жалобу
func ReviewReportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moderatorID := r.Context().Value("userID").(int)

		pathParts := strings.Split(r.URL.Path, "/")
		reportID, err := strconv.Atoi(pathParts[len(pathParts)-1])
		if err != nil {
			http.Error(w, "Invalid report ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Action  string `json:"action"`
			Comment string `json:"comment"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем жалобу
		_, err = db.Exec(`
			UPDATE reports 
			SET status = 'resolved', 
			    moderator_id = ?, 
			    moderator_action = ?, 
			    moderator_comment = ?,
			    reviewed_at = ?
			WHERE id = ?
		`, moderatorID, req.Action, req.Comment, time.Now(), reportID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Report %d reviewed by moderator %d", reportID, moderatorID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Report reviewed successfully",
		})
	}
}

// GetModerationStatsHandler - статистика модерации
func GetModerationStatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := make(map[string]interface{})

		// Всего жалоб
		var totalReports int
		db.QueryRow("SELECT COUNT(*) FROM reports").Scan(&totalReports)
		stats["total_reports"] = totalReports

		// Ожидают рассмотрения
		var pendingReports int
		db.QueryRow("SELECT COUNT(*) FROM reports WHERE status = 'pending'").Scan(&pendingReports)
		stats["pending_reports"] = pendingReports

		// Рассмотрено
		var resolvedReports int
		db.QueryRow("SELECT COUNT(*) FROM reports WHERE status = 'resolved'").Scan(&resolvedReports)
		stats["resolved_reports"] = resolvedReports

		// По типам
		rows, _ := db.Query("SELECT target_type, COUNT(*) FROM reports GROUP BY target_type")
		defer rows.Close()

		byType := make(map[string]int)
		for rows.Next() {
			var targetType string
			var count int
			rows.Scan(&targetType, &count)
			byType[targetType] = count
		}
		stats["by_type"] = byType

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    stats,
		})
	}
}
