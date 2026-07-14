package customer

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidCardNumber = errors.New("card number is required for card payment method")
	ErrInvalidExpiry     = errors.New("expiration month and year are required")
	ErrInvalidWallet     = errors.New("wallet type and phone are required for wallet payment method")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCustomer(ctx context.Context, merchantID string, req CreateCustomerRequest) (*Customer, error) {
	c := &Customer{
		ID:         uuid.New().String(),
		MerchantID: merchantID,
		Email:      req.Email,
		Phone:      req.Phone,
		Name:       req.Name,
		Metadata:   req.Metadata,
	}

	if c.Metadata == nil {
		c.Metadata = map[string]string{}
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}

	return c, nil
}

func (s *Service) GetCustomer(ctx context.Context, id, merchantID string) (*Customer, error) {
	return s.repo.GetByID(ctx, id, merchantID)
}

func (s *Service) ListCustomers(ctx context.Context, merchantID string, page, perPage int) ([]Customer, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	return s.repo.List(ctx, merchantID, perPage, offset)
}

func (s *Service) UpdateCustomer(ctx context.Context, id, merchantID string, req UpdateCustomerRequest) (*Customer, error) {
	return s.repo.Update(ctx, id, merchantID, req)
}

func detectBrand(cardNumber string) string {
	if len(cardNumber) == 0 {
		return "unknown"
	}
	switch cardNumber[0] {
	case '4':
		return "visa"
	case '5':
		return "mastercard"
	case '3':
		return "amex"
	default:
		return "unknown"
	}
}

func generateFingerprint(cardNumber string) string {
	h := sha256.Sum256([]byte(cardNumber))
	return fmt.Sprintf("%x", h)
}

func (s *Service) AttachPaymentMethod(ctx context.Context, merchantID, customerID string, req AttachPaymentMethodRequest) (*PaymentMethod, error) {
	if _, err := s.repo.GetByID(ctx, customerID, merchantID); err != nil {
		return nil, err
	}

	pm := &PaymentMethod{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Type:       req.Type,
		IsActive:   true,
	}

	switch req.Type {
	case "card":
		if req.CardNumber == nil || *req.CardNumber == "" {
			return nil, ErrInvalidCardNumber
		}
		if req.ExpMonth == nil || req.ExpYear == nil {
			return nil, ErrInvalidExpiry
		}

		cardNum := strings.TrimSpace(*req.CardNumber)
		last4 := cardNum[len(cardNum)-4:]
		fingerprint := generateFingerprint(cardNum)

		token := uuid.New().String()

		pm.Last4 = &last4
		pm.Brand = strPtr(detectBrand(cardNum))
		pm.Fingerprint = &fingerprint
		pm.Token = &token
		pm.ExpMonth = req.ExpMonth
		pm.ExpYear = req.ExpYear
		pm.CardholderName = req.CardholderName

	case "wallet":
		if req.WalletType == nil || *req.WalletType == "" || req.WalletPhone == nil || *req.WalletPhone == "" {
			return nil, ErrInvalidWallet
		}
		pm.WalletType = req.WalletType
		pm.WalletPhone = req.WalletPhone

	default:
		return nil, fmt.Errorf("unsupported payment method type: %s", req.Type)
	}

	if req.SetAsDefault {
		if err := s.repo.UnsetDefaultPaymentMethods(ctx, customerID, merchantID); err != nil {
			return nil, err
		}
		pm.IsDefault = true
	}

	if err := s.repo.CreatePaymentMethod(ctx, merchantID, pm); err != nil {
		return nil, fmt.Errorf("create payment method: %w", err)
	}

	return pm, nil
}

func (s *Service) ListPaymentMethods(ctx context.Context, customerID, merchantID string) ([]PaymentMethod, error) {
	if _, err := s.repo.GetByID(ctx, customerID, merchantID); err != nil {
		return nil, err
	}

	return s.repo.ListPaymentMethods(ctx, customerID, merchantID)
}

func (s *Service) DetachPaymentMethod(ctx context.Context, customerID, merchantID, pmID string) error {
	if _, err := s.repo.GetByID(ctx, customerID, merchantID); err != nil {
		return err
	}

	return s.repo.DeactivatePaymentMethod(ctx, pmID, customerID, merchantID)
}

func strPtr(s string) *string { return &s }
