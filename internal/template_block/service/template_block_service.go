package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type TemplateBlockService struct {
	repository interfaces.TemplateBlockRepository
}

func NewTemplateBlockService(repository interfaces.TemplateBlockRepository) *TemplateBlockService {
	return &TemplateBlockService{repository: repository}
}

func (s *TemplateBlockService) CreateTemplateBlock(block *dto.CreateTemplateBlockDTO) (*string, error) {
	if block == nil {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "Block payload is nil", nil)
	}
	existingBlock, err := s.repository.GetTemplateBlockByTemplateIDAndName(block.TemplateID, block.BlockName)
	if err == nil && existingBlock != nil && existingBlock.ID != "" {
		return nil, apierror.New(errorcode_enum.CodeConflict, fmt.Sprintf("Template block with name '%s' already exists in template", block.BlockName), nil)
	}
	id, err := s.repository.CreateTemplateBlock(block)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("[TemplateBlockService] CreateTemplateBlock DB error: %v\n", err)
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create template block (DB error)", err)
	}
	return id, nil
}

func (s *TemplateBlockService) GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error) {
	block, err := s.repository.GetTemplateBlockByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Template block not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve template block by ID", err)
	}
	return block, nil
}

func (s *TemplateBlockService) ListTemplateBlocksByTemplateID(templateID string) ([]*dto.TemplateBlockDTO, error) {
	blocks, err := s.repository.GetTemplateBlocksByTemplateID(templateID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list template blocks", err)
	}
	return blocks, nil
}

func (s *TemplateBlockService) UpdateTemplateBlock(id string, update *dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error) {
	updatedBlock, err := s.repository.UpdateTemplateBlock(id, update)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to update template block", err)
	}
	return updatedBlock, nil
}

func (s *TemplateBlockService) DeleteTemplateBlock(id string) error {
	err := s.repository.DeleteTemplateBlock(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete template block", err)
	}
	return nil
}
