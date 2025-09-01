package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_template/dto"

type CustomWorkoutTemplateRepository interface {
	Create(gymID string, template dto.CreateCustomWorkoutTemplateDTO) (string, error)
	GetByID(gymID, id string) (dto.ResponseCustomWorkoutTemplateDTO, error)
	GetByName(gymID, name string) (dto.ResponseCustomWorkoutTemplateDTO, error)
	List(gymID string) ([]dto.ResponseCustomWorkoutTemplateDTO, error)
	Update(gymID string, template dto.UpdateCustomWorkoutTemplateDTO) error
	Delete(gymID, id string) error
}
