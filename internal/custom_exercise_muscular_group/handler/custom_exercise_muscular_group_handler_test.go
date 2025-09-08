package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
)

type mockService struct {
	CreateLinkFn                 func(gymID string, req *dto.CustomExerciseMuscularGroupCreationDTO) (*string, error)
	DeleteLinkFn                 func(gymID, id string) error
	RemoveAllLinksForExerciseFn  func(gymID, customExerciseID string) error
	GetLinkByIDFn                func(gymID, id string) (*dto.CustomExerciseMuscularGroup, error)
	GetLinksByExerciseIDFn       func(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error)
	GetLinksByMuscularGroupIDFn  func(gymID, muscularGroupID string) ([]*dto.CustomExerciseMuscularGroup, error)
	GetLinksByCustomExerciseIDFn func(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error) // for interface compatibility
}

func (m *mockService) CreateLink(gymID string, req *dto.CustomExerciseMuscularGroupCreationDTO) (*string, error) {
	return m.CreateLinkFn(gymID, req)
}
func (m *mockService) DeleteLink(gymID, id string) error {
	return m.DeleteLinkFn(gymID, id)
}
func (m *mockService) RemoveAllLinksForExercise(gymID, customExerciseID string) error {
	return m.RemoveAllLinksForExerciseFn(gymID, customExerciseID)
}
func (m *mockService) GetLinkByID(gymID, id string) (*dto.CustomExerciseMuscularGroup, error) {
	return m.GetLinkByIDFn(gymID, id)
}
func (m *mockService) GetLinksByExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error) {
	return m.GetLinksByExerciseIDFn(gymID, customExerciseID)
}
func (m *mockService) GetLinksByCustomExerciseID(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error) {
	return m.GetLinksByCustomExerciseIDFn(gymID, customExerciseID)
}
func (m *mockService) GetLinksByMuscularGroupID(gymID, muscularGroupID string) ([]*dto.CustomExerciseMuscularGroup, error) {
	return m.GetLinksByMuscularGroupIDFn(gymID, muscularGroupID)
}

func setTestGymID(ctx context.Context, gymID string) context.Context {
	return context.WithValue(ctx, middleware.GymIDKey, gymID)
}

func TestCreateLink_Success(t *testing.T) {
	service := &mockService{
		CreateLinkFn: func(gymID string, req *dto.CustomExerciseMuscularGroupCreationDTO) (*string, error) {
			id := "123"
			return &id, nil
		},
	}
	h := &CustomExerciseMuscularGroupHandler{service: service}
	body, _ := json.Marshal(dto.CustomExerciseMuscularGroupCreationDTO{CustomExerciseID: "ex1", MuscularGroupID: "mg1"})
	r := httptest.NewRequest("POST", "/custom_exercise_muscular_group", bytes.NewReader(body))
	r = r.WithContext(setTestGymID(r.Context(), "tenant1"))
	w := httptest.NewRecorder()

	h.CreateLink(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// ...additional tests for DeleteLink, RemoveAllLinksForExercise, GetLinkByID, GetLinksByExerciseID, GetLinksByMuscularGroupID...
func TestDeleteLink_Success(t *testing.T) {
	service := &mockService{
		DeleteLinkFn: func(gymID, id string) error { return nil },
	}
	h := &CustomExerciseMuscularGroupHandler{service: service}
	r := httptest.NewRequest("DELETE", "/custom_exercise_muscular_group/id1", nil)
	r = r.WithContext(setTestGymID(r.Context(), "tenant1"))
	w := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "id1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	h.DeleteLink(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRemoveAllLinksForExercise_Success(t *testing.T) {
	service := &mockService{
		RemoveAllLinksForExerciseFn: func(gymID, customExerciseID string) error { return nil },
	}
	h := &CustomExerciseMuscularGroupHandler{service: service}
	r := httptest.NewRequest("DELETE", "/custom_exercise_muscular_group/exercise/ex1", nil)
	r = r.WithContext(setTestGymID(r.Context(), "tenant1"))
	w := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("customExerciseID", "ex1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	h.RemoveAllLinksForExercise(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetLinkByID_Success(t *testing.T) {
	service := &mockService{
		GetLinkByIDFn: func(gymID, id string) (*dto.CustomExerciseMuscularGroup, error) {
			return &dto.CustomExerciseMuscularGroup{ID: id, CustomExerciseID: "ex1", MuscularGroupID: "mg1"}, nil
		},
	}
	h := &CustomExerciseMuscularGroupHandler{service: service}
	r := httptest.NewRequest("GET", "/custom_exercise_muscular_group/id1", nil)
	r = r.WithContext(setTestGymID(r.Context(), "tenant1"))
	w := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "id1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	h.GetLinkByID(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetLinksByExerciseID_Success(t *testing.T) {
	service := &mockService{
		GetLinksByExerciseIDFn: func(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error) {
			return []*dto.CustomExerciseMuscularGroup{{ID: "id1", CustomExerciseID: customExerciseID, MuscularGroupID: "mg1"}}, nil
		},
		GetLinksByCustomExerciseIDFn: func(gymID, customExerciseID string) ([]*dto.CustomExerciseMuscularGroup, error) {
			return []*dto.CustomExerciseMuscularGroup{{ID: "id1", CustomExerciseID: customExerciseID, MuscularGroupID: "mg1"}}, nil
		},
	}
	h := &CustomExerciseMuscularGroupHandler{service: service}
	r := httptest.NewRequest("GET", "/custom_exercise_muscular_group/exercise/ex1", nil)
	r = r.WithContext(setTestGymID(r.Context(), "tenant1"))
	w := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("customExerciseID", "ex1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	h.GetLinksByExerciseID(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetLinksByMuscularGroupID_Success(t *testing.T) {
	service := &mockService{
		GetLinksByMuscularGroupIDFn: func(gymID, muscularGroupID string) ([]*dto.CustomExerciseMuscularGroup, error) {
			return []*dto.CustomExerciseMuscularGroup{{ID: "id1", CustomExerciseID: "ex1", MuscularGroupID: muscularGroupID}}, nil
		},
	}
	h := &CustomExerciseMuscularGroupHandler{service: service}
	r := httptest.NewRequest("GET", "/custom_exercise_muscular_group/muscular_group/mg1", nil)
	r = r.WithContext(setTestGymID(r.Context(), "tenant1"))
	w := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("muscularGroupID", "mg1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	h.GetLinksByMuscularGroupID(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// ...additional tests for DeleteLink, RemoveAllLinksForExercise, GetLinkByID, GetLinksByExerciseID, GetLinksByMuscularGroupID...
