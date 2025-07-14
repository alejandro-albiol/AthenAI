package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/go-chi/chi/v5"
)

// NewAuthRouter creates and configures auth routes
func NewAuthRouter(authHandler interfaces.AuthHandler) http.Handler {
	r := chi.NewRouter()

	// Authentication endpoints
	r.Post("/login", authHandler.Login)
	r.Post("/refresh", authHandler.RefreshToken)
	r.Post("/logout", authHandler.Logout)
	r.Get("/validate", authHandler.ValidateToken)

	return r
}
