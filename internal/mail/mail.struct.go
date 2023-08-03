package mail

import "time"

type Recipient struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
}

type Mail struct {
	Id        *string     `json:"id,omitempty"`
	Body      *string     `json:"body,omitempty"`
	HtmlBody  *string     `json:"html_body,omitempty"`
	UserId    *string     `json:"user_id,omitempty"`
	Date      *time.Time  `json:"date,omitempty"`
	Snippet   *string     `json:"snippet,omitempty"`
	Subject   *string     `json:"subject,omitempty"`
	ThreadId  *string     `json:"thread_id,omitempty"`
	Starred   *bool       `json:"starred,omitempty"`
	Unread    *bool       `json:"unread,omitempty"`
	Cc        []Recipient `json:"cc,omitempty"`
	Bcc       []Recipient `json:"bcc,omitempty"`
	From      []Recipient `json:"from,omitempty"`
	ReplyTo   []Recipient `json:"reply_to,omitempty"`
	Labels    []string    `json:"labels,omitempty"`
	InReplyTo string      `json:"in_reply_to,omitempty"`
}

type MailResponse struct {
	Data []Mail `json:"data"`
	Meta Meta   `json:"meta"`
}
