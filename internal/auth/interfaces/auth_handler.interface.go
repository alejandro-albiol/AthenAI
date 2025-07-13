package interfaces

import "net/http"

// AuthHandler defines the HTTP layer interface
type AuthHandler interface {
	// Login handles POST /auth/login (routes based on subdomain)
	Login(w http.ResponseWriter, r *http.Request)

	// RefreshToken handles POST /auth/refresh
	RefreshToken(w http.ResponseWriter, r *http.Request)

	// Logout handles POST /auth/logout
	Logout(w http.ResponseWriter, r *http.Request)

	// ValidateToken handles GET /auth/validate
	ValidateToken(w http.ResponseWriter, r *http.Request)
}
