package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api/middleware"
)

func RegisterAuthRoutes(r chi.Router, authService *AuthService, rateLimiter middleware.Limiter) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(rateLimiter))
		r.Post("/auth/login", LoginHandler(authService))
		r.Post("/auth/register", RegisterHandler(authService))
		r.Post("/auth/refresh", RefreshHandler(authService))
		r.Post("/auth/forgot-password", ForgotPasswordHandler(authService))
		r.Post("/auth/reset-password", ResetPasswordHandler(authService))
	})
	r.With(AuthMiddleware(authService)).Get("/auth/me", MeHandler(authService))
}
