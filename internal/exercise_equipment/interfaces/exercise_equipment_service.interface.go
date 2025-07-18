package interfaces

import "github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"

type ExerciseEquipmentService interface {
	CreateLink(link dto.ExerciseEquipment) error
	DeleteLink(id string) error
	GetLinkByID(id string) (dto.ExerciseEquipment, error)
	GetLinksByExerciseID(exerciseID string) ([]dto.ExerciseEquipment, error)
	GetLinksByEquipmentID(equipmentID string) ([]dto.ExerciseEquipment, error)
}
