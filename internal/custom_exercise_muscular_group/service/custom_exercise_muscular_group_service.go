package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomExerciseMuscularGroupService struct {
	repository interfaces.CustomExerciseMuscularGroupRepository
}

func NewCustomExerciseMuscularGroupService(repository interfaces.CustomExerciseMuscularGroupRepository) *CustomExerciseMuscularGroupService {
	return &CustomExerciseMuscularGroupService{repository: repository}
}

func (s *CustomExerciseMuscularGroupService) CreateLink(link dto.CustomExerciseMuscularGroup) (string, error) {
	linkID, err := s.repository.CreateLink(link)
	if err != nil {
		return "", apierror.New(errorcode_enum.CodeInternal, "Failed to create custom_exercise-muscular group link", err)
	}
	return linkID, nil
}

func (s *CustomExerciseMuscularGroupService) DeleteLink(id string) error {
	err := s.repository.DeleteLink(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom_exercise-muscular group link", err)
	}
	return nil
}

func (s *CustomExerciseMuscularGroupService) GetLinkByID(id string) (dto.CustomExerciseMuscularGroup, error) {
	link, err := s.repository.FindByID(id)
	if err != nil {
		return dto.CustomExerciseMuscularGroup{}, apierror.New(errorcode_enum.CodeNotFound, "Link not found", err)
	}
	return link, nil
}

func (s *CustomExerciseMuscularGroupService) GetLinksByCustomExerciseID(customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error) {
	links, err := s.repository.FindByCustomExerciseID(customExerciseID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by custom exercise ID", err)
	}
	return links, nil
}

func (s *CustomExerciseMuscularGroupService) GetLinksByMuscularGroupID(muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error) {
	links, err := s.repository.FindByMuscularGroupID(muscularGroupID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by muscular group ID", err)
	}
	return links, nil
}
