package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	authdto "github.com/alejandro-albiol/athenai/internal/auth/dto"
	authinterfaces "github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

type InvitationHandler struct {
	invitationService authinterfaces.InvitationServiceInterface
}

func NewInvitationHandler(invitationService authinterfaces.InvitationServiceInterface) authinterfaces.InvitationHandler {
	return &InvitationHandler{
		invitationService: invitationService,
	}
}

// CreateInvitation handles POST /api/v1/invitations
func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	// Only platform admins can create invitations
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Only platform administrators can create invitations",
			nil,
		))
		return
	}

	var req authdto.InvitationCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	// Get the creator's user ID from JWT
	creatorID := middleware.GetUserID(r)

	invitation, apiErr := h.invitationService.CreateInvitation(&req, creatorID)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Invitation created successfully", invitation)
}

// GetGymInvitations handles GET /api/v1/gyms/{gymId}/invitations
func (h *InvitationHandler) GetGymInvitations(w http.ResponseWriter, r *http.Request) {
	// Only platform admins can view all invitations
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Only platform administrators can view invitations",
			nil,
		))
		return
	}

	gymID := chi.URLParam(r, "gymId")
	if gymID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Gym ID is required",
			nil,
		))
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	status := r.URL.Query().Get("status") // pending, accepted, expired

	limit := 20 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	invitations, apiErr := h.invitationService.GetGymInvitations(gymID, limit, offset, status)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Invitations retrieved successfully", invitations)
}

// DecodeInvitation handles GET /api/v1/invitations/decode/{token}
func (h *InvitationHandler) DecodeInvitation(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invitation token is required",
			nil,
		))
		return
	}

	decoded, apiErr := h.invitationService.DecodeInvitation(token)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Invitation decoded successfully", decoded)
}

// AcceptInvitation handles POST /api/v1/invitations/accept
func (h *InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	var req authdto.InvitationAcceptRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invalid request payload",
			err,
		))
		return
	}

	user, apiErr := h.invitationService.AcceptInvitation(&req)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Invitation accepted successfully", user)
}

// ResendInvitation handles POST /api/v1/invitations/{id}/resend
func (h *InvitationHandler) ResendInvitation(w http.ResponseWriter, r *http.Request) {
	// Only platform admins can resend invitations
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Only platform administrators can resend invitations",
			nil,
		))
		return
	}

	invitationID := chi.URLParam(r, "id")
	if invitationID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invitation ID is required",
			nil,
		))
		return
	}

	apiErr := h.invitationService.ResendInvitation(invitationID)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Invitation resent successfully", nil)
}

// DeleteInvitation handles DELETE /api/v1/invitations/{id}
func (h *InvitationHandler) DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	// Only platform admins can delete invitations
	if !middleware.IsPlatformAdmin(r) {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeForbidden,
			"Only platform administrators can delete invitations",
			nil,
		))
		return
	}

	invitationID := chi.URLParam(r, "id")
	if invitationID == "" {
		response.WriteAPIError(w, apierror.New(
			errorcode_enum.CodeBadRequest,
			"Invitation ID is required",
			nil,
		))
		return
	}

	apiErr := h.invitationService.DeleteInvitation(invitationID)
	if apiErr != nil {
		response.WriteAPIError(w, apiErr)
		return
	}

	response.WriteAPISuccess(w, "Invitation deleted successfully", nil)
}
