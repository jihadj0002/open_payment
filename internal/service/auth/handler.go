package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/openpayment/gateway/internal/api"
	"github.com/rs/zerolog/log"
)

func LoginHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if req.Email == "" || req.Password == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_fields", "email and password are required")
			return
		}

		result, err := authService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			if err == ErrInvalidCredentials {
				api.RespondError(w, http.StatusUnauthorized, "auth_error", "invalid_credentials", "invalid email or password")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, result)
	}
}

func RegisterHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_fields", "name, email, and password are required")
			return
		}

		result, err := authService.Register(r.Context(), req)
		if err != nil {
			if err == ErrEmailAlreadyExists {
				api.RespondError(w, http.StatusConflict, "invalid_request", "email_exists", "a merchant with this email already exists")
				return
			}
			log.Error().Err(err).Msg("register failed")
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusCreated, result)
	}
}

func ForgotPasswordHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if body.Email == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_fields", "email is required")
			return
		}

		if err := authService.ForgotPassword(r.Context(), body.Email); err != nil {
			log.Error().Err(err).Msg("forgot password failed")
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"message": "if the email exists, a reset link has been sent"})
	}
}

func ResetPasswordHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if body.Token == "" || body.Password == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_fields", "token and password are required")
			return
		}

		if len(body.Password) < 6 {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "weak_password", "password must be at least 6 characters")
			return
		}

		if err := authService.ResetPassword(r.Context(), body.Token, body.Password); err != nil {
			if errors.Is(err, ErrInvalidResetToken) || errors.Is(err, ErrResetTokenUsed) {
				api.RespondError(w, http.StatusBadRequest, "auth_error", "invalid_token", err.Error())
				return
			}
			log.Error().Err(err).Msg("reset password failed")
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, map[string]string{"message": "password has been reset successfully"})
	}
}

func RefreshHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		if body.RefreshToken == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_field", "refresh_token is required")
			return
		}

		tokenPair, err := authService.RefreshAccessToken(r.Context(), body.RefreshToken)
		if err != nil {
			if errors.Is(err, ErrInvalidToken) {
				api.RespondError(w, http.StatusUnauthorized, "auth_error", "invalid_token", "invalid or expired refresh token")
				return
			}
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		api.RespondJSON(w, http.StatusOK, tokenPair)
	}
}

func MeHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		api.RespondJSON(w, http.StatusOK, claims)
	}
}
