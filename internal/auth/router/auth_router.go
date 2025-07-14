package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/go-chi/chi/v5"
)

// NewAuthRouter creates a new router for authentication endpoints
func NewAuthRouter(handler interfaces.AuthHandler) http.Handler {
	r := chi.NewRouter()

	// POST /auth/login - Single login endpoint (checks X-Gym-ID header)
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handler.Login(w, r)
	})

	// POST /auth/refresh - Refresh access token using refresh token
	r.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		handler.RefreshToken(w, r)
	})

	// POST /auth/logout - Logout and revoke refresh token
	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		handler.Logout(w, r)
	})

	// GET /auth/validate - Validate JWT token
	r.Get("/validate", func(w http.ResponseWriter, r *http.Request) {
		handler.ValidateToken(w, r)
	})

	return r
}
