package interfaces

import "github.com/alejandro-albiol/athenai/internal/workout_template/dto"

type WorkoutTemplateRepository interface {
	Create(template dto.WorkoutTemplateDTO) (string, error)
	GetByID(id string) (dto.WorkoutTemplateDTO, error)
	GetByName(name string) (dto.WorkoutTemplateDTO, error)
	GetByDifficulty(difficulty string) ([]dto.WorkoutTemplateDTO, error)
	GetByTargetAudience(targetAudience string) ([]dto.WorkoutTemplateDTO, error)
	GetAll() ([]dto.WorkoutTemplateDTO, error)
	Update(id string, template dto.WorkoutTemplateDTO) (dto.WorkoutTemplateDTO, error)
	Delete(id string) error
}
