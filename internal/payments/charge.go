package payments

import "time"

type Address struct {
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

type BillingDetails struct {
	Address Address `json:"address,omitempty"`
	Email   string  `json:"email,omitempty"`
	Name    string  `json:"name,omitempty"`
	Phone   string  `json:"phone,omitempty"`
}

type PaymentMethod struct {
	Type    PaymentMethodType `json:"type,omitempty"`
	Details PaymentDetails    `json:"details,omitempty"`
}

type PaymentDetails struct {
	Details map[string]interface{} `json:"details,omitempty"`
}

type Outcome struct {
	NetworkStatus string `json:"network_status,omitempty"`
	Reason        string `json:"reason,omitempty"`
	RiskLevel     string `json:"risk_level,omitempty"`
	SellerMessage string `json:"seller_message,omitempty"`
	Type          string `json:"type,omitempty"`
}

type Charge struct {
	Id                   string                 `json:"id,omitempty"`
	Amount               string                 `json:"amount,omitempty"`
	AmountRefunded       string                 `json:"amount_refunded,omitempty"`
	Application          string                 `json:"application,omitempty"`
	ApplicationFeeAmount string                 `json:"application_fee_amount,omitempty"`
	BillingDetails       BillingDetails         `json:"billing_details,omitempty"`
	Captured             bool                   `json:"captured,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	Currency             string                 `json:"currency,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Disputed             bool                   `json:"disputed,omitempty"`
	FailureCode          string                 `json:"failure_code,omitempty"`
	FailureMessage       string                 `json:"failure_message,omitempty"`
	Invoice              string                 `json:"invoice,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	Outcome              Outcome                `json:"outcome,omitempty"`
	Paid                 bool                   `json:"paid,omitempty"`
	PaymentMethod        PaymentMethod          `json:"payment_method,omitempty"`
	Email                string                 `json:"email,omitempty"`
	Contact              string                 `json:"contact,omitempty"`
	Status               string                 `json:"status,omitempty"`
	Raw                  map[string]interface{} `json:"raw,omitempty"`
	IntegrationAccountId *string                `json:"integration_account_id,omitempty"`
	CqCreatedAt          *time.Time             `json:"cq_created_at,omitempty"`
	CqUpdatedAt          *time.Time             `json:"cq_updated_at,omitempty"`
}

type PaymentMethodType string
type PaymentMethodStatus string

const (
	AchCreditTransfer PaymentMethodType = "ach_credit_transfer"
	AchDebit          PaymentMethodType = "ach_debit"
	AcssDebit         PaymentMethodType = "acss_debit"
	Alipay            PaymentMethodType = "alipay"
	AuBecsDebit       PaymentMethodType = "au_becs_debit"
	Bancontact        PaymentMethodType = "bancontact"
	Card              PaymentMethodType = "card"
	CardPresent       PaymentMethodType = "card_present"
	Eps               PaymentMethodType = "eps"
	Giropay           PaymentMethodType = "giropay"
	Ideal             PaymentMethodType = "ideal"
	Klarna            PaymentMethodType = "klarna"
	Multibanco        PaymentMethodType = "multibanco"
	P24               PaymentMethodType = "p24"
	SepaDebit         PaymentMethodType = "sepa_debit"
	Sofort            PaymentMethodType = "sofort"
	StripeAccount     PaymentMethodType = "stripe_account"
	Wechat            PaymentMethodType = "wechat"
	Netbank           PaymentMethodType = "netbank"
	Wallet            PaymentMethodType = "wallet"
	Emi               PaymentMethodType = "emi"
	Upi               PaymentMethodType = "upi"
)

const (
	Created    PaymentMethodStatus = "created"
	Authorized PaymentMethodStatus = "authorized"
	Succeeded  PaymentMethodStatus = "succeeded"
	Refunded   PaymentMethodStatus = "refunded"
	Failed     PaymentMethodStatus = "failed"
)

type ChargesResponse struct {
	Data []Charge `json:"data"`
	Meta Meta     `json:"meta"`
}
