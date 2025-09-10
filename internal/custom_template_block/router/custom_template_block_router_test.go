package router_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/router"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomTemplateBlockHandler is a mock for testing
type MockCustomTemplateBlockHandler struct {
	mock.Mock
}

func (m *MockCustomTemplateBlockHandler) Create(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	response.WriteAPICreated(w, "Custom template block created successfully", "block123")
}

func (m *MockCustomTemplateBlockHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	block := dto.ResponseCustomTemplateBlockDTO{
		ID:                       "block123",
		TemplateID:               "template123",
		BlockName:                "Warm-up",
		BlockType:                "warmup",
		BlockOrder:               1,
		ExerciseCount:            3,
		EstimatedDurationMinutes: intPtr(10),
		Instructions:             "Start with light exercises",
		Reps:                     intPtr(15),
		Series:                   intPtr(3),
		RestTimeSeconds:          intPtr(60),
		IsActive:                 true,
		CreatedBy:                "user123",
	}
	response.WriteAPISuccess(w, "Custom template block retrieved successfully", block)
}

func (m *MockCustomTemplateBlockHandler) ListByTemplateID(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	blocks := []*dto.ResponseCustomTemplateBlockDTO{
		{
			ID:                       "block123",
			TemplateID:               "template123",
			BlockName:                "Warm-up",
			BlockType:                "warmup",
			BlockOrder:               1,
			ExerciseCount:            3,
			EstimatedDurationMinutes: intPtr(10),
			Instructions:             "Start with light exercises",
			Reps:                     intPtr(15),
			Series:                   intPtr(3),
			RestTimeSeconds:          intPtr(60),
			IsActive:                 true,
			CreatedBy:                "user123",
		},
	}
	response.WriteAPISuccess(w, "Custom template blocks retrieved successfully", blocks)
}

func (m *MockCustomTemplateBlockHandler) Update(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	updatedBlock := dto.ResponseCustomTemplateBlockDTO{
		ID:                       "block123",
		TemplateID:               "template123",
		BlockName:                "Updated Warm-up",
		BlockType:                "warmup",
		BlockOrder:               1,
		ExerciseCount:            5,
		EstimatedDurationMinutes: intPtr(15),
		Instructions:             "Updated instructions",
		Reps:                     intPtr(20),
		Series:                   intPtr(4),
		RestTimeSeconds:          intPtr(45),
		IsActive:                 true,
		CreatedBy:                "user123",
	}
	response.WriteAPISuccess(w, "Custom template block updated successfully", updatedBlock)
}

func (m *MockCustomTemplateBlockHandler) Delete(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	response.WriteAPISuccess(w, "Custom template block deleted successfully", nil)
}

func TestNewCustomTemplateBlockRouter(t *testing.T) {
	mockHandler := new(MockCustomTemplateBlockHandler)
	r := router.NewCustomTemplateBlockRouter(mockHandler)
	assert.NotNil(t, r)
}

func TestCreateCustomTemplateBlockRoute(t *testing.T) {
	mockHandler := new(MockCustomTemplateBlockHandler)
	mockHandler.On("Create", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request"))

	r := router.NewCustomTemplateBlockRouter(mockHandler)

	create := dto.CreateCustomTemplateBlockDTO{
		TemplateID:               "template123",
		BlockName:                "Warm-up",
		BlockType:                "warmup",
		BlockOrder:               1,
		ExerciseCount:            3,
		EstimatedDurationMinutes: intPtr(10),
		Instructions:             "Start with light exercises",
		Reps:                     intPtr(15),
		Series:                   intPtr(3),
		RestTimeSeconds:          intPtr(60),
		CreatedBy:                "user123",
	}
	body, _ := json.Marshal(create)
	req := httptest.NewRequest(http.MethodPost, "/custom-template-block", bytes.NewBuffer(body))

	// Add context with gymId parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gymId", "gym123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestGetCustomTemplateBlockByIDRoute(t *testing.T) {
	mockHandler := new(MockCustomTemplateBlockHandler)
	mockHandler.On("GetByID", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request"))

	r := router.NewCustomTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/custom-template-block/block123", nil)

	// Add context with gymId and id parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gymId", "gym123")
	rctx.URLParams.Add("id", "block123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestListCustomTemplateBlocksByTemplateIDRoute(t *testing.T) {
	mockHandler := new(MockCustomTemplateBlockHandler)
	mockHandler.On("ListByTemplateID", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request"))

	r := router.NewCustomTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/custom-template-block/template/template123", nil)

	// Add context with gymId and templateID parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gymId", "gym123")
	rctx.URLParams.Add("templateID", "template123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestUpdateCustomTemplateBlockRoute(t *testing.T) {
	mockHandler := new(MockCustomTemplateBlockHandler)
	mockHandler.On("Update", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request"))

	r := router.NewCustomTemplateBlockRouter(mockHandler)

	update := dto.UpdateCustomTemplateBlockDTO{
		BlockName:       stringPtr("Updated Warm-up"),
		ExerciseCount:   intPtr(5),
		Reps:            intPtr(20),
		Series:          intPtr(4),
		RestTimeSeconds: intPtr(45),
	}
	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPut, "/custom-template-block/block123", bytes.NewBuffer(body))

	// Add context with gymId and id parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gymId", "gym123")
	rctx.URLParams.Add("id", "block123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestDeleteCustomTemplateBlockRoute(t *testing.T) {
	mockHandler := new(MockCustomTemplateBlockHandler)
	mockHandler.On("Delete", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request"))

	r := router.NewCustomTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodDelete, "/custom-template-block/block123", nil)

	// Add context with gymId and id parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gymId", "gym123")
	rctx.URLParams.Add("id", "block123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

// ErrorMockHandler for testing error cases
type ErrorMockHandler struct {
	mock.Mock
}

func (m *ErrorMockHandler) Create(w http.ResponseWriter, r *http.Request)           {}
func (m *ErrorMockHandler) ListByTemplateID(w http.ResponseWriter, r *http.Request) {}
func (m *ErrorMockHandler) Update(w http.ResponseWriter, r *http.Request)           {}
func (m *ErrorMockHandler) Delete(w http.ResponseWriter, r *http.Request)           {}

func (m *ErrorMockHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	apiErr := apierror.New(errorcode_enum.CodeNotFound, "Custom template block not found", nil)
	response.WriteAPIError(w, apiErr)
}

func TestRouterErrorHandling(t *testing.T) {
	mockHandler := new(ErrorMockHandler)
	mockHandler.On("GetByID", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request"))

	r := router.NewCustomTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/custom-template-block/nonexistent", nil)

	// Add context with gymId and id parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gymId", "gym123")
	rctx.URLParams.Add("id", "nonexistent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.NotEmpty(t, resp.Message)
	mockHandler.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

