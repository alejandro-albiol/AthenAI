package repository

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
	"github.com/stretchr/testify/assert"
)

func setupMockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	return db, mock, err
}

func TestCreateCustomTemplateBlock(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCustomTemplateBlockRepository(db)

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

	expectedQuery := `INSERT INTO "gym123"\.custom_template_block \(template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, reps, series, rest_time_seconds, created_by\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11\) RETURNING id`
	mock.ExpectQuery(expectedQuery).
		WithArgs(block.TemplateID, block.BlockName, block.BlockType, block.BlockOrder, block.ExerciseCount, block.EstimatedDurationMinutes, block.Instructions, block.Reps, block.Series, block.RestTimeSeconds, block.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("block123"))

	id, err := repo.CreateCustomTemplateBlock(gymID, block)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "block123", *id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCustomTemplateBlockRepository(db)

	gymID := "gym123"
	blockID := "block123"
	expectedQuery := `SELECT id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, reps, series, rest_time_seconds, created_at, updated_at, deleted_at, is_active, created_by FROM "gym123"\.custom_template_block WHERE id = \$1`

	now := time.Now()
	mock.ExpectQuery(expectedQuery).
		WithArgs(blockID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "template_id", "block_name", "block_type", "block_order", "exercise_count", "estimated_duration_minutes", "instructions", "reps", "series", "rest_time_seconds", "created_at", "updated_at", "deleted_at", "is_active", "created_by"}).
			AddRow("block123", "template123", "Warm-up", "warmup", 1, 3, 10, "Start with light exercises", 15, 3, 60, now, now, nil, true, "user123"))

	result, err := repo.GetCustomTemplateBlockByID(gymID, blockID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "block123", result.ID)
	assert.Equal(t, "template123", result.TemplateID)
	assert.Equal(t, "Warm-up", result.BlockName)
	assert.Equal(t, "warmup", result.BlockType)
	assert.Equal(t, 1, result.BlockOrder)
	assert.Equal(t, 3, result.ExerciseCount)
	assert.Equal(t, 10, *result.EstimatedDurationMinutes)
	assert.Equal(t, "Start with light exercises", result.Instructions)
	assert.Equal(t, 15, *result.Reps)
	assert.Equal(t, 3, *result.Series)
	assert.Equal(t, 60, *result.RestTimeSeconds)
	assert.Equal(t, "user123", result.CreatedBy)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByTemplateID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCustomTemplateBlockRepository(db)

	gymID := "gym123"
	templateID := "template123"
	expectedQuery := `SELECT id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, reps, series, rest_time_seconds, created_at, updated_at, deleted_at, is_active, created_by FROM "gym123"\.custom_template_block WHERE template_id = \$1 AND deleted_at IS NULL ORDER BY block_order ASC`

	now := time.Now()
	mock.ExpectQuery(expectedQuery).
		WithArgs(templateID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "template_id", "block_name", "block_type", "block_order", "exercise_count", "estimated_duration_minutes", "instructions", "reps", "series", "rest_time_seconds", "created_at", "updated_at", "deleted_at", "is_active", "created_by"}).
			AddRow("block123", "template123", "Warm-up", "warmup", 1, 3, 10, "Start with light exercises", 15, 3, 60, now, now, nil, true, "user123").
			AddRow("block124", "template123", "Main Set", "main", 2, 5, 20, "Focus on strength", 8, 4, 90, now, now, nil, true, "user123"))

	result, err := repo.ListCustomTemplateBlocksByTemplateID(gymID, templateID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "block123", result[0].ID)
	assert.Equal(t, "block124", result[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCustomTemplateBlock(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCustomTemplateBlockRepository(db)

	gymID := "gym123"
	blockID := "block123"
	update := &dto.UpdateCustomTemplateBlockDTO{
		BlockName:                stringPtr("Updated Warm-up"),
		BlockType:                stringPtr("warmup"),
		BlockOrder:               intPtr(1),
		ExerciseCount:            intPtr(5),
		EstimatedDurationMinutes: intPtr(15),
		Instructions:             stringPtr("Updated instructions"),
		Reps:                     intPtr(20),
		Series:                   intPtr(4),
		RestTimeSeconds:          intPtr(45),
		IsActive:                 boolPtr(true),
	}

	expectedQuery := `UPDATE "gym123"\.custom_template_block SET block_name = \$1, block_type = \$2, block_order = \$3, exercise_count = \$4, estimated_duration_minutes = \$5, instructions = \$6, reps = \$7, series = \$8, rest_time_seconds = \$9, is_active = \$10, updated_at = CURRENT_TIMESTAMP WHERE id = \$11 AND deleted_at IS NULL`
	mock.ExpectExec(expectedQuery).
		WithArgs(update.BlockName, update.BlockType, update.BlockOrder, update.ExerciseCount, update.EstimatedDurationMinutes, update.Instructions, update.Reps, update.Series, update.RestTimeSeconds, update.IsActive, blockID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateCustomTemplateBlock(gymID, blockID, update)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCustomTemplateBlock(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCustomTemplateBlockRepository(db)

	gymID := "gym123"
	blockID := "block123"
	expectedQuery := `UPDATE "gym123"\.custom_template_block SET deleted_at = CURRENT_TIMESTAMP WHERE id = \$1 AND deleted_at IS NULL`
	mock.ExpectExec(expectedQuery).
		WithArgs(blockID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteCustomTemplateBlock(gymID, blockID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
