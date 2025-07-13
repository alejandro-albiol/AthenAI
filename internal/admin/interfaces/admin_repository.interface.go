package interfaces

import "github.com/alejandro-albiol/athenai/internal/admin/dto"

type AdminRepository interface {
	// Exercise management
	CreateExercise(exercise dto.ExerciseCreationDTO) (string, error)
	GetExerciseByID(id string) (dto.ExerciseResponseDTO, error)
	GetAllExercises() ([]dto.ExerciseResponseDTO, error)
	UpdateExercise(id string, exercise dto.ExerciseUpdateDTO) (dto.ExerciseResponseDTO, error)
	DeleteExercise(id string) error
	GetExercisesByMuscularGroup(muscularGroups []string) ([]dto.ExerciseResponseDTO, error)
	GetExercisesByEquipment(equipment []string) ([]dto.ExerciseResponseDTO, error)

	// Equipment management
	CreateEquipment(equipment dto.EquipmentCreationDTO) (string, error)
	GetEquipmentByID(id string) (dto.EquipmentResponseDTO, error)
	GetAllEquipment() ([]dto.EquipmentResponseDTO, error)
	UpdateEquipment(id string, equipment dto.EquipmentUpdateDTO) (dto.EquipmentResponseDTO, error)
	DeleteEquipment(id string) error

	// Muscular group management
	CreateMuscularGroup(group dto.MuscularGroupCreationDTO) (string, error)
	GetMuscularGroupByID(id string) (dto.MuscularGroupResponseDTO, error)
	GetAllMuscularGroups() ([]dto.MuscularGroupResponseDTO, error)
	UpdateMuscularGroup(id string, group dto.MuscularGroupUpdateDTO) (dto.MuscularGroupResponseDTO, error)
	DeleteMuscularGroup(id string) error
}
