package admin

import "time"

type MerchantListItem struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Status            string    `json:"status"`
	VolumeThisMonth   int64     `json:"volume_this_month"`
	TransactionCount  int       `json:"transaction_count"`
	FeeConfigID       *string   `json:"fee_config_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type TransactionListItem struct {
	ID           string `json:"id"`
	PaymentID    string `json:"payment_id"`
	MerchantID   string `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	ProcessorRef string `json:"processor_ref,omitempty"`
	CreatedAt    int64  `json:"created"`
}

type Dispute struct {
	ID         string     `json:"id"`
	PaymentID  string     `json:"payment_id"`
	MerchantID string     `json:"merchant_id"`
	Amount     int64      `json:"amount"`
	Currency   string     `json:"currency"`
	Status     string     `json:"status"`
	Reason     string     `json:"reason"`
	RespondBy  time.Time  `json:"respond_by"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type FeeConfig struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Rate             float64   `json:"rate"`
	FixedFee         int64     `json:"fixed_fee"`
	CrossBorderRate  *float64  `json:"cross_border_rate,omitempty"`
	MonthlyFee       int64     `json:"monthly_fee"`
	MinMonthlyVolume int64     `json:"min_monthly_volume"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type SystemConfig struct {
	ID                   string    `json:"id"`
	SupportedCurrencies  []string  `json:"supported_currencies"`
	SupportedCountries   []string  `json:"supported_countries"`
	MaxTransactionAmount int64     `json:"max_transaction_amount"`
	MaintenanceMode      bool      `json:"maintenance_mode"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type AuditLogEntry struct {
	ID           string    `json:"id"`
	ActorID      string    `json:"actor_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id,omitempty"`
	Details      string    `json:"details,omitempty"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}
