package interfaces

import "github.com/alejandro-albiol/athenai/internal/template_block/dto"

type TemplateBlockRepository interface {
	Create(block dto.CreateTemplateBlockDTO) error
	GetByID(id string) (dto.TemplateBlockDTO, error)
	GetByTemplateID(templateID string) ([]dto.TemplateBlockDTO, error)
	GetByTemplateIDAndName(templateID, name string) (dto.TemplateBlockDTO, error)
	Update(id string, block dto.UpdateTemplateBlockDTO) (dto.TemplateBlockDTO, error)
	Delete(id string) error
}
