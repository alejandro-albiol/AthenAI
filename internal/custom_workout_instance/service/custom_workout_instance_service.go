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

func (s *CustomWorkoutInstanceService) CreateCustomWorkoutInstance(gymID string, instance dto.CreateCustomWorkoutInstanceDTO) (string, error) {
	// No duplicate check needed for instance name, but could be added if required
	id, err := s.Repo.Create(gymID, instance)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create workout instance", err)
	}
	return id, nil
}

func (s *CustomWorkoutInstanceService) GetCustomWorkoutInstanceByID(gymID, id string) (dto.ResponseCustomWorkoutInstanceDTO, error) {
	instance, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.ResponseCustomWorkoutInstanceDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", nil)
		}
		return dto.ResponseCustomWorkoutInstanceDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get workout instance", err)
	}
	return instance, nil
}

func (s *CustomWorkoutInstanceService) ListCustomWorkoutInstances(gymID string) ([]dto.ResponseCustomWorkoutInstanceDTO, error) {
	instances, err := s.Repo.List(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout instances", err)
	}
	if len(instances) == 0 {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "No workout instances found", nil)
	}
	return instances, nil
}

func (s *CustomWorkoutInstanceService) UpdateCustomWorkoutInstance(gymID string, instance dto.UpdateCustomWorkoutInstanceDTO) error {
	// Check existence first
	_, err := s.Repo.GetByID(gymID, instance.ID)
	if err != nil {
		return apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", err)
	}
	err = s.Repo.Update(gymID, instance)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update workout instance", err)
	}
	return nil
}

func (s *CustomWorkoutInstanceService) DeleteCustomWorkoutInstance(gymID, id string) error {
	// Check existence first
	_, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeNotFound, "Workout instance not found", err)
	}
	err = s.Repo.Delete(gymID, id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete workout instance", err)
	}
	return nil
}
