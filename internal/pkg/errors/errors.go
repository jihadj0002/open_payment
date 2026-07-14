package errors

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrValidation      = errors.New("validation error")
	ErrConflict        = errors.New("resource conflict")
	ErrRateLimited     = errors.New("rate limited")
	ErrInternal        = errors.New("internal server error")
	ErrIdempotency     = errors.New("idempotency conflict")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrDuplicateRequest = errors.New("duplicate request")
)

type DomainError struct {
	Err     error
	Message string
	Code    string
	Status  int
}

func (e *DomainError) Error() string {
	return e.Message
}
