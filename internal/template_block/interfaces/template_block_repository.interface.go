package interfaces

import "github.com/alejandro-albiol/athenai/internal/template_block/dto"

type TemplateBlockRepository interface {
	CreateTemplateBlock(block *dto.CreateTemplateBlockDTO) (string, error)
	GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error)
	GetTemplateBlocksByTemplateID(templateID string) ([]*dto.TemplateBlockDTO, error)
	GetTemplateBlockByTemplateIDAndName(templateID string, name string) (*dto.TemplateBlockDTO, error)
	UpdateTemplateBlock(id string, block *dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error)
	DeleteTemplateBlock(id string) error
}
