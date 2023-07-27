package ticketing

import "time"

type User struct {
	Id                   *string                `json:"id,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Avatar               *string                `json:"avatar,omitempty"`
	EmailAddress         *string                `json:"email,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integrtion_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type UsersResponse struct {
	Data []User `json:"data"`
	Meta Meta   `json:"meta"`
}

type UserResponse struct {
	Data User `json:"data"`
}
