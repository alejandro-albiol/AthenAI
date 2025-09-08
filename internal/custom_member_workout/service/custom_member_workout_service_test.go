package service

import (
	"database/sql"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	CreateFn         func(string, *dto.CreateCustomMemberWorkoutDTO) (*string, error)
	GetByIDFn        func(string, string) (*dto.ResponseCustomMemberWorkoutDTO, error)
	ListByMemberIDFn func(string, string) ([]*dto.ResponseCustomMemberWorkoutDTO, error)
	UpdateFn         func(string, *dto.UpdateCustomMemberWorkoutDTO) error
	DeleteFn         func(string, string) error
}

func (m *mockRepo) Create(gymID string, d *dto.CreateCustomMemberWorkoutDTO) (*string, error) {
	return m.CreateFn(gymID, d)
}
func (m *mockRepo) GetByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
	return m.GetByIDFn(gymID, id)
}
func (m *mockRepo) ListByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error) {
	return m.ListByMemberIDFn(gymID, memberID)
}
func (m *mockRepo) Update(gymID string, d *dto.UpdateCustomMemberWorkoutDTO) error {
	return m.UpdateFn(gymID, d)
}
func (m *mockRepo) Delete(gymID, id string) error {
	return m.DeleteFn(gymID, id)
}

func TestCreateCustomMemberWorkout_Validation(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{})
	cases := []struct {
		name    string
		input   dto.CreateCustomMemberWorkoutDTO
		wantErr string
	}{
		{"missing all", dto.CreateCustomMemberWorkoutDTO{}, errorcode_enum.CodeBadRequest},
		{"missing memberID", dto.CreateCustomMemberWorkoutDTO{WorkoutInstanceID: "w", ScheduledDate: "d"}, errorcode_enum.CodeBadRequest},
		{"bad rating", dto.CreateCustomMemberWorkoutDTO{MemberID: "m", WorkoutInstanceID: "w", ScheduledDate: "d", Rating: intPtr(10)}, errorcode_enum.CodeBadRequest},
	}
	for _, c := range cases {
		_, err := svc.CreateCustomMemberWorkout("gym", &c.input)
		assert.Error(t, err, c.name)
		apiErr := err.(*apierror.APIError)
		assert.Equal(t, c.wantErr, apiErr.Code, c.name)
	}
}

func intPtr(i int) *int { return &i }
func TestCreateCustomMemberWorkout_Success(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{
		CreateFn: func(gymID string, d *dto.CreateCustomMemberWorkoutDTO) (*string, error) {
			id := "okid"
			return &id, nil
		},
	})
	input := &dto.CreateCustomMemberWorkoutDTO{MemberID: "m", WorkoutInstanceID: "w", ScheduledDate: "d"}
	id, err := svc.CreateCustomMemberWorkout("gym", input)
	assert.NoError(t, err)
	assert.Equal(t, "okid", *id)
}
func TestGetCustomMemberWorkoutByID_Success(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{
		GetByIDFn: func(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
			return &dto.ResponseCustomMemberWorkoutDTO{ID: id}, nil
		},
	})
	res, err := svc.GetCustomMemberWorkoutByID("gym", "id1")
	assert.NoError(t, err)
	assert.Equal(t, "id1", res.ID)
}
func TestUpdateCustomMemberWorkout_Success(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{
		UpdateFn: func(gymID string, d *dto.UpdateCustomMemberWorkoutDTO) error { return nil },
	})
	err := svc.UpdateCustomMemberWorkout("gym", &dto.UpdateCustomMemberWorkoutDTO{ID: "id"})
	assert.NoError(t, err)
}
func TestDeleteCustomMemberWorkout_Success(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{
		DeleteFn: func(gymID, id string) error { return nil },
	})
	err := svc.DeleteCustomMemberWorkout("gym", "id")
	assert.NoError(t, err)
}

func TestGetCustomMemberWorkoutByID_NotFound(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{
		GetByIDFn: func(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
			return nil, sql.ErrNoRows
		},
	})
	_, err := svc.GetCustomMemberWorkoutByID("gym", "notfound")
	assert.Error(t, err)
	apiErr := err.(*apierror.APIError)
	assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
}

func TestUpdateCustomMemberWorkout_InvalidStatus(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{})
	badStatus := "bad"
	err := svc.UpdateCustomMemberWorkout("gym", &dto.UpdateCustomMemberWorkoutDTO{ID: "id", Status: &badStatus})
	assert.Error(t, err)
	apiErr := err.(*apierror.APIError)
	assert.Equal(t, errorcode_enum.CodeBadRequest, apiErr.Code)
}

func TestDeleteCustomMemberWorkout_NotFound(t *testing.T) {
	svc := NewCustomMemberWorkoutService(&mockRepo{
		DeleteFn: func(gymID, id string) error { return sql.ErrNoRows },
	})
	err := svc.DeleteCustomMemberWorkout("gym", "notfound")
	assert.Error(t, err)
}
