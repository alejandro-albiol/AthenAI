package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/equipment/dto"
	"github.com/alejandro-albiol/athenai/internal/equipment/router"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEquipmentHandler struct {
	mock.Mock
}

func (m *MockEquipmentHandler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockEquipmentHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockEquipmentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockEquipmentHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockEquipmentHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestNewEquipmentRouter(t *testing.T) {
	mockHandler := new(MockEquipmentHandler)
	r := router.NewEquipmentRouter(mockHandler)

	assert.NotNil(t, r)
}

func TestCreateEquipmentRoute(t *testing.T) {
	mockHandler := new(MockEquipmentHandler)
	mockHandler.On("CreateEquipment", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		response.WriteAPISuccess(w, "Equipment created successfully", map[string]string{"id": "equipment123"})
	})

	r := router.NewEquipmentRouter(mockHandler)

	input := dto.EquipmentCreationDTO{
		Name:        "Treadmill",
		Description: "Cardio equipment for running",
		Category:    "Cardio",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestGetEquipmentRoute(t *testing.T) {
	mockHandler := new(MockEquipmentHandler)
	mockHandler.On("GetEquipment", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		equipment := dto.EquipmentResponseDTO{
			ID:          "equipment123",
			Name:        "Treadmill",
			Description: "Cardio equipment for running",
			Category:    "Cardio",
			IsActive:    true,
			CreatedAt:   "2023-01-01T00:00:00Z",
			UpdatedAt:   "2023-01-01T00:00:00Z",
		}
		response.WriteAPISuccess(w, "Equipment retrieved successfully", equipment)
	})

	r := router.NewEquipmentRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/equipment123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestListEquipmentRoute(t *testing.T) {
	mockHandler := new(MockEquipmentHandler)
	mockHandler.On("ListEquipment", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		equipment := []dto.EquipmentResponseDTO{
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
		response.WriteAPISuccess(w, "Equipment list retrieved successfully", equipment)
	})

	r := router.NewEquipmentRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestUpdateEquipmentRoute(t *testing.T) {
	mockHandler := new(MockEquipmentHandler)
	mockHandler.On("UpdateEquipment", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		updatedEquipment := dto.EquipmentResponseDTO{
			ID:          "equipment123",
			Name:        "Updated Treadmill",
			Description: "Updated description",
			Category:    "Cardio",
			IsActive:    true,
			CreatedAt:   "2023-01-01T00:00:00Z",
			UpdatedAt:   "2023-01-01T01:00:00Z",
		}
		response.WriteAPISuccess(w, "Equipment updated successfully", updatedEquipment)
	})

	r := router.NewEquipmentRouter(mockHandler)

	update := dto.EquipmentUpdateDTO{
		Name:        stringPtr("Updated Treadmill"),
		Description: stringPtr("Updated description"),
	}
	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPut, "/equipment123", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	mockHandler.AssertExpectations(t)
}

func TestDeleteEquipmentRoute(t *testing.T) {
	mockHandler := new(MockEquipmentHandler)
	mockHandler.On("DeleteEquipment", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		response.WriteAPISuccess(w, "Equipment deleted successfully", nil)
	})

	r := router.NewEquipmentRouter(mockHandler)

	req := httptest.NewRequest(http.MethodDelete, "/equipment123", nil)
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
	mockHandler := new(MockEquipmentHandler)
	mockHandler.On("GetEquipment", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		apiErr := apierror.New(errorcode_enum.CodeNotFound, "Equipment not found", nil)
		response.WriteAPIError(w, apiErr)
	})

	r := router.NewEquipmentRouter(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
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

// Helper function for pointer types
func stringPtr(s string) *string {
	return &s
}
