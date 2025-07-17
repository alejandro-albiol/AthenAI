package service

import (
	"github.com/alejandro-albiol/athenai/internal/exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise/repository"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type ExerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}

}

func (s *ExerciseService) CreateExercise(exercise dto.ExerciseCreationDTO) (string, error) {
	existingExercise, err := s.repo.GetExerciseByName(exercise.Name)
	if err == nil && existingExercise.ID != "" {
		return "", apierror.New(errorcode_enum.CodeConflict, "Exercise with this name already exists", err)
	}
	id, err := s.repo.CreateExercise(exercise)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create exercise", err)
	}
	return id, nil
}

func (s *ExerciseService) GetExerciseByID(id string) (dto.ExerciseResponseDTO, error) {
	exercise, err := s.repo.GetExerciseByID(id)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Exercise with ID "+id+" not found", err)
	}
	return exercise, nil
}

func (s *ExerciseService) GetExerciseByName(name string) (dto.ExerciseResponseDTO, error) {
	exercise, err := s.repo.GetExerciseByName(name)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Exercise with name '"+name+"' not found", err)
	}
	return exercise, nil
}

func (s *ExerciseService) GetAllExercises() ([]dto.ExerciseResponseDTO, error) {
	exercises, err := s.repo.GetAllExercises()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve exercises", err)
	}
	return exercises, nil
}

func (s *ExerciseService) UpdateExercise(id string, exercise dto.ExerciseUpdateDTO) (dto.ExerciseResponseDTO, error) {
	// Check existence first
	_, err := s.repo.GetExerciseByID(id)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Exercise with ID "+id+" not found", err)
	}
	updatedExercise, err := s.repo.UpdateExercise(id, exercise)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update exercise", err)
	}
	return updatedExercise, nil
}

func (s *ExerciseService) DeleteExercise(id string) error {
	// Check existence first
	_, err := s.repo.GetExerciseByID(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeNotFound, "Exercise with ID "+id+" not found", err)
	}
	err = s.repo.DeleteExercise(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete exercise", err)
	}
	return nil
}
