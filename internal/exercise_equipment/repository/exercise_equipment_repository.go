 package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/interfaces"
)

type ExerciseEquipmentRepositoryImpl struct {
	db *sql.DB
}

func NewExerciseEquipmentRepository(db *sql.DB) interfaces.ExerciseEquipmentRepository {
	return &ExerciseEquipmentRepositoryImpl{db: db}
}

func (r *ExerciseEquipmentRepositoryImpl) CreateLink(link dto.ExerciseEquipment) (string, error) {
	// TODO: Implement DB insert
	return "", nil
}

func (r *ExerciseEquipmentRepositoryImpl) DeleteLink(id string) error {
	// TODO: Implement DB delete
	return nil
}

func (r *ExerciseEquipmentRepositoryImpl) FindByID(id string) (dto.ExerciseEquipment, error) {
	// TODO: Implement DB select by ID
	return dto.ExerciseEquipment{}, nil
}

func (r *ExerciseEquipmentRepositoryImpl) FindByExerciseID(exerciseID string) ([]dto.ExerciseEquipment, error) {
	// TODO: Implement DB select by exerciseID
	return nil, nil
}

func (r *ExerciseEquipmentRepositoryImpl) FindByEquipmentID(equipmentID string) ([]dto.ExerciseEquipment, error) {
	// TODO: Implement DB select by equipmentID
	return nil, nil
}
