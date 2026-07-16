package nagad

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Credentials struct {
	MerchantID         string
	MerchantPrivateKey string
	PGPublicKey        string
	BaseURL            string
}

type InitializePaymentResponse struct {
	MerchantID       string `json:"merchantId"`
	OrderID          string `json:"orderId"`
	PaymentRefID     string `json:"paymentRefId"`
	Amount           string `json:"amount"`
	ClientMobileNo   string `json:"clientMobileNo"`
	MerchantMobileNo string `json:"merchantMobileNo"`
	OrderDateTime    string `json:"orderDateTime"`
	IssuerPaymentDateTime string `json:"issuerPaymentDateTime"`
	IssuerPaymentRefNo    string `json:"issuerPaymentRefNo"`
	AdditionalMerchantInfo string `json:"additionalMerchantInfo"`
	Status           string `json:"status"`
	StatusCode       string `json:"statusCode"`
	ServiceType      string `json:"serviceType"`
	CallbackURL      string `json:"callbackURL"`
}

type CompletePaymentResponse struct {
	MerchantID       string `json:"merchantId"`
	OrderID          string `json:"orderId"`
	PaymentRefID     string `json:"paymentRefId"`
	Amount           string `json:"amount"`
	ClientMobileNo   string `json:"clientMobileNo"`
	MerchantMobileNo string `json:"merchantMobileNo"`
	OrderDateTime    string `json:"orderDateTime"`
	IssuerPaymentDateTime string `json:"issuerPaymentDateTime"`
	IssuerPaymentRefNo    string `json:"issuerPaymentRefNo"`
	AdditionalMerchantInfo string `json:"additionalMerchantInfo"`
	Status           string `json:"status"`
	StatusCode       string `json:"statusCode"`
	CancelIssuerDateTime string `json:"cancelIssuerDateTime"`
	CancelIssuerRefNo    string `json:"cancelIssuerRefNo"`
	ServiceType      string `json:"serviceType"`
}

type PaymentRequest struct {
	MerchantID      string `json:"merchantId"`
	OrderID         string `json:"orderId"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	CallbackURL     string `json:"callbackURL"`
	MerchantName    string `json:"merchantName,omitempty"`
	CustomerName    string `json:"customerName,omitempty"`
	CustomerMobile  string `json:"customerMobile,omitempty"`
	CustomerEmail   string `json:"customerEmail,omitempty"`
	AdditionalInfo  string `json:"additionalInfo,omitempty"`
}

type Adapter struct {
	creds      Credentials
	baseURL    string
	httpClient *http.Client
	callbackBase string
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewAdapter(creds Credentials, callbackBase string) (*Adapter, error) {
	privateKey, err := parsePrivateKey(creds.MerchantPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse merchant private key: %w", err)
	}

	publicKey, err := parsePublicKey(creds.PGPublicKey)
	if err != nil {
		return nil, fmt.Errorf("parse pg public key: %w", err)
	}

	baseURL := strings.TrimRight(creds.BaseURL, "/")

	return &Adapter{
		creds:        creds,
		baseURL:      baseURL,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		callbackBase: strings.TrimRight(callbackBase, "/"),
		privateKey:   privateKey,
		publicKey:    publicKey,
	}, nil
}

func parsePrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		keyBytes = []byte(keyStr)
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		parsedKey2, err2 := x509.ParsePKCS1PrivateKey(keyBytes)
		if err2 != nil {
			return nil, fmt.Errorf("parse private key (PKCS8: %v, PKCS1: %v)", err, err2)
		}
		return parsedKey2, nil
	}
	return parsedKey.(*rsa.PrivateKey), nil
}

func parsePublicKey(keyStr string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		keyBytes = []byte(keyStr)
	}
	parsedKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	pub, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA public key")
	}
	return pub, nil
}

func (a *Adapter) sign(data string) (string, error) {
	hash := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, a.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("rsa sign: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (a *Adapter) verifySignature(data, signature string) error {
	hash := sha256.Sum256([]byte(data))
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	return rsa.VerifyPKCS1v15(a.publicKey, crypto.SHA256, hash[:], sigBytes)
}

func generateOrderID() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}

func (a *Adapter) InitializePayment(ctx context.Context, req PaymentRequest) (*InitializePaymentResponse, error) {
	orderID := generateOrderID()

	payload := map[string]interface{}{
		"merchantId":  a.creds.MerchantID,
		"orderId":     orderID,
		"amount":      req.Amount,
		"currency":    req.Currency,
		"callbackURL": a.callbackBase + "/v1/nagad/callback",
	}

	if req.CustomerMobile != "" {
		payload["customerMobile"] = req.CustomerMobile
	}
	if req.CustomerName != "" {
		payload["customerName"] = req.CustomerName
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal init payment: %w", err)
	}

	signature, err := a.sign(string(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("sign init payment: %w", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(payloadBytes)

	initReq := map[string]string{
		"merchantId": a.creds.MerchantID,
		"orderId":    orderID,
		"payload":    encodedPayload,
		"signature":  signature,
	}

	body, err := json.Marshal(initReq)
	if err != nil {
		return nil, fmt.Errorf("marshal init request wrapper: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.baseURL+"/api/check-out/initialize/"+a.creds.MerchantID+"/"+orderID,
		bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create init request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Merchant-ID", a.creds.MerchantID)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do init payment: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read init response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("nagad init payment failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var initResp InitializePaymentResponse
	if err := json.Unmarshal(respBody, &initResp); err != nil {
		return nil, fmt.Errorf("parse init response: %w", err)
	}

	return &initResp, nil
}

func (a *Adapter) CompletePayment(ctx context.Context, paymentRefID string) (*CompletePaymentResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.baseURL+"/api/check-out/complete/"+paymentRefID,
		nil)
	if err != nil {
		return nil, fmt.Errorf("create complete request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Merchant-ID", a.creds.MerchantID)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do complete payment: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read complete response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("nagad complete payment failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var completeResp CompletePaymentResponse
	if err := json.Unmarshal(respBody, &completeResp); err != nil {
		return nil, fmt.Errorf("parse complete response: %w", err)
	}

	return &completeResp, nil
}

func (a *Adapter) VerifyCallbackSignature(payload, signature string) error {
	return a.verifySignature(payload, signature)
}

type RefundRequest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	TrxID    string `json:"trxId"`
	Reason   string `json:"reason,omitempty"`
}

type RefundResponse struct {
	MerchantID    string `json:"merchantId"`
	OrderID       string `json:"orderId"`
	RefundTrxID   string `json:"refundTrxId"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	StatusCode    string `json:"statusCode"`
	CompletedTime string `json:"completedTime"`
}

func (a *Adapter) RefundPayment(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	payload := map[string]interface{}{
		"merchantId": a.creds.MerchantID,
		"amount":     req.Amount,
		"currency":   req.Currency,
		"trxId":      req.TrxID,
		"reason":     req.Reason,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal refund payload: %w", err)
	}

	signature, err := a.sign(string(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("sign refund: %w", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(payloadBytes)

	refundReq := map[string]string{
		"merchantId": a.creds.MerchantID,
		"payload":    encodedPayload,
		"signature":  signature,
	}

	body, err := json.Marshal(refundReq)
	if err != nil {
		return nil, fmt.Errorf("marshal refund request wrapper: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.baseURL+"/api/check-out/refund", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create refund request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Merchant-ID", a.creds.MerchantID)

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do refund: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read refund response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("nagad refund failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result RefundResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse refund response: %w", err)
	}

	return &result, nil
}
