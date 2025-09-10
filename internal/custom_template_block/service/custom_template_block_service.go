package service

import (
	"database/sql"
	"errors"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomTemplateBlockService struct {
	Repo interfaces.CustomTemplateBlockRepository
}

func NewCustomTemplateBlockService(repo interfaces.CustomTemplateBlockRepository) *CustomTemplateBlockService {
	return &CustomTemplateBlockService{Repo: repo}
}

func (s *CustomTemplateBlockService) CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (*string, error) {
	if block == nil {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "Block payload is nil", nil)
	}

	id, err := s.Repo.CreateCustomTemplateBlock(gymID, block)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create custom template block", err)
	}
	return id, nil
}

func (s *CustomTemplateBlockService) GetCustomTemplateBlockByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error) {
	res, err := s.Repo.GetCustomTemplateBlockByID(gymID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Custom template block not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get custom template block", err)
	}
	return res, nil
}

func (s *CustomTemplateBlockService) ListCustomTemplateBlocksByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	res, err := s.Repo.ListCustomTemplateBlocksByTemplateID(gymID, templateID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom template blocks by template ID", err)
	}
	return res, nil
}

func (s *CustomTemplateBlockService) ListCustomTemplateBlocks(gymID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	res, err := s.Repo.ListCustomTemplateBlocks(gymID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to list custom template blocks", err)
	}
	return res, nil
}

func (s *CustomTemplateBlockService) UpdateCustomTemplateBlock(gymID, id string, update *dto.UpdateCustomTemplateBlockDTO) (*dto.ResponseCustomTemplateBlockDTO, error) {
	if err := s.Repo.UpdateCustomTemplateBlock(gymID, id, update); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Custom template block not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to update custom template block", err)
	}
	// Return the updated block
	return s.GetCustomTemplateBlockByID(gymID, id)
}

func (s *CustomTemplateBlockService) DeleteCustomTemplateBlock(gymID, id string) error {
	if err := s.Repo.DeleteCustomTemplateBlock(gymID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierror.New(errorcode_enum.CodeNotFound, "Custom template block not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom template block", err)
	}
	return nil
}
