package service

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/enum"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type CustomMemberWorkoutService struct {
	repository interfaces.CustomMemberWorkoutRepository
}

func NewCustomMemberWorkoutService(repo interfaces.CustomMemberWorkoutRepository) *CustomMemberWorkoutService {
	return &CustomMemberWorkoutService{repository: repo}
}

func (s *CustomMemberWorkoutService) CreateCustomMemberWorkout(gymID string, memberWorkout *dto.CreateCustomMemberWorkoutDTO) (*string, error) {
	// Validate required fields
	if memberWorkout.MemberID == "" || memberWorkout.WorkoutInstanceID == "" || memberWorkout.ScheduledDate == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "Missing required fields", nil)
	}
	// Validate rating if present
	if memberWorkout.Rating != nil {
		if *memberWorkout.Rating < 1 || *memberWorkout.Rating > 5 {
			return nil, apierror.New(errorcode_enum.CodeBadRequest, "Rating must be between 1 and 5", nil)
		}
	}
	// Status is always scheduled on create
	return s.repository.Create(gymID, memberWorkout)
}

func (s *CustomMemberWorkoutService) GetCustomMemberWorkoutByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
	if id == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "ID is required", nil)
	}
	res, err := s.repository.GetByID(gymID, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, apierror.New(errorcode_enum.CodeNotFound, "Custom member workout not found", err)
        }
        return nil, apierror.New(errorcode_enum.CodeInternal, "Failed to get custom member workout", err)
    }
	return res, nil
}

func (s *CustomMemberWorkoutService) ListCustomMemberWorkoutsByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error) {
	if memberID == "" {
		return nil, apierror.New(errorcode_enum.CodeBadRequest, "MemberID is required", nil)
	}
	return s.repository.ListByMemberID(gymID, memberID)
}

func (s *CustomMemberWorkoutService) UpdateCustomMemberWorkout(gymID string, memberWorkout *dto.UpdateCustomMemberWorkoutDTO) error {
	if memberWorkout.ID == "" {
		return apierror.New(errorcode_enum.CodeBadRequest, "ID is required", nil)
	}
	if memberWorkout.Status != nil {
		status := enum.CustomMemberWorkoutStatus(*memberWorkout.Status)
		if !status.IsValid() {
			return apierror.New(errorcode_enum.CodeBadRequest, "Invalid status value", nil)
		}
	}
	if memberWorkout.Rating != nil {
		if *memberWorkout.Rating < 1 || *memberWorkout.Rating > 5 {
			return apierror.New(errorcode_enum.CodeBadRequest, "Rating must be between 1 and 5", nil)
		}
	}
	return s.repository.Update(gymID, memberWorkout)
}

func (s *CustomMemberWorkoutService) DeleteCustomMemberWorkout(gymID, id string) error {
	if id == "" {
		return apierror.New(errorcode_enum.CodeBadRequest, "ID is required", nil)
	}
	return s.repository.Delete(gymID, id)
}
