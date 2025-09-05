package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/exercise/dto"
	"github.com/go-chi/chi/v5"
)

type mockService struct {
	CreateExerciseFunc    func(ex *dto.ExerciseCreationDTO) (*string, error)
	GetExerciseByIDFunc   func(id string) (*dto.ExerciseResponseDTO, error)
	GetAllExercisesFunc   func() ([]*dto.ExerciseResponseDTO, error)
	UpdateExerciseFunc    func(id string, exercise *dto.ExerciseUpdateDTO) (*dto.ExerciseResponseDTO, error)
	GetExerciseByNameFunc func(name string) (*dto.ExerciseResponseDTO, error)

	GetExercisesByMuscularGroupFunc             func(muscularGroups []string) ([]*dto.ExerciseResponseDTO, error)
	GetExercisesByEquipmentFunc                 func(equipment []string) ([]*dto.ExerciseResponseDTO, error)
	GetExercisesByMuscularGroupAndEquipmentFunc func(muscularGroups []string, equipment []string) ([]*dto.ExerciseResponseDTO, error)
	DeleteExerciseFunc                          func(id string) error
}

func (m *mockService) CreateExercise(ex *dto.ExerciseCreationDTO) (*string, error) {
	return m.CreateExerciseFunc(ex)
}
func (m *mockService) GetExerciseByID(id string) (*dto.ExerciseResponseDTO, error) {
	return m.GetExerciseByIDFunc(id)
}
func (m *mockService) GetExerciseByName(name string) (*dto.ExerciseResponseDTO, error) {
	return m.GetExerciseByNameFunc(name)
}
func (m *mockService) GetExercisesByMuscularGroup(muscularGroups []string) ([]*dto.ExerciseResponseDTO, error) {
	return m.GetExercisesByMuscularGroupFunc(muscularGroups)
}
func (m *mockService) GetExercisesByEquipment(equipment []string) ([]*dto.ExerciseResponseDTO, error) {
	return m.GetExercisesByEquipmentFunc(equipment)
}
func (m *mockService) GetExercisesByMuscularGroupAndEquipment(muscularGroups []string, equipment []string) ([]*dto.ExerciseResponseDTO, error) {
	return m.GetExercisesByMuscularGroupAndEquipmentFunc(muscularGroups, equipment)
}
func (m *mockService) GetAllExercises() ([]*dto.ExerciseResponseDTO, error) {
	return m.GetAllExercisesFunc()
}
func (m *mockService) UpdateExercise(id string, exercise *dto.ExerciseUpdateDTO) (*dto.ExerciseResponseDTO, error) {
	return m.UpdateExerciseFunc(id, exercise)
}
func (m *mockService) DeleteExercise(id string) error {
	return m.DeleteExerciseFunc(id)
}

func TestExerciseHandler_CreateExercise(t *testing.T) {
	ts := []struct {
		name       string
		input      string
		mockFunc   func(*dto.ExerciseCreationDTO) (*string, error)
		wantStatus int
	}{
		{
			name:       "success",
			input:      `{"name":"Pushup","difficulty_level":"beginner","exercise_type":"strength","synonyms":["press-up"],"equipment":["eq1"],"muscular_groups":["mg1"],"instructions":"Do a pushup","created_by":"tester"}`,
			mockFunc:   func(ex *dto.ExerciseCreationDTO) (*string, error) { id := "id1"; return &id, nil },
			wantStatus: 201,
		},
		{
			name:       "bad request",
			input:      `{"name":`,
			mockFunc:   func(ex *dto.ExerciseCreationDTO) (*string, error) { return nil, nil },
			wantStatus: 400,
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockService{CreateExerciseFunc: tc.mockFunc}
			h := &ExerciseHandler{service: service}
			req := httptest.NewRequest("POST", "/exercises", strings.NewReader(tc.input))
			req.Header.Set("Content-Type", "application/json")
			rw := httptest.NewRecorder()
			h.CreateExercise(rw, req)
			if rw.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, rw.Code)
			}
		})
	}
}

func TestExerciseHandler_GetExerciseByID(t *testing.T) {
	service := &mockService{
		GetExerciseByIDFunc: func(id string) (*dto.ExerciseResponseDTO, error) {
			if id == "notfound" {
				return nil, errors.New("not found")
			}
			return &dto.ExerciseResponseDTO{ID: id, Name: "Pushup"}, nil
		},
	}
	h := &ExerciseHandler{service: service}
	req := httptest.NewRequest("GET", "/exercises/id1", nil)
	req = muxSetParam(req, "id", "id1")
	rw := httptest.NewRecorder()
	h.GetExerciseByID(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}

	req = httptest.NewRequest("GET", "/exercises/notfound", nil)
	req = muxSetParam(req, "id", "notfound")
	rw = httptest.NewRecorder()
	h.GetExerciseByID(rw, req)
	if rw.Code == 200 {
		t.Error("expected error status for not found, got 200")
	}
}

func TestExerciseHandler_DeleteExercise(t *testing.T) {
	service := &mockService{
		DeleteExerciseFunc: func(id string) error {
			if id == "fail" {
				return errors.New("delete error")
			}
			return nil
		},
	}
	h := &ExerciseHandler{service: service}
	req := httptest.NewRequest("DELETE", "/exercises/id1", nil)
	req = muxSetParam(req, "id", "id1")
	rw := httptest.NewRecorder()
	h.DeleteExercise(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}

	req = httptest.NewRequest("DELETE", "/exercises/fail", nil)
	req = muxSetParam(req, "id", "fail")
	rw = httptest.NewRecorder()
	h.DeleteExercise(rw, req)
	if rw.Code == 200 {
		t.Error("expected error status for delete error, got 200")
	}
}

