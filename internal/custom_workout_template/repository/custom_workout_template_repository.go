package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_template/dto"
)

type CustomWorkoutTemplateRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomWorkoutTemplateRepository(db *sql.DB) *CustomWorkoutTemplateRepositoryImpl {
	return &CustomWorkoutTemplateRepositoryImpl{DB: db}
}

func (r *CustomWorkoutTemplateRepositoryImpl) Create(gymID string, template *dto.CreateCustomWorkoutTemplateDTO) (string, error) {
	// SQL: INSERT INTO <gymID>.custom_workout_template ...
	return "", nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) GetByID(gymID, id string) (*dto.ResponseCustomWorkoutTemplateDTO, error) {
	// SQL: SELECT ... FROM <gymID>.custom_workout_template WHERE id = $1
	return nil, nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) List(gymID string) ([]*dto.ResponseCustomWorkoutTemplateDTO, error) {
	// SQL: SELECT ... FROM <gymID>.custom_workout_template
	return nil, nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) Update(gymID string, template *dto.UpdateCustomWorkoutTemplateDTO) error {
	// SQL: UPDATE <gymID>.custom_workout_template SET ... WHERE id = $1
	return nil
}

func (r *CustomWorkoutTemplateRepositoryImpl) Delete(gymID, id string) error {
	// SQL: DELETE FROM <gymID>.custom_workout_template WHERE id = $1
	return nil
}
