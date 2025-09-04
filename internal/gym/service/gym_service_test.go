package service_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGymRepository struct {
	mock.Mock
}

func (m *MockGymRepository) CreateGym(gym *dto.GymCreationDTO) (*string, error) {
	args := m.Called(gym)
	str := args.String(0)
	return &str, args.Error(1)
}

func (m *MockGymRepository) GetGymByID(id string) (*dto.GymResponseDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymRepository) GetGymByName(name string) (*dto.GymResponseDTO, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymRepository) GetAllGyms() ([]*dto.GymResponseDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymRepository) UpdateGym(id string, gym *dto.GymUpdateDTO) (*dto.GymResponseDTO, error) {
	args := m.Called(id, gym)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GymResponseDTO), args.Error(1)
}

func (m *MockGymRepository) SetGymActive(id string, active bool) error {
	args := m.Called(id, active)
	return args.Error(0)
}

func (m *MockGymRepository) DeleteGym(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestMain(m *testing.M) {
	// Set test environment to skip database operations
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	// Run tests
	m.Run()
}

func TestCreateGym(t *testing.T) {
	gymDTO := &dto.GymCreationDTO{
		Name:    "Test Gym",
		Email:   "test@gym.com",
		Address: "123 Test St",
		Phone:   "+1234567890",
	}

	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(MockGymRepository)
		svc := service.NewGymService(mockRepo)

		mockRepo.On("GetGymByName", gymDTO.Name).Return(nil, sql.ErrNoRows)
		mockRepo.On("CreateGym", gymDTO).Return("gym123", nil)

		id, err := svc.CreateGym(gymDTO)
		assert.NoError(t, err)
		assert.Equal(t, "gym123", id)
	})

	t.Run("domain already exists", func(t *testing.T) {
		mockRepo := new(MockGymRepository)
		svc := service.NewGymService(mockRepo)

		existingGym := &dto.GymResponseDTO{ID: "gym123", Name: gymDTO.Name}
		mockRepo.On("GetGymByName", gymDTO.Name).Return(existingGym, nil)

		_, err := svc.CreateGym(gymDTO)
		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, errorcode_enum.CodeConflict, apiErr.Code)
	})
}

func TestGetGym(t *testing.T) {
	mockRepo := new(MockGymRepository)
	svc := service.NewGymService(mockRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedGym := &dto.GymResponseDTO{
			ID:       "gym123",
			Name:     "Test Gym",
			Email:    "test@gym.com",
			IsActive: true,
		}
		mockRepo.On("GetGymByID", "gym123").Return(expectedGym, nil)

		gym, err := svc.GetGymByID("gym123")
		assert.NoError(t, err)
		assert.Equal(t, expectedGym, gym)
	})

	t.Run("gym not found", func(t *testing.T) {
		mockRepo.On("GetGymByID", "nonexistent").Return(nil, sql.ErrNoRows)

		_, err := svc.GetGymByID("nonexistent")
		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
	})
}

func TestGetAllGyms(t *testing.T) {
	mockRepo := new(MockGymRepository)
	svc := service.NewGymService(mockRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedGyms := []*dto.GymResponseDTO{
			{ID: "gym123", Name: "Test Gym 1"},
			{ID: "gym456", Name: "Test Gym 2"},
		}
		mockRepo.On("GetAllGyms").Return(expectedGyms, nil)

		gyms, err := svc.GetAllGyms()
		assert.NoError(t, err)
		assert.Equal(t, expectedGyms, gyms)
	})
}

func TestUpdateGym(t *testing.T) {
	mockRepo := new(MockGymRepository)
	svc := service.NewGymService(mockRepo)

	// Use local variables for pointer fields
	name := "Updated Gym"
	email := "updated@gym.com"
	updateDTO := &dto.GymUpdateDTO{
		Name:  &name,
		Email: &email,
	}

	t.Run("successful update", func(t *testing.T) {
		mockRepo.On("GetGymByID", "gym123").Return(&dto.GymResponseDTO{ID: "gym123"}, nil)
		mockRepo.On("UpdateGym", "gym123", updateDTO).Return(&dto.GymResponseDTO{ID: "gym123"}, nil)

		updatedGym, err := svc.UpdateGym("gym123", updateDTO)
		assert.NoError(t, err)
		assert.Equal(t, "gym123", updatedGym.ID)
	})

	t.Run("gym not found", func(t *testing.T) {
		mockRepo.On("GetGymByID", "nonexistent").Return(nil, sql.ErrNoRows)
		mockRepo.On("UpdateGym", "nonexistent", updateDTO).Return(nil, sql.ErrNoRows)

		_, err := svc.UpdateGym("nonexistent", updateDTO)
		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
	})
}

func TestDeleteGym(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockRepo := new(MockGymRepository)
		svc := service.NewGymService(mockRepo)

		mockRepo.On("GetGymByID", "gym123").Return(&dto.GymResponseDTO{ID: "gym123"}, nil)
		mockRepo.On("DeleteGym", "gym123").Return(nil)

		err := svc.DeleteGym("gym123")
		assert.NoError(t, err)
	})

	t.Run("gym not found", func(t *testing.T) {
		mockRepo := new(MockGymRepository)
		svc := service.NewGymService(mockRepo)

		mockRepo.On("GetGymByID", "nonexistent").Return(nil, sql.ErrNoRows)

		err := svc.DeleteGym("nonexistent")
		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
	})
}

func TestSetGymActive(t *testing.T) {
	mockRepo := new(MockGymRepository)
	svc := service.NewGymService(mockRepo)

	t.Run("successful activation", func(t *testing.T) {
		mockRepo.On("GetGymByID", "gym123").Return(&dto.GymResponseDTO{ID: "gym123"}, nil)
		mockRepo.On("SetGymActive", "gym123", true).Return(nil)

		err := svc.SetGymActive("gym123", true)
		assert.NoError(t, err)
	})

	t.Run("gym not found", func(t *testing.T) {
		mockRepo.On("GetGymByID", "nonexistent").Return(nil, sql.ErrNoRows)

		err := svc.SetGymActive("nonexistent", true)
		assert.Error(t, err)
		var apiErr *apierror.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
	})
}
