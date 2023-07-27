package ticketing

import "time"

type Tag struct {
	Id                   *string                `json:"id,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Color                *string                `json:"color,omitempty"`
	Description          *string                `json:"description,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integration_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type TagResponse struct {
	Data []Tag `json:"data"`
	Meta Meta  `json:"meta"`
}
