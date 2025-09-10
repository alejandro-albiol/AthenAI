package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomTemplateBlockService is a mock for testing
type MockCustomTemplateBlockService struct {
	mock.Mock
}

func (m *MockCustomTemplateBlockService) CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (*string, error) {
	args := m.Called(gymID, block)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockCustomTemplateBlockService) GetCustomTemplateBlockByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error) {
	args := m.Called(gymID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseCustomTemplateBlockDTO), args.Error(1)
}

func (m *MockCustomTemplateBlockService) ListCustomTemplateBlocksByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	args := m.Called(gymID, templateID)
	return args.Get(0).([]*dto.ResponseCustomTemplateBlockDTO), args.Error(1)
}

func (m *MockCustomTemplateBlockService) ListCustomTemplateBlocks(gymID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	args := m.Called(gymID)
	return args.Get(0).([]*dto.ResponseCustomTemplateBlockDTO), args.Error(1)
}

func (m *MockCustomTemplateBlockService) UpdateCustomTemplateBlock(gymID, id string, update *dto.UpdateCustomTemplateBlockDTO) (*dto.ResponseCustomTemplateBlockDTO, error) {
	args := m.Called(gymID, id, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseCustomTemplateBlockDTO), args.Error(1)
}

func (m *MockCustomTemplateBlockService) DeleteCustomTemplateBlock(gymID, id string) error {
	args := m.Called(gymID, id)
	return args.Error(0)
}

func TestCreateCustomTemplateBlock(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mockService *MockCustomTemplateBlockService)
		requestBody    interface{}
		expectedStatus int
		gymID          string
	}{
		{
			name: "successful_creation",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				expectedID := "block123"
				mockService.On("CreateCustomTemplateBlock", "gym123", mock.AnythingOfType("*dto.CreateCustomTemplateBlockDTO")).Return(&expectedID, nil)
			},
			requestBody: dto.CreateCustomTemplateBlockDTO{
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
			},
			expectedStatus: http.StatusCreated,
			gymID:          "gym123",
		},
		{
			name:           "missing_required_fields",
			setupMock:      func(mockService *MockCustomTemplateBlockService) {},
			requestBody:    dto.CreateCustomTemplateBlockDTO{BlockName: "Incomplete"},
			expectedStatus: http.StatusBadRequest,
			gymID:          "gym123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCustomTemplateBlockService)
			handler := NewCustomTemplateBlockHandler(mockService)

			tt.setupMock(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/custom-template-block", bytes.NewBuffer(body))

			// Add URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gymId", tt.gymID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.Create(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetCustomTemplateBlock(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mockService *MockCustomTemplateBlockService)
		expectedStatus int
		gymID          string
		blockID        string
	}{
		{
			name: "successful_retrieval",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				expectedBlock := &dto.ResponseCustomTemplateBlockDTO{
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
				mockService.On("GetCustomTemplateBlockByID", "gym123", "block123").Return(expectedBlock, nil)
			},
			expectedStatus: http.StatusOK,
			gymID:          "gym123",
			blockID:        "block123",
		},
		{
			name: "block_not_found",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				mockService.On("GetCustomTemplateBlockByID", "gym123", "nonexistent").Return(nil, apierror.New(errorcode_enum.CodeNotFound, "Custom template block not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			gymID:          "gym123",
			blockID:        "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCustomTemplateBlockService)
			handler := NewCustomTemplateBlockHandler(mockService)

			tt.setupMock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/custom-template-block/"+tt.blockID, nil)

			// Add URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gymId", tt.gymID)
			rctx.URLParams.Add("id", tt.blockID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.GetByID(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestListCustomTemplateBlocksByTemplateID(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mockService *MockCustomTemplateBlockService)
		expectedStatus int
		gymID          string
		templateID     string
	}{
		{
			name: "successful_list",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				expectedBlocks := []*dto.ResponseCustomTemplateBlockDTO{
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
				mockService.On("ListCustomTemplateBlocksByTemplateID", "gym123", "template123").Return(expectedBlocks, nil)
			},
			expectedStatus: http.StatusOK,
			gymID:          "gym123",
			templateID:     "template123",
		},
		{
			name: "empty_list",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				mockService.On("ListCustomTemplateBlocksByTemplateID", "gym123", "template456").Return([]*dto.ResponseCustomTemplateBlockDTO{}, nil)
			},
			expectedStatus: http.StatusOK,
			gymID:          "gym123",
			templateID:     "template456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCustomTemplateBlockService)
			handler := NewCustomTemplateBlockHandler(mockService)

			tt.setupMock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/custom-template-block/template/"+tt.templateID, nil)

			// Add URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gymId", tt.gymID)
			rctx.URLParams.Add("templateId", tt.templateID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.ListByTemplateID(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateCustomTemplateBlock(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mockService *MockCustomTemplateBlockService)
		requestBody    interface{}
		expectedStatus int
		gymID          string
		blockID        string
	}{
		{
			name: "successful_update",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				expectedBlock := &dto.ResponseCustomTemplateBlockDTO{
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
				mockService.On("UpdateCustomTemplateBlock", "gym123", "block123", mock.AnythingOfType("*dto.UpdateCustomTemplateBlockDTO")).Return(expectedBlock, nil)
			},
			requestBody: dto.UpdateCustomTemplateBlockDTO{
				BlockName:       stringPtr("Updated Warm-up"),
				ExerciseCount:   intPtr(5),
				Reps:            intPtr(20),
				Series:          intPtr(4),
				RestTimeSeconds: intPtr(45),
			},
			expectedStatus: http.StatusOK,
			gymID:          "gym123",
			blockID:        "block123",
		},
		{
			name: "block_not_found",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				mockService.On("UpdateCustomTemplateBlock", "gym123", "nonexistent", mock.AnythingOfType("*dto.UpdateCustomTemplateBlockDTO")).Return(nil, apierror.New(errorcode_enum.CodeNotFound, "Custom template block not found", nil))
			},
			requestBody: dto.UpdateCustomTemplateBlockDTO{
				BlockName: stringPtr("Updated Name"),
			},
			expectedStatus: http.StatusNotFound,
			gymID:          "gym123",
			blockID:        "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCustomTemplateBlockService)
			handler := NewCustomTemplateBlockHandler(mockService)

			tt.setupMock(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/custom-template-block/"+tt.blockID, bytes.NewBuffer(body))

			// Add URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gymId", tt.gymID)
			rctx.URLParams.Add("id", tt.blockID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.Update(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteCustomTemplateBlock(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mockService *MockCustomTemplateBlockService)
		expectedStatus int
		gymID          string
		blockID        string
	}{
		{
			name: "successful_deletion",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				mockService.On("DeleteCustomTemplateBlock", "gym123", "block123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			gymID:          "gym123",
			blockID:        "block123",
		},
		{
			name: "block_not_found",
			setupMock: func(mockService *MockCustomTemplateBlockService) {
				mockService.On("DeleteCustomTemplateBlock", "gym123", "nonexistent").Return(apierror.New(errorcode_enum.CodeNotFound, "Custom template block not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			gymID:          "gym123",
			blockID:        "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCustomTemplateBlockService)
			handler := NewCustomTemplateBlockHandler(mockService)

			tt.setupMock(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/custom-template-block/"+tt.blockID, nil)

			// Add URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gymId", tt.gymID)
			rctx.URLParams.Add("id", tt.blockID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.Delete(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

