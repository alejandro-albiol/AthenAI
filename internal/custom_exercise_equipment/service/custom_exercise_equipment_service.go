package service

import (
	"database/sql"

	customEquipmentInterfaces "github.com/alejandro-albiol/athenai/internal/custom_equipment/interfaces"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/interfaces"
	publicEquipmentInterfaces "github.com/alejandro-albiol/athenai/internal/equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomExerciseEquipmentService struct {
	repository          interfaces.CustomExerciseEquipmentRepository
	customEquipmentRepo customEquipmentInterfaces.CustomEquipmentRepository
	publicEquipmentRepo publicEquipmentInterfaces.EquipmentRepository
}

func NewCustomExerciseEquipmentService(
	repo interfaces.CustomExerciseEquipmentRepository,
	customEquipmentRepo customEquipmentInterfaces.CustomEquipmentRepository,
	publicEquipmentRepo publicEquipmentInterfaces.EquipmentRepository,
) interfaces.CustomExerciseEquipmentService {
	return &CustomExerciseEquipmentService{
		repository:          repo,
		customEquipmentRepo: customEquipmentRepo,
		publicEquipmentRepo: publicEquipmentRepo,
	}
}

func (s *CustomExerciseEquipmentService) CreateLink(gymID string, link *dto.CustomExerciseEquipment) (*string, error) {
	// Validate equipment exists in either tenant or public table
	var found bool
	if s.customEquipmentRepo != nil {
		eq, err := s.customEquipmentRepo.GetByID(gymID, link.EquipmentID)
		if err == nil && eq != nil {
			found = true
		}
	}
	if !found && s.publicEquipmentRepo != nil {
		eq, err := s.publicEquipmentRepo.GetEquipmentByID(link.EquipmentID)
		if err == nil && eq.ID != "" {
			found = true
		}
	}
	if !found {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "The equipment is nonexistent", nil)
	}

	id, err := s.repository.CreateLink(gymID, link)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to create custom exercise equipment link", err)
	}
	return id, nil
}

func (s *CustomExerciseEquipmentService) DeleteLink(gymID, id string) error {
	err := s.repository.DeleteLink(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Custom exercise equipment link not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom exercise equipment link", err)
	}
	return nil
}

func (s *CustomExerciseEquipmentService) RemoveAllLinksForExercise(gymID, customExerciseID string) error {
	err := s.repository.RemoveAllLinksForExercise(gymID, customExerciseID)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to remove all equipment links for custom exercise", err)
	}
	return nil
}

func (s *CustomExerciseEquipmentService) FindByID(gymID, id string) (*dto.CustomExerciseEquipment, error) {
	link, err := s.repository.FindByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(errorcode_enum.CodeNotFound, "Custom exercise equipment link not found", err)
		}
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to find custom exercise equipment link", err)
	}
	return link, nil
}

func (s *CustomExerciseEquipmentService) FindByCustomExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseEquipment, error) {
	links, err := s.repository.FindByCustomExerciseID(gymID, customExerciseID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to find equipment links for custom exercise", err)
	}
	return links, nil
}

func (s *CustomExerciseEquipmentService) FindByEquipmentID(gymID, equipmentID string) ([]*dto.CustomExerciseEquipment, error) {
	links, err := s.repository.FindByEquipmentID(gymID, equipmentID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to find custom exercise links for equipment", err)
	}
	return links, nil
}
