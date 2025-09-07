package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/equipment/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEquipmentService struct {
	mock.Mock
}

func (m *MockEquipmentService) CreateEquipment(equipment *dto.EquipmentCreationDTO) (*string, error) {
	args := m.Called(equipment)
	if args.Get(0) != nil {
		id := args.Get(0).(string)
		return &id, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockEquipmentService) GetEquipmentByID(id string) (*dto.EquipmentResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.EquipmentResponseDTO), args.Error(1)
}

func (m *MockEquipmentService) GetAllEquipment() ([]*dto.EquipmentResponseDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.EquipmentResponseDTO), args.Error(1)
}

func (m *MockEquipmentService) UpdateEquipment(id string, update *dto.EquipmentUpdateDTO) (*dto.EquipmentResponseDTO, error) {
	args := m.Called(id, update)
	return args.Get(0).(*dto.EquipmentResponseDTO), args.Error(1)
}

func (m *MockEquipmentService) DeleteEquipment(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateEquipment(t *testing.T) {
	testCases := []struct {
		name       string
		input      dto.EquipmentCreationDTO
		setupMock  func(*MockEquipmentService)
		wantStatus int
	}{
		{
			name: "successful creation",
			input: dto.EquipmentCreationDTO{
				Name:        "Treadmill",
				Description: "Cardio equipment for running",
				Category:    "Cardio",
			},
			setupMock: func(mockService *MockEquipmentService) {
				mockService.On("CreateEquipment", mock.AnythingOfType("*dto.EquipmentCreationDTO")).Return("equipment123", nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate name conflict",
			input: dto.EquipmentCreationDTO{
				Name:        "Treadmill",
				Description: "Cardio equipment for running",
				Category:    "Cardio",
			},
			setupMock: func(mockService *MockEquipmentService) {
				mockService.On("CreateEquipment", mock.AnythingOfType("*dto.EquipmentCreationDTO")).Return("",
					apierror.New(errorcode_enum.CodeConflict, "Equipment with name 'Treadmill' already exists", nil),
				)
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockEquipmentService)
			tc.setupMock(mockService)

			h := handler.NewEquipmentHandler(mockService)

			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/equipment", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			h.CreateEquipment(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			if tc.wantStatus == http.StatusOK || tc.wantStatus == http.StatusCreated {
				assert.Equal(t, "success", resp.Status)
			} else {
				assert.Equal(t, "error", resp.Status)
				assert.NotEmpty(t, resp.Message)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetEquipment(t *testing.T) {
	testCases := []struct {
		name        string
		equipmentID string
		setupMock   func(*MockEquipmentService)
		wantStatus  int
	}{
		{
			name:        "successful retrieval",
			equipmentID: "equipment123",
			setupMock: func(mockService *MockEquipmentService) {
				equipment := &dto.EquipmentResponseDTO{
					ID:          "equipment123",
					Name:        "Treadmill",
					Description: "Cardio equipment for running",
					Category:    "Cardio",
					IsActive:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				}
				mockService.On("GetEquipmentByID", "equipment123").Return(equipment, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "equipment not found",
			equipmentID: "nonexistent",
			setupMock: func(mockService *MockEquipmentService) {
				mockService.On("GetEquipmentByID", "nonexistent").Return((*dto.EquipmentResponseDTO)(nil),
					apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockEquipmentService)
			tc.setupMock(mockService)

			h := handler.NewEquipmentHandler(mockService)

			// Use chi router to set URL param for id
			router := chi.NewRouter()
			router.Get("/{id}", h.GetEquipment)

			req := httptest.NewRequest(http.MethodGet, "/"+tc.equipmentID, nil)
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

func TestListEquipment(t *testing.T) {
	testCases := []struct {
		name       string
		setupMock  func(*MockEquipmentService)
		wantStatus int
	}{
		{
			name: "successful list",
			setupMock: func(mockService *MockEquipmentService) {
				equipment := []*dto.EquipmentResponseDTO{
					{
						ID:          "equipment1",
						Name:        "Treadmill",
						Description: "Cardio equipment for running",
						Category:    "Cardio",
						IsActive:    true,
						CreatedAt:   "2023-01-01T00:00:00Z",
						UpdatedAt:   "2023-01-01T00:00:00Z",
					},
					{
						ID:          "equipment2",
						Name:        "Dumbbells",
						Description: "Free weights for strength training",
						Category:    "Strength",
						IsActive:    true,
						CreatedAt:   "2023-01-01T00:00:00Z",
						UpdatedAt:   "2023-01-01T00:00:00Z",
					},
				}
				mockService.On("GetAllEquipment").Return(equipment, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "empty list",
			setupMock: func(mockService *MockEquipmentService) {
				mockService.On("GetAllEquipment").Return([]*dto.EquipmentResponseDTO{}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockEquipmentService)
			tc.setupMock(mockService)

			h := handler.NewEquipmentHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/equipment", nil)
			w := httptest.NewRecorder()

			h.ListEquipment(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, "success", resp.Status)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateEquipment(t *testing.T) {
	testCases := []struct {
		name        string
		equipmentID string
		input       dto.EquipmentUpdateDTO
		setupMock   func(*MockEquipmentService)
		wantStatus  int
	}{
		{
			name:        "successful update",
			equipmentID: "equipment123",
			input: dto.EquipmentUpdateDTO{
				Name:        stringPtr("Updated Treadmill"),
				Description: stringPtr("Updated description"),
			},
			setupMock: func(mockService *MockEquipmentService) {
				updatedEquipment := &dto.EquipmentResponseDTO{
					ID:          "equipment123",
					Name:        "Updated Treadmill",
					Description: "Updated description",
					Category:    "Cardio",
					IsActive:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T01:00:00Z",
				}
				mockService.On("UpdateEquipment", "equipment123", mock.AnythingOfType("*dto.EquipmentUpdateDTO")).Return(updatedEquipment, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "equipment not found",
			equipmentID: "nonexistent",
			input: dto.EquipmentUpdateDTO{
				Name: stringPtr("Updated Name"),
			},
			setupMock: func(mockService *MockEquipmentService) {
				mockService.On("UpdateEquipment", "nonexistent", mock.AnythingOfType("*dto.EquipmentUpdateDTO")).Return(
					(*dto.EquipmentResponseDTO)(nil), apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockEquipmentService)
			tc.setupMock(mockService)

			h := handler.NewEquipmentHandler(mockService)

			// Use chi router to set URL param for id
			router := chi.NewRouter()
			router.Put("/{id}", h.UpdateEquipment)

			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPut, "/"+tc.equipmentID, bytes.NewBuffer(body))
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

func TestDeleteEquipment(t *testing.T) {
	testCases := []struct {
		name        string
		equipmentID string
		setupMock   func(*MockEquipmentService)
		wantStatus  int
	}{
		{
			name:        "successful deletion",
			equipmentID: "equipment123",
			setupMock: func(mockService *MockEquipmentService) {
				mockService.On("DeleteEquipment", "equipment123").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "equipment not found",
			equipmentID: "nonexistent",
			setupMock: func(mockService *MockEquipmentService) {
				mockService.On("DeleteEquipment", "nonexistent").Return(
					apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockEquipmentService)
			tc.setupMock(mockService)

			h := handler.NewEquipmentHandler(mockService)

			// Use chi router to set URL param for id
			router := chi.NewRouter()
			router.Delete("/{id}", h.DeleteEquipment)

			req := httptest.NewRequest(http.MethodDelete, "/"+tc.equipmentID, nil)
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

// Helper function for pointer types
func stringPtr(s string) *string {
	return &s
}
