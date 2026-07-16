package payment

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/openpayment/gateway/internal/api"
)

func HandleBkashCallback(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bkashPaymentID := r.URL.Query().Get("paymentID")
		status := r.URL.Query().Get("status")

		if bkashPaymentID == "" || status == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_params",
				"paymentID and status are required")
			return
		}

		pi, err := svc.repo.GetPaymentIntentByProviderRef(r.Context(), bkashPaymentID)
		if err != nil {
			api.RespondError(w, http.StatusNotFound, "not_found", "payment_not_found",
				"payment not found for bKash ref: "+bkashPaymentID)
			return
		}

		if status == "success" {
			if err := svc.HandleWalletCallback(r.Context(), pi.ID, status, bkashPaymentID); err != nil {
				redirectURL := pi.CancelURL
				if redirectURL != nil && *redirectURL != "" {
					http.Redirect(w, r, *redirectURL+"?status=failed&error="+err.Error(), http.StatusTemporaryRedirect)
					return
				}
				api.RespondError(w, http.StatusInternalServerError, "payment_error", "callback_failed", err.Error())
				return
			}

			if pi.ReturnURL != nil && *pi.ReturnURL != "" {
				http.Redirect(w, r, *pi.ReturnURL+"?status=success&payment_intent="+pi.ID, http.StatusTemporaryRedirect)
				return
			}
			http.Redirect(w, r, "/checkout/"+pi.ID+"/success", http.StatusTemporaryRedirect)
			return
		}

		if status == "cancel" {
			svc.HandleWalletCallback(r.Context(), pi.ID, "cancel", "")
			if pi.CancelURL != nil && *pi.CancelURL != "" {
				http.Redirect(w, r, *pi.CancelURL+"?status=canceled&payment_intent="+pi.ID, http.StatusTemporaryRedirect)
				return
			}
			http.Redirect(w, r, "/checkout/"+pi.ID+"/cancel", http.StatusTemporaryRedirect)
			return
		}

		if pi.CancelURL != nil && *pi.CancelURL != "" {
			http.Redirect(w, r, *pi.CancelURL+"?status=failed&payment_intent="+pi.ID, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, "/checkout/"+pi.ID+"/error", http.StatusTemporaryRedirect)
	}
}

type NagadCallbackPayload struct {
	MerchantID       string `json:"merchantId"`
	OrderID          string `json:"orderId"`
	PaymentRefID     string `json:"paymentRefId"`
	Amount           string `json:"amount"`
	ClientMobileNo   string `json:"clientMobileNo"`
	Status           string `json:"status"`
	StatusCode       string `json:"statusCode"`
	IssuerPaymentDateTime string `json:"issuerPaymentDateTime,omitempty"`
	IssuerPaymentRefNo    string `json:"issuerPaymentRefNo,omitempty"`
	CancelIssuerDateTime  string `json:"cancelIssuerDateTime,omitempty"`
	CancelIssuerRefNo     string `json:"cancelIssuerRefNo,omitempty"`
	Signature        string `json:"signature,omitempty"`
}

