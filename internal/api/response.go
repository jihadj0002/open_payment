package api

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
	Status  int    `json:"status"`
}

type PaginatedResponse struct {
	Data       []interface{} `json:"data"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
	HasMore    bool          `json:"has_more"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{Data: data})
}

func RespondError(w http.ResponseWriter, status int, errType, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Error: &APIError{
			Type:    errType,
			Code:    code,
			Message: message,
			Status:  status,
		},
	})
}

func RespondPaginated(w http.ResponseWriter, data []interface{}, total, page, perPage int) {
	totalPages := (total + perPage - 1) / perPage
	RespondJSON(w, http.StatusOK, PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	})
}
