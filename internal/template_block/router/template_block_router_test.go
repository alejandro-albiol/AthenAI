package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/template_block/router"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTemplateBlockHandler struct {
	mock.Mock
}

func (m *MockTemplateBlockHandler) CreateTemplateBlock(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockTemplateBlockHandler) GetTemplateBlockByID(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockTemplateBlockHandler) ListTemplateBlocksByTemplateID(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// Add missing interface method for filter
func (m *MockTemplateBlockHandler) GetTemplateBlocksByFilter(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockTemplateBlockHandler) UpdateTemplateBlock(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockTemplateBlockHandler) DeleteTemplateBlock(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestNewTemplateBlockRouter(t *testing.T) {
	mockHandler := new(MockTemplateBlockHandler)
	r := router.NewTemplateBlockRouter(mockHandler)

	assert.NotNil(t, r)
}

func TestCreateTemplateBlockRoute(t *testing.T) {
	mockHandler := new(MockTemplateBlockHandler)
	mockHandler.On("CreateTemplateBlock", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		response.WriteAPISuccess(w, "Template block created successfully", nil)
	})

	r := router.NewTemplateBlockRouter(mockHandler)

	input := dto.CreateTemplateBlockDTO{
		TemplateID:    "template123",
		BlockName:          "Warm-up",
		BlockType:          "warmup",
		BlockOrder:         1,
		ExerciseCount: 3,
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/template-blocks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestGetTemplateBlockByIDRoute(t *testing.T) {
	mockHandler := new(MockTemplateBlockHandler)
	mockHandler.On("GetTemplateBlockByID", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		block := dto.TemplateBlockDTO{
			ID:            "block123",
			TemplateID:    "template123",
			BlockName:          "Warm-up",
			BlockType:          "warmup",
			BlockOrder:         1,
			ExerciseCount: 3,
		}
		response.WriteAPISuccess(w, "Template block retrieved successfully", block)
	})

	r := router.NewTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/template-blocks/block123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestListTemplateBlocksByTemplateIDRoute(t *testing.T) {
	mockHandler := new(MockTemplateBlockHandler)
	mockHandler.On("ListTemplateBlocksByTemplateID", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		blocks := []dto.TemplateBlockDTO{
			{
				ID:            "block1",
				TemplateID:    "template123",
				BlockName:          "Warm-up",
				BlockType:          "warmup",
				BlockOrder:         1,
				ExerciseCount: 3,
			},
			{
				ID:            "block2",
				TemplateID:    "template123",
				BlockName:          "Main Set",
				BlockType:          "main",
				BlockOrder:         2,
				ExerciseCount: 5,
			},
		}
		response.WriteAPISuccess(w, "Template blocks retrieved successfully", blocks)
	})

	r := router.NewTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/templates/template123/blocks", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestUpdateTemplateBlockRoute(t *testing.T) {
	mockHandler := new(MockTemplateBlockHandler)
	mockHandler.On("UpdateTemplateBlock", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		updatedBlock := dto.TemplateBlockDTO{
			ID:            "block123",
			TemplateID:    "template123",
			BlockName:          "Updated Warm-up",
			BlockType:          "warmup",
			BlockOrder:         1,
			ExerciseCount: 5,
		}
		response.WriteAPISuccess(w, "Template block updated successfully", updatedBlock)
	})

	r := router.NewTemplateBlockRouter(mockHandler)

	update := dto.UpdateTemplateBlockDTO{
		BlockName:          stringPtr("Updated Warm-up"),
		ExerciseCount: intPtr(5),
	}
	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPut, "/template-blocks/block123", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestDeleteTemplateBlockRoute(t *testing.T) {
	mockHandler := new(MockTemplateBlockHandler)
	mockHandler.On("DeleteTemplateBlock", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		response.WriteAPISuccess(w, "Template block deleted successfully", nil)
	})

	r := router.NewTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodDelete, "/template-blocks/block123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestRouterErrorHandling(t *testing.T) {
	mockHandler := new(MockTemplateBlockHandler)
	mockHandler.On("GetTemplateBlockByID", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		apiErr := apierror.New(errorcode_enum.CodeNotFound, "Template block not found", nil)
		response.WriteAPIError(w, apiErr)
	})

	r := router.NewTemplateBlockRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/template-blocks/nonexistent", nil)
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

// Helper functions for pointer types
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
