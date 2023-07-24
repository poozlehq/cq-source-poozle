package ticketing

import "time"

type Creator struct {
	Id       *string `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
}

type Comment struct {
	Id                   *string                `json:"id,omitempty"`
	TicketId             *string                `json:"ticket_id,omitempty"`
	Body                 *string                `json:"body,omitempty"`
	HtmlBody             *string                `json:"html_body,omitempty"`
	CreatedById          *string                `json:"created_by_id,omitempty"`
	CreatedBy            Creator                `json:"created_by,omitempty"`
	IsPrivate            *bool                  `json:"is_private,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integrtion_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type CommentResponse struct {
	Data []Comment `json:"data"`
	Meta Meta      `json:"meta"`
}
