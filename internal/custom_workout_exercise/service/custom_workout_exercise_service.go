package service

import (
	"database/sql"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomWorkoutExerciseService struct {
	Repo interfaces.CustomWorkoutExerciseRepository
}

func NewCustomWorkoutExerciseService(repo interfaces.CustomWorkoutExerciseRepository) *CustomWorkoutExerciseService {
	return &CustomWorkoutExerciseService{Repo: repo}
}

func (s *CustomWorkoutExerciseService) CreateCustomWorkoutExercise(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (*string, error) {
	// Validate required fields
	if exercise.CreatedBy == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "CreatedBy is required", nil)
	}
	if exercise.WorkoutInstanceID == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "WorkoutInstanceID is required", nil)
	}
	if exercise.BlockName == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "BlockName is required", nil)
	}
	if exercise.ExerciseOrder <= 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "ExerciseOrder must be greater than 0", nil)
	}

	// Validate exercise source and IDs
	if exercise.ExerciseSource != "public" && exercise.ExerciseSource != "gym" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "ExerciseSource must be 'public' or 'gym'", nil)
	}

	if exercise.ExerciseSource == "public" {
		if exercise.PublicExerciseID == nil || *exercise.PublicExerciseID == "" {
			return nil, apierror.New(errorcode_enum.CodeBadRequest, "PublicExerciseID is required when ExerciseSource is 'public'", nil)
		}
		if exercise.GymExerciseID != nil && *exercise.GymExerciseID != "" {
			return nil, apierror.New(errorcode_enum.CodeBadRequest, "GymExerciseID must be empty when ExerciseSource is 'public'", nil)
		}
	} else {
		if exercise.GymExerciseID == nil || *exercise.GymExerciseID == "" {
			return nil, apierror.New(errorcode_enum.CodeBadRequest, "GymExerciseID is required when ExerciseSource is 'gym'", nil)
		}
		if exercise.PublicExerciseID != nil && *exercise.PublicExerciseID != "" {
			return nil, apierror.New(errorcode_enum.CodeBadRequest, "PublicExerciseID must be empty when ExerciseSource is 'gym'", nil)
		}
	}

	// Validate numeric values
	if exercise.Sets != nil && *exercise.Sets <= 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "Sets must be greater than 0", nil)
	}
	if exercise.RepsMin != nil && *exercise.RepsMin <= 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "RepsMin must be greater than 0", nil)
	}
	if exercise.RepsMax != nil && *exercise.RepsMax <= 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "RepsMax must be greater than 0", nil)
	}
	if exercise.RepsMin != nil && exercise.RepsMax != nil && *exercise.RepsMin > *exercise.RepsMax {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "RepsMin cannot be greater than RepsMax", nil)
	}
	if exercise.WeightKg != nil && *exercise.WeightKg < 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "WeightKg cannot be negative", nil)
	}
	if exercise.DurationSeconds != nil && *exercise.DurationSeconds <= 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "DurationSeconds must be greater than 0", nil)
	}
	if exercise.RestSeconds != nil && *exercise.RestSeconds < 0 {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "RestSeconds cannot be negative", nil)
	}

	// Check for duplicate exercise order in the same block and workout instance
	existingExercises, err := s.Repo.ListByWorkoutInstanceID(gymID, exercise.WorkoutInstanceID)
	if err != nil && err != sql.ErrNoRows {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to check existing exercises", err)
	}

	for _, existing := range existingExercises {
		if existing.BlockName == exercise.BlockName && existing.ExerciseOrder == exercise.ExerciseOrder {
			return nil, apierror.New(errorcode_enum.CodeConflict,
				fmt.Sprintf("Exercise order %d already exists in block '%s'", exercise.ExerciseOrder, exercise.BlockName), nil)
		}
	}

	// Create the exercise
	id, err := s.Repo.Create(gymID, exercise)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create workout exercise", err)
	}

	return id, nil
}

func (s *CustomWorkoutExerciseService) GetCustomWorkoutExerciseByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if id == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "ID is required", nil)
	}

	exercise, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Workout exercise not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get workout exercise", err)
	}

	return exercise, nil
}

func (s *CustomWorkoutExerciseService) ListCustomWorkoutExercisesByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if workoutInstanceID == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "WorkoutInstanceID is required", nil)
	}

	exercises, err := s.Repo.ListByWorkoutInstanceID(gymID, workoutInstanceID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout exercises", err)
	}

	return exercises, nil
}

func (s *CustomWorkoutExerciseService) UpdateCustomWorkoutExercise(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error {
	if exercise.ID == "" {
		return apierror.New(errorcode_enum.CodeBadRequest, "ID is required", nil)
	}

	// Validate numeric values if provided
	if exercise.Sets != nil && *exercise.Sets <= 0 {
		return apierror.New(errorcode_enum.CodeBadRequest, "Sets must be greater than 0", nil)
	}
	if exercise.RepsMin != nil && *exercise.RepsMin <= 0 {
		return apierror.New(errorcode_enum.CodeBadRequest, "RepsMin must be greater than 0", nil)
	}
	if exercise.RepsMax != nil && *exercise.RepsMax <= 0 {
		return apierror.New(errorcode_enum.CodeBadRequest, "RepsMax must be greater than 0", nil)
	}
	if exercise.RepsMin != nil && exercise.RepsMax != nil && *exercise.RepsMin > *exercise.RepsMax {
		return apierror.New(errorcode_enum.CodeBadRequest, "RepsMin cannot be greater than RepsMax", nil)
	}
	if exercise.WeightKg != nil && *exercise.WeightKg < 0 {
		return apierror.New(errorcode_enum.CodeBadRequest, "WeightKg cannot be negative", nil)
	}
	if exercise.DurationSeconds != nil && *exercise.DurationSeconds <= 0 {
		return apierror.New(errorcode_enum.CodeBadRequest, "DurationSeconds must be greater than 0", nil)
	}
	if exercise.RestSeconds != nil && *exercise.RestSeconds < 0 {
		return apierror.New(errorcode_enum.CodeBadRequest, "RestSeconds cannot be negative", nil)
	}

	// Check if exercise exists before updating
	_, err := s.Repo.GetByID(gymID, exercise.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Workout exercise not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check existing workout exercise", err)
	}

	err = s.Repo.Update(gymID, exercise)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Workout exercise not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update workout exercise", err)
	}

	return nil
}

func (s *CustomWorkoutExerciseService) DeleteCustomWorkoutExercise(gymID, id string) error {
	if id == "" {
		return apierror.New(errorcode_enum.CodeBadRequest, "ID is required", nil)
	}

	// Check if exercise exists before deleting
	_, err := s.Repo.GetByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Workout exercise not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check existing workout exercise", err)
	}

	err = s.Repo.Delete(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Workout exercise not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete workout exercise", err)
	}

	return nil
}

func (s *CustomWorkoutExerciseService) ListCustomWorkoutExercisesByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if muscularGroupID == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "MuscularGroupID is required", nil)
	}

	exercises, err := s.Repo.ListByMuscularGroupID(gymID, muscularGroupID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout exercises by muscular group", err)
	}

	return exercises, nil
}

func (s *CustomWorkoutExerciseService) ListCustomWorkoutExercisesByEquipmentID(gymID, equipmentID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	if equipmentID == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "EquipmentID is required", nil)
	}

	exercises, err := s.Repo.ListByEquipmentID(gymID, equipmentID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list workout exercises by equipment", err)
	}

	return exercises, nil
}
