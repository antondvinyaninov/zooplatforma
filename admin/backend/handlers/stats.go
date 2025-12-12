package handlers

import (
	"database"
	"net/http"
)

// AdminStatsOverviewHandler - общая статистика
func AdminStatsOverviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var usersCount, postsCount, petsCount int

	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&usersCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&postsCount)
	database.DB.QueryRow("SELECT COUNT(*) FROM pets").Scan(&petsCount)

	sendSuccess(w, map[string]interface{}{
		"users_count": usersCount,
		"posts_count": postsCount,
		"pets_count":  petsCount,
	})
}
