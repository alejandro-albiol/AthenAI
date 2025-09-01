package service

import (
	"database/sql"
	"errors"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomWorkoutTemplateServiceImpl struct {
	Repo interfaces.CustomWorkoutTemplateRepository
}

func NewCustomWorkoutTemplateService(repo interfaces.CustomWorkoutTemplateRepository) *CustomWorkoutTemplateServiceImpl {
	return &CustomWorkoutTemplateServiceImpl{Repo: repo}
}

func (s *CustomWorkoutTemplateServiceImpl) CreateCustomWorkoutTemplate(gymID string, template dto.CreateCustomWorkoutTemplateDTO) (string, error) {
	existing, err := s.Repo.GetByName(gymID, template.Name)
	if err == nil && existing.ID != "" {
		return "", apierror.New(errorcode_enum.CodeConflict, "Custom workout template with this name already exists", nil)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to check template name existence", err)
	}
	id, err := s.Repo.Create(gymID, template)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create custom workout template", err)
	}
	return id, nil
}

func (s *CustomWorkoutTemplateServiceImpl) GetCustomWorkoutTemplateByID(gymID, id string) (dto.ResponseCustomWorkoutTemplateDTO, error) {
	template, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || template.ID == "" {
			return dto.ResponseCustomWorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Custom workout template not found", err)
		}
		return dto.ResponseCustomWorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get custom workout template", err)
	}
	return template, nil
}

func (s *CustomWorkoutTemplateServiceImpl) ListCustomWorkoutTemplates(gymID string) ([]dto.ResponseCustomWorkoutTemplateDTO, error) {
	templates, err := s.Repo.List(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom workout templates", err)
	}
	if len(templates) == 0 {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "No custom workout templates found", nil)
	}
	return templates, nil
}

func (s *CustomWorkoutTemplateServiceImpl) UpdateCustomWorkoutTemplate(gymID string, template dto.UpdateCustomWorkoutTemplateDTO) error {
	existing, err := s.Repo.GetByID(gymID, template.ID)
	if err != nil || existing.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, "Custom workout template not found", err)
	}
	// Optionally check for duplicate name if updating name
	if template.Name != nil && *template.Name != existing.Name {
		dup, err := s.Repo.GetByName(gymID, *template.Name)
		if err == nil && dup.ID != "" {
			return apierror.New(errorcode_enum.CodeConflict, "Custom workout template with this name already exists", nil)
		}
	}
	err = s.Repo.Update(gymID, template)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update custom workout template", err)
	}
	return nil
}

func (s *CustomWorkoutTemplateServiceImpl) DeleteCustomWorkoutTemplate(gymID, id string) error {
	existing, err := s.Repo.GetByID(gymID, id)
	if err != nil || existing.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, "Custom workout template not found", err)
	}
	err = s.Repo.Delete(gymID, id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom workout template", err)
	}
	return nil
}
