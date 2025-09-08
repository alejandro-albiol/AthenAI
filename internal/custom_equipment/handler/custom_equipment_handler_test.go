package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/interfaces"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	interfaces.CustomEquipmentService
	CreateFn  func(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error)
	GetByIDFn func(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error)
	ListFn    func(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error)
	UpdateFn  func(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error
	DeleteFn  func(gymID, id string) error
}

func (m *mockService) CreateCustomEquipment(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error) {
	return m.CreateFn(gymID, equipment)
}
func (m *mockService) GetCustomEquipmentByID(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error) {
	return m.GetByIDFn(gymID, id)
}
func (m *mockService) ListCustomEquipment(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
	return m.ListFn(gymID)
}
func (m *mockService) UpdateCustomEquipment(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error {
	return m.UpdateFn(gymID, equipment)
}
func (m *mockService) DeleteCustomEquipment(gymID, id string) error {
	return m.DeleteFn(gymID, id)
}

func TestCreateCustomEquipmentHandler(t *testing.T) {
	service := &mockService{
		CreateFn: func(gymID string, equipment *dto.CreateCustomEquipmentDTO) (*string, error) {
			if equipment == nil || equipment.Name == "" {
				return nil, assert.AnError
			}
			id := "eq-1"
			return &id, nil
		},
	}
	h := NewCustomEquipmentHandler(service)
	body := `{"created_by":"user123","name":"Dumbbell","description":"A dumbbell","category":"weight","is_active":true}`
	req := httptest.NewRequest(http.MethodPost, "/custom-equipment", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateCustomEquipment(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateCustomEquipmentHandler_InvalidBody(t *testing.T) {
	service := &mockService{}
	h := NewCustomEquipmentHandler(service)
	body := `{"created_by":"user123","name":"Dumbbell","is_active":true}`
	req := httptest.NewRequest(http.MethodPost, "/custom-equipment", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateCustomEquipment(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetCustomEquipmentByIDHandler(t *testing.T) {
	service := &mockService{
		GetByIDFn: func(gymID, id string) (*dto.ResponseCustomEquipmentDTO, error) {
			if id == "notfound" {
				return nil, assert.AnError
			}
			return &dto.ResponseCustomEquipmentDTO{ID: id, Name: "Dumbbell"}, nil
		},
	}
	h := NewCustomEquipmentHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/custom-equipment/eq-1", nil)
	w := httptest.NewRecorder()

	h.GetCustomEquipmentByID(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestListCustomEquipmentHandler(t *testing.T) {
	service := &mockService{
		ListFn: func(gymID string) ([]*dto.ResponseCustomEquipmentDTO, error) {
			return []*dto.ResponseCustomEquipmentDTO{{ID: "eq-1", Name: "Dumbbell"}}, nil
		},
	}
	h := NewCustomEquipmentHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/custom-equipment", nil)
	w := httptest.NewRecorder()

	h.ListCustomEquipment(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateCustomEquipmentHandler(t *testing.T) {
	service := &mockService{
		UpdateFn: func(gymID string, equipment *dto.UpdateCustomEquipmentDTO) error {
			if equipment == nil || equipment.ID == "notfound" {
				return assert.AnError
			}
			return nil
		},
	}
	h := NewCustomEquipmentHandler(service)
	body := `{"id":"eq-1","name":"Barbell","description":"A barbell","category":"weight","is_active":true}`
	req := httptest.NewRequest(http.MethodPut, "/custom-equipment/eq-1", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.UpdateCustomEquipment(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteCustomEquipmentHandler(t *testing.T) {
	service := &mockService{
		DeleteFn: func(gymID, id string) error {
			if id == "notfound" {
				return assert.AnError
			}
			return nil
		},
	}
	h := NewCustomEquipmentHandler(service)
	req := httptest.NewRequest(http.MethodDelete, "/custom-equipment/eq-1", nil)
	w := httptest.NewRecorder()

	h.DeleteCustomEquipment(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
