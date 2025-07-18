package interfaces

import "github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"

type ExerciseEquipmentRepository interface {
	CreateLink(link dto.ExerciseEquipment) (string, error)
	DeleteLink(id string) error
	FindByID(id string) (dto.ExerciseEquipment, error)
	FindByExerciseID(exerciseID string) ([]dto.ExerciseEquipment, error)
	FindByEquipmentID(equipmentID string) ([]dto.ExerciseEquipment, error)
}
