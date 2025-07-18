package service

import (
	"github.com/alejandro-albiol/athenai/internal/exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type ExerciseService struct {
	repository interfaces.ExerciseRepository
}

func NewExerciseService(repo interfaces.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repository: repo}

}

func (s *ExerciseService) CreateExercise(exercise dto.ExerciseCreationDTO) (string, error) {
	existingExercise, err := s.repository.GetExerciseByName(exercise.Name)
	if err == nil && existingExercise.ID != "" {
		return "", apierror.New(errorcode_enum.CodeConflict, "Exercise with this name already exists", err)
	}
	id, err := s.repository.CreateExercise(exercise)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create exercise", err)
	}
	return id, nil
}

func (s *ExerciseService) GetExerciseByID(id string) (dto.ExerciseResponseDTO, error) {
	exercise, err := s.repository.GetExerciseByID(id)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Exercise with ID "+id+" not found", err)
	}
	return exercise, nil
}

func (s *ExerciseService) GetExerciseByName(name string) (dto.ExerciseResponseDTO, error) {
	exercise, err := s.repository.GetExerciseByName(name)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Exercise with name '"+name+"' not found", err)
	}
	if exercise.ID == "" {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Exercise with name '"+name+"' not found", nil)
	}
	return exercise, nil
}

func (s *ExerciseService) GetExercisesByMuscularGroup(muscularGroups []string) ([]dto.ExerciseResponseDTO, error) {
	exercises, err := s.repository.GetExercisesByMuscularGroup(muscularGroups)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve exercises by muscular group", err)
	}
	if len(exercises) == 0 {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "No exercises found for the specified muscular groups", nil)
	}
	return exercises, nil
}

func (s *ExerciseService) GetExercisesByEquipment(equipment []string) ([]dto.ExerciseResponseDTO, error) {
	exercises, err := s.repository.GetExercisesByEquipment(equipment)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve exercises by equipment", err)
	}
	if len(exercises) == 0 {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "No exercises found for the specified equipment", nil)
	}
	return exercises, nil
}

func (s *ExerciseService) GetAllExercises() ([]dto.ExerciseResponseDTO, error) {
	exercises, err := s.repository.GetAllExercises()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve exercises", err)
	}
	if len(exercises) == 0 {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "No exercises found", nil)
	}
	return exercises, nil
}

func (s *ExerciseService) UpdateExercise(id string, exercise dto.ExerciseUpdateDTO) (dto.ExerciseResponseDTO, error) {
	// Check existence first
	_, err := s.repository.GetExerciseByID(id)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Exercise with ID "+id+" not found", err)
	}
	updatedExercise, err := s.repository.UpdateExercise(id, exercise)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update exercise", err)
	}
	return updatedExercise, nil
}

func (s *ExerciseService) DeleteExercise(id string) error {
	// Check existence first
	_, err := s.repository.GetExerciseByID(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeNotFound, "Exercise with ID "+id+" not found", err)
	}
	err = s.repository.DeleteExercise(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete exercise", err)
	}
	return nil
}

func (s *ExerciseService) GetExercisesByMuscularGroupAndEquipment(muscularGroups []string, equipment []string) ([]dto.ExerciseResponseDTO, error) {
	if len(muscularGroups) == 0 && len(equipment) == 0 {
		return s.GetAllExercises()
	}
	var exercisesByGroup, exercisesByEquipment []dto.ExerciseResponseDTO
	var err error

	if len(muscularGroups) > 0 {
		exercisesByGroup, err = s.repository.GetExercisesByMuscularGroup(muscularGroups)
		if err != nil {
			return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve exercises by muscular group", err)
		}
	}
	if len(equipment) > 0 {
		exercisesByEquipment, err = s.repository.GetExercisesByEquipment(equipment)
		if err != nil {
			return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve exercises by equipment", err)
		}
	}

	// If both filters are present, return intersection
	if len(muscularGroups) > 0 && len(equipment) > 0 {
		idSet := make(map[string]struct{})
		for _, e := range exercisesByGroup {
			idSet[e.ID] = struct{}{}
		}
		var intersection []dto.ExerciseResponseDTO
		for _, e := range exercisesByEquipment {
			if _, found := idSet[e.ID]; found {
				intersection = append(intersection, e)
			}
		}
		if len(intersection) == 0 {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "No exercises found for the specified filters", nil)
		}
		return intersection, nil
	}

	// If only one filter, return its result
	if len(muscularGroups) > 0 {
		if len(exercisesByGroup) == 0 {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "No exercises found for the specified muscular groups", nil)
		}
		return exercisesByGroup, nil
	}
	if len(exercisesByEquipment) > 0 {
		if len(exercisesByEquipment) == 0 {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "No exercises found for the specified equipment", nil)
		}
		return exercisesByEquipment, nil
	}
	return nil, apierror.New(errorcode_enum.CodeNotFound, "No exercises found for the specified filters", nil)
}
