package clients

import (
	"fmt"
)

// MainClient - клиент для Main Service
type MainClient struct {
	*BaseClient
}

// NewMainClient создает новый клиент для Main Service
func NewMainClient(baseURL string) *MainClient {
	return &MainClient{
		BaseClient: NewBaseClient(baseURL, 0),
	}
}

// Organization - структура организации
type Organization struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
	Cover       string `json:"cover"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Verified    bool   `json:"verified"`
	CreatedAt   string `json:"created_at"`
}

// OrganizationMember - участник организации
type OrganizationMember struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	OrganizationID int    `json:"organization_id"`
	Role           string `json:"role"`
	Position       string `json:"position"`
	UserName       string `json:"user_name"`
	UserEmail      string `json:"user_email"`
	UserAvatar     string `json:"user_avatar"`
	JoinedAt       string `json:"joined_at"`
}

// GetOrganization получает информацию об организации
func (c *MainClient) GetOrganization(orgID int, token string) (*Organization, error) {
	endpoint := fmt.Sprintf("/api/organizations/%d", orgID)
	resp, err := c.Get(endpoint, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get organization: %s", resp.Error)
	}

	orgMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return mapToOrganization(orgMap), nil
}

// GetUserOrganizations получает организации пользователя
func (c *MainClient) GetUserOrganizations(userID int, token string) ([]*Organization, error) {
	resp, err := c.Get("/api/organizations/my", token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get user organizations: %s", resp.Error)
	}

	orgsData, ok := resp.Data.([]interface{})
	if !ok {
		return []*Organization{}, nil
	}

	orgs := make([]*Organization, 0, len(orgsData))
	for _, item := range orgsData {
		orgMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		orgs = append(orgs, mapToOrganization(orgMap))
	}

	return orgs, nil
}

// GetOrganizationMembers получает участников организации
func (c *MainClient) GetOrganizationMembers(orgID int, token string) ([]*OrganizationMember, error) {
	endpoint := fmt.Sprintf("/api/organizations/%d/members", orgID)
	resp, err := c.Get(endpoint, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization members: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get organization members: %s", resp.Error)
	}

	membersData, ok := resp.Data.([]interface{})
	if !ok {
		return []*OrganizationMember{}, nil
	}

	members := make([]*OrganizationMember, 0, len(membersData))
	for _, item := range membersData {
		memberMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		members = append(members, mapToOrganizationMember(memberMap))
	}

	return members, nil
}

// CheckMembership проверяет является ли пользователь участником организации
func (c *MainClient) CheckMembership(userID, orgID int, token string) (bool, error) {
	members, err := c.GetOrganizationMembers(orgID, token)
	if err != nil {
		return false, err
	}

	for _, member := range members {
		if member.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}

// Helper функции
func mapToOrganization(m map[string]interface{}) *Organization {
	return &Organization{
		ID:          getIntOrZero(m, "id"),
		Name:        getStringOrEmpty(m, "name"),
		Type:        getStringOrEmpty(m, "type"),
		Description: getStringOrEmpty(m, "description"),
		Address:     getStringOrEmpty(m, "address"),
		Phone:       getStringOrEmpty(m, "phone"),
		Email:       getStringOrEmpty(m, "email"),
		Website:     getStringOrEmpty(m, "website"),
		Logo:        getStringOrEmpty(m, "logo"),
		Cover:       getStringOrEmpty(m, "cover"),
		Region:      getStringOrEmpty(m, "region"),
		City:        getStringOrEmpty(m, "city"),
		Verified:    getBoolOrFalse(m, "verified"),
		CreatedAt:   getStringOrEmpty(m, "created_at"),
	}
}

func mapToOrganizationMember(m map[string]interface{}) *OrganizationMember {
	return &OrganizationMember{
		ID:             getIntOrZero(m, "id"),
		UserID:         getIntOrZero(m, "user_id"),
		OrganizationID: getIntOrZero(m, "organization_id"),
		Role:           getStringOrEmpty(m, "role"),
		Position:       getStringOrEmpty(m, "position"),
		UserName:       getStringOrEmpty(m, "user_name"),
		UserEmail:      getStringOrEmpty(m, "user_email"),
		UserAvatar:     getStringOrEmpty(m, "user_avatar"),
		JoinedAt:       getStringOrEmpty(m, "joined_at"),
	}
}

func getBoolOrFalse(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok && val != nil {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
