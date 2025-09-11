package service

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomWorkoutInstanceService struct {
	Repo interfaces.CustomWorkoutInstanceRepository
}

func NewCustomWorkoutInstanceService(repo interfaces.CustomWorkoutInstanceRepository) *CustomWorkoutInstanceService {
	return &CustomWorkoutInstanceService{Repo: repo}
}

func (s *CustomWorkoutInstanceService) CreateCustomWorkoutInstance(gymID, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error) {
	// Validate template source and IDs
	if instance.TemplateSource == "public" && instance.PublicTemplateID == nil {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "public_template_id is required when template_source is 'public'", nil)
	}
	if instance.TemplateSource == "gym" && instance.GymTemplateID == nil {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "gym_template_id is required when template_source is 'gym'", nil)
	}

	id, err := s.Repo.Create(gymID, createdBy, instance)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create workout instance", err)
	}
	return id, nil
}

func (s *CustomWorkoutInstanceService) GetCustomWorkoutInstanceByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
	instance, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get workout instance", err)
	}
	return instance, nil
}

func (s *CustomWorkoutInstanceService) GetCustomWorkoutInstanceSummaryByID(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error) {
	instance, err := s.Repo.GetSummaryByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get workout instance", err)
	}
	return instance, nil
}

func (s *CustomWorkoutInstanceService) GetCustomWorkoutInstancesByUserID(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	instances, err := s.Repo.GetByUserID(gymID, userID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get workout instances", err)
	}
	return instances, nil
}

func (s *CustomWorkoutInstanceService) GetCustomWorkoutInstanceSummariesByUserID(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	instances, err := s.Repo.GetSummariesByUserID(gymID, userID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get workout instances", err)
	}
	return instances, nil
}

func (s *CustomWorkoutInstanceService) GetLastCustomWorkoutInstancesByUserID(gymID, userID string, numberOfWorkouts int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	if numberOfWorkouts <= 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "numberOfWorkouts must be greater than 0", nil)
	}
	if numberOfWorkouts > 100 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "numberOfWorkouts cannot exceed 100", nil)
	}

	instances, err := s.Repo.GetLastsByUserID(gymID, userID, numberOfWorkouts)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get workout instances", err)
	}
	return instances, nil
}

func (s *CustomWorkoutInstanceService) ListCustomWorkoutInstances(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	instances, err := s.Repo.List(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout instances", err)
	}
	return instances, nil
}

func (s *CustomWorkoutInstanceService) ListCustomWorkoutInstanceSummaries(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	instances, err := s.Repo.ListSummaries(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout instances", err)
	}
	return instances, nil
}

func (s *CustomWorkoutInstanceService) UpdateCustomWorkoutInstance(gymID, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error {
	// Validate template source and IDs if being updated
	if instance.TemplateSource != nil {
		if *instance.TemplateSource == "public" && instance.PublicTemplateID == nil {
			return apierror.New(errorcode_enum.CodeBadRequest, "public_template_id is required when template_source is 'public'", nil)
		}
		if *instance.TemplateSource == "gym" && instance.GymTemplateID == nil {
			return apierror.New(errorcode_enum.CodeBadRequest, "gym_template_id is required when template_source is 'gym'", nil)
		}
	}

	err := s.Repo.Update(gymID, id, instance)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update workout instance", err)
	}
	return nil
}

func (s *CustomWorkoutInstanceService) DeleteCustomWorkoutInstance(gymID, id string) error {
	err := s.Repo.Delete(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete workout instance", err)
	}
	return nil
}
