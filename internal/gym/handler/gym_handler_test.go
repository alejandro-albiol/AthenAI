package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGymService struct {
	mock.Mock
}

func (m *MockGymService) CreateGym(gym *dto.GymCreationDTO) (string, error) {
	args := m.Called(gym)
	if len(args) > 0 {
		return args.String(0), args.Error(1)
	}
	return "", args.Error(0)
}

func (m *MockGymService) GetGymByID(id string) (*dto.GymResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymService) GetGymByName(name string) (*dto.GymResponseDTO, error) {
	args := m.Called(name)
	return args.Get(0).(*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymService) GetAllGyms() ([]*dto.GymResponseDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymService) UpdateGym(id string, gym *dto.GymUpdateDTO) (*dto.GymResponseDTO, error) {
	args := m.Called(id, gym)
	return args.Get(0).(*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymService) SetGymActive(id string, active bool) error {
	args := m.Called(id, active)
	return args.Error(0)
}

func (m *MockGymService) DeleteGym(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateGym(t *testing.T) {
	testCases := []struct {
		name        string
		input       dto.GymCreationDTO
		setupMock   func(*MockGymService)
		wantStatus  int
		gymIDHeader string // For multi-tenancy, if needed
	}{
		{
			name: "successful creation",
			input: dto.GymCreationDTO{
				Name:    "Test Gym",
				Email:   "test@test.com",
				Address: "123 Test St",
				Phone:   "+1234567890",
			},
			setupMock: func(mockService *MockGymService) {
				mockService.On("CreateGym", mock.AnythingOfType("*dto.GymCreationDTO")).Return("gym-uuid-123", nil)
			},
			wantStatus:  http.StatusCreated,
			gymIDHeader: "gym-uuid-123",
		},
		{
			name: "invalid input",
			input: dto.GymCreationDTO{
				Name: "", // missing required field
			},
			setupMock: func(mockService *MockGymService) {
				mockService.On("CreateGym", mock.AnythingOfType("*dto.GymCreationDTO")).Return("", apierror.New(errorcode_enum.CodeBadRequest, "Invalid input", nil))
			},
			wantStatus:  http.StatusBadRequest,
			gymIDHeader: "gym-uuid-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockGymService)
			if tc.setupMock != nil {
				tc.setupMock(mockService)
			}

			h := handler.NewGymHandler(mockService)

			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/gyms", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			// Inject platform_admin into context for security check
			ctx := req.Context()
			ctx = context.WithValue(ctx, middleware.UserTypeKey, "platform_admin")
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			h.CreateGym(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			if tc.wantStatus == http.StatusCreated {
				assert.NotEmpty(t, resp.Data)
			} else {
				assert.Equal(t, "error", resp.Status)
				assert.NotEmpty(t, resp.Message)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetGymByID(t *testing.T) {
	testCases := []struct {
		name        string
		gymID       string
		setupMock   func(*MockGymService)
		wantStatus  int
		gymIDHeader string
	}{
		{
			name:  "successful retrieval",
			gymID: "gym-uuid-123",
			setupMock: func(mockService *MockGymService) {
				mockService.On("GetGymByID", "gym-uuid-123").Return(&dto.GymResponseDTO{
					ID:       "gym-uuid-123",
					Name:     "Test Gym",
					Email:    "test@test.com",
					IsActive: true,
				}, nil)
			},
			wantStatus:  http.StatusOK,
			gymIDHeader: "gym-uuid-123",
		},
		{
			name:  "gym not found",
			gymID: "nonexistent",
			setupMock: func(mockService *MockGymService) {
				mockService.On("GetGymByID", "nonexistent").Return(&dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found", nil))
			},
			wantStatus:  http.StatusNotFound,
			gymIDHeader: "gym-uuid-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockGymService)
			tc.setupMock(mockService)

			h := handler.NewGymHandler(mockService)

			// Use chi router to set URL param for id
			router := chi.NewRouter()
			router.Get("/{id}", h.GetGymByID)

			req := httptest.NewRequest(http.MethodGet, "/"+tc.gymID, nil)
			// Inject platform_admin into context for security check
			ctx := req.Context()
			ctx = context.WithValue(ctx, middleware.UserTypeKey, "platform_admin")
			ctx = context.WithValue(ctx, middleware.GymIDKey, tc.gymID)
			req = req.WithContext(ctx)
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

// Removed: TestGetGymByDomain. All gym lookups are now UUID-only.
