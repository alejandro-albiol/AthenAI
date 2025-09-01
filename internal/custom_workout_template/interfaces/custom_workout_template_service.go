package interfaces

import "github.com/alejandro-albiol/athenai/internal/custom_workout_template/dto"

type CustomWorkoutTemplateService interface {
	CreateCustomWorkoutTemplate(gymID string, template dto.CreateCustomWorkoutTemplateDTO) error
	GetCustomWorkoutTemplateByID(gymID, id string) (dto.ResponseCustomWorkoutTemplateDTO, error)
	ListCustomWorkoutTemplates(gymID string) ([]dto.ResponseCustomWorkoutTemplateDTO, error)
	UpdateCustomWorkoutTemplate(gymID string, template dto.UpdateCustomWorkoutTemplateDTO) error
	DeleteCustomWorkoutTemplate(gymID, id string) error
}
