package service

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomExerciseEquipmentServiceImpl struct {
	repository interfaces.CustomExerciseEquipmentRepository
}

func NewCustomExerciseEquipmentService(repo interfaces.CustomExerciseEquipmentRepository) interfaces.CustomExerciseEquipmentService {
	return &CustomExerciseEquipmentServiceImpl{repository: repo}
}

func (s *CustomExerciseEquipmentServiceImpl) CreateLink(gymID string, link dto.CustomExerciseEquipment) (string, error) {
	id, err := s.repository.CreateLink(gymID, link)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create custom exercise equipment link", err)
	}
	return id, nil
}

func (s *CustomExerciseEquipmentServiceImpl) DeleteLink(gymID, id string) error {
	err := s.repository.DeleteLink(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(errorcode_enum.CodeNotFound, "Custom exercise equipment link not found", err)
		}
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom exercise equipment link", err)
	}
	return nil
}

func (s *CustomExerciseEquipmentServiceImpl) RemoveAllLinksForExercise(gymID, customExerciseID string) error {
	err := s.repository.RemoveAllLinksForExercise(gymID, customExerciseID)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to remove all equipment links for custom exercise", err)
	}
	return nil
}

func (s *CustomExerciseEquipmentServiceImpl) FindByID(gymID, id string) (dto.CustomExerciseEquipment, error) {
	link, err := s.repository.FindByID(gymID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.CustomExerciseEquipment{}, apierror.New(errorcode_enum.CodeNotFound, "Custom exercise equipment link not found", err)
		}
		return dto.CustomExerciseEquipment{}, apierror.New(errorcode_enum.CodeInternal, "Failed to find custom exercise equipment link", err)
	}
	return link, nil
}

func (s *CustomExerciseEquipmentServiceImpl) FindByCustomExerciseID(gymID, customExerciseID string) ([]dto.CustomExerciseEquipment, error) {
	links, err := s.repository.FindByCustomExerciseID(gymID, customExerciseID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to find equipment links for custom exercise", err)
	}
	return links, nil
}

func (s *CustomExerciseEquipmentServiceImpl) FindByEquipmentID(gymID, equipmentID string) ([]dto.CustomExerciseEquipment, error) {
	links, err := s.repository.FindByEquipmentID(gymID, equipmentID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to find custom exercise links for equipment", err)
	}
	return links, nil
}