func TestExerciseHandler_GetAllExercises(t *testing.T) {
	service := &mockService{
		GetAllExercisesFunc: func() ([]*dto.ExerciseResponseDTO, error) {
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	h := &ExerciseHandler{service: service}
	req := httptest.NewRequest("GET", "/exercises", nil)
	rw := httptest.NewRecorder()
	h.GetAllExercises(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestExerciseHandler_GetExerciseByName(t *testing.T) {
	service := &mockService{
		GetExerciseByNameFunc: func(name string) (*dto.ExerciseResponseDTO, error) {
			if name == "notfound" {
				return nil, errors.New("not found")
			}
			return &dto.ExerciseResponseDTO{ID: "id1", Name: name}, nil
		},
	}
	h := &ExerciseHandler{service: service}
	req := httptest.NewRequest("GET", "/exercises/name/pushup", nil)
	req = muxSetParam(req, "name", "pushup")
	rw := httptest.NewRecorder()
	h.GetExerciseByName(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}

	req = httptest.NewRequest("GET", "/exercises/name/notfound", nil)
	req = muxSetParam(req, "name", "notfound")
	rw = httptest.NewRecorder()
	h.GetExerciseByName(rw, req)
	if rw.Code == 200 {
		t.Error("expected error status for not found, got 200")
	}
}

func TestExerciseHandler_GetExercisesByEquipment(t *testing.T) {
	service := &mockService{
		GetExercisesByEquipmentFunc: func(equipment []string) ([]*dto.ExerciseResponseDTO, error) {
			if len(equipment) > 0 && equipment[0] == "fail" {
				return nil, errors.New("equipment error")
			}
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	h := &ExerciseHandler{service: service}
	req := httptest.NewRequest("GET", "/exercises?equipment=eq1", nil)
	rw := httptest.NewRecorder()
	h.GetExerciseByEquipment(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}

	req = httptest.NewRequest("GET", "/exercises?equipment=fail", nil)
	rw = httptest.NewRecorder()
	h.GetExerciseByEquipment(rw, req)
	if rw.Code == 200 {
		t.Error("expected error status for equipment error, got 200")
	}
}

func TestExerciseHandler_GetExercisesByMuscularGroup(t *testing.T) {
	service := &mockService{
		GetExercisesByMuscularGroupFunc: func(groups []string) ([]*dto.ExerciseResponseDTO, error) {
			if len(groups) > 0 && groups[0] == "fail" {
				return nil, errors.New("group error")
			}
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	h := &ExerciseHandler{service: service}
	req := httptest.NewRequest("GET", "/exercises?group=mg1", nil)
	rw := httptest.NewRecorder()
	h.GetExerciseByMuscularGroup(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}

	req = httptest.NewRequest("GET", "/exercises?group=fail", nil)
	rw = httptest.NewRecorder()
	h.GetExerciseByMuscularGroup(rw, req)
	if rw.Code == 200 {
		t.Error("expected error status for group error, got 200")
	}
}

func TestExerciseHandler_GetExercisesByMuscularGroupAndEquipment(t *testing.T) {
	service := &mockService{
		GetExercisesByMuscularGroupAndEquipmentFunc: func(groups, equipment []string) ([]*dto.ExerciseResponseDTO, error) {
			if (len(groups) > 0 && groups[0] == "fail") || (len(equipment) > 0 && equipment[0] == "fail") {
				return nil, errors.New("filter error")
			}
			return []*dto.ExerciseResponseDTO{{ID: "id1"}}, nil
		},
	}
	h := &ExerciseHandler{service: service}
	req := httptest.NewRequest("GET", "/exercises?group=mg1&equipment=eq1", nil)
	rw := httptest.NewRecorder()
	h.GetExercisesByFilters(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}

	req = httptest.NewRequest("GET", "/exercises?group=fail&equipment=eq1", nil)
	rw = httptest.NewRecorder()
	h.GetExercisesByFilters(rw, req)
	if rw.Code == 200 {
		t.Error("expected error status for filter error, got 200")
	}
}

func TestExerciseHandler_UpdateExercise(t *testing.T) {
	service := &mockService{
		UpdateExerciseFunc: func(id string, ex *dto.ExerciseUpdateDTO) (*dto.ExerciseResponseDTO, error) {
			if id == "fail" {
				return nil, errors.New("update error")
			}
			return &dto.ExerciseResponseDTO{ID: id, Name: "Updated"}, nil
		},
	}
	h := &ExerciseHandler{service: service}
	upd := `{"name":"Updated"}`
	req := httptest.NewRequest("PUT", "/exercises/id1", strings.NewReader(upd))
	req = muxSetParam(req, "id", "id1")
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	h.UpdateExercise(rw, req)
	if rw.Code != 200 {
		t.Errorf("expected status 200, got %d", rw.Code)
	}

	req = httptest.NewRequest("PUT", "/exercises/fail", strings.NewReader(upd))
	req = muxSetParam(req, "id", "fail")
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	h.UpdateExercise(rw, req)
	if rw.Code == 200 {
		t.Error("expected error status for update error, got 200")
	}
}

// muxSetParam is a helper to set chi/mux URL params in tests
func muxSetParam(r *http.Request, key, val string) *http.Request {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}
