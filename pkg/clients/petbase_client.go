package clients

import (
	"fmt"
)

// PetBaseClient - клиент для PetBase Service
type PetBaseClient struct {
	*BaseClient
}

// NewPetBaseClient создает новый клиент для PetBase Service
func NewPetBaseClient(baseURL string) *PetBaseClient {
	return &PetBaseClient{
		BaseClient: NewBaseClient(baseURL, 0),
	}
}

// Pet - структура питомца
type Pet struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Species           string `json:"species"`
	Breed             string `json:"breed"`
	BirthDate         string `json:"birth_date"`
	Gender            string `json:"gender"`
	Color             string `json:"color"`
	Size              string `json:"size"`
	Photo             string `json:"photo"`
	ChipNumber        string `json:"chip_number"`
	TagNumber         string `json:"tag_number"`
	Status            string `json:"status"`
	OwnerID           int    `json:"owner_id"`
	OrganizationID    int    `json:"organization_id"`
	CuratorID         int    `json:"curator_id"`
	Location          string `json:"location"`
	Story             string `json:"story"`
	HealthStatus      string `json:"health_status"`
	SterilizationDate string `json:"sterilization_date"`
	CreatedAt         string `json:"created_at"`
}

// MedicalRecord - медицинская запись
type MedicalRecord struct {
	ID           int    `json:"id"`
	PetID        int    `json:"pet_id"`
	Date         string `json:"date"`
	ClinicName   string `json:"clinic_name"`
	Diagnosis    string `json:"diagnosis"`
	Treatment    string `json:"treatment"`
	Veterinarian string `json:"veterinarian"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"created_at"`
}

// Vaccination - вакцинация
type Vaccination struct {
	ID           int    `json:"id"`
	PetID        int    `json:"pet_id"`
	VaccineName  string `json:"vaccine_name"`
	Date         string `json:"date"`
	NextDate     string `json:"next_date"`
	ClinicName   string `json:"clinic_name"`
	Veterinarian string `json:"veterinarian"`
	CreatedAt    string `json:"created_at"`
}

// GetPet получает информацию о питомце
func (c *PetBaseClient) GetPet(petID int, token string) (*Pet, error) {
	endpoint := fmt.Sprintf("/api/pets/%d", petID)
	resp, err := c.Get(endpoint, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get pet: %s", resp.Error)
	}

	// Парсим данные
	petMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	pet := mapToPet(petMap)
	return pet, nil
}

// GetUserPets получает питомцев пользователя
func (c *PetBaseClient) GetUserPets(userID int, token string) ([]*Pet, error) {
	endpoint := fmt.Sprintf("/api/pets/user/%d", userID)
	resp, err := c.Get(endpoint, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user pets: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get user pets: %s", resp.Error)
	}

	// Парсим массив питомцев
	petsData, ok := resp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	pets := make([]*Pet, 0, len(petsData))
	for _, item := range petsData {
		petMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		pets = append(pets, mapToPet(petMap))
	}

	return pets, nil
}

// CreatePet создает нового питомца
func (c *PetBaseClient) CreatePet(pet *Pet, token string) (*Pet, error) {
	resp, err := c.Post("/api/pets", pet, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to create pet: %s", resp.Error)
	}

	petMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return mapToPet(petMap), nil
}

// UpdatePet обновляет питомца
func (c *PetBaseClient) UpdatePet(petID int, updates map[string]interface{}, token string) (*Pet, error) {
	endpoint := fmt.Sprintf("/api/pets/%d", petID)
	resp, err := c.Put(endpoint, updates, token)
	if err != nil {
		return nil, fmt.Errorf("failed to update pet: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to update pet: %s", resp.Error)
	}

	petMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return mapToPet(petMap), nil
}

