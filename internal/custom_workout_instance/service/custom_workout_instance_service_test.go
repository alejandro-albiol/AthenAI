package service_test

import (
	"database/sql"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	createErr               error
	getByIDErr              error
	getSummaryByIDErr       error
	getByUserIDErr          error
	getSummariesByUserIDErr error
	getLastsByUserIDErr     error
	listErr                 error
	listSummariesErr        error
	updateErr               error
	deleteErr               error
	instances               []*dto.ResponseCustomWorkoutInstanceDTO
	summaries               []*dto.SummaryCustomWorkoutInstanceDTO
	lastCreatedID           string
}

func (m *mockRepository) Create(gymID string, createdBy string, instance *dto.CreateCustomWorkoutInstanceDTO) (*string, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &m.lastCreatedID, nil
}

func (m *mockRepository) GetByID(gymID, id string) (*dto.ResponseCustomWorkoutInstanceDTO, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	for _, instance := range m.instances {
		if instance.ID == id {
			return instance, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) GetSummaryByID(gymID, id string) (*dto.SummaryCustomWorkoutInstanceDTO, error) {
	if m.getSummaryByIDErr != nil {
		return nil, m.getSummaryByIDErr
	}
	for _, summary := range m.summaries {
		if summary.ID == id {
			return summary, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) GetByUserID(gymID, userID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	if m.getByUserIDErr != nil {
		return nil, m.getByUserIDErr
	}
	var result []*dto.ResponseCustomWorkoutInstanceDTO
	// Note: Since CreatedBy is not in the DTO, we'll assume all instances belong to the user for testing
	result = append(result, m.instances...)
	return result, nil
}

func (m *mockRepository) GetSummariesByUserID(gymID, userID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	if m.getSummariesByUserIDErr != nil {
		return nil, m.getSummariesByUserIDErr
	}
	return m.summaries, nil
}

func (m *mockRepository) GetLastsByUserID(gymID, userID string, count int) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	if m.getLastsByUserIDErr != nil {
		return nil, m.getLastsByUserIDErr
	}
	var result []*dto.ResponseCustomWorkoutInstanceDTO
	for i, instance := range m.instances {
		if i < count {
			result = append(result, instance)
		}
	}
	return result, nil
}

func (m *mockRepository) List(gymID string) ([]*dto.ResponseCustomWorkoutInstanceDTO, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.instances, nil
}

func (m *mockRepository) ListSummaries(gymID string) ([]*dto.SummaryCustomWorkoutInstanceDTO, error) {
	if m.listSummariesErr != nil {
		return nil, m.listSummariesErr
	}
	return m.summaries, nil
}

func (m *mockRepository) Update(gymID string, id string, instance *dto.UpdateCustomWorkoutInstanceDTO) error {
	return m.updateErr
}

func (m *mockRepository) Delete(gymID, id string) error {
	return m.deleteErr
}

func TestCreateCustomWorkoutInstance_Success(t *testing.T) {
	mockRepo := &mockRepository{
		lastCreatedID: "instance123",
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	createDTO := &dto.CreateCustomWorkoutInstanceDTO{
		Name:           "Test Workout",
		Description:    "Test description",
		TemplateSource: "gym",
		GymTemplateID:  stringPtr("template123"),
	}

	id, err := service.CreateCustomWorkoutInstance("gym123", "user123", createDTO)

	assert.NoError(t, err)
	assert.Equal(t, "instance123", *id)
}

func TestCreateCustomWorkoutInstance_ValidationError_EmptyName(t *testing.T) {
	mockRepo := &mockRepository{}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	createDTO := &dto.CreateCustomWorkoutInstanceDTO{
		Name:           "", // Empty name should cause validation error
		TemplateSource: "gym",
	}

	id, err := service.CreateCustomWorkoutInstance("gym123", "user123", createDTO)

	assert.Error(t, err)
	assert.Nil(t, id)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeBadRequest, apiErr.Code)
}

func TestCreateCustomWorkoutInstance_ValidationError_MissingPublicTemplateID(t *testing.T) {
	mockRepo := &mockRepository{}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	createDTO := &dto.CreateCustomWorkoutInstanceDTO{
		Name:           "Test Workout",
		TemplateSource: "public", // public template source but no public_template_id
	}

	id, err := service.CreateCustomWorkoutInstance("gym123", "user123", createDTO)

	assert.Error(t, err)
	assert.Nil(t, id)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeBadRequest, apiErr.Code)
}

func TestCreateCustomWorkoutInstance_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		createErr: apierror.New(errorcode_enum.CodeInternal, "Database error", nil),
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	createDTO := &dto.CreateCustomWorkoutInstanceDTO{
		Name:           "Test Workout",
		TemplateSource: "gym",
		GymTemplateID:  stringPtr("template123"),
	}

	id, err := service.CreateCustomWorkoutInstance("gym123", "user123", createDTO)

	assert.Error(t, err)
	assert.Nil(t, id)
}

func TestGetCustomWorkoutInstanceByID_Success(t *testing.T) {
	instance := &dto.ResponseCustomWorkoutInstanceDTO{
		ID:                       "instance123",
		Name:                     "Test Workout",
		Description:              "Test description",
		TemplateSource:           "gym",
		GymTemplateID:            stringPtr("template456"),
		DifficultyLevel:          "intermediate",
		EstimatedDurationMinutes: 45,
		TotalExercises:           5,
	}

	mockRepo := &mockRepository{
		instances: []*dto.ResponseCustomWorkoutInstanceDTO{instance},
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	result, err := service.GetCustomWorkoutInstanceByID("gym123", "instance123")

	assert.NoError(t, err)
	assert.Equal(t, instance, result)
}

func TestGetCustomWorkoutInstanceByID_NotFound(t *testing.T) {
	mockRepo := &mockRepository{
		getByIDErr: sql.ErrNoRows,
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	result, err := service.GetCustomWorkoutInstanceByID("gym123", "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeNotFound, apiErr.Code)
}

func TestGetCustomWorkoutInstanceSummaryByID_Success(t *testing.T) {
	summary := &dto.SummaryCustomWorkoutInstanceDTO{
		ID:                       "instance123",
		Name:                     "Test Workout",
		Description:              "Test description",
		TemplateSource:           "gym",
		DifficultyLevel:          "intermediate",
		EstimatedDurationMinutes: 45,
		TotalExercises:           5,
	}

	mockRepo := &mockRepository{
		summaries: []*dto.SummaryCustomWorkoutInstanceDTO{summary},
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	result, err := service.GetCustomWorkoutInstanceSummaryByID("gym123", "instance123")

	assert.NoError(t, err)
	assert.Equal(t, summary, result)
}

func TestGetCustomWorkoutInstancesByUserID_Success(t *testing.T) {
	instances := []*dto.ResponseCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Workout 1",
			TemplateSource: "gym",
		},
		{
			ID:             "instance2",
			Name:           "Workout 2",
			TemplateSource: "public",
		},
	}

	mockRepo := &mockRepository{
		instances: instances,
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	result, err := service.GetCustomWorkoutInstancesByUserID("gym123", "user123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetLastCustomWorkoutInstancesByUserID_Success(t *testing.T) {
	instances := []*dto.ResponseCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Recent Workout",
			TemplateSource: "gym",
		},
	}

	mockRepo := &mockRepository{
		instances: instances,
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	result, err := service.GetLastCustomWorkoutInstancesByUserID("gym123", "user123", 5)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "instance1", result[0].ID)
}

func TestUpdateCustomWorkoutInstance_Success(t *testing.T) {
	mockRepo := &mockRepository{}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	updateDTO := &dto.UpdateCustomWorkoutInstanceDTO{
		Name:        stringPtr("Updated Workout"),
		Description: stringPtr("Updated description"),
	}

	err := service.UpdateCustomWorkoutInstance("gym123", "instance123", updateDTO)

	assert.NoError(t, err)
}

func TestUpdateCustomWorkoutInstance_ValidationError_MissingGymTemplateID(t *testing.T) {
	mockRepo := &mockRepository{}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	updateDTO := &dto.UpdateCustomWorkoutInstanceDTO{
		TemplateSource: stringPtr("gym"), // gym template source but no gym_template_id
	}

	err := service.UpdateCustomWorkoutInstance("gym123", "instance123", updateDTO)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, errorcode_enum.CodeBadRequest, apiErr.Code)
}

func TestDeleteCustomWorkoutInstance_Success(t *testing.T) {
	mockRepo := &mockRepository{}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	err := service.DeleteCustomWorkoutInstance("gym123", "instance123")

	assert.NoError(t, err)
}

func TestDeleteCustomWorkoutInstance_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		deleteErr: apierror.New(errorcode_enum.CodeInternal, "Database error", nil),
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	err := service.DeleteCustomWorkoutInstance("gym123", "instance123")

	assert.Error(t, err)
}

func TestListCustomWorkoutInstances_Success(t *testing.T) {
	instances := []*dto.ResponseCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Workout 1",
			TemplateSource: "gym",
		},
		{
			ID:             "instance2",
			Name:           "Workout 2",
			TemplateSource: "public",
		},
	}

	mockRepo := &mockRepository{
		instances: instances,
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	result, err := service.ListCustomWorkoutInstances("gym123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListCustomWorkoutInstanceSummaries_Success(t *testing.T) {
	summaries := []*dto.SummaryCustomWorkoutInstanceDTO{
		{
			ID:             "instance1",
			Name:           "Workout 1",
			TemplateSource: "gym",
		},
		{
			ID:             "instance2",
			Name:           "Workout 2",
			TemplateSource: "public",
		},
	}

	mockRepo := &mockRepository{
		summaries: summaries,
	}
	service := service.NewCustomWorkoutInstanceService(mockRepo)

	result, err := service.ListCustomWorkoutInstanceSummaries("gym123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
