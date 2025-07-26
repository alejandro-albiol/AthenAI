package interfaces

import "github.com/alejandro-albiol/athenai/internal/template_block/dto"

type TemplateBlockRepository interface {
	Create(block dto.TemplateBlockDTO) (string, error)
	GetByID(id string) (dto.TemplateBlockDTO, error)
	GetByTemplateID(templateID string) ([]dto.TemplateBlockDTO, error)
	Update(id string, block dto.TemplateBlockDTO) (dto.TemplateBlockDTO, error)
	Delete(id string) error
}
