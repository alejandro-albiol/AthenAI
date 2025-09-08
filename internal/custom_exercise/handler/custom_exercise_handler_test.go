package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/interfaces"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	interfaces.CustomExerciseService
	CreateFn  func(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error)
	GetByIDFn func(gymID, id string) (*dto.CustomExerciseResponseDTO, error)
	ListFn    func(gymID string) ([]*dto.CustomExerciseResponseDTO, error)
	DeleteFn  func(gymID, id string) error
}

func (m *mockService) CreateCustomExercise(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error) {
	return m.CreateFn(gymID, exercise)
}
func (m *mockService) GetCustomExerciseByID(gymID, id string) (*dto.CustomExerciseResponseDTO, error) {
	return m.GetByIDFn(gymID, id)
}
func (m *mockService) ListCustomExercises(gymID string) ([]*dto.CustomExerciseResponseDTO, error) {
	return m.ListFn(gymID)
}
func (m *mockService) DeleteCustomExercise(gymID, id string) error {
	return m.DeleteFn(gymID, id)
}

func TestCreateCustomExerciseHandler(t *testing.T) {
	service := &mockService{
		CreateFn: func(gymID string, exercise *dto.CustomExerciseCreationDTO) (*string, error) {
			if exercise == nil || exercise.Name == "" {
				return nil, assert.AnError
			}
			id := "ex-1"
			return &id, nil
		},
	}
	h := NewCustomExerciseHandler(service)
	body := `{"created_by":"user1","name":"Push Up","description":"A bodyweight exercise","difficulty_level":"easy","exercise_type":"bodyweight","instructions":"Do a push up","video_url":"http://example.com/video","image_url":"http://example.com/image","muscular_groups":["chest","triceps"],"is_active":true}`
	req := httptest.NewRequest(http.MethodPost, "/custom-exercise", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.Create(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateCustomExerciseHandler_InvalidBody(t *testing.T) {
	service := &mockService{}
	h := NewCustomExerciseHandler(service)
	body := `{"created_by":"user1"}`
	req := httptest.NewRequest(http.MethodPost, "/custom-exercise", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.Create(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetCustomExerciseByIDHandler(t *testing.T) {
	service := &mockService{
		GetByIDFn: func(gymID, id string) (*dto.CustomExerciseResponseDTO, error) {
			if id == "notfound" {
				return nil, assert.AnError
			}
			return &dto.CustomExerciseResponseDTO{ID: id, Name: "Push Up"}, nil
		},
	}
	h := NewCustomExerciseHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/custom-exercise/ex-1", nil)
	w := httptest.NewRecorder()

	h.GetByID(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestListCustomExercisesHandler(t *testing.T) {
	service := &mockService{
		ListFn: func(gymID string) ([]*dto.CustomExerciseResponseDTO, error) {
			return []*dto.CustomExerciseResponseDTO{{ID: "ex-1", Name: "Push Up"}}, nil
		},
	}
	h := NewCustomExerciseHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/custom-exercise", nil)
	w := httptest.NewRecorder()

	h.List(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteCustomExerciseHandler(t *testing.T) {
	service := &mockService{
		DeleteFn: func(gymID, id string) error {
			if id == "notfound" {
				return assert.AnError
			}
			return nil
		},
	}
	h := NewCustomExerciseHandler(service)
	req := httptest.NewRequest(http.MethodDelete, "/custom-exercise/ex-1", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
