package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/workout_template/dto"
	"github.com/alejandro-albiol/athenai/internal/workout_template/repository"
	"github.com/stretchr/testify/assert"
)

func getTestResponseWorkoutTemplate() *dto.ResponseWorkoutTemplateDTO {
	desc := "A test template"
	diff := "easy"
	dur := 30
	target := "beginner"
	return &dto.ResponseWorkoutTemplateDTO{
		ID:                       "1",
		Name:                     "Test Template",
		Description:              desc,
		DifficultyLevel:          diff,
		EstimatedDurationMinutes: dur,
		TargetAudience:           target,
		CreatedBy:                "admin",
		IsActive:                 true,
		IsPublic:                 false,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}
}

func TestCreateWorkoutTemplate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewWorkoutTemplateRepository(db)
	desc := "A test template"
	diff := "easy"
	dur := 30
	target := "beginner"
	createDTO := &dto.CreateWorkoutTemplateDTO{
		Name:                     "Test Template",
		Description:              desc,
		DifficultyLevel:          diff,
		EstimatedDurationMinutes: dur,
		TargetAudience:           target,
		CreatedBy:                "admin",
	}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO public.workout_template (name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)).
		WithArgs(createDTO.Name, createDTO.Description, createDTO.DifficultyLevel, createDTO.EstimatedDurationMinutes, createDTO.TargetAudience, createDTO.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	id, err := repo.CreateWorkoutTemplate(createDTO)
	assert.NoError(t, err)
	if assert.NotNil(t, id) {
		assert.Equal(t, "1", *id)
	}
}

func TestGetWorkoutTemplateByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewWorkoutTemplateRepository(db)
	resp := getTestResponseWorkoutTemplate()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE id = $1`)).
		WithArgs(resp.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "difficulty_level", "estimated_duration_minutes", "target_audience", "created_by", "is_active", "is_public", "created_at", "updated_at"}).
			AddRow(resp.ID, resp.Name, resp.Description, resp.DifficultyLevel, resp.EstimatedDurationMinutes, resp.TargetAudience, resp.CreatedBy, resp.IsActive, resp.IsPublic, resp.CreatedAt, resp.UpdatedAt))

	result, err := repo.GetWorkoutTemplateByID(resp.ID)
	assert.NoError(t, err)
	assert.Equal(t, resp.ID, result.ID)
	assert.Equal(t, resp.Name, result.Name)
}

func TestGetWorkoutTemplateByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewWorkoutTemplateRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE id = $1`)).
		WithArgs("999").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "difficulty_level", "estimated_duration_minutes", "target_audience", "created_by", "is_active", "is_public", "created_at", "updated_at"}))

	result, err := repo.GetWorkoutTemplateByID("999")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateWorkoutTemplate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewWorkoutTemplateRepository(db)
	desc := "Updated Description"
	diff := "medium"
	dur := 45
	target := "intermediate"
	name := "Updated Name"
	isActive := true
	isPublic := true
	updateDTO := &dto.UpdateWorkoutTemplateDTO{
		Name:                     &name,
		Description:              &desc,
		DifficultyLevel:          &diff,
		EstimatedDurationMinutes: &dur,
		TargetAudience:           &target,
		IsActive:                 &isActive,
		IsPublic:                 &isPublic,
	}
	resp := getTestResponseWorkoutTemplate()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE public.workout_template SET name = $1, description = $2, difficulty_level = $3, estimated_duration_minutes = $4, target_audience = $5, is_active = $6, is_public = $7, updated_at = NOW() WHERE id = $8 RETURNING id`)).
		WithArgs(updateDTO.Name, updateDTO.Description, updateDTO.DifficultyLevel, updateDTO.EstimatedDurationMinutes, updateDTO.TargetAudience, updateDTO.IsActive, updateDTO.IsPublic, resp.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE id = $1`)).
		WithArgs(resp.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "difficulty_level", "estimated_duration_minutes", "target_audience", "created_by", "is_active", "is_public", "created_at", "updated_at"}).
			AddRow(resp.ID, resp.Name, resp.Description, resp.DifficultyLevel, resp.EstimatedDurationMinutes, resp.TargetAudience, resp.CreatedBy, resp.IsActive, resp.IsPublic, resp.CreatedAt, resp.UpdatedAt))

	result, err := repo.UpdateWorkoutTemplate(resp.ID, updateDTO)
	assert.NoError(t, err)
	assert.Equal(t, resp.ID, result.ID)
}

func TestDeleteWorkoutTemplate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewWorkoutTemplateRepository(db)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM public.workout_template WHERE id = $1`)).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteWorkoutTemplate("1")
	assert.NoError(t, err)
}

func TestGetWorkoutTemplatesByDifficulty_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewWorkoutTemplateRepository(db)
	resp := getTestResponseWorkoutTemplate()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, difficulty_level, estimated_duration_minutes, target_audience, created_by, is_active, is_public, created_at, updated_at FROM public.workout_template WHERE difficulty_level = $1`)).
		WithArgs(resp.DifficultyLevel).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "difficulty_level", "estimated_duration_minutes", "target_audience", "created_by", "is_active", "is_public", "created_at", "updated_at"}).
			AddRow(resp.ID, resp.Name, resp.Description, resp.DifficultyLevel, resp.EstimatedDurationMinutes, resp.TargetAudience, resp.CreatedBy, resp.IsActive, resp.IsPublic, resp.CreatedAt, resp.UpdatedAt))

	results, err := repo.GetWorkoutTemplatesByDifficulty(resp.DifficultyLevel)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, resp.ID, results[0].ID)
}
