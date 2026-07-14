package customer

import "time"

type Customer struct {
	ID         string            `json:"id"`
	MerchantID string            `json:"merchant_id"`
	Email      *string           `json:"email,omitempty"`
	Phone      *string           `json:"phone,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type CreateCustomerRequest struct {
	Email    *string           `json:"email,omitempty"`
	Phone    *string           `json:"phone,omitempty"`
	Name     *string           `json:"name,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type UpdateCustomerRequest struct {
	Email    *string           `json:"email,omitempty"`
	Phone    *string           `json:"phone,omitempty"`
	Name     *string           `json:"name,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type PaymentMethod struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id"`
	Type           string    `json:"type"`
	Token          *string   `json:"token,omitempty"`
	Last4          *string   `json:"last4,omitempty"`
	Brand          *string   `json:"brand,omitempty"`
	ExpMonth       *int      `json:"exp_month,omitempty"`
	ExpYear        *int      `json:"exp_year,omitempty"`
	CardholderName *string   `json:"cardholder_name,omitempty"`
	WalletType      *string   `json:"wallet_type,omitempty"`
	WalletPhone     *string   `json:"wallet_phone,omitempty"`
	Fingerprint     *string   `json:"fingerprint,omitempty"`
	IsDefault       bool      `json:"is_default"`
	IsActive        bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AttachPaymentMethodRequest struct {
	Type           string  `json:"type"`
	CardNumber     *string `json:"card_number,omitempty"`
	ExpMonth       *int    `json:"exp_month,omitempty"`
	ExpYear        *int    `json:"exp_year,omitempty"`
	CVC            *string `json:"cvc,omitempty"`
	CardholderName *string `json:"cardholder_name,omitempty"`
	WalletType     *string `json:"wallet_type,omitempty"`
	WalletPhone    *string `json:"wallet_phone,omitempty"`
	SetAsDefault   bool    `json:"set_as_default"`
}
