package handler

import (
	"errors"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomExerciseMuscularGroupHandler struct {
	service interfaces.CustomExerciseMuscularGroupService
}

func NewCustomExerciseMuscularGroupHandler(service interfaces.CustomExerciseMuscularGroupService) *CustomExerciseMuscularGroupHandler {
	return &CustomExerciseMuscularGroupHandler{service: service}
}

func (h *CustomExerciseMuscularGroupHandler) CreateLink(gymID string, link dto.CustomExerciseMuscularGroup) error {
	err := h.service.CreateLink(gymID, link)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			return apiErr
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create custom_exercise-muscular group link", err)
	}
	return nil
}

func (h *CustomExerciseMuscularGroupHandler) DeleteLink(gymID, id string) error {
	err := h.service.DeleteLink(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			return apiErr
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom_exercise-muscular group link", err)
	}
	return nil
}

func (h *CustomExerciseMuscularGroupHandler) RemoveAllLinksForExercise(gymID, id string) error {
	err := h.service.RemoveAllLinksForExercise(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			return apiErr
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to remove all links for exercise", err)
	}
	return nil
}

func (h *CustomExerciseMuscularGroupHandler) GetLinkByID(gymID, id string) (dto.CustomExerciseMuscularGroup, error) {
	link, err := h.service.GetLinkByID(gymID, id)
	if err != nil {
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			return dto.CustomExerciseMuscularGroup{}, apiErr
		}
		return dto.CustomExerciseMuscularGroup{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get link by ID", err)
	}
	return link, nil
}

func (h *CustomExerciseMuscularGroupHandler) GetLinksByCustomExerciseID(gymID, customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error) {
	return h.service.GetLinksByCustomExerciseID(gymID, customExerciseID)
}

func (h *CustomExerciseMuscularGroupHandler) GetLinksByMuscularGroupID(gymID, muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error) {
	return h.service.GetLinksByMuscularGroupID(gymID, muscularGroupID)
}
