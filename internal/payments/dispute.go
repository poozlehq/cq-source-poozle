package payments

import "time"

type Dispute struct {
	Id                   string                 `json:"id,omitempty"`
	Amount               string                 `json:"amount,omitempty"`
	ChargeID             string                 `json:"charge_id,omitempty"`
	Currency             string                 `json:"currency,omitempty"`
	Reason               string                 `json:"reason,omitempty"`
	Evidence             Evidence               `json:"evidence,omitempty"`
	Status               string                 `json:"status,omitempty"`
	Priority             string                 `json:"priority,omitempty"`
	IsChargeRefundable   bool                   `json:"is_charge_refundable,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integration_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type Evidence struct {
	AccessActivityLog            string   `json:"access_activity_log,omitempty"`
	BillingAddress               string   `json:"billing_address,omitempty"`
	CancellationPolicy           []string `json:"cancellation_policy,omitempty"`
	CancellationPolicyDisclosure string   `json:"cancellation_policy_disclosure,omitempty"`
	CancellationRebuttal         string   `json:"cancellation_rebuttal,omitempty"`
	CustomerCommunication        []string `json:"customer_communication,omitempty"`
	CustomerEmailAddress         string   `json:"customer_email_address,omitempty"`
	CustomerName                 string   `json:"customer_name,omitempty"`
	CustomerPurchaseIP           string   `json:"customer_purchase_ip,omitempty"`
	CustomerSignature            []string `json:"customer_signature,omitempty"`
	DuplicateChargeDocumentation []string `json:"duplicate_charge_documentation,omitempty"`
	DuplicateChargeExplanation   string   `json:"duplicate_charge_explanation,omitempty"`
	DuplicateChargeID            string   `json:"duplicate_charge_id,omitempty"`
	ProductDescription           string   `json:"product_description,omitempty"`
	Receipt                      []string `json:"receipt,omitempty"`
	RefundPolicy                 []string `json:"refund_policy,omitempty"`
	RefundPolicyDisclosure       string   `json:"refund_policy_disclosure,omitempty"`
	RefundRefusalExplanation     string   `json:"refund_refusal_explanation,omitempty"`
	ServiceDate                  string   `json:"service_date,omitempty"`
	ServiceDocumentation         []string `json:"service_documentation,omitempty"`
	ShippingAddress              string   `json:"shipping_address,omitempty"`
	ShippingCarrier              string   `json:"shipping_carrier,omitempty"`
	ShippingDate                 string   `json:"shipping_date,omitempty"`
	ShippingDocumentation        []string `json:"shipping_documentation,omitempty"`
	ShippingTrackingNumber       string   `json:"shipping_tracking_number,omitempty"`
	UncategorizedFile            []string `json:"uncategorized_file,omitempty"`
	UncategorizedText            string   `json:"uncategorized_text,omitempty"`
	DueBy                        string   `json:"due_by,omitempty"`
	HasEvidence                  bool     `json:"has_evidence,omitempty"`
	PastDue                      bool     `json:"past_due,omitempty"`
	SubmissionCount              string   `json:"submission_count,omitempty"`
}

type DisputesResponse struct {
	Data []Dispute `json:"data"`
	Meta Meta      `json:"meta"`
}
