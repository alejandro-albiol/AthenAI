package router

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/go-chi/chi/v5"
)

// NewAuthRouter creates a new router for authentication endpoints
func NewAuthRouter(handler interfaces.AuthHandler, invitationHandler interfaces.InvitationHandler) http.Handler {
	r := chi.NewRouter()

	// Authentication endpoints
	r.Post("/login", handler.Login)           // POST /auth/login - Single login endpoint (checks X-Gym-ID header)
	r.Post("/refresh", handler.RefreshToken)  // POST /auth/refresh - Refresh access token using refresh token
	r.Post("/logout", handler.Logout)         // POST /auth/logout - Logout and revoke refresh token
	r.Get("/validate", handler.ValidateToken) // GET /auth/validate - Validate JWT token

	return r
}

// NewInvitationRouter creates a new router for invitation endpoints
func NewInvitationRouter(invitationHandler interfaces.InvitationHandler) http.Handler {
	r := chi.NewRouter()

	// Invitation endpoints
	r.Post("/invitation", invitationHandler.CreateInvitation)                // POST /invitation - Create new invitation
	r.Get("/gym/{gymId}/invitations", invitationHandler.GetGymInvitations)   // GET /gym/{gymId}/invitations - Get all invitations for a gym
	r.Get("/invitation/decode/{token}", invitationHandler.DecodeInvitation)  // GET /invitation/decode/{token} - Decode invitation token
	r.Post("/invitation/accept/{token}", invitationHandler.AcceptInvitation) // POST /invitation/accept/{token} - Accept invitation
	r.Post("/invitation/{id}/resend", invitationHandler.ResendInvitation)    // POST /invitation/{id}/resend - Resend invitation
	r.Delete("/invitation/{id}", invitationHandler.DeleteInvitation)         // DELETE /invitation/{id} - Delete invitation

	return r
}