// GetMedicalRecords получает медицинские записи питомца
func (c *PetBaseClient) GetMedicalRecords(petID int, token string) ([]*MedicalRecord, error) {
	endpoint := fmt.Sprintf("/api/pets/%d/medical", petID)
	resp, err := c.Get(endpoint, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical records: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get medical records: %s", resp.Error)
	}

	recordsData, ok := resp.Data.([]interface{})
	if !ok {
		return []*MedicalRecord{}, nil
	}

	records := make([]*MedicalRecord, 0, len(recordsData))
	for _, item := range recordsData {
		recordMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		records = append(records, mapToMedicalRecord(recordMap))
	}

	return records, nil
}

// GetVaccinations получает вакцинации питомца
func (c *PetBaseClient) GetVaccinations(petID int, token string) ([]*Vaccination, error) {
	endpoint := fmt.Sprintf("/api/pets/%d/vaccinations", petID)
	resp, err := c.Get(endpoint, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get vaccinations: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get vaccinations: %s", resp.Error)
	}

	vaccsData, ok := resp.Data.([]interface{})
	if !ok {
		return []*Vaccination{}, nil
	}

	vaccs := make([]*Vaccination, 0, len(vaccsData))
	for _, item := range vaccsData {
		vaccMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		vaccs = append(vaccs, mapToVaccination(vaccMap))
	}

	return vaccs, nil
}

// Helper функции для маппинга
func mapToPet(m map[string]interface{}) *Pet {
	return &Pet{
		ID:                getIntOrZero(m, "id"),
		Name:              getStringOrEmpty(m, "name"),
		Species:           getStringOrEmpty(m, "species"),
		Breed:             getStringOrEmpty(m, "breed"),
		BirthDate:         getStringOrEmpty(m, "birth_date"),
		Gender:            getStringOrEmpty(m, "gender"),
		Color:             getStringOrEmpty(m, "color"),
		Size:              getStringOrEmpty(m, "size"),
		Photo:             getStringOrEmpty(m, "photo"),
		ChipNumber:        getStringOrEmpty(m, "chip_number"),
		TagNumber:         getStringOrEmpty(m, "tag_number"),
		Status:            getStringOrEmpty(m, "status"),
		OwnerID:           getIntOrZero(m, "owner_id"),
		OrganizationID:    getIntOrZero(m, "organization_id"),
		CuratorID:         getIntOrZero(m, "curator_id"),
		Location:          getStringOrEmpty(m, "location"),
		Story:             getStringOrEmpty(m, "story"),
		HealthStatus:      getStringOrEmpty(m, "health_status"),
		SterilizationDate: getStringOrEmpty(m, "sterilization_date"),
		CreatedAt:         getStringOrEmpty(m, "created_at"),
	}
}

func mapToMedicalRecord(m map[string]interface{}) *MedicalRecord {
	return &MedicalRecord{
		ID:           getIntOrZero(m, "id"),
		PetID:        getIntOrZero(m, "pet_id"),
		Date:         getStringOrEmpty(m, "date"),
		ClinicName:   getStringOrEmpty(m, "clinic_name"),
		Diagnosis:    getStringOrEmpty(m, "diagnosis"),
		Treatment:    getStringOrEmpty(m, "treatment"),
		Veterinarian: getStringOrEmpty(m, "veterinarian"),
		Notes:        getStringOrEmpty(m, "notes"),
		CreatedAt:    getStringOrEmpty(m, "created_at"),
	}
}

func mapToVaccination(m map[string]interface{}) *Vaccination {
	return &Vaccination{
		ID:           getIntOrZero(m, "id"),
		PetID:        getIntOrZero(m, "pet_id"),
		VaccineName:  getStringOrEmpty(m, "vaccine_name"),
		Date:         getStringOrEmpty(m, "date"),
		NextDate:     getStringOrEmpty(m, "next_date"),
		ClinicName:   getStringOrEmpty(m, "clinic_name"),
		Veterinarian: getStringOrEmpty(m, "veterinarian"),
		CreatedAt:    getStringOrEmpty(m, "created_at"),
	}
}

func getIntOrZero(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok && val != nil {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}
