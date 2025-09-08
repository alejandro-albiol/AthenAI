package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/interfaces"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	interfaces.CustomExerciseEquipmentService
	CreateLinkFn func(gymID string, link *dto.CustomExerciseEquipment) (*string, error)
	DeleteLinkFn func(gymID, id string) error
	FindByIDFn   func(gymID, id string) (*dto.CustomExerciseEquipment, error)
}

func (m *mockService) CreateLink(gymID string, link *dto.CustomExerciseEquipment) (*string, error) {
	return m.CreateLinkFn(gymID, link)
}
func (m *mockService) DeleteLink(gymID, id string) error {
	return m.DeleteLinkFn(gymID, id)
}
func (m *mockService) FindByID(gymID, id string) (*dto.CustomExerciseEquipment, error) {
	return m.FindByIDFn(gymID, id)
}

func TestCreateLinkHandler(t *testing.T) {
	service := &mockService{
		CreateLinkFn: func(gymID string, link *dto.CustomExerciseEquipment) (*string, error) {
			if link == nil || link.CustomExerciseID == "" || link.EquipmentID == "" {
				return nil, assert.AnError
			}
			id := "link-1"
			return &id, nil
		},
	}
	h := &CustomExerciseEquipmentHandler{service: service}
	body := `{"custom_exercise_id":"ex-1","equipment_id":"eq-1"}`
	req := httptest.NewRequest(http.MethodPost, "/custom-exercise-equipment", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateLinkHandler_InvalidBody(t *testing.T) {
	service := &mockService{}
	h := &CustomExerciseEquipmentHandler{service: service}
	body := `{"custom_exercise_id":""}`
	req := httptest.NewRequest(http.MethodPost, "/custom-exercise-equipment", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFindByIDHandler(t *testing.T) {
	service := &mockService{
		FindByIDFn: func(gymID, id string) (*dto.CustomExerciseEquipment, error) {
			if id == "notfound" {
				return nil, assert.AnError
			}
			return &dto.CustomExerciseEquipment{ID: id, CustomExerciseID: "ex-1", EquipmentID: "eq-1"}, nil
		},
	}
	h := &CustomExerciseEquipmentHandler{service: service}
	req := httptest.NewRequest(http.MethodGet, "/custom-exercise-equipment/link-1", nil)
	w := httptest.NewRecorder()

	h.FindByID(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
