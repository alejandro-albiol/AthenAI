package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
)

type CustomMemberWorkoutRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomMemberWorkoutRepository(db *sql.DB) *CustomMemberWorkoutRepositoryImpl {
	return &CustomMemberWorkoutRepositoryImpl{DB: db}
}

func (r *CustomMemberWorkoutRepositoryImpl) Create(gymID string, memberWorkout *dto.CreateCustomMemberWorkoutDTO) (string, error) {
	// TODO: Implement
	return "", nil
}

func (r *CustomMemberWorkoutRepositoryImpl) GetByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
	// TODO: Implement
	return nil, nil
}

func (r *CustomMemberWorkoutRepositoryImpl) ListByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error) {
	// TODO: Implement
	return nil, nil
}

func (r *CustomMemberWorkoutRepositoryImpl) Update(gymID string, memberWorkout *dto.UpdateCustomMemberWorkoutDTO) error {
	// TODO: Implement
	return nil
}

func (r *CustomMemberWorkoutRepositoryImpl) Delete(gymID, id string) error {
	// TODO: Implement
	return nil
}
