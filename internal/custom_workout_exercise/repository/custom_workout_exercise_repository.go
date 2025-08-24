package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
)

type CustomWorkoutExerciseRepository struct {
	DB *sql.DB
}

func NewCustomWorkoutExerciseRepository(db *sql.DB) *CustomWorkoutExerciseRepository {
	return &CustomWorkoutExerciseRepository{DB: db}
}

func (r *CustomWorkoutExerciseRepository) Create(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (string, error) {
	// TODO: Implement DB insert logic
	return "", nil
}

func (r *CustomWorkoutExerciseRepository) GetByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
	// TODO: Implement DB select logic
	return &dto.ResponseCustomWorkoutExerciseDTO{}, nil
}

func (r *CustomWorkoutExerciseRepository) ListByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	// TODO: Implement DB select logic
	return []*dto.ResponseCustomWorkoutExerciseDTO{}, nil
}

func (r *CustomWorkoutExerciseRepository) Update(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error {
	// TODO: Implement DB update logic
	return nil
}

func (r *CustomWorkoutExerciseRepository) Delete(gymID, id string) error {
	// TODO: Implement DB delete logic
	return nil
}
