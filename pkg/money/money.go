package money

type Currency string

const (
	BDT Currency = "BDT"
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

type Money struct {
	Amount   int64    `json:"amount"`
	Currency Currency `json:"currency"`
}

func New(amount int64, currency Currency) Money {
	return Money{Amount: amount, Currency: currency}
}
