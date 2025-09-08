package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"

type CustomExerciseEquipmentRepository interface {
	CreateLink(gymID string, link *dto.CustomExerciseEquipment) (*string, error)
	DeleteLink(gymID, id string) error
	RemoveAllLinksForExercise(gymID, customExerciseID string) error
	FindByID(gymID, id string) (*dto.CustomExerciseEquipment, error)
	FindByCustomExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseEquipment, error)
	FindByEquipmentID(gymID, equipmentID string) ([]*dto.CustomExerciseEquipment, error)
}
