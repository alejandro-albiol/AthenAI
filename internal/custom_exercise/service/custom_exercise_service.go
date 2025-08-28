package service

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomExerciseService struct {
	Repo interfaces.CustomExerciseRepository
}

func NewCustomExerciseService(repo interfaces.CustomExerciseRepository) *CustomExerciseService {
	return &CustomExerciseService{Repo: repo}
}

func (s *CustomExerciseService) CreateCustomExercise(gymID string, exercise dto.CustomExerciseCreationDTO) error {
	if err := s.Repo.CreateCustomExercise(gymID, exercise); err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create custom exercise", err)
	}
	return nil
}

func (s *CustomExerciseService) GetCustomExerciseByID(gymID, id string) (dto.CustomExerciseResponseDTO, error) {
	res, err := s.Repo.GetCustomExerciseByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.CustomExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Custom exercise not found", err)
		}
		return dto.CustomExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to get custom exercise", err)
	}
	return res, nil
}

func (s *CustomExerciseService) ListCustomExercises(gymID string) ([]dto.CustomExerciseResponseDTO, error) {
	res, err := s.Repo.ListCustomExercises(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom exercises", err)
	}
	return res, nil
}

func (s *CustomExerciseService) UpdateCustomExercise(gymID, id string, update dto.CustomExerciseUpdateDTO) error {
	if err := s.Repo.UpdateCustomExercise(gymID, id, update); err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update custom exercise", err)
	}
	return nil
}

func (s *CustomExerciseService) DeleteCustomExercise(gymID, id string) error {
	if err := s.Repo.DeleteCustomExercise(gymID, id); err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom exercise", err)
	}
	return nil
}
