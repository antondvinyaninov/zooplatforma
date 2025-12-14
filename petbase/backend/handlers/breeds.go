package handlers

import (
	"database"
	"encoding/json"
	"net/http"
	"petbase/models"
	"strconv"
	"strings"
)

func BreedsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getBreeds(w, r)
	case "POST":
		createBreed(w, r)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func BreedDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/breeds/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getBreedDetail(w, r, id)
	case "PUT":
		updateBreed(w, r, id)
	case "DELETE":
		deleteBreed(w, r, id)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func BreedsBySpeciesHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/breeds/species/")
	speciesID, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "Invalid species ID", http.StatusBadRequest)
		return
	}

	getBreedsBySpecies(w, r, speciesID)
}

func getBreeds(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT b.id, b.species_id, s.name as species_name, b.name, b.name_en, 
		       b.description, b.origin, b.size, b.weight_min, b.weight_max,
		       b.lifespan_min, b.lifespan_max, b.temperament, b.care_level, 
		       b.photo, b.created_at
		FROM breeds b
		LEFT JOIN species s ON b.species_id = s.id
		ORDER BY b.name
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var breeds []models.Breed
	for rows.Next() {
		var b models.Breed
		if err := rows.Scan(
			&b.ID, &b.SpeciesID, &b.SpeciesName, &b.Name, &b.NameEn,
			&b.Description, &b.Origin, &b.Size, &b.WeightMin, &b.WeightMax,
			&b.LifespanMin, &b.LifespanMax, &b.Temperament, &b.CareLevel,
			&b.Photo, &b.CreatedAt,
		); err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breeds = append(breeds, b)
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    breeds,
	})
}

func getBreedsBySpecies(w http.ResponseWriter, r *http.Request, speciesID int) {
	query := `
		SELECT b.id, b.species_id, s.name as species_name, b.name, b.name_en, 
		       b.description, b.origin, b.size, b.weight_min, b.weight_max,
		       b.lifespan_min, b.lifespan_max, b.temperament, b.care_level, 
		       b.photo, b.created_at
		FROM breeds b
		LEFT JOIN species s ON b.species_id = s.id
		WHERE b.species_id = ?
		ORDER BY b.name
	`
	rows, err := database.DB.Query(query, speciesID)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var breeds []models.Breed
	for rows.Next() {
		var b models.Breed
		if err := rows.Scan(
			&b.ID, &b.SpeciesID, &b.SpeciesName, &b.Name, &b.NameEn,
			&b.Description, &b.Origin, &b.Size, &b.WeightMin, &b.WeightMax,
			&b.LifespanMin, &b.LifespanMax, &b.Temperament, &b.CareLevel,
			&b.Photo, &b.CreatedAt,
		); err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		breeds = append(breeds, b)
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    breeds,
	})
}

func getBreedDetail(w http.ResponseWriter, r *http.Request, id int) {
	query := `
		SELECT b.id, b.species_id, s.name as species_name, b.name, b.name_en, 
		       b.description, b.origin, b.size, b.weight_min, b.weight_max,
		       b.lifespan_min, b.lifespan_max, b.temperament, b.care_level, 
		       b.photo, b.created_at
		FROM breeds b
		LEFT JOIN species s ON b.species_id = s.id
		WHERE b.id = ?
	`
	var b models.Breed
	err := database.DB.QueryRow(query, id).Scan(
		&b.ID, &b.SpeciesID, &b.SpeciesName, &b.Name, &b.NameEn,
		&b.Description, &b.Origin, &b.Size, &b.WeightMin, &b.WeightMax,
		&b.LifespanMin, &b.LifespanMax, &b.Temperament, &b.CareLevel,
		&b.Photo, &b.CreatedAt,
	)

	if err != nil {
		sendError(w, "Breed not found", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    b,
	})
}

func createBreed(w http.ResponseWriter, r *http.Request) {
	var req models.Breed

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO breeds (
			species_id, name, name_en, description, origin, size,
			weight_min, weight_max, lifespan_min, lifespan_max,
			temperament, care_level, photo
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		req.SpeciesID, req.Name, req.NameEn, req.Description, req.Origin, req.Size,
		req.WeightMin, req.WeightMax, req.LifespanMin, req.LifespanMax,
		req.Temperament, req.CareLevel, req.Photo,
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

func updateBreed(w http.ResponseWriter, r *http.Request, id int) {
	var req models.Breed

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`
		UPDATE breeds SET
			species_id = ?, name = ?, name_en = ?, description = ?, origin = ?, size = ?,
			weight_min = ?, weight_max = ?, lifespan_min = ?, lifespan_max = ?,
			temperament = ?, care_level = ?, photo = ?
		WHERE id = ?
	`,
		req.SpeciesID, req.Name, req.NameEn, req.Description, req.Origin, req.Size,
		req.WeightMin, req.WeightMax, req.LifespanMin, req.LifespanMax,
		req.Temperament, req.CareLevel, req.Photo, id,
	)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
	})
}

func deleteBreed(w http.ResponseWriter, r *http.Request, id int) {
	_, err := database.DB.Exec("DELETE FROM breeds WHERE id = ?", id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
	})
}
