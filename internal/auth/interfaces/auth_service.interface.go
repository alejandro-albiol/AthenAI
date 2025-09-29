package interfaces

import (
	"net/http"

	"github.com/alejandro-albiol/athenai/internal/auth/dto"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
)

// AuthServiceInterface defines authentication business logic with simplified approach
type AuthServiceInterface interface {
	// Single login method that checks X-Gym-ID header to determine authentication type
	Login(r *http.Request, loginReq *dto.LoginRequestDTO) (*dto.LoginResponseDTO, *apierror.APIError)

	// Token operations
	ValidateToken(token string) (*dto.TokenValidationResponseDTO, *apierror.APIError)
	RefreshToken(refreshReq *dto.RefreshTokenRequestDTO) (*dto.LoginResponseDTO, *apierror.APIError)
	Logout(logoutReq *dto.LogoutRequestDTO) *apierror.APIError
}

// InvitationServiceInterface defines invitation business logic
type InvitationServiceInterface interface {
	// CreateInvitation generates a new gym invitation
	CreateInvitation(req *dto.InvitationCreateRequestDTO, creatorID string) (*dto.InvitationResponseDTO, *apierror.APIError)

	// GetGymInvitations retrieves invitations for a specific gym
	GetGymInvitations(gymID string, limit, offset int, status string) (*dto.InvitationListResponseDTO, *apierror.APIError)

	// DecodeInvitation validates and decodes an invitation token
	DecodeInvitation(token string) (*dto.InvitationDecodeResponseDTO, *apierror.APIError)

	// AcceptInvitation processes invitation acceptance and creates user account
	AcceptInvitation(req *dto.InvitationAcceptRequestDTO) (*dto.LoginResponseDTO, *apierror.APIError)

	// ResendInvitation sends invitation email again
	ResendInvitation(invitationID string) *apierror.APIError

	// DeleteInvitation removes an invitation
	DeleteInvitation(invitationID string) *apierror.APIError
}
