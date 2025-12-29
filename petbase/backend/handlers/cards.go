package handlers

import (
	"database"
	"encoding/json"
	"net/http"
	"petbase/models"
	"strconv"
	"strings"
)

func CardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getCards(w, r)
	case "POST":
		createCard(w, r)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CardDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/cards/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getCardDetail(w, r, id)
	case "PUT":
		updateCard(w, r, id)
	case "DELETE":
		deleteCard(w, r, id)
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CardsByBreedHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/cards/breed/")
	breedID, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, "Invalid breed ID", http.StatusBadRequest)
		return
	}

	getCardsByBreed(w, r, breedID)
}

func getCards(w http.ResponseWriter, _ *http.Request) {
	query := `
		SELECT c.id, c.breed_id, b.name as breed_name, c.title, c.description,
		       c.characteristics, c.care_tips, c.health_info, c.nutrition,
		       c.photos, c.is_published, c.created_at, c.updated_at
		FROM pet_cards c
		LEFT JOIN breeds b ON c.breed_id = b.id
		ORDER BY c.created_at DESC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cards []models.PetCard
	for rows.Next() {
		var c models.PetCard
		if err := rows.Scan(
			&c.ID, &c.BreedID, &c.BreedName, &c.Title, &c.Description,
			&c.Characteristics, &c.CareTips, &c.HealthInfo, &c.Nutrition,
			&c.Photos, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, c)
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    cards,
	})
}

func getCardsByBreed(w http.ResponseWriter, _ *http.Request, breedID int) {
	query := `
		SELECT c.id, c.breed_id, b.name as breed_name, c.title, c.description,
		       c.characteristics, c.care_tips, c.health_info, c.nutrition,
		       c.photos, c.is_published, c.created_at, c.updated_at
		FROM pet_cards c
		LEFT JOIN breeds b ON c.breed_id = b.id
		WHERE c.breed_id = ?
		ORDER BY c.created_at DESC
	`
	rows, err := database.DB.Query(query, breedID)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cards []models.PetCard
	for rows.Next() {
		var c models.PetCard
		if err := rows.Scan(
			&c.ID, &c.BreedID, &c.BreedName, &c.Title, &c.Description,
			&c.Characteristics, &c.CareTips, &c.HealthInfo, &c.Nutrition,
			&c.Photos, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, c)
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    cards,
	})
}

func getCardDetail(w http.ResponseWriter, _ *http.Request, id int) {
	query := `
		SELECT c.id, c.breed_id, b.name as breed_name, c.title, c.description,
		       c.characteristics, c.care_tips, c.health_info, c.nutrition,
		       c.photos, c.is_published, c.created_at, c.updated_at
		FROM pet_cards c
		LEFT JOIN breeds b ON c.breed_id = b.id
		WHERE c.id = ?
	`
	var c models.PetCard
	err := database.DB.QueryRow(query, id).Scan(
		&c.ID, &c.BreedID, &c.BreedName, &c.Title, &c.Description,
		&c.Characteristics, &c.CareTips, &c.HealthInfo, &c.Nutrition,
		&c.Photos, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		sendError(w, "Card not found", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"data":    c,
	})
}

func createCard(w http.ResponseWriter, r *http.Request) {
	var req models.PetCard

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO pet_cards (
			breed_id, title, description, characteristics, care_tips,
			health_info, nutrition, photos, is_published
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		req.BreedID, req.Title, req.Description, req.Characteristics, req.CareTips,
		req.HealthInfo, req.Nutrition, req.Photos, req.IsPublished,
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

func updateCard(w http.ResponseWriter, r *http.Request, id int) {
	var req models.PetCard

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`
		UPDATE pet_cards SET
			breed_id = ?, title = ?, description = ?, characteristics = ?, care_tips = ?,
			health_info = ?, nutrition = ?, photos = ?, is_published = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		req.BreedID, req.Title, req.Description, req.Characteristics, req.CareTips,
		req.HealthInfo, req.Nutrition, req.Photos, req.IsPublished, id,
	)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
	})
}

func deleteCard(w http.ResponseWriter, _ *http.Request, id int) {
	_, err := database.DB.Exec("DELETE FROM pet_cards WHERE id = ?", id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
	})
}
