package ledger

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordPayment(ctx context.Context, merchantID, transactionID, currency string, amount, fee int64) error {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
		}
	}()

	bal, err := s.repo.GetCurrentBalanceTx(ctx, tx, merchantID, currency)
	if err != nil {
		return fmt.Errorf("get current balance: %w", err)
	}

	paymentIn := &Entry{
		TransactionID: transactionID,
		MerchantID:    merchantID,
		EntryType:     EntryTypePaymentIn,
		Amount:        amount,
		Currency:      currency,
		BalanceBefore: bal,
		BalanceAfter:  bal + amount,
		Description:   "payment received",
	}
	if err := s.repo.CreateEntryTx(ctx, tx, paymentIn); err != nil {
		return fmt.Errorf("create payment entry: %w", err)
	}

	if fee > 0 {
		feeEntry := &Entry{
			TransactionID: transactionID,
			MerchantID:    merchantID,
			EntryType:     EntryTypeFee,
			Amount:        fee,
			Currency:      currency,
			BalanceBefore: paymentIn.BalanceAfter,
			BalanceAfter:  paymentIn.BalanceAfter - fee,
			Description:   "processing fee",
		}
		if err := s.repo.CreateEntryTx(ctx, tx, feeEntry); err != nil {
			return fmt.Errorf("create fee entry: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}

func (s *Service) RecordRefund(ctx context.Context, merchantID, transactionID, currency string, amount int64) error {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
		}
	}()

	bal, err := s.repo.GetCurrentBalanceTx(ctx, tx, merchantID, currency)
	if err != nil {
		return fmt.Errorf("get current balance: %w", err)
	}

	refund := &Entry{
		TransactionID: transactionID,
		MerchantID:    merchantID,
		EntryType:     EntryTypeRefundOut,
		Amount:        amount,
		Currency:      currency,
		BalanceBefore: bal,
		BalanceAfter:  bal - amount,
		Description:   "refund issued",
	}
	if err := s.repo.CreateEntryTx(ctx, tx, refund); err != nil {
		return fmt.Errorf("create refund entry: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}

func (s *Service) RecordChargeback(ctx context.Context, merchantID, transactionID, currency string, amount int64) error {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
		}
	}()

	bal, err := s.repo.GetCurrentBalanceTx(ctx, tx, merchantID, currency)
	if err != nil {
		return fmt.Errorf("get current balance: %w", err)
	}

	chargeback := &Entry{
		TransactionID: transactionID,
		MerchantID:    merchantID,
		EntryType:     EntryTypeChargeback,
		Amount:        amount,
		Currency:      currency,
		BalanceBefore: bal,
		BalanceAfter:  bal - amount,
		Description:   "chargeback",
	}
	if err := s.repo.CreateEntryTx(ctx, tx, chargeback); err != nil {
		return fmt.Errorf("create chargeback entry: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}

func (s *Service) RecordSettlement(ctx context.Context, merchantID, currency string, amount int64) error {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
		}
	}()

	bal, err := s.repo.GetCurrentBalanceTx(ctx, tx, merchantID, currency)
	if err != nil {
		return fmt.Errorf("get current balance: %w", err)
	}

	settlement := &Entry{
		MerchantID:    merchantID,
		EntryType:     EntryTypeSettlement,
		Amount:        amount,
		Currency:      currency,
		BalanceBefore: bal,
		BalanceAfter:  bal - amount,
		Description:   "settlement payout",
	}
	if err := s.repo.CreateEntryTx(ctx, tx, settlement); err != nil {
		return fmt.Errorf("create settlement entry: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}

func (s *Service) GetBalance(ctx context.Context, merchantID, currency string) (*Balance, error) {
	available, err := s.repo.GetCurrentBalance(ctx, merchantID, currency)
	if err != nil {
		return nil, fmt.Errorf("get available balance: %w", err)
	}

	pending, err := s.repo.GetPendingBalance(ctx, merchantID, currency)
	if err != nil {
		return nil, fmt.Errorf("get pending balance: %w", err)
	}

	reserve := available * 5 / 100

	return &Balance{
		MerchantID: merchantID,
		Currency:   currency,
		Available:  available,
		Pending:    pending,
		Reserve:    reserve,
	}, nil
}

func (s *Service) GetBalanceTransactions(ctx context.Context, merchantID, currency string, limit, offset int) ([]BalanceTransaction, error) {
	entries, err := s.repo.GetEntries(ctx, merchantID, currency, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get entries: %w", err)
	}

	txs := make([]BalanceTransaction, len(entries))
	for i, e := range entries {
		txs[i] = BalanceTransaction{
			ID:          e.ID,
			Type:        e.EntryType,
			Amount:      e.Amount,
			Currency:    e.Currency,
			Description: e.Description,
			CreatedAt:   e.CreatedAt.Unix(),
		}
	}
	return txs, nil
}
