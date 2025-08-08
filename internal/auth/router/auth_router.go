package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/go-chi/chi/v5"
)

// NewAuthRouter creates a new router for authentication endpoints
func NewAuthRouter(handler interfaces.AuthHandler) http.Handler {
	r := chi.NewRouter()

	// Authentication endpoints
	r.Post("/login", handler.Login)           // POST /auth/login - Single login endpoint (checks X-Gym-ID header)
	r.Post("/refresh", handler.RefreshToken)  // POST /auth/refresh - Refresh access token using refresh token
	r.Post("/logout", handler.Logout)         // POST /auth/logout - Logout and revoke refresh token
	r.Get("/validate", handler.ValidateToken) // GET /auth/validate - Validate JWT token

	return r
}
