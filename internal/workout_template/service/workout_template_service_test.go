package service_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/workout_template/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWorkoutTemplateRepository struct {
	mock.Mock
}

func (m *MockWorkoutTemplateRepository) CreateWorkoutTemplate(input *dto.CreateWorkoutTemplateDTO) (*string, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	id := args.Get(0).(string)
	return &id, args.Error(1)
}
func (m *MockWorkoutTemplateRepository) GetWorkoutTemplateByID(id string) (*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateRepository) GetWorkoutTemplateByName(name string) (*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateRepository) GetWorkoutTemplatesByDifficulty(difficulty string) ([]*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(difficulty)
	return args.Get(0).([]*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateRepository) GetWorkoutTemplatesByTargetAudience(targetAudience string) ([]*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(targetAudience)
	return args.Get(0).([]*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateRepository) GetAllWorkoutTemplates() ([]*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateRepository) UpdateWorkoutTemplate(id string, input *dto.UpdateWorkoutTemplateDTO) (*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateRepository) DeleteWorkoutTemplate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateWorkoutTemplate_Success(t *testing.T) {
	repo := new(MockWorkoutTemplateRepository)
	svc := service.NewWorkoutTemplateService(repo)
	createDTO := &dto.CreateWorkoutTemplateDTO{Name: "Test Template"}
	repo.On("GetWorkoutTemplateByName", createDTO.Name).Return(nil, sql.ErrNoRows)
	repo.On("CreateWorkoutTemplate", createDTO).Return("123", nil)

	id, err := svc.CreateWorkoutTemplate(createDTO)
	assert.NoError(t, err)
	if assert.NotNil(t, id) {
		assert.Equal(t, "123", *id)
	}
}

func TestCreateWorkoutTemplate_Conflict(t *testing.T) {
	repo := new(MockWorkoutTemplateRepository)
	svc := service.NewWorkoutTemplateService(repo)
	createDTO := &dto.CreateWorkoutTemplateDTO{Name: "Test Template"}
	repo.On("GetWorkoutTemplateByName", createDTO.Name).Return(&dto.ResponseWorkoutTemplateDTO{ID: "123"}, nil)

	id, err := svc.CreateWorkoutTemplate(createDTO)
	assert.Error(t, err)
	assert.Nil(t, id)
}

func TestGetWorkoutTemplateByID_Success(t *testing.T) {
	repo := new(MockWorkoutTemplateRepository)
	svc := service.NewWorkoutTemplateService(repo)
	repo.On("GetWorkoutTemplateByID", "123").Return(&dto.ResponseWorkoutTemplateDTO{ID: "123", Name: "Test"}, nil)

	template, err := svc.GetWorkoutTemplateByID("123")
	assert.NoError(t, err)
	assert.Equal(t, "123", template.ID)
}

func TestGetWorkoutTemplateByID_NotFound(t *testing.T) {
	repo := new(MockWorkoutTemplateRepository)
	svc := service.NewWorkoutTemplateService(repo)
	repo.On("GetWorkoutTemplateByID", "123").Return(nil, errors.New("not found"))

	template, err := svc.GetWorkoutTemplateByID("123")
	assert.Error(t, err)
	assert.Nil(t, template)
}

// Additional tests for update, delete, and list methods can be added similarly.
