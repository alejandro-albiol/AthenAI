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

// InvitationHandler defines the invitation HTTP layer interface
type InvitationHandler interface {
	// CreateInvitation handles POST /invitations
	CreateInvitation(w http.ResponseWriter, r *http.Request)

	// GetGymInvitations handles GET /gyms/{gymId}/invitations
	GetGymInvitations(w http.ResponseWriter, r *http.Request)

	// DecodeInvitation handles GET /invitations/decode/{token}
	DecodeInvitation(w http.ResponseWriter, r *http.Request)

	// AcceptInvitation handles POST /invitations/accept
	AcceptInvitation(w http.ResponseWriter, r *http.Request)

	// ResendInvitation handles POST /invitations/{id}/resend
	ResendInvitation(w http.ResponseWriter, r *http.Request)

	// DeleteInvitation handles DELETE /invitations/{id}
	DeleteInvitation(w http.ResponseWriter, r *http.Request)
}
