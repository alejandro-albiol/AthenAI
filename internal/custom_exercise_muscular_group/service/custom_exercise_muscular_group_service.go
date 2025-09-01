package service

import (
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/interfaces"
	publicMuscularGroupInterfaces "github.com/alejandro-albiol/athenai/internal/muscular_group/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomExerciseMuscularGroupService struct {
	repository              interfaces.CustomExerciseMuscularGroupRepository
	publicMuscularGroupRepo publicMuscularGroupInterfaces.MuscularGroupRepository
}

func NewCustomExerciseMuscularGroupService(repository interfaces.CustomExerciseMuscularGroupRepository, publicMuscularGroupRepo publicMuscularGroupInterfaces.MuscularGroupRepository) *CustomExerciseMuscularGroupService {
	return &CustomExerciseMuscularGroupService{
		repository:              repository,
		publicMuscularGroupRepo: publicMuscularGroupRepo,
	}
}

func (s *CustomExerciseMuscularGroupService) CreateLink(gymID string, link dto.CustomExerciseMuscularGroup) error {
	// Validate muscular group exists in public table and is active
	var found bool
	if s.publicMuscularGroupRepo != nil {
		mg, err := s.publicMuscularGroupRepo.GetMuscularGroupByID(link.MuscularGroupID)
		if err == nil && mg.ID != "" && mg.IsActive {
			found = true
		}
	}
	if !found {
		return apierror.New(errorcode_enum.CodeNotFound, "Muscular group not found or not allowed", nil)
	}
	err := s.repository.CreateLink(gymID, link)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to create custom_exercise-muscular group link", err)
	}
	return nil
}

func (s *CustomExerciseMuscularGroupService) DeleteLink(gymID, id string) error {
	err := s.repository.DeleteLink(gymID, id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to delete custom_exercise-muscular group link", err)
	}
	return nil
}

func (s *CustomExerciseMuscularGroupService) RemoveAllLinksForExercise(gymID, id string) error {
	err := s.repository.DeleteLink(gymID, id)
	if err != nil {
		return apierror.New(errorcode_enum.CodeInternal, "Failed to remove all links for exercise", err)
	}
	return nil
}

func (s *CustomExerciseMuscularGroupService) GetLinkByID(gymID, id string) (dto.CustomExerciseMuscularGroup, error) {
	link, err := s.repository.FindByID(gymID, id)
	if err != nil {
		return dto.CustomExerciseMuscularGroup{}, apierror.New(errorcode_enum.CodeNotFound, "Link not found", err)
	}
	return link, nil
}

func (s *CustomExerciseMuscularGroupService) GetLinksByCustomExerciseID(gymID, customExerciseID string) ([]dto.CustomExerciseMuscularGroup, error) {
	links, err := s.repository.FindByCustomExerciseID(gymID, customExerciseID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by custom exercise ID", err)
	}
	return links, nil
}

func (s *CustomExerciseMuscularGroupService) GetLinksByMuscularGroupID(gymID, muscularGroupID string) ([]dto.CustomExerciseMuscularGroup, error) {
	links, err := s.repository.FindByMuscularGroupID(gymID, muscularGroupID)
	if err != nil {
		return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get links by muscular group ID", err)
	}
	return links, nil
}
