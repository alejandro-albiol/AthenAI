package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/workout_template/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWorkoutTemplateService struct {
	mock.Mock
}

func (m *MockWorkoutTemplateService) CreateWorkoutTemplate(dto *dto.CreateWorkoutTemplateDTO) (string, error) {
	args := m.Called(dto)
	return args.String(0), args.Error(1)
}
func (m *MockWorkoutTemplateService) GetWorkoutTemplateByID(id string) (*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateService) GetWorkoutTemplateByName(name string) (*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}

func (m *MockWorkoutTemplateService) GetWorkoutTemplatesByDifficulty(difficulty string) ([]*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(difficulty)
	return args.Get(0).([]*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateService) GetWorkoutTemplatesByTargetAudience(targetAudience string) ([]*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(targetAudience)
	return args.Get(0).([]*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateService) GetAllWorkoutTemplates() ([]*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}
func (m *MockWorkoutTemplateService) UpdateWorkoutTemplate(id string, updateDTO *dto.UpdateWorkoutTemplateDTO) (*dto.ResponseWorkoutTemplateDTO, error) {
	args := m.Called(id, updateDTO)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseWorkoutTemplateDTO), args.Error(1)
}

func (m *MockWorkoutTemplateService) DeleteWorkoutTemplate(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateWorkoutTemplate_Success(t *testing.T) {
	mockService := new(MockWorkoutTemplateService)
	handler := handler.NewWorkoutTemplateHandler(mockService)
	createDTO := &dto.CreateWorkoutTemplateDTO{Name: "Test Template"}
	mockService.On("CreateWorkoutTemplate", createDTO).Return("123", nil)

	body, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/workout-template", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateWorkoutTemplate(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateWorkoutTemplate_BadRequest(t *testing.T) {
	mockService := new(MockWorkoutTemplateService)
	handler := handler.NewWorkoutTemplateHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/workout-template", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.CreateWorkoutTemplate(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateWorkoutTemplate_Conflict(t *testing.T) {
	mockService := new(MockWorkoutTemplateService)
	handler := handler.NewWorkoutTemplateHandler(mockService)
	createDTO := &dto.CreateWorkoutTemplateDTO{Name: "Test Template"}
	mockService.On("CreateWorkoutTemplate", createDTO).Return("", apierror.New(errorcode_enum.CodeConflict, "conflict", nil))

	body, _ := json.Marshal(createDTO)
	req := httptest.NewRequest(http.MethodPost, "/workout-template", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateWorkoutTemplate(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetWorkoutTemplateByID_Success(t *testing.T) {
	mockService := new(MockWorkoutTemplateService)
	handler := handler.NewWorkoutTemplateHandler(mockService)
	workoutTemplate := &dto.ResponseWorkoutTemplateDTO{ID: "123", Name: "Test"}
	mockService.On("GetWorkoutTemplateByID", "123").Return(workoutTemplate, nil)

	req := httptest.NewRequest(http.MethodGet, "/workout-template/123", nil)
	w := httptest.NewRecorder()

	handler.GetWorkoutTemplateByID(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetWorkoutTemplateByID_NotFound(t *testing.T) {
	mockService := new(MockWorkoutTemplateService)
	handler := handler.NewWorkoutTemplateHandler(mockService)
	mockService.On("GetWorkoutTemplateByID", "123").Return(nil, apierror.New(errorcode_enum.CodeNotFound, "not found", nil))

	req := httptest.NewRequest(http.MethodGet, "/workout-template/123", nil)
	w := httptest.NewRecorder()

	handler.GetWorkoutTemplateByID(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
