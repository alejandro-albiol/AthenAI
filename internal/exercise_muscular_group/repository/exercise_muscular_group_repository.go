package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/interfaces"
)

type ExerciseMuscularGroupRepositoryImpl struct {
	db *sql.DB
}

func NewExerciseMuscularGroupRepository(db *sql.DB) interfaces.ExerciseMuscularGroupRepository {
	return &ExerciseMuscularGroupRepositoryImpl{db: db}
}

func (r *ExerciseMuscularGroupRepositoryImpl) CreateLink(link dto.ExerciseMuscularGroup) (string, error) {
	// TODO: Implement DB insert
	return "", nil
}

func (r *ExerciseMuscularGroupRepositoryImpl) DeleteLink(id string) error {
	// TODO: Implement DB delete
	return nil
}

func (r *ExerciseMuscularGroupRepositoryImpl) FindByID(id string) (dto.ExerciseMuscularGroup, error) {
	// TODO: Implement DB select by ID
	return dto.ExerciseMuscularGroup{}, nil
}

func (r *ExerciseMuscularGroupRepositoryImpl) FindByExerciseID(exerciseID string) ([]dto.ExerciseMuscularGroup, error) {
	// TODO: Implement DB select by exerciseID
	return nil, nil
}

func (r *ExerciseMuscularGroupRepositoryImpl) FindByMuscularGroupID(muscularGroupID string) ([]dto.ExerciseMuscularGroup, error) {
	// TODO: Implement DB select by muscularGroupID
	return nil, nil
}
