package service

import (
	"database/sql"
	"errors"

	"github.com/alejandro-albiol/athenai/internal/muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type MuscularGroupService struct {
	repository interfaces.MuscularGroupRepository
}

func NewMuscularGroupService(repository interfaces.MuscularGroupRepository) *MuscularGroupService {
	return &MuscularGroupService{repository: repository}
}

func (s *MuscularGroupService) CreateMuscularGroup(mg *dto.CreateMuscularGroupDTO) (*string, error) {
	// Optimized: Check for name uniqueness using GetMuscularGroupByName
	existing, err := s.repository.GetMuscularGroupByName(mg.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to check muscular group name uniqueness", err)
	}
	if existing != nil {
		return nil, apierror.New(errorcode_enum.CodeConflict, "Muscular group with this name already exists", nil)
	}

	muscularGroupID, err := s.repository.CreateMuscularGroup(mg)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create muscular group", err)
	}
	return muscularGroupID, nil
}

func (s *MuscularGroupService) GetMuscularGroupByID(id string) (*dto.MuscularGroupResponseDTO, error) {
	mg, err := s.repository.GetMuscularGroupByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Muscular group not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve muscular group by ID", err)
	}
	return mg, nil
}

func (s *MuscularGroupService) GetAllMuscularGroups() ([]*dto.MuscularGroupResponseDTO, error) {
	muscularGroups, err := s.repository.GetAllMuscularGroups()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve muscular groups", err)
	}
	return muscularGroups, nil
}

func (s *MuscularGroupService) UpdateMuscularGroup(id string, dto *dto.UpdateMuscularGroupDTO) (*dto.MuscularGroupResponseDTO, error) {
	// Check if muscular group exists
	existing, err := s.repository.GetMuscularGroupByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Muscular group not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve muscular group", err)
	}
	if existing == nil {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "Muscular group not found", nil)
	}

	// If name is being updated, check for uniqueness
	if dto.Name != nil {
		dup, err := s.repository.GetMuscularGroupByName(*dto.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to check name uniqueness", err)
		}
		if dup != nil && dup.ID != id {
			return nil, apierror.New(errorcode_enum.CodeConflict, "Muscular group with this name already exists", nil)
		}
	}

	updated, err := s.repository.UpdateMuscularGroup(id, dto)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to update muscular group", err)
	}
	return updated, nil
}

func (s *MuscularGroupService) DeleteMuscularGroup(id string) error {
	// Check if muscular group exists
	_, err := s.repository.GetMuscularGroupByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierror.New(errorcode_enum.CodeNotFound, "Muscular group not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve muscular group", err)
	}

	err = s.repository.DeleteMuscularGroup(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete muscular group", err)
	}
	return nil
}
