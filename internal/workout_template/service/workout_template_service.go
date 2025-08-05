package service

import (
	"github.com/alejandro-albiol/athenai/internal/workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/workout_template/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type WorkoutTemplateService struct {
	Repo interfaces.WorkoutTemplateRepository
}

// CreateWorkoutTemplate creates a new workout template in the system.
func (s *WorkoutTemplateService) CreateWorkoutTemplate(input dto.CreateWorkoutTemplateDTO) error {
	existingTemplate, err := s.Repo.GetByName(input.Name)
	if err == nil && existingTemplate.ID != "" {
		return apierror.New(errorcode_enum.CodeConflict, "Workout template already exists", nil)
	}
	if err != nil && err.Error() != "not found" {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check existing workout template", err)
	}

	err = s.Repo.Create(input)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create workout template", err)
	}
	return nil
}

// GetWorkoutTemplateByID retrieves a workout template by its unique ID.
func (s *WorkoutTemplateService) GetWorkoutTemplateByID(id string) (dto.WorkoutTemplateDTO, error) {
	template, err := s.Repo.GetByID(id)
	if template.ID == "" && err != nil {
		return dto.WorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Workout template not found", err)
	}
	if err != nil {
		return dto.WorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve workout template by ID", err)
	}
	return template, nil
}

// GetWorkoutTemplateByName retrieves a workout template by its name.
func (s *WorkoutTemplateService) GetWorkoutTemplateByName(name string) (dto.WorkoutTemplateDTO, error) {
	template, err := s.Repo.GetByName(name)
	if template.ID == "" && err != nil {
		return dto.WorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Workout template not found", err)
	}
	if err != nil {
		return dto.WorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve workout template by name", err)
	}
	return template, nil
}

// GetWorkoutTemplatesByDifficulty retrieves all workout templates matching a given difficulty level.
func (s *WorkoutTemplateService) GetWorkoutTemplatesByDifficulty(difficulty string) ([]dto.WorkoutTemplateDTO, error) {
	templates, err := s.Repo.GetByDifficulty(difficulty)
	if templates == nil && err != nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "There are no workout templates found for the given difficulty", err)
	}
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout templates by difficulty", err)
	}
	return templates, nil
}

// GetWorkoutTemplatesByTargetAudience retrieves all workout templates for a specific target audience.
func (s *WorkoutTemplateService) GetWorkoutTemplatesByTargetAudience(targetAudience string) ([]dto.WorkoutTemplateDTO, error) {
	templates, err := s.Repo.GetByTargetAudience(targetAudience)
	if templates == nil && err != nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "There are no workout templates found for the given target audience", err)
	}
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout templates by target audience", err)
	}
	return templates, nil
}

// GetAllWorkoutTemplates retrieves all workout templates in the system.
func (s *WorkoutTemplateService) GetAllWorkoutTemplates() ([]dto.WorkoutTemplateDTO, error) {
	templates, err := s.Repo.GetAll()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list all workout templates", err)
	}
	return templates, nil
}

// UpdateWorkoutTemplate updates an existing workout template by ID.
func (s *WorkoutTemplateService) UpdateWorkoutTemplate(id string, input dto.UpdateWorkoutTemplateDTO) (dto.WorkoutTemplateDTO, error) {
	findExistingTemplate, err := s.Repo.GetByID(id)
	if findExistingTemplate.ID == "" && err != nil {
		return dto.WorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Workout template not found", err)
	}

	template, err := s.Repo.Update(id, input)
	if err != nil {
		return dto.WorkoutTemplateDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update workout template", err)
	}
	return template, nil
}

// DeleteWorkoutTemplate deletes a workout template by ID.
func (s *WorkoutTemplateService) DeleteWorkoutTemplate(id string) error {
	err := s.Repo.Delete(id)
	if err != nil {
		return apierror.New("WT_DELETE_FAILED", "Failed to delete workout template", err)
	}
	return nil
}
