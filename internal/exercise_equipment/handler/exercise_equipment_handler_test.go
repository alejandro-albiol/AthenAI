package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
)

type mockService struct {
	CreateLinkFunc                func(link *dto.ExerciseEquipment) (*string, error)
	DeleteLinkFunc                func(id string) error
	GetLinkByIDFunc               func(id string) (*dto.ExerciseEquipment, error)
	GetLinksByExerciseIDFunc      func(exerciseID string) ([]*dto.ExerciseEquipment, error)
	RemoveAllLinksForExerciseFunc func(exerciseID string) error
	GetLinksByEquipmentIDFunc     func(equipmentID string) ([]*dto.ExerciseEquipment, error)
}

func (m *mockService) CreateLink(link *dto.ExerciseEquipment) (*string, error) {
	return m.CreateLinkFunc(link)
}
func (m *mockService) DeleteLink(id string) error {
	return m.DeleteLinkFunc(id)
}
func (m *mockService) GetLinkByID(id string) (*dto.ExerciseEquipment, error) {
	return m.GetLinkByIDFunc(id)
}
func (m *mockService) GetLinksByExerciseID(exerciseID string) ([]*dto.ExerciseEquipment, error) {
	return m.GetLinksByExerciseIDFunc(exerciseID)
}
func (m *mockService) RemoveAllLinksForExercise(exerciseID string) error {
	if m.RemoveAllLinksForExerciseFunc != nil {
		return m.RemoveAllLinksForExerciseFunc(exerciseID)
	}
	return nil
}
func (m *mockService) GetLinksByEquipmentID(equipmentID string) ([]*dto.ExerciseEquipment, error) {
	return m.GetLinksByEquipmentIDFunc(equipmentID)
}

func TestExerciseEquipmentHandler_CreateLink(t *testing.T) {
	mockSvc := &mockService{
		CreateLinkFunc: func(link *dto.ExerciseEquipment) (*string, error) {
			id := "id1"
			return &id, nil
		},
	}
	h := &ExerciseEquipmentHandler{service: mockSvc}

	body, _ := json.Marshal(&dto.ExerciseEquipment{ExerciseID: "ex1", EquipmentID: "eq1"})
	req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateLink(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Error case
	mockSvc.CreateLinkFunc = func(link *dto.ExerciseEquipment) (*string, error) {
		return nil, errors.New("service error")
	}
	w = httptest.NewRecorder()
	h.CreateLink(w, req)
	resp = w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Error("expected error status, got 200")
	}
}

func TestExerciseEquipmentHandler_DeleteLink(t *testing.T) {
	mockSvc := &mockService{
		DeleteLinkFunc: func(id string) error {
			return nil
		},
	}
	h := &ExerciseEquipmentHandler{service: mockSvc}

	req := httptest.NewRequest(http.MethodDelete, "/link/id1", nil)
	w := httptest.NewRecorder()
	h.DeleteLink(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Error case
	mockSvc.DeleteLinkFunc = func(id string) error {
		return errors.New("service error")
	}
	w = httptest.NewRecorder()
	h.DeleteLink(w, req)
	resp = w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Error("expected error status, got 200")
	}
}

func TestExerciseEquipmentHandler_GetLinkByID(t *testing.T) {
	mockSvc := &mockService{
		GetLinkByIDFunc: func(id string) (*dto.ExerciseEquipment, error) {
			return &dto.ExerciseEquipment{ExerciseID: "ex1", EquipmentID: "eq1"}, nil
		},
	}
	h := &ExerciseEquipmentHandler{service: mockSvc}

	req := httptest.NewRequest(http.MethodGet, "/link/id1", nil)
	w := httptest.NewRecorder()
	h.GetLinkByID(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("ex1")) {
		t.Errorf("expected response body to contain 'ex1', got %s", string(body))
	}

	// Error case
	mockSvc.GetLinkByIDFunc = func(id string) (*dto.ExerciseEquipment, error) {
		return nil, errors.New("service error")
	}
	w = httptest.NewRecorder()
	h.GetLinkByID(w, req)
	resp = w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Error("expected error status, got 200")
	}
}

// ...add more tests for other handler methods as needed
