package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
	"github.com/alejandro-albiol/athenai/internal/user/handler"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create a test request with JWT context using middleware constants
func createTestRequest(method, url string, body []byte, gymID string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))

	// Add JWT context to the request using the same keys as middleware
	ctx := context.WithValue(req.Context(), middleware.GymIDKey, gymID)
	ctx = context.WithValue(ctx, middleware.UserTypeKey, "tenant_user")
	ctx = context.WithValue(ctx, middleware.UserRoleKey, "admin")
	ctx = context.WithValue(ctx, middleware.UserIDKey, "test-user-123")
	req = req.WithContext(ctx)

	return req
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(gymID string, user dto.UserCreationDTO) error {
	args := m.Called(gymID, user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(gymID, id string) (dto.UserResponseDTO, error) {
	args := m.Called(gymID, id)
	return args.Get(0).(dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(gymID, username string) (dto.UserResponseDTO, error) {
	args := m.Called(gymID, username)
	return args.Get(0).(dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(gymID, email string) (dto.UserResponseDTO, error) {
	args := m.Called(gymID, email)
	return args.Get(0).(dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserService) GetAllUsers(gymID string) ([]dto.UserResponseDTO, error) {
	args := m.Called(gymID)
	return args.Get(0).([]dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserService) GetPasswordHashByUsername(gymID, username string) (string, error) {
	args := m.Called(gymID, username)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) UpdateUser(gymID string, id string, user dto.UserUpdateDTO) error {
	args := m.Called(gymID, id, user)
	return args.Error(0)
}

func (m *MockUserService) UpdatePassword(gymID, id string, newPasswordHash string) error {
	args := m.Called(gymID, id, newPasswordHash)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(gymID, id string) error {
	args := m.Called(gymID, id)
	return args.Error(0)
}

func (m *MockUserService) VerifyUser(gymID, userID string) error {
	args := m.Called(gymID, userID)
	return args.Error(0)
}

func (m *MockUserService) SetUserActive(gymID, userID string, active bool) error {
	args := m.Called(gymID, userID, active)
	return args.Error(0)
}

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		input      dto.UserCreationDTO
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:  "successful registration",
			gymID: "gym123",
			input: dto.UserCreationDTO{
				Username: "testuser",
				Email:    "test@test.com",
				Password: "password123",
				Role:     userrole_enum.User,
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("RegisterUser", "gym123", mock.AnythingOfType("dto.UserCreationDTO")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:  "username conflict",
			gymID: "gym123",
			input: dto.UserCreationDTO{
				Username: "existing",
				Email:    "test@test.com",
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("RegisterUser", "gym123", mock.AnythingOfType("dto.UserCreationDTO")).Return(
					apierror.New(errorcode_enum.CodeConflict, "Username already exists", nil),
				)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:  "invalid request body",
			gymID: "gym123",
			input: dto.UserCreationDTO{}, // This will be overridden with invalid JSON
			setupMock: func(mockService *MockUserService) {
				// No mock setup needed as handler should fail before calling service
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			var body []byte
			var err error
			if tc.name == "invalid request body" {
				body = []byte(`{"invalid": json}`) // Invalid JSON
			} else {
				body, err = json.Marshal(tc.input)
				assert.NoError(t, err)
			}

			req := createTestRequest(http.MethodPost, "/users", body, tc.gymID)
			w := httptest.NewRecorder()

			handler.RegisterUser(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestVerifyUser(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		userID     string
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:   "successful verification",
			gymID:  "gym123",
			userID: "user123",
			setupMock: func(mockService *MockUserService) {
				mockService.On("VerifyUser", "gym123", "user123").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			setupMock: func(mockService *MockUserService) {
				mockService.On("VerifyUser", "gym123", "nonexistent").Return(
					apierror.New(errorcode_enum.CodeNotFound, "User not found", nil),
				)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "user already verified",
			gymID:  "gym123",
			userID: "user123",
			setupMock: func(mockService *MockUserService) {
				mockService.On("VerifyUser", "gym123", "user123").Return(
					apierror.New(errorcode_enum.CodeConflict, "User is already verified", nil),
				)
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodPost, "/users/"+tc.userID+"/verify", nil, tc.gymID)

			handler.VerifyUser(w, req, tc.userID)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:  "successful fetch",
			gymID: "gym123",
			setupMock: func(mockService *MockUserService) {
				users := []dto.UserResponseDTO{
					{
						ID:       "user1",
						Username: "testuser1",
						Email:    "test1@test.com",
						Role:     userrole_enum.User,
					},
					{
						ID:       "user2",
						Username: "testuser2",
						Email:    "test2@test.com",
						Role:     userrole_enum.Admin,
					},
				}
				mockService.On("GetAllUsers", "gym123").Return(users, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "empty result",
			gymID: "gym123",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetAllUsers", "gym123").Return([]dto.UserResponseDTO{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "service error",
			gymID: "gym123",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetAllUsers", "gym123").Return([]dto.UserResponseDTO{},
					apierror.New(errorcode_enum.CodeInternal, "Failed to retrieve users", nil))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodGet, "/user", nil, tc.gymID)

			handler.GetAllUsers(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		userID     string
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:   "successful fetch",
			gymID:  "gym123",
			userID: "user123",
			setupMock: func(mockService *MockUserService) {
				user := dto.UserResponseDTO{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@test.com",
					Role:     userrole_enum.User,
				}
				mockService.On("GetUserByID", "gym123", "user123").Return(user, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetUserByID", "gym123", "nonexistent").Return(dto.UserResponseDTO{},
					apierror.New(errorcode_enum.CodeNotFound, "User not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "internal server error",
			gymID:  "gym123",
			userID: "user123",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodGet, "/user/"+tc.userID, nil, tc.gymID)

			handler.GetUserByID(w, req, tc.userID)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		username   string
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:     "successful fetch",
			gymID:    "gym123",
			username: "testuser",
			setupMock: func(mockService *MockUserService) {
				user := dto.UserResponseDTO{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@test.com",
					Role:     userrole_enum.User,
				}
				mockService.On("GetUserByUsername", "gym123", "testuser").Return(user, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "user not found",
			gymID:    "gym123",
			username: "nonexistent",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetUserByUsername", "gym123", "nonexistent").Return(dto.UserResponseDTO{},
					apierror.New(errorcode_enum.CodeNotFound, "User not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "internal server error",
			gymID:    "gym123",
			username: "testuser",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetUserByUsername", "gym123", "testuser").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodGet, "/user/username/"+tc.username, nil, tc.gymID)

			handler.GetUserByUsername(w, req, tc.username)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		email      string
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:  "successful fetch",
			gymID: "gym123",
			email: "test@test.com",
			setupMock: func(mockService *MockUserService) {
				user := dto.UserResponseDTO{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@test.com",
					Role:     userrole_enum.User,
				}
				mockService.On("GetUserByEmail", "gym123", "test@test.com").Return(user, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "user not found",
			gymID: "gym123",
			email: "nonexistent@test.com",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetUserByEmail", "gym123", "nonexistent@test.com").Return(dto.UserResponseDTO{},
					apierror.New(errorcode_enum.CodeNotFound, "User not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "internal server error",
			gymID: "gym123",
			email: "test@test.com",
			setupMock: func(mockService *MockUserService) {
				mockService.On("GetUserByEmail", "gym123", "test@test.com").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodGet, "/user/email/"+tc.email, nil, tc.gymID)

			handler.GetUserByEmail(w, req, tc.email)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		userID     string
		input      dto.UserUpdateDTO
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:   "successful update",
			gymID:  "gym123",
			userID: "user123",
			input: dto.UserUpdateDTO{
				Username: "updateduser",
				Email:    "updated@test.com",
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("UpdateUser", "gym123", "user123", mock.AnythingOfType("dto.UserUpdateDTO")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			input: dto.UserUpdateDTO{
				Username: "updateduser",
				Email:    "updated@test.com",
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("UpdateUser", "gym123", "nonexistent", mock.AnythingOfType("dto.UserUpdateDTO")).Return(
					apierror.New(errorcode_enum.CodeNotFound, "User not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "username conflict",
			gymID:  "gym123",
			userID: "user123",
			input: dto.UserUpdateDTO{
				Username: "existinguser",
				Email:    "updated@test.com",
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("UpdateUser", "gym123", "user123", mock.AnythingOfType("dto.UserUpdateDTO")).Return(
					apierror.New(errorcode_enum.CodeConflict, "Username already exists", nil))
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:   "internal server error",
			gymID:  "gym123",
			userID: "user123",
			input: dto.UserUpdateDTO{
				Username: "updateduser",
				Email:    "updated@test.com",
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("UpdateUser", "gym123", "user123", mock.AnythingOfType("dto.UserUpdateDTO")).Return(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodPut, "/user/"+tc.userID, nil, tc.gymID)

			handler.UpdateUser(w, req, tc.userID, tc.input)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		userID     string
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:   "successful deletion",
			gymID:  "gym123",
			userID: "user123",
			setupMock: func(mockService *MockUserService) {
				mockService.On("DeleteUser", "gym123", "user123").Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			setupMock: func(mockService *MockUserService) {
				mockService.On("DeleteUser", "gym123", "nonexistent").Return(
					apierror.New(errorcode_enum.CodeNotFound, "User not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "internal server error",
			gymID:  "gym123",
			userID: "user123",
			setupMock: func(mockService *MockUserService) {
				mockService.On("DeleteUser", "gym123", "user123").Return(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodDelete, "/user/"+tc.userID, nil, tc.gymID)

			handler.DeleteUser(w, req, tc.userID)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSetUserActive(t *testing.T) {
	testCases := []struct {
		name       string
		gymID      string
		userID     string
		active     bool
		setupMock  func(*MockUserService)
		wantStatus int
	}{
		{
			name:   "successful activation",
			gymID:  "gym123",
			userID: "user123",
			active: true,
			setupMock: func(mockService *MockUserService) {
				mockService.On("SetUserActive", "gym123", "user123", true).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "successful deactivation",
			gymID:  "gym123",
			userID: "user123",
			active: false,
			setupMock: func(mockService *MockUserService) {
				mockService.On("SetUserActive", "gym123", "user123", false).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			active: true,
			setupMock: func(mockService *MockUserService) {
				mockService.On("SetUserActive", "gym123", "nonexistent", true).Return(
					apierror.New(errorcode_enum.CodeNotFound, "User not found", nil))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "user already active",
			gymID:  "gym123",
			userID: "user123",
			active: true,
			setupMock: func(mockService *MockUserService) {
				mockService.On("SetUserActive", "gym123", "user123", true).Return(
					apierror.New(errorcode_enum.CodeConflict, "User is already active", nil))
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:   "internal server error",
			gymID:  "gym123",
			userID: "user123",
			active: true,
			setupMock: func(mockService *MockUserService) {
				mockService.On("SetUserActive", "gym123", "user123", true).Return(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockUserService)
			handler := handler.NewUsersHandler(mockService)
			tc.setupMock(mockService)

			w := httptest.NewRecorder()
			req := createTestRequest(http.MethodPatch, "/user/"+tc.userID+"/active", nil, tc.gymID)

			handler.SetUserActive(w, req, tc.userID, tc.active)

			assert.Equal(t, tc.wantStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRegisterUserInternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUsersHandler(mockService)

	input := dto.UserCreationDTO{
		Username: "testuser",
		Email:    "test@test.com",
		Password: "password123",
		Role:     userrole_enum.User,
	}

	// Mock service returns a non-APIError
	mockService.On("RegisterUser", "gym123", mock.AnythingOfType("dto.UserCreationDTO")).Return(assert.AnError)

	body, err := json.Marshal(input)
	assert.NoError(t, err)

	req := createTestRequest(http.MethodPost, "/user", body, "gym123")
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllUsersInternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUsersHandler(mockService)

	// Mock service returns a non-APIError
	mockService.On("GetAllUsers", "gym123").Return([]dto.UserResponseDTO{}, assert.AnError)

	w := httptest.NewRecorder()
	req := createTestRequest(http.MethodGet, "/user", nil, "gym123")

	handler.GetAllUsers(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserByIDInternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUsersHandler(mockService)

	// Mock service returns a non-APIError
	mockService.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{}, assert.AnError)

	w := httptest.NewRecorder()
	req := createTestRequest(http.MethodGet, "/user/user123", nil, "gym123")

	handler.GetUserByID(w, req, "user123")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateUserInternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUsersHandler(mockService)

	input := dto.UserUpdateDTO{
		Username: "updateduser",
		Email:    "updated@test.com",
	}

	// Mock service returns a non-APIError
	mockService.On("UpdateUser", "gym123", "user123", mock.AnythingOfType("dto.UserUpdateDTO")).Return(assert.AnError)

	w := httptest.NewRecorder()
	req := createTestRequest(http.MethodPut, "/user/user123", nil, "gym123")

	handler.UpdateUser(w, req, "user123", input)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteUserInternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUsersHandler(mockService)

	// Mock service returns a non-APIError
	mockService.On("DeleteUser", "gym123", "user123").Return(assert.AnError)

	w := httptest.NewRecorder()
	req := createTestRequest(http.MethodDelete, "/user/user123", nil, "gym123")

	handler.DeleteUser(w, req, "user123")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestVerifyUserInternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUsersHandler(mockService)

	// Mock service returns a non-APIError
	mockService.On("VerifyUser", "gym123", "user123").Return(assert.AnError)

	w := httptest.NewRecorder()
	req := createTestRequest(http.MethodPost, "/user/user123/verify", nil, "gym123")

	handler.VerifyUser(w, req, "user123")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSetUserActiveInternalError(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUsersHandler(mockService)

	// Mock service returns a non-APIError
	mockService.On("SetUserActive", "gym123", "user123", true).Return(assert.AnError)

	w := httptest.NewRecorder()
	req := createTestRequest(http.MethodPatch, "/user/user123/active", nil, "gym123")

	handler.SetUserActive(w, req, "user123", true)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
