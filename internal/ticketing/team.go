package ticketing

import "time"

type Member struct {
	Id       *string `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
}

type Team struct {
	Id                   *string                `json:"id,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Member               Member                 `json:"color,omitempty"`
	Description          *string                `json:"description,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integration_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type TeamResponse struct {
	Data []Team `json:"data"`
	Meta Meta   `json:"meta"`
}
