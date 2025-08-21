package service

import (
	"github.com/alejandro-albiol/athenai/internal/exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise/interfaces"
	exerciseEquipmentDTO "github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
	equipmentIF "github.com/alejandro-albiol/athenai/internal/exercise_equipment/interfaces"
	exerciseMuscularGroupDTO "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
	muscularGroupIF "github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type ExerciseService struct {
	repository                   interfaces.ExerciseRepository
	exerciseEquipmentService     equipmentIF.ExerciseEquipmentService
	exerciseMuscularGroupService muscularGroupIF.ExerciseMuscularGroupService
}

func NewExerciseService(repo interfaces.ExerciseRepository, equipmentService equipmentIF.ExerciseEquipmentService, muscularGroupService muscularGroupIF.ExerciseMuscularGroupService) *ExerciseService {
	return &ExerciseService{
		repository:                   repo,
		exerciseEquipmentService:     equipmentService,
		exerciseMuscularGroupService: muscularGroupService,
	}

}

func (s *ExerciseService) CreateExercise(exercise dto.ExerciseCreationDTO) (string, error) {
	// Validate enums
	if !exercise.DifficultyLevel.IsValid() {
		return "", apierror.New(errorcode_enum.CodeBadRequest, "Invalid difficulty level", nil)
	}
	if !exercise.ExerciseType.IsValid() {
		return "", apierror.New(errorcode_enum.CodeBadRequest, "Invalid exercise type", nil)
	}
	// Validate synonyms
	if exercise.Synonyms == nil {
		return "", apierror.New(errorcode_enum.CodeBadRequest, "Synonyms must be provided as an array", nil)
	}
	seen := make(map[string]struct{})
	for _, syn := range exercise.Synonyms {
		if syn == "" {
			return "", apierror.New(errorcode_enum.CodeBadRequest, "Synonyms cannot contain empty strings", nil)
		}
		if _, exists := seen[syn]; exists {
			return "", apierror.New(errorcode_enum.CodeBadRequest, "Synonyms must be unique", nil)
		}
		seen[syn] = struct{}{}
	}
	existingExercise, err := s.repository.GetExerciseByName(exercise.Name)
	if err == nil && existingExercise.ID != "" {
		return "", apierror.New(errorcode_enum.CodeConflict, "Exercise with this name already exists", err)
	}
	id, err := s.repository.CreateExercise(exercise)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create exercise", err)
	}

	// Create links in join tables
	if s.exerciseEquipmentService != nil && len(exercise.Equipment) > 0 {
		for _, eqID := range exercise.Equipment {
			link := exerciseEquipmentDTO.ExerciseEquipment{
				ExerciseID:  id,
				EquipmentID: eqID,
			}
			_, err := s.exerciseEquipmentService.CreateLink(link)
			if err != nil {
				return "", apierror.New(errorcode_enum.CodeInternal, "Failed to link equipment to exercise", err)
			}
		}
	}
	if s.exerciseMuscularGroupService != nil && len(exercise.MuscularGroups) > 0 {
		for _, mgID := range exercise.MuscularGroups {
			link := exerciseMuscularGroupDTO.ExerciseMuscularGroup{
				ExerciseID:      id,
				MuscularGroupID: mgID,
			}
			_, err := s.exerciseMuscularGroupService.CreateLink(link)
			if err != nil {
				return "", apierror.New(errorcode_enum.CodeInternal, "Failed to link muscular group to exercise", err)
			}
		}
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
	// Validate enums if provided
	if exercise.DifficultyLevel != "" && !exercise.DifficultyLevel.IsValid() {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeBadRequest, "Invalid difficulty level", nil)
	}
	if exercise.ExerciseType != "" && !exercise.ExerciseType.IsValid() {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeBadRequest, "Invalid exercise type", nil)
	}
	// Validate synonyms if provided
	if exercise.Synonyms != nil {
		seen := make(map[string]struct{})
		for _, syn := range exercise.Synonyms {
			if syn == "" {
				return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeBadRequest, "Synonyms cannot contain empty strings", nil)
			}
			if _, exists := seen[syn]; exists {
				return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeBadRequest, "Synonyms must be unique", nil)
			}
			seen[syn] = struct{}{}
		}
	}
	updatedExercise, err := s.repository.UpdateExercise(id, exercise)
	if err != nil {
		return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update exercise", err)
	}
	// Update join tables: remove all links, then add new ones if provided
	if s.exerciseEquipmentService != nil {
		_ = s.exerciseEquipmentService.RemoveAllLinksForExercise(id)
		if exercise.Equipment != nil {
			for _, eqID := range exercise.Equipment {
				link := exerciseEquipmentDTO.ExerciseEquipment{
					ExerciseID:  id,
					EquipmentID: eqID,
				}
				_, err := s.exerciseEquipmentService.CreateLink(link)
				if err != nil {
					return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update equipment links", err)
				}
			}
		}
	}
	if s.exerciseMuscularGroupService != nil {
		_ = s.exerciseMuscularGroupService.RemoveAllLinksForExercise(id)
		if exercise.MuscularGroups != nil {
			for _, mgID := range exercise.MuscularGroups {
				link := exerciseMuscularGroupDTO.ExerciseMuscularGroup{
					ExerciseID:      id,
					MuscularGroupID: mgID,
				}
				_, err := s.exerciseMuscularGroupService.CreateLink(link)
				if err != nil {
					return dto.ExerciseResponseDTO{}, apierror.New(errorcode_enum.CodeInternal, "Failed to update muscular group links", err)
				}
			}
		}
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
	// Remove all join table links
	if s.exerciseEquipmentService != nil {
		_ = s.exerciseEquipmentService.RemoveAllLinksForExercise(id)
	}
	if s.exerciseMuscularGroupService != nil {
		_ = s.exerciseMuscularGroupService.RemoveAllLinksForExercise(id)
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