func HandleNagadCallback(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "read_error", "failed to read body")
			return
		}
		defer r.Body.Close()

		var callback NagadCallbackPayload
		if err := json.Unmarshal(body, &callback); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "parse_error", "invalid JSON")
			return
		}

		pi, err := svc.GetPaymentByID(r.Context(), callback.OrderID)
		if err != nil {
			api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
			return
		}

		if callback.Status == "Success" || callback.StatusCode == "Success" {
			providerRef := callback.PaymentRefID
			if providerRef == "" {
				providerRef = callback.IssuerPaymentRefNo
			}
			if err := svc.HandleWalletCallback(r.Context(), pi.ID, "success", providerRef); err != nil {
				api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
				return
			}
		} else if callback.Status == "Cancel" || callback.StatusCode == "Cancel" {
			svc.HandleWalletCallback(r.Context(), pi.ID, "cancel", "")
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

type BkashSNSNotification struct {
	Type             string `json:"Type"`
	Message          string `json:"Message"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Subject          string `json:"Subject,omitempty"`
	SubscribeURL     string `json:"SubscribeURL,omitempty"`
	Timestamp        string `json:"Timestamp"`
	Token            string `json:"Token,omitempty"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL,omitempty"`
}

func buildSNSStringToSign(notif *BkashSNSNotification) string {
	fields := []struct {
		name  string
		value string
	}{
		{"Message", notif.Message},
		{"MessageId", notif.MessageId},
		{"Subject", notif.Subject},
		{"SubscribeURL", notif.SubscribeURL},
		{"Timestamp", notif.Timestamp},
		{"Token", notif.Token},
		{"TopicArn", notif.TopicArn},
		{"Type", notif.Type},
	}

	var parts []string
	for _, f := range fields {
		if f.value != "" {
			parts = append(parts, f.name+"\n"+f.value)
		}
	}
	return strings.Join(parts, "\n") + "\n"
}

func verifySNSMessage(notif *BkashSNSNotification) error {
	certResp, err := http.Get(notif.SigningCertURL)
	if err != nil {
		return fmt.Errorf("fetch signing cert: %w", err)
	}
	defer certResp.Body.Close()

	certBytes, err := io.ReadAll(certResp.Body)
	if err != nil {
		return fmt.Errorf("read signing cert: %w", err)
	}

	block, _ := pem.Decode(certBytes)
	if block == nil {
		certBytes = []byte("-----BEGIN CERTIFICATE-----\n" + base64.StdEncoding.EncodeToString(certBytes) + "\n-----END CERTIFICATE-----")
		block, _ = pem.Decode(certBytes)
		if block == nil {
			return fmt.Errorf("failed to decode certificate PEM")
		}
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse certificate: %w", err)
	}

	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("certificate public key is not RSA")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(notif.Signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	stringToSign := buildSNSStringToSign(notif)

	if notif.SignatureVersion == "1" {
		hash := sha1.Sum([]byte(stringToSign))
		if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hash[:], sigBytes); err != nil {
			return fmt.Errorf("SNS signature verification failed (SHA1): %w", err)
		}
	} else if notif.SignatureVersion == "2" {
		hash := sha256.Sum256([]byte(stringToSign))
		if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes); err != nil {
			return fmt.Errorf("SNS signature verification failed (SHA256): %w", err)
		}
	} else {
		return fmt.Errorf("unsupported SNS signature version: %s", notif.SignatureVersion)
	}

	return nil
}

func HandleBkashWebhook(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "read_error", "failed to read body")
			return
		}
		defer r.Body.Close()

		var notif BkashSNSNotification
		if err := json.Unmarshal(body, &notif); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "parse_error", "invalid SNS JSON")
			return
		}

		if err := verifySNSMessage(&notif); err != nil {
			api.RespondError(w, http.StatusUnauthorized, "verification_error", "sns_signature_invalid", err.Error())
			return
		}

		snsMsgType := r.Header.Get("x-amz-sns-message-type")

		if snsMsgType == "SubscriptionConfirmation" {
			if notif.SubscribeURL != "" {
				go func() {
					http.Get(notif.SubscribeURL)
				}()
			}
			api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
			return
		}

		if snsMsgType == "Notification" {
			var msgData map[string]interface{}
			if err := json.Unmarshal([]byte(notif.Message), &msgData); err != nil {
				api.RespondError(w, http.StatusBadRequest, "invalid_request", "parse_error", "invalid message JSON")
				return
			}

			merchantInvoiceNumber, _ := msgData["merchantInvoiceNumber"].(string)
			transactionStatus, _ := msgData["transactionStatus"].(string)
			trxID, _ := msgData["trxID"].(string)

			if merchantInvoiceNumber == "" {
				api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
				return
			}

			if transactionStatus == "Completed" {
				svc.HandleWalletCallback(r.Context(), merchantInvoiceNumber, "success", trxID)
			}
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

type BkashCallbackError struct {
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

func HandleBkashCallbackURL(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID := r.URL.Query().Get("paymentID")
		status := r.URL.Query().Get("status")

		fmt.Printf("bKash callback: paymentID=%s, status=%s\n", paymentID, status)

		if paymentID == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_payment_id", "paymentID is required")
			return
		}

		HandleBkashCallback(svc).ServeHTTP(w, r)
	}
}
