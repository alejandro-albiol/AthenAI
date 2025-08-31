package service

import (
	"testing"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(gymID string, user dto.UserCreationDTO) (string, error) {
	args := m.Called(gymID, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(gymID, id string) (dto.UserResponseDTO, error) {
	args := m.Called(gymID, id)
	return args.Get(0).(dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(gymID, username string) (dto.UserResponseDTO, error) {
	args := m.Called(gymID, username)
	return args.Get(0).(dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(gymID, email string) (dto.UserResponseDTO, error) {
	args := m.Called(gymID, email)
	return args.Get(0).(dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(gymID string) ([]dto.UserResponseDTO, error) {
	args := m.Called(gymID)
	return args.Get(0).([]dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(gymID, id string, user dto.UserUpdateDTO) error {
	args := m.Called(gymID, id, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(gymID, id string) error {
	args := m.Called(gymID, id)
	return args.Error(0)
}

func (m *MockUserRepository) VerifyUser(gymID, userID string) error {
	args := m.Called(gymID, userID)
	return args.Error(0)
}

func (m *MockUserRepository) SetUserActive(gymID, userID string, active bool) error {
	args := m.Called(gymID, userID, active)
	return args.Error(0)
}

func (m *MockUserRepository) GetPasswordHashByUsername(gymID, username string) (string, error) {
	args := m.Called(gymID, username)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(gymID, userID, newPasswordHash string) error {
	args := m.Called(gymID, userID, newPasswordHash)
	return args.Error(0)
}

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		userDTO   dto.UserCreationDTO
		mockSetup func(*MockUserRepository)
		wantErr   bool
	}{
		{
			name:  "successful registration",
			gymID: "gym123",
			userDTO: dto.UserCreationDTO{
				Username: "testuser",
				Email:    "test@test.com",
				Password: "password123",
				Role:     userrole_enum.User,
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByUsername", "gym123", "testuser").Return(dto.UserResponseDTO{}, nil)
				mockRepo.On("GetUserByEmail", "gym123", "test@test.com").Return(dto.UserResponseDTO{}, nil)
				mockRepo.On("CreateUser", "gym123", mock.AnythingOfType("dto.UserCreationDTO")).Return("user-123", nil)
			},
			wantErr: false,
		},
		{
			name:  "username already exists",
			gymID: "gym123",
			userDTO: dto.UserCreationDTO{
				Username: "existing",
				Email:    "test@test.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByUsername", "gym123", "existing").Return(dto.UserResponseDTO{ID: "1"}, nil)
			},
			wantErr: true,
		},
		{
			name:  "email already exists",
			gymID: "gym123",
			userDTO: dto.UserCreationDTO{
				Username: "newuser",
				Email:    "existing@test.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByUsername", "gym123", "newuser").Return(dto.UserResponseDTO{}, nil)
				mockRepo.On("GetUserByEmail", "gym123", "existing@test.com").Return(dto.UserResponseDTO{ID: "1"}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			userID, err := service.RegisterUser(tc.gymID, tc.userDTO)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, userID)
			}
		})
	}
}

func TestVerifyUser(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		userID    string
		mockSetup func(*MockUserRepository)
		wantErr   bool
	}{
		{
			name:   "successful verification",
			gymID:  "gym123",
			userID: "user123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					Verified: false,
				}, nil)
				mockRepo.On("VerifyUser", "gym123", "user123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "already verified",
			gymID:  "gym123",
			userID: "user123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					Verified: true,
				}, nil)
			},
			wantErr: true,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "nonexistent").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			err := service.VerifyUser(tc.gymID, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		userID    string
		mockSetup func(*MockUserRepository)
		wantErr   bool
		want      dto.UserResponseDTO
	}{
		{
			name:   "successful fetch",
			gymID:  "gym123",
			userID: "user123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@test.com",
					Role:     userrole_enum.User,
					GymID:    "gym123",
				}, nil)
			},
			wantErr: false,
			want: dto.UserResponseDTO{
				ID:       "user123",
				Username: "testuser",
				Email:    "test@test.com",
				Role:     userrole_enum.User,
				GymID:    "gym123",
			},
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "nonexistent").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
			want:    dto.UserResponseDTO{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			got, err := service.GetUserByID(tc.gymID, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		username  string
		mockSetup func(*MockUserRepository)
		wantErr   bool
		want      dto.UserResponseDTO
	}{
		{
			name:     "successful fetch",
			gymID:    "gym123",
			username: "testuser",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByUsername", "gym123", "testuser").Return(dto.UserResponseDTO{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@test.com",
					Role:     userrole_enum.User,
					GymID:    "gym123",
				}, nil)
			},
			wantErr: false,
			want: dto.UserResponseDTO{
				ID:       "user123",
				Username: "testuser",
				Email:    "test@test.com",
				Role:     userrole_enum.User,
				GymID:    "gym123",
			},
		},
		{
			name:     "user not found",
			gymID:    "gym123",
			username: "nonexistent",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByUsername", "gym123", "nonexistent").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
			want:    dto.UserResponseDTO{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			got, err := service.GetUserByUsername(tc.gymID, tc.username)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		email     string
		mockSetup func(*MockUserRepository)
		wantErr   bool
		want      dto.UserResponseDTO
	}{
		{
			name:  "successful fetch",
			gymID: "gym123",
			email: "test@test.com",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByEmail", "gym123", "test@test.com").Return(dto.UserResponseDTO{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@test.com",
					Role:     userrole_enum.User,
					GymID:    "gym123",
				}, nil)
			},
			wantErr: false,
			want: dto.UserResponseDTO{
				ID:       "user123",
				Username: "testuser",
				Email:    "test@test.com",
				Role:     userrole_enum.User,
				GymID:    "gym123",
			},
		},
		{
			name:  "user not found",
			gymID: "gym123",
			email: "nonexistent@test.com",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByEmail", "gym123", "nonexistent@test.com").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
			want:    dto.UserResponseDTO{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			got, err := service.GetUserByEmail(tc.gymID, tc.email)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetPasswordHashByUsername(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		username  string
		mockSetup func(*MockUserRepository)
		wantErr   bool
		want      string
	}{
		{
			name:     "successful fetch",
			gymID:    "gym123",
			username: "testuser",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetPasswordHashByUsername", "gym123", "testuser").Return("hashedpassword123", nil)
			},
			wantErr: false,
			want:    "hashedpassword123",
		},
		{
			name:     "user not found",
			gymID:    "gym123",
			username: "nonexistent",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetPasswordHashByUsername", "gym123", "nonexistent").Return("", assert.AnError)
			},
			wantErr: true,
			want:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			got, err := service.GetPasswordHashByUsername(tc.gymID, tc.username)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		mockSetup func(*MockUserRepository)
		wantErr   bool
		want      []dto.UserResponseDTO
	}{
		{
			name:  "successful fetch with multiple users",
			gymID: "gym123",
			mockSetup: func(mockRepo *MockUserRepository) {
				users := []dto.UserResponseDTO{
					{
						ID:       "user1",
						Username: "testuser1",
						Email:    "test1@test.com",
						Role:     userrole_enum.User,
						GymID:    "gym123",
					},
					{
						ID:       "user2",
						Username: "testuser2",
						Email:    "test2@test.com",
						Role:     userrole_enum.Admin,
						GymID:    "gym123",
					},
				}
				mockRepo.On("GetAllUsers", "gym123").Return(users, nil)
			},
			wantErr: false,
			want: []dto.UserResponseDTO{
				{
					ID:       "user1",
					Username: "testuser1",
					Email:    "test1@test.com",
					Role:     userrole_enum.User,
					GymID:    "gym123",
				},
				{
					ID:       "user2",
					Username: "testuser2",
					Email:    "test2@test.com",
					Role:     userrole_enum.Admin,
					GymID:    "gym123",
				},
			},
		},
		{
			name:  "successful fetch with empty result",
			gymID: "gym123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetAllUsers", "gym123").Return([]dto.UserResponseDTO{}, nil)
			},
			wantErr: false,
			want:    []dto.UserResponseDTO{},
		},
		{
			name:  "repository error",
			gymID: "gym123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetAllUsers", "gym123").Return([]dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
			want:    []dto.UserResponseDTO(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			got, err := service.GetAllUsers(tc.gymID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		userID    string
		userDTO   dto.UserUpdateDTO
		mockSetup func(*MockUserRepository)
		wantErr   bool
	}{
		{
			name:   "successful update",
			gymID:  "gym123",
			userID: "user123",
			userDTO: dto.UserUpdateDTO{
				Username: "updateduser",
				Email:    "updated@test.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					Username: "olduser",
					Email:    "old@test.com",
					Role:     userrole_enum.User,
					GymID:    "gym123",
				}, nil)
				// Username doesn't already exist
				mockRepo.On("GetUserByUsername", "gym123", "updateduser").Return(dto.UserResponseDTO{}, assert.AnError)
				// Email doesn't already exist
				mockRepo.On("GetUserByEmail", "gym123", "updated@test.com").Return(dto.UserResponseDTO{}, assert.AnError)
				// Update succeeds
				mockRepo.On("UpdateUser", "gym123", "user123", mock.AnythingOfType("dto.UserUpdateDTO")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			userDTO: dto.UserUpdateDTO{
				Username: "updateduser",
				Email:    "updated@test.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "nonexistent").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:   "username already exists",
			gymID:  "gym123",
			userID: "user123",
			userDTO: dto.UserUpdateDTO{
				Username: "existinguser",
				Email:    "updated@test.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					Username: "olduser",
					Email:    "old@test.com",
					Role:     userrole_enum.User,
					GymID:    "gym123",
				}, nil)
				// Username already exists (different user)
				mockRepo.On("GetUserByUsername", "gym123", "existinguser").Return(dto.UserResponseDTO{
					ID: "otheruser",
				}, nil)
			},
			wantErr: true,
		},
		{
			name:   "email already exists",
			gymID:  "gym123",
			userID: "user123",
			userDTO: dto.UserUpdateDTO{
				Username: "updateduser",
				Email:    "existing@test.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					Username: "olduser",
					Email:    "old@test.com",
					Role:     userrole_enum.User,
					GymID:    "gym123",
				}, nil)
				// Username doesn't already exist
				mockRepo.On("GetUserByUsername", "gym123", "updateduser").Return(dto.UserResponseDTO{}, assert.AnError)
				// Email already exists (different user)
				mockRepo.On("GetUserByEmail", "gym123", "existing@test.com").Return(dto.UserResponseDTO{
					ID: "otheruser",
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			err := service.UpdateUser(tc.gymID, tc.userID, tc.userDTO)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	testCases := []struct {
		name        string
		gymID       string
		userID      string
		newPassword string
		mockSetup   func(*MockUserRepository)
		wantErr     bool
	}{
		{
			name:        "successful password update",
			gymID:       "gym123",
			userID:      "user123",
			newPassword: "newpassword123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:    "user123",
					GymID: "gym123",
				}, nil)
				// Update succeeds
				mockRepo.On("UpdatePassword", "gym123", "user123", mock.AnythingOfType("string")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "user not found",
			gymID:       "gym123",
			userID:      "nonexistent",
			newPassword: "newpassword123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "nonexistent").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			err := service.UpdatePassword(tc.gymID, tc.userID, tc.newPassword)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		userID    string
		mockSetup func(*MockUserRepository)
		wantErr   bool
	}{
		{
			name:   "successful deletion",
			gymID:  "gym123",
			userID: "user123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:    "user123",
					GymID: "gym123",
				}, nil)
				// Delete succeeds
				mockRepo.On("DeleteUser", "gym123", "user123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "nonexistent").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			err := service.DeleteUser(tc.gymID, tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetUserActive(t *testing.T) {
	testCases := []struct {
		name      string
		gymID     string
		userID    string
		active    bool
		mockSetup func(*MockUserRepository)
		wantErr   bool
	}{
		{
			name:   "successful activation",
			gymID:  "gym123",
			userID: "user123",
			active: true,
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists and is currently inactive
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					GymID:    "gym123",
					IsActive: false,
				}, nil)
				// Activation succeeds
				mockRepo.On("SetUserActive", "gym123", "user123", true).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "successful deactivation",
			gymID:  "gym123",
			userID: "user123",
			active: false,
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists and is currently active
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					GymID:    "gym123",
					IsActive: true,
				}, nil)
				// Deactivation succeeds
				mockRepo.On("SetUserActive", "gym123", "user123", false).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			gymID:  "gym123",
			userID: "nonexistent",
			active: true,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetUserByID", "gym123", "nonexistent").Return(dto.UserResponseDTO{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:   "user already active",
			gymID:  "gym123",
			userID: "user123",
			active: true,
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists and is already active
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					GymID:    "gym123",
					IsActive: true,
				}, nil)
			},
			wantErr: true,
		},
		{
			name:   "user already inactive",
			gymID:  "gym123",
			userID: "user123",
			active: false,
			mockSetup: func(mockRepo *MockUserRepository) {
				// User exists and is already inactive
				mockRepo.On("GetUserByID", "gym123", "user123").Return(dto.UserResponseDTO{
					ID:       "user123",
					GymID:    "gym123",
					IsActive: false,
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUsersService(mockRepo)
			tc.mockSetup(mockRepo)
			err := service.SetUserActive(tc.gymID, tc.userID, tc.active)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
