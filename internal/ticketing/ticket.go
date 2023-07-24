package ticketing

import "time"

type Assignee struct {
	Id       *string `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
}

type TicketTag struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type Ticket struct {
	Id                   *string                `json:"id,omitempty"`
	ParentId             *string                `json:"parent_id,omitempty"`
	CollectionId         *string                `json:"collection_id,omitempty"`
	Type                 *string                `json:"type,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	Status               *string                `json:"status,omitempty"`
	Priority             *string                `json:"priority,omitempty"`
	TicketUrl            *string                `json:"ticket_url,omitempty"`
	Assignees            []Assignee             `json:"assignees,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	CreatedBy            *string                `json:"created_by,omitempty"`
	DueDate              *string                `json:"due_date,omitempty"`
	CompletedAt          *string                `json:"completed_at,omitempty"`
	Tags                 []TicketTag            `json:"tags,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integrtion_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type TicketResponse struct {
	Data []Ticket `json:"data"`
	Meta Meta     `json:"meta"`
}
