package service

import (
	"database/sql"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomTemplateBlockRepository is a mock for testing
type MockCustomTemplateBlockRepository struct {
	mock.Mock
}

func (m *MockCustomTemplateBlockRepository) CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (*string, error) {
	args := m.Called(gymID, block)
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockCustomTemplateBlockRepository) GetCustomTemplateBlockByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error) {
	args := m.Called(gymID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseCustomTemplateBlockDTO), args.Error(1)
}

func (m *MockCustomTemplateBlockRepository) ListCustomTemplateBlocksByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	args := m.Called(gymID, templateID)
	return args.Get(0).([]*dto.ResponseCustomTemplateBlockDTO), args.Error(1)
}

func (m *MockCustomTemplateBlockRepository) ListCustomTemplateBlocks(gymID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	args := m.Called(gymID)
	return args.Get(0).([]*dto.ResponseCustomTemplateBlockDTO), args.Error(1)
}

func (m *MockCustomTemplateBlockRepository) UpdateCustomTemplateBlock(gymID, id string, update *dto.UpdateCustomTemplateBlockDTO) error {
	args := m.Called(gymID, id, update)
	return args.Error(0)
}

func (m *MockCustomTemplateBlockRepository) DeleteCustomTemplateBlock(gymID, id string) error {
	args := m.Called(gymID, id)
	return args.Error(0)
}

func TestCustomTemplateBlockService_CreateCustomTemplateBlock(t *testing.T) {
	mockRepo := new(MockCustomTemplateBlockRepository)
	service := NewCustomTemplateBlockService(mockRepo)

	gymID := "gym123"
	block := &dto.CreateCustomTemplateBlockDTO{
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
	}

	expectedID := "block123"
	mockRepo.On("CreateCustomTemplateBlock", gymID, block).Return(&expectedID, nil)

	id, err := service.CreateCustomTemplateBlock(gymID, block)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, expectedID, *id)
	mockRepo.AssertExpectations(t)
}

func TestCustomTemplateBlockService_CreateCustomTemplateBlock_NilPayload(t *testing.T) {
	mockRepo := new(MockCustomTemplateBlockRepository)
	service := NewCustomTemplateBlockService(mockRepo)

	gymID := "gym123"

	id, err := service.CreateCustomTemplateBlock(gymID, nil)
	assert.Nil(t, id)
	assert.Error(t, err)

	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, "BAD_REQUEST", apiErr.Code)
	assert.Contains(t, apiErr.Message, "Block payload is nil")

	// Should not call repository if payload is nil
	mockRepo.AssertNotCalled(t, "CreateCustomTemplateBlock")
}

func TestCustomTemplateBlockService_GetCustomTemplateBlockByID(t *testing.T) {
	mockRepo := new(MockCustomTemplateBlockRepository)
	service := NewCustomTemplateBlockService(mockRepo)

	gymID := "gym123"
	blockID := "block123"
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

	mockRepo.On("GetCustomTemplateBlockByID", gymID, blockID).Return(expectedBlock, nil)

	result, err := service.GetCustomTemplateBlockByID(gymID, blockID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBlock, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomTemplateBlockService_GetCustomTemplateBlockByID_NotFound(t *testing.T) {
	mockRepo := new(MockCustomTemplateBlockRepository)
	service := NewCustomTemplateBlockService(mockRepo)

	gymID := "gym123"
	blockID := "nonexistent"

	mockRepo.On("GetCustomTemplateBlockByID", gymID, blockID).Return(nil, sql.ErrNoRows)

	result, err := service.GetCustomTemplateBlockByID(gymID, blockID)
	assert.Nil(t, result)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Contains(t, apiErr.Message, "not found")
	mockRepo.AssertExpectations(t)
}

func TestCustomTemplateBlockService_ListCustomTemplateBlocksByTemplateID(t *testing.T) {
	mockRepo := new(MockCustomTemplateBlockRepository)
	service := NewCustomTemplateBlockService(mockRepo)

	gymID := "gym123"
	templateID := "template123"
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
		{
			ID:                       "block124",
			TemplateID:               "template123",
			BlockName:                "Main Set",
			BlockType:                "main",
			BlockOrder:               2,
			ExerciseCount:            5,
			EstimatedDurationMinutes: intPtr(20),
			Instructions:             "Focus on strength",
			Reps:                     intPtr(8),
			Series:                   intPtr(4),
			RestTimeSeconds:          intPtr(90),
			IsActive:                 true,
			CreatedBy:                "user123",
		},
	}

	mockRepo.On("ListCustomTemplateBlocksByTemplateID", gymID, templateID).Return(expectedBlocks, nil)

	result, err := service.ListCustomTemplateBlocksByTemplateID(gymID, templateID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBlocks, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomTemplateBlockService_UpdateCustomTemplateBlock(t *testing.T) {
	mockRepo := new(MockCustomTemplateBlockRepository)
	service := NewCustomTemplateBlockService(mockRepo)

	gymID := "gym123"
	blockID := "block123"
	update := &dto.UpdateCustomTemplateBlockDTO{
		BlockName:       stringPtr("Updated Warm-up"),
		ExerciseCount:   intPtr(5),
		Reps:            intPtr(20),
		Series:          intPtr(4),
		RestTimeSeconds: intPtr(45),
	}

	expectedBlock := &dto.ResponseCustomTemplateBlockDTO{
		ID:                       "block123",
		TemplateID:               "template123",
		BlockName:                "Updated Warm-up",
		BlockType:                "warmup",
		BlockOrder:               1,
		ExerciseCount:            5,
		EstimatedDurationMinutes: intPtr(10),
		Instructions:             "Start with light exercises",
		Reps:                     intPtr(20),
		Series:                   intPtr(4),
		RestTimeSeconds:          intPtr(45),
		IsActive:                 true,
		CreatedBy:                "user123",
	}

	mockRepo.On("UpdateCustomTemplateBlock", gymID, blockID, update).Return(nil)
	mockRepo.On("GetCustomTemplateBlockByID", gymID, blockID).Return(expectedBlock, nil)

	result, err := service.UpdateCustomTemplateBlock(gymID, blockID, update)
	assert.NoError(t, err)
	assert.Equal(t, expectedBlock, result)
	mockRepo.AssertExpectations(t)
}

func TestCustomTemplateBlockService_DeleteCustomTemplateBlock(t *testing.T) {
	mockRepo := new(MockCustomTemplateBlockRepository)
	service := NewCustomTemplateBlockService(mockRepo)

	gymID := "gym123"
	blockID := "block123"

	mockRepo.On("DeleteCustomTemplateBlock", gymID, blockID).Return(nil)

	err := service.DeleteCustomTemplateBlock(gymID, blockID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
