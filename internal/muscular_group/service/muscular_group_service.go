package service

import (
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

func (s *MuscularGroupService) CreateMuscularGroup(mg dto.MuscularGroup) error {
	groups, err := s.repository.GetAllMuscularGroups()
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to check muscular group name uniqueness", err)
	}
	for _, g := range groups {
		if g.Name == mg.Name {
			return apierror.New(errorcode_enum.CodeConflict, "Muscular group with this name already exists", nil)
		}
	}
	_, err = s.repository.CreateMuscularGroup(mg)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create muscular group", err)
	}
	return nil
}

func (s *MuscularGroupService) GetMuscularGroupByID(id string) (dto.MuscularGroup, error) {
	muscularGroup, err := s.repository.FindByID(id)
	if err != nil {
		return dto.MuscularGroup{}, apierror.New(errorcode_enum.CodeNotFound, "Muscular group with ID "+id+" not found", err)
	}
	return muscularGroup, nil
}

func (s *MuscularGroupService) GetAllMuscularGroups() ([]dto.MuscularGroup, error) {
	groups, err := s.repository.GetAllMuscularGroups()
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve muscular groups", err)
	}
	return groups, nil
}

func (s *MuscularGroupService) UpdateMuscularGroup(id string, mg dto.MuscularGroup) error {
	existing, err := s.repository.FindByID(id)
	if err != nil || existing.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, "Muscular group with ID "+id+" not found", err)
	}
	_, err = s.repository.UpdateMuscularGroup(id, mg)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to update muscular group", err)
	}
	return nil
}

func (s *MuscularGroupService) DeleteMuscularGroup(id string) error {
	existing, err := s.repository.FindByID(id)
	if err != nil || existing.ID == "" {
		return apierror.New(errorcode_enum.CodeNotFound, "Muscular group with ID "+id+" not found", err)
	}
	err = s.repository.DeleteMuscularGroup(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete muscular group", err)
	}
	return nil
}
