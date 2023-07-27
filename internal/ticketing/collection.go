package ticketing

import "time"

type Collection struct {
	Id                   *string                `json:"id,omitempty"`
	ParentId             *string                `json:"parent_id,omitempty"`
	Type                 *string                `json:"type,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	UpdatedAt            *string                `json:"updated_at,omitempty"`
	CreatedAt            *string                `json:"created_at,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integrtion_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type CollectionResponse struct {
	Data []Collection `json:"data"`
	Meta Meta         `json:"meta"`
}
