package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGymService struct {
	mock.Mock
}

func (m *MockGymService) CreateGym(gym dto.GymCreationDTO) (string, error) {
	args := m.Called(gym)
	if len(args) > 0 {
		return args.String(0), args.Error(1)
	}
	return "", args.Error(0)
}

func (m *MockGymService) GetGymByID(id string) (dto.GymResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymService) GetGymByDomain(domain string) (dto.GymResponseDTO, error) {
	args := m.Called(domain)
	return args.Get(0).(dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymService) GetAllGyms() ([]dto.GymResponseDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymService) UpdateGym(id string, gym dto.GymUpdateDTO) (dto.GymResponseDTO, error) {
	args := m.Called(id, gym)
	return args.Get(0).(dto.GymResponseDTO), args.Error(1)
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
		name       string
		input      dto.GymCreationDTO
		setupMock  func(*MockGymService)
		wantStatus int
	}{
		{
			name: "successful creation",
			input: dto.GymCreationDTO{
				Name:        "Test Gym",
				Domain:      "test-gym",
				Email:       "test@test.com",
				Address:     "123 Test St",
				ContactName: "John Doe",
				Phone:       "+1234567890",
			},
			setupMock: func(mockService *MockGymService) {
				mockService.On("CreateGym", mock.Anything).Return("gym123", nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate domain",
			input: dto.GymCreationDTO{
				Name:        "Test Gym",
				Domain:      "existing-gym",
				Email:       "test@test.com",
				Address:     "123 Test St",
				ContactName: "John Doe",
				Phone:       "+1234567890",
			},
			setupMock: func(mockService *MockGymService) {
				mockService.On("CreateGym", mock.AnythingOfType("dto.GymCreationDTO")).Return("", apierror.New(errorcode_enum.CodeConflict, "Domain already exists"))
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "invalid input",
			input: dto.GymCreationDTO{
				Name:   "", // missing required field
				Domain: "test-gym",
			},
			setupMock: func(mockService *MockGymService) {
				mockService.On("CreateGym", mock.AnythingOfType("dto.GymCreationDTO")).Return("", apierror.New(errorcode_enum.CodeBadRequest, "Invalid input"))
			},
			wantStatus: http.StatusBadRequest,
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
			w := httptest.NewRecorder()

			h.CreateGym(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			var resp response.APIResponse[any]
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			if tc.wantStatus == http.StatusCreated {
				_, ok := resp.Data.(string)
				assert.True(t, ok || resp.Data != nil) // Data should be a string or not nil
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
		name       string
		gymID      string
		setupMock  func(*MockGymService)
		wantStatus int
	}{
		{
			name:  "successful retrieval",
			gymID: "gym123",
			setupMock: func(mockService *MockGymService) {
				mockService.On("GetGymByID", "gym123").Return(dto.GymResponseDTO{
					ID:       "gym123",
					Name:     "Test Gym",
					Domain:   "test-gym",
					Email:    "test@test.com",
					IsActive: true,
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "gym not found",
			gymID: "nonexistent",
			setupMock: func(mockService *MockGymService) {
				mockService.On("GetGymByID", "nonexistent").Return(dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockGymService)
			tc.setupMock(mockService)

			h := handler.NewGymHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/gyms/"+tc.gymID, nil)
			w := httptest.NewRecorder()

			h.GetGymByID(w, req, tc.gymID)

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

func TestGetGymByDomain(t *testing.T) {
	testCases := []struct {
		name       string
		domain     string
		setupMock  func(*MockGymService)
		wantStatus int
	}{
		{
			name:   "successful retrieval",
			domain: "test-gym",
			setupMock: func(mockService *MockGymService) {
				mockService.On("GetGymByDomain", "test-gym").Return(dto.GymResponseDTO{
					ID:       "gym123",
					Name:     "Test Gym",
					Domain:   "test-gym",
					Email:    "test@test.com",
					IsActive: true,
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "gym not found",
			domain: "nonexistent",
			setupMock: func(mockService *MockGymService) {
				mockService.On("GetGymByDomain", "nonexistent").Return(dto.GymResponseDTO{}, apierror.New(errorcode_enum.CodeNotFound, "Gym not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockGymService)
			tc.setupMock(mockService)

			h := handler.NewGymHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/gyms/domain/"+tc.domain, nil)
			w := httptest.NewRecorder()

			h.GetGymByDomain(w, req, tc.domain)

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
