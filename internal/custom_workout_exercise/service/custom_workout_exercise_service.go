package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/interfaces"
)

type CustomWorkoutExerciseService struct {
	Repo interfaces.CustomWorkoutExerciseRepository
}

func NewCustomWorkoutExerciseService(repo interfaces.CustomWorkoutExerciseRepository) *CustomWorkoutExerciseService {
	return &CustomWorkoutExerciseService{Repo: repo}
}

func (s *CustomWorkoutExerciseService) CreateCustomWorkoutExercise(gymID string, exercise *dto.CreateCustomWorkoutExerciseDTO) (string, error) {
	return s.Repo.Create(gymID, exercise)
}

func (s *CustomWorkoutExerciseService) GetCustomWorkoutExerciseByID(gymID, id string) (*dto.ResponseCustomWorkoutExerciseDTO, error) {
	return s.Repo.GetByID(gymID, id)
}

func (s *CustomWorkoutExerciseService) ListCustomWorkoutExercisesByWorkoutInstanceID(gymID, workoutInstanceID string) ([]*dto.ResponseCustomWorkoutExerciseDTO, error) {
	return s.Repo.ListByWorkoutInstanceID(gymID, workoutInstanceID)
}

func (s *CustomWorkoutExerciseService) UpdateCustomWorkoutExercise(gymID string, exercise *dto.UpdateCustomWorkoutExerciseDTO) error {
	return s.Repo.Update(gymID, exercise)
}

func (s *CustomWorkoutExerciseService) DeleteCustomWorkoutExercise(gymID, id string) error {
	return s.Repo.Delete(gymID, id)
}
