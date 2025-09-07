package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/template_block/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTemplateBlockService struct {
	mock.Mock
}

func (m *MockTemplateBlockService) CreateTemplateBlock(block *dto.CreateTemplateBlockDTO) (*string, error) {
	args := m.Called(block)
	if id, ok := args.Get(0).(*string); ok {
		return id, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTemplateBlockService) GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.TemplateBlockDTO), args.Error(1)
}

func (m *MockTemplateBlockService) ListTemplateBlocksByTemplateID(templateID string) ([]*dto.TemplateBlockDTO, error) {
	args := m.Called(templateID)
	return args.Get(0).([]*dto.TemplateBlockDTO), args.Error(1)
}

func (m *MockTemplateBlockService) UpdateTemplateBlock(id string, update *dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error) {
	args := m.Called(id, update)
	return args.Get(0).(*dto.TemplateBlockDTO), args.Error(1)
}

func (m *MockTemplateBlockService) DeleteTemplateBlock(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTemplateBlock(t *testing.T) {
	testCases := []struct {
		name       string
		input      dto.CreateTemplateBlockDTO
		setupMock  func(*MockTemplateBlockService)
		wantStatus int
	}{
		{
			name: "successful creation",
			input: dto.CreateTemplateBlockDTO{
				TemplateID:    "template123",
				BlockName:          "Warm-up",
				BlockType:          "warmup",
				BlockOrder:         1,
				ExerciseCount: 3,
			},
			setupMock: func(mockService *MockTemplateBlockService) {
				mockService.On("CreateTemplateBlock", mock.AnythingOfType("*dto.CreateTemplateBlockDTO")).Return("block-123", nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate name conflict",
			input: dto.CreateTemplateBlockDTO{
				TemplateID:    "template123",
				BlockName:          "Warm-up",
				BlockType:          "warmup",
				BlockOrder:         1,
				ExerciseCount: 3,
			},
			setupMock: func(mockService *MockTemplateBlockService) {
				mockService.On("CreateTemplateBlock", mock.AnythingOfType("*dto.CreateTemplateBlockDTO")).Return(
					"", apierror.New(errorcode_enum.CodeConflict, "Template block with name 'Warm-up' already exists in template", nil),
				)
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTemplateBlockService)
			tc.setupMock(mockService)

			h := handler.NewTemplateBlockHandler(mockService)

			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/template-blocks", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			h.CreateTemplateBlock(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			if tc.wantStatus == http.StatusCreated {
				assert.Equal(t, "success", resp.Status)
			} else {
				assert.Equal(t, "error", resp.Status)
				assert.NotEmpty(t, resp.Message)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetTemplateBlock(t *testing.T) {
	testCases := []struct {
		name       string
		blockID    string
		setupMock  func(*MockTemplateBlockService)
		wantStatus int
	}{
		{
			name:    "successful retrieval",
			blockID: "block123",
			setupMock: func(mockService *MockTemplateBlockService) {
				block := &dto.TemplateBlockDTO{
					ID:            "block123",
					TemplateID:    "template123",
					BlockName:          "Warm-up",
					BlockType:          "warmup",
					BlockOrder:         1,
					ExerciseCount: 3,
				}
				mockService.On("GetTemplateBlockByID", "block123").Return(block, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "block not found",
			blockID: "nonexistent",
			setupMock: func(mockService *MockTemplateBlockService) {
				mockService.On("GetTemplateBlockByID", "nonexistent").Return((*dto.TemplateBlockDTO)(nil),
					apierror.New(errorcode_enum.CodeNotFound, "Template block not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTemplateBlockService)
			tc.setupMock(mockService)

			h := handler.NewTemplateBlockHandler(mockService)

			// Use chi router to set URL param for id
			router := chi.NewRouter()
			router.Get("/{id}", h.GetTemplateBlockByID)

			req := httptest.NewRequest(http.MethodGet, "/"+tc.blockID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, "success", resp.Status)
				assert.NotEmpty(t, resp.Data)
			} else {
				assert.Equal(t, "error", resp.Status)
				assert.NotEmpty(t, resp.Message)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestListTemplateBlocksByTemplateID(t *testing.T) {
	testCases := []struct {
		name       string
		templateID string
		setupMock  func(*MockTemplateBlockService)
		wantStatus int
	}{
		{
			name:       "successful list",
			templateID: "template123",
			setupMock: func(mockService *MockTemplateBlockService) {
				blocks := []*dto.TemplateBlockDTO{
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
				mockService.On("ListTemplateBlocksByTemplateID", "template123").Return(blocks, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty list",
			templateID: "template456",
			setupMock: func(mockService *MockTemplateBlockService) {
				mockService.On("ListTemplateBlocksByTemplateID", "template456").Return([]*dto.TemplateBlockDTO{}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTemplateBlockService)
			tc.setupMock(mockService)

			h := handler.NewTemplateBlockHandler(mockService)

			// Use chi router to set URL param for templateId
			router := chi.NewRouter()
			router.Get("/{templateId}/blocks", h.ListTemplateBlocksByTemplateID)

			req := httptest.NewRequest(http.MethodGet, "/"+tc.templateID+"/blocks", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, "success", resp.Status)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateTemplateBlock(t *testing.T) {
	testCases := []struct {
		name       string
		blockID    string
		input      dto.UpdateTemplateBlockDTO
		setupMock  func(*MockTemplateBlockService)
		wantStatus int
	}{
		{
			name:    "successful update",
			blockID: "block123",
			input: dto.UpdateTemplateBlockDTO{
				BlockName:      stringPtr("Updated Warm-up"),
				ExerciseCount:  intPtr(5),
			},
			setupMock: func(mockService *MockTemplateBlockService) {
				updatedBlock := &dto.TemplateBlockDTO{
					ID:            "block123",
					TemplateID:    "template123",
					BlockName:      "Updated Warm-up",
					BlockType:      "warmup",
					BlockOrder:     1,
					ExerciseCount:  5,
				}
				mockService.On("UpdateTemplateBlock", "block123", mock.AnythingOfType("*dto.UpdateTemplateBlockDTO")).Return(updatedBlock, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "block not found",
			blockID: "nonexistent",
			input: dto.UpdateTemplateBlockDTO{
				BlockName: stringPtr("Updated Name"),
			},
			setupMock: func(mockService *MockTemplateBlockService) {
				mockService.On("UpdateTemplateBlock", "nonexistent", mock.AnythingOfType("*dto.UpdateTemplateBlockDTO")).Return(
					(*dto.TemplateBlockDTO)(nil), apierror.New(errorcode_enum.CodeNotFound, "Template block not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTemplateBlockService)
			tc.setupMock(mockService)

			h := handler.NewTemplateBlockHandler(mockService)

			// Use chi router to set URL param for id
			router := chi.NewRouter()
			router.Put("/{id}", h.UpdateTemplateBlock)

			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPut, "/"+tc.blockID, bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, "success", resp.Status)
				assert.NotEmpty(t, resp.Data)
			} else {
				assert.Equal(t, "error", resp.Status)
				assert.NotEmpty(t, resp.Message)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteTemplateBlock(t *testing.T) {
	testCases := []struct {
		name       string
		blockID    string
		setupMock  func(*MockTemplateBlockService)
		wantStatus int
	}{
		{
			name:    "successful deletion",
			blockID: "block123",
			setupMock: func(mockService *MockTemplateBlockService) {
				mockService.On("DeleteTemplateBlock", "block123").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "block not found",
			blockID: "nonexistent",
			setupMock: func(mockService *MockTemplateBlockService) {
				mockService.On("DeleteTemplateBlock", "nonexistent").Return(
					apierror.New(errorcode_enum.CodeNotFound, "Template block not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockTemplateBlockService)
			tc.setupMock(mockService)

			h := handler.NewTemplateBlockHandler(mockService)

			// Use chi router to set URL param for id
			router := chi.NewRouter()
			router.Delete("/{id}", h.DeleteTemplateBlock)

			req := httptest.NewRequest(http.MethodDelete, "/"+tc.blockID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, "success", resp.Status)
			} else {
				assert.Equal(t, "error", resp.Status)
				assert.NotEmpty(t, resp.Message)
			}
			mockService.AssertExpectations(t)
		})
	}
}

// Helper functions for pointer types
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
