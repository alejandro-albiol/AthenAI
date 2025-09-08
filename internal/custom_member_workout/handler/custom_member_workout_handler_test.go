package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/handler"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	CreateFn         func(string, *dto.CreateCustomMemberWorkoutDTO) (*string, error)
	GetByIDFn        func(string, string) (*dto.ResponseCustomMemberWorkoutDTO, error)
	ListByMemberIDFn func(string, string) ([]*dto.ResponseCustomMemberWorkoutDTO, error)
	UpdateFn         func(string, *dto.UpdateCustomMemberWorkoutDTO) error
	DeleteFn         func(string, string) error
}

func (m *mockService) CreateCustomMemberWorkout(gymID string, d *dto.CreateCustomMemberWorkoutDTO) (*string, error) {
	return m.CreateFn(gymID, d)
}
func (m *mockService) GetCustomMemberWorkoutByID(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
	return m.GetByIDFn(gymID, id)
}
func (m *mockService) ListCustomMemberWorkoutsByMemberID(gymID, memberID string) ([]*dto.ResponseCustomMemberWorkoutDTO, error) {
	return m.ListByMemberIDFn(gymID, memberID)
}
func (m *mockService) UpdateCustomMemberWorkout(gymID string, d *dto.UpdateCustomMemberWorkoutDTO) error {
	return m.UpdateFn(gymID, d)
}
func (m *mockService) DeleteCustomMemberWorkout(gymID, id string) error {
	return m.DeleteFn(gymID, id)
}

func TestCreateCustomMemberWorkoutHandler_BadRequest(t *testing.T) {
	h := handler.NewCustomMemberWorkoutHandler(&mockService{})
	req := httptest.NewRequest(http.MethodPost, "/custom-member-workout", bytes.NewBuffer([]byte("bad json")))
	req.Header.Set("Authorization", "Bearer testtoken") // simulate JWT for middleware
	w := httptest.NewRecorder()
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCustomMemberWorkoutHandler_NotFound(t *testing.T) {
	h := handler.NewCustomMemberWorkoutHandler(&mockService{
		GetByIDFn: func(gymID, id string) (*dto.ResponseCustomMemberWorkoutDTO, error) {
			return nil, errors.New("not found")
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/custom-member-workout/123", nil)
	req.Header.Set("Authorization", "Bearer testtoken")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.GetByID(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateCustomMemberWorkoutHandler_BadRequest(t *testing.T) {
	h := handler.NewCustomMemberWorkoutHandler(&mockService{})
	req := httptest.NewRequest(http.MethodPut, "/custom-member-workout/123", bytes.NewBuffer([]byte("bad json")))
	req.Header.Set("Authorization", "Bearer testtoken")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.Update(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
