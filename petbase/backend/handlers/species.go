package handlers

import (
	"database"
	"encoding/json"
	"net/http"
	"petbase/models"
	"strconv"
	"strings"
)

func SpeciesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getSpecies(w, r)
	case "POST":
		createSpecies(w, r)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func SpeciesDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/species/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getSpeciesDetail(w, r, id)
	case "PUT":
		updateSpecies(w, r, id)
	case "DELETE":
		deleteSpecies(w, r, id)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getSpecies(w http.ResponseWriter, _ *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, name_en, description, icon, created_at FROM species ORDER BY name")
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var species []models.Species
	for rows.Next() {
		var s models.Species
		if err := rows.Scan(&s.ID, &s.Name, &s.NameEn, &s.Description, &s.Icon, &s.CreatedAt); err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		species = append(species, s)
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    species,
	})
}

func getSpeciesDetail(w http.ResponseWriter, _ *http.Request, id int) {
	var s models.Species
	err := database.DB.QueryRow(
		"SELECT id, name, name_en, description, icon, created_at FROM species WHERE id = ?",
		id,
	).Scan(&s.ID, &s.Name, &s.NameEn, &s.Description, &s.Icon, &s.CreatedAt)

	if err != nil {
		sendError(w, "Species not found", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    s,
	})
}

func createSpecies(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		NameEn      string `json:"name_en"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO species (name, name_en, description, icon) VALUES (?, ?, ?, ?)",
		req.Name, req.NameEn, req.Description, req.Icon,
	)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	sendJSON(w, map[string]interface{}{
		"success": true,
		"id":      id,
	})
}

func updateSpecies(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		Name        string `json:"name"`
		NameEn      string `json:"name_en"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(
		"UPDATE species SET name = ?, name_en = ?, description = ?, icon = ? WHERE id = ?",
		req.Name, req.NameEn, req.Description, req.Icon, id,
	)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
	})
}

func deleteSpecies(w http.ResponseWriter, _ *http.Request, id int) {
	_, err := database.DB.Exec("DELETE FROM species WHERE id = ?", id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
	})
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
