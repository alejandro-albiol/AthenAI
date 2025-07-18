package repository

import (
	"database/sql"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/interfaces"
)

type CustomExerciseMuscularGroupRepositoryImpl struct {
	db *sql.DB
}

func NewCustomExerciseMuscularGroupRepository(db *sql.DB) interfaces.CustomExerciseMuscularGroupRepository {
	return &CustomExerciseMuscularGroupRepositoryImpl{db: db}
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) CreateLink(link dto.CustomExerciseMuscularGroup) (string, error) {
	// TODO: Implement DB insert
	return "", nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) DeleteLink(id string) error {
	// TODO: Implement DB delete
	return nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) FindByID(id string) (dto.CustomExerciseMuscularGroup, error) {
	// TODO: Implement DB select by ID
	return dto.CustomExerciseMuscularGroup{}, nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) FindByCustomExerciseID(customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error) {
	// TODO: Implement DB select by customExerciseID
	return nil, nil
}

func (r *CustomExerciseMuscularGroupRepositoryImpl) FindByMuscularGroupID(muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error) {
	// TODO: Implement DB select by muscularGroupID
	return nil, nil
}
