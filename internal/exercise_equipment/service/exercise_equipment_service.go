package service

import (
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type ExerciseEquipmentService struct {
	repository interfaces.ExerciseEquipmentRepository
}

func NewExerciseEquipmentService(repository interfaces.ExerciseEquipmentRepository) *ExerciseEquipmentService {
	return &ExerciseEquipmentService{repository: repository}
}

func (s *ExerciseEquipmentService) CreateLink(link dto.ExerciseEquipment) (string, error) {
	linkID, err := s.repository.CreateLink(link)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create exercise-equipment link", err)
	}
	return linkID, nil
}

func (s *ExerciseEquipmentService) DeleteLink(id string) error {
	err := s.repository.DeleteLink(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete exercise-equipment link", err)
	}
	return nil
}

func (s *ExerciseEquipmentService) RemoveAllLinksForExercise(exerciseID string) error {
	err := s.repository.RemoveAllLinksForExercise(exerciseID)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to remove all links for exercise", err)
	}
	return nil
}

func (s *ExerciseEquipmentService) GetLinkByID(id string) (dto.ExerciseEquipment, error) {
	link, err := s.repository.FindByID(id)
	if err != nil {
		return dto.ExerciseEquipment{}, apierror.New(errorcode_enum.CodeNotFound, "Link not found", err)
	}
	return link, nil
}

func (s *ExerciseEquipmentService) GetLinksByExerciseID(exerciseID string) ([]dto.ExerciseEquipment, error) {
	links, err := s.repository.FindByExerciseID(exerciseID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by exercise ID", err)
	}
	return links, nil
}

func (s *ExerciseEquipmentService) GetLinksByEquipmentID(equipmentID string) ([]dto.ExerciseEquipment, error) {
	links, err := s.repository.FindByEquipmentID(equipmentID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by equipment ID", err)
	}
	if len(links) == 0 {
		return nil, apierror.New(errorcode_enum.CodeNotFound, "No links found for the given equipment ID", nil)
	}
	return links, nil
}
