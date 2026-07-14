package auth

import "github.com/go-chi/chi/v5"

func RegisterAuthRoutes(r chi.Router, authService *AuthService) {
	r.Post("/auth/login", LoginHandler(authService))
	r.Post("/auth/register", RegisterHandler(authService))
	r.Post("/auth/refresh", RefreshHandler(authService))
	r.With(AuthMiddleware(authService)).Get("/auth/me", MeHandler(authService))
}
