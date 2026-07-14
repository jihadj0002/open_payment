package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type CardRequest struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	CVV            string `json:"cvv"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	MerchantRef    string `json:"merchant_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

type WalletRequest struct {
	WalletID       string `json:"wallet_id"`
	PIN            string `json:"pin"`
	OTP            string `json:"otp"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	MerchantRef    string `json:"merchant_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

type BankRequest struct {
	AccountNumber  string `json:"account_number"`
	RoutingNumber  string `json:"routing_number"`
	AccountType    string `json:"account_type"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	MerchantRef    string `json:"merchant_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

type ProcessorResponse struct {
	Success      bool   `json:"success"`
	ProcessorRef string `json:"processor_ref"`
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	Fee          int64  `json:"fee"`
	ProcessedAt  string `json:"processed_at"`
}

type ProcessorInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Healthy bool   `json:"healthy"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /process/card", handleCardProcess)
	mux.HandleFunc("POST /process/wallet", handleWalletProcess)
	mux.HandleFunc("POST /process/bank", handleBankProcess)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Mock Processor Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(ProcessorInfo{
		Name:    "mock-processor",
		Type:    "payment",
		Healthy: true,
	})
}

func handleCardProcess(w http.ResponseWriter, r *http.Request) {
	var req CardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)

	// Simulate card brand detection
	brand := "visa"
	if len(req.CardNumber) > 0 {
		switch req.CardNumber[0] {
		case '4':
			brand = "visa"
		case '5':
			brand = "mastercard"
		case '3':
			brand = "amex"
		default:
			brand = "unknown"
		}
	}

	simulateFailure := rand.Float64() < 0.1
	status := "succeeded"
	msg := "Payment approved"

	if simulateFailure {
		status = "failed"
		msg = "Card declined - insufficient funds"
	}

	resp := ProcessorResponse{
		Success:      !simulateFailure,
		ProcessorRef: fmt.Sprintf("card_%s_%d", brand, time.Now().UnixMilli()),
		Status:       status,
		Message:      msg,
		Fee:          int64(float64(req.Amount) * 0.022),
		ProcessedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}

	statusCode := http.StatusOK
	if simulateFailure {
		statusCode = http.StatusPaymentRequired
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func handleWalletProcess(w http.ResponseWriter, r *http.Request) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(100+rand.Intn(300)) * time.Millisecond)

	provider := "bkash"
	if req.WalletID != "" {
		if req.WalletID[0] == '0' {
			provider = "nagad"
		} else if req.WalletID[0] == '1' {
			provider = "rocket"
		}
	}

	simulateFailure := rand.Float64() < 0.05
	status := "succeeded"
	msg := "Wallet payment approved"

	if simulateFailure {
		status = "failed"
		msg = "Insufficient balance in wallet"
	}

	resp := ProcessorResponse{
		Success:      !simulateFailure,
		ProcessorRef: fmt.Sprintf("wallet_%s_%d", provider, time.Now().UnixMilli()),
		Status:       status,
		Message:      msg,
		Fee:          int64(float64(req.Amount) * 0.015),
		ProcessedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}

	statusCode := http.StatusOK
	if simulateFailure {
		statusCode = http.StatusPaymentRequired
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func handleBankProcess(w http.ResponseWriter, r *http.Request) {
	var req BankRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(200+rand.Intn(500)) * time.Millisecond)

	simulateFailure := rand.Float64() < 0.08
	status := "succeeded"
	msg := "Bank transfer initiated"

	if simulateFailure {
		status = "failed"
		msg = "Account verification failed"
	}

	resp := ProcessorResponse{
		Success:      !simulateFailure,
		ProcessorRef: fmt.Sprintf("bank_%d", time.Now().UnixMilli()),
		Status:       status,
		Message:      msg,
		Fee:          int64(float64(req.Amount) * 0.005),
		ProcessedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}

	statusCode := http.StatusOK
	if simulateFailure {
		statusCode = http.StatusPaymentRequired
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
