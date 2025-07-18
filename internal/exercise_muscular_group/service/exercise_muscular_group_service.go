package service

import (
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type ExerciseMuscularGroupService struct {
	repository interfaces.ExerciseMuscularGroupRepository
}

func NewExerciseMuscularGroupService(repository interfaces.ExerciseMuscularGroupRepository) *ExerciseMuscularGroupService {
	return &ExerciseMuscularGroupService{repository: repository}
}

func (s *ExerciseMuscularGroupService) CreateLink(link dto.ExerciseMuscularGroup) error {
	_, err := s.repository.CreateLink(link)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create exercise-muscular group link", err)
	}
	return nil
}

func (s *ExerciseMuscularGroupService) DeleteLink(id string) error {
	err := s.repository.DeleteLink(id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete exercise-muscular group link", err)
	}
	return nil
}

func (s *ExerciseMuscularGroupService) GetLinkByID(id string) (dto.ExerciseMuscularGroup, error) {
	link, err := s.repository.FindByID(id)
	if err != nil {
		return dto.ExerciseMuscularGroup{}, apierror.New(errorcode_enum.CodeNotFound, "Link not found", err)
	}
	return link, nil
}

func (s *ExerciseMuscularGroupService) GetLinksByExerciseID(exerciseID string) ([]dto.ExerciseMuscularGroup, error) {
	links, err := s.repository.FindByExerciseID(exerciseID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by exercise ID", err)
	}
	return links, nil
}

func (s *ExerciseMuscularGroupService) GetLinksByMuscularGroupID(muscularGroupID string) ([]dto.ExerciseMuscularGroup, error) {
	links, err := s.repository.FindByMuscularGroupID(muscularGroupID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by muscular group ID", err)
	}
	return links, nil
}
