package api

import (
	"encoding/json"
	"net/http"
)

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type StructuredError struct {
	Type    string        `json:"type"`
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Status  int           `json:"status"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ValidationError struct {
	Message string
	Details []ErrorDetail
}

func (e *ValidationError) Error() string { return e.Message }

type NotFoundError struct {
	Resource string
	Message  string
}

func (e *NotFoundError) Error() string { return e.Message }

type AuthError struct {
	Message string
	Code    string
}

func (e *AuthError) Error() string { return e.Message }

type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string { return e.Message }

func (e *StructuredError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func RespondStructuredError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var se StructuredError

	switch e := err.(type) {
	case *ValidationError:
		se = StructuredError{
			Type:    "validation_error",
			Code:    "validation_failed",
			Message: e.Message,
			Status:  http.StatusBadRequest,
			Details: e.Details,
		}
	case *NotFoundError:
		se = StructuredError{
			Type:    "not_found",
			Code:    e.Resource + "_not_found",
			Message: e.Message,
			Status:  http.StatusNotFound,
		}
	case *AuthError:
		code := e.Code
		if code == "" {
			code = "unauthorized"
		}
		se = StructuredError{
			Type:    "auth_error",
			Code:    code,
			Message: e.Message,
			Status:  http.StatusUnauthorized,
		}
	case *RateLimitError:
		se = StructuredError{
			Type:    "rate_limit_error",
			Code:    "rate_limited",
			Message: e.Message,
			Status:  http.StatusTooManyRequests,
		}
	default:
		se = StructuredError{
			Type:    "server_error",
			Code:    "internal_error",
			Message: "an unexpected error occurred",
			Status:  http.StatusInternalServerError,
		}
	}

	w.WriteHeader(se.Status)
	json.NewEncoder(w).Encode(APIResponse{Error: &APIError{
		Type:    se.Type,
		Code:    se.Code,
		Message: se.Message,
		Status:  se.Status,
	}})
}
