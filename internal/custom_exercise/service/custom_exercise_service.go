package service

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/repository"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomExerciseService struct {
	Repo *repository.CustomExerciseRepository
}

func NewCustomExerciseService(repo *repository.CustomExerciseRepository) *CustomExerciseService {
	return &CustomExerciseService{Repo: repo}
}

func (s *CustomExerciseService) Create(gymID string, exercise *dto.CustomExerciseCreationDTO) *apierror.APIError {
	if err := s.Repo.Create(gymID, exercise); err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create custom exercise", err)
	}
	return nil
}

func (s *CustomExerciseService) GetByID(gymID, id string) (*dto.CustomExerciseResponseDTO, *apierror.APIError) {
	res, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Custom exercise not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get custom exercise", err)
	}
	return res, nil
}

func (s *CustomExerciseService) List(gymID string) ([]*dto.CustomExerciseResponseDTO, *apierror.APIError) {
	res, err := s.Repo.List(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom exercises", err)
	}
	return res, nil
}

func (s *CustomExerciseService) Update(gymID string, exercise *dto.CustomExerciseUpdateDTO) *apierror.APIError {
	if err := s.Repo.Update(gymID, exercise); err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update custom exercise", err)
	}
	return nil
}

func (s *CustomExerciseService) Delete(gymID, id string) *apierror.APIError {
	if err := s.Repo.Delete(gymID, id); err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom exercise", err)
	}
	return nil
}
