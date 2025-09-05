package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
)

type mockService struct {
	CreateLinkFunc                func(link *dto.ExerciseMuscularGroup) (*string, error)
	DeleteLinkFunc                func(id string) error
	RemoveAllLinksForExerciseFunc func(exerciseID string) error
	GetLinkByIDFunc               func(id string) (*dto.ExerciseMuscularGroup, error)
	GetLinksByExerciseIDFunc      func(exerciseID string) ([]*dto.ExerciseMuscularGroup, error)
	GetLinksByMuscularGroupIDFunc func(muscularGroupID string) ([]*dto.ExerciseMuscularGroup, error)
}

func (m *mockService) CreateLink(link *dto.ExerciseMuscularGroup) (*string, error) {
	return m.CreateLinkFunc(link)
}
func (m *mockService) DeleteLink(id string) error {
	return m.DeleteLinkFunc(id)
}
func (m *mockService) RemoveAllLinksForExercise(exerciseID string) error {
	return m.RemoveAllLinksForExerciseFunc(exerciseID)
}

func (m *mockService) GetLinkByID(id string) (*dto.ExerciseMuscularGroup, error) {
	return m.GetLinkByIDFunc(id)
}
func (m *mockService) GetLinksByExerciseID(exerciseID string) ([]*dto.ExerciseMuscularGroup, error) {
	return m.GetLinksByExerciseIDFunc(exerciseID)
}
func (m *mockService) GetLinksByMuscularGroupID(muscularGroupID string) ([]*dto.ExerciseMuscularGroup, error) {
	return m.GetLinksByMuscularGroupIDFunc(muscularGroupID)
}

func TestExerciseMuscularGroupHandler_CreateLink(t *testing.T) {
	mockSvc := &mockService{
		CreateLinkFunc: func(link *dto.ExerciseMuscularGroup) (*string, error) {
			id := "ex1"
			return &id, nil
		},
	}
	h := &ExerciseMuscularGroupHandler{service: mockSvc}

	body, _ := json.Marshal(&dto.ExerciseMuscularGroup{ExerciseID: "ex1", MuscularGroupID: "mg1"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateLink(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	// Error case
	mockSvc.CreateLinkFunc = func(link *dto.ExerciseMuscularGroup) (*string, error) {
		return nil, errors.New("service error")
	}
	w = httptest.NewRecorder()
	h.CreateLink(w, req)
	resp = w.Result()
	if resp.StatusCode == http.StatusCreated {
		t.Error("expected error status, got 201")
	}
}

func TestExerciseMuscularGroupHandler_DeleteLink(t *testing.T) {
	mockSvc := &mockService{
		DeleteLinkFunc: func(id string) error {
			return nil
		},
	}
	h := &ExerciseMuscularGroupHandler{service: mockSvc}

	req := httptest.NewRequest(http.MethodDelete, "/link/ex1", nil)
	req = muxSetURLParam(req, "id", "ex1")
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

func TestExerciseMuscularGroupHandler_GetLinkByID(t *testing.T) {
	mockSvc := &mockService{
		GetLinkByIDFunc: func(id string) (*dto.ExerciseMuscularGroup, error) {
			return &dto.ExerciseMuscularGroup{ExerciseID: "ex1", MuscularGroupID: "mg1"}, nil
		},
	}
	h := &ExerciseMuscularGroupHandler{service: mockSvc}

	req := httptest.NewRequest(http.MethodGet, "/link/ex1", nil)
	req = muxSetURLParam(req, "id", "ex1")
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
	mockSvc.GetLinkByIDFunc = func(id string) (*dto.ExerciseMuscularGroup, error) {
		return nil, errors.New("service error")
	}
	w = httptest.NewRecorder()
	h.GetLinkByID(w, req)
	resp = w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Error("expected error status, got 200")
	}
}

// Helper to set chi URL param in test requests
func muxSetURLParam(r *http.Request, key, value string) *http.Request {
	ctx := r.Context()
	ctx = contextWithURLParam(ctx, key, value)
	return r.WithContext(ctx)
}

// contextWithURLParam is a minimal chi-like param setter for tests
func contextWithURLParam(ctx context.Context, key, value string) context.Context {
	type urlParamsKey struct{}
	params := map[string]string{key: value}
	return context.WithValue(ctx, urlParamsKey{}, params)
}

// ...add more tests for other handler methods as needed
