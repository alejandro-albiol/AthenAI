package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/template_block/repository"
)

func TestCreateTemplateBlock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	repo := repository.NewTemplateBlockRepository(mockDB)

	instructions := "Do this"
	estimatedDurationMinutes := 30
	reps := 12
	series := 3
	restTimeSeconds := 60
	block := &dto.CreateTemplateBlockDTO{
		TemplateID:               "template-uuid",
		BlockName:                "Block Name",
		BlockType:                "Type",
		BlockOrder:               1,
		ExerciseCount:            5,
		EstimatedDurationMinutes: &estimatedDurationMinutes,
		Instructions:             &instructions,
		Reps:                     &reps,
		Series:                   &series,
		RestTimeSeconds:          &restTimeSeconds,
		CreatedBy:                "admin-uuid",
	}

	mock.ExpectQuery("INSERT INTO public.template_block").
		WithArgs(block.TemplateID, block.BlockName, block.BlockType, block.BlockOrder, block.ExerciseCount, block.EstimatedDurationMinutes, block.Instructions, block.Reps, block.Series, block.RestTimeSeconds, block.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("block-uuid"))

	id, err := repo.CreateTemplateBlock(block)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if id == nil || *id != "block-uuid" {
		t.Errorf("expected id 'block-uuid', got '%v'", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	repo := repository.NewTemplateBlockRepository(mockDB)

	row := sqlmock.NewRows([]string{"id", "template_id", "block_name", "block_type", "block_order", "exercise_count", "estimated_duration_minutes", "instructions", "reps", "series", "rest_time_seconds", "created_at", "created_by"}).
		AddRow("block-uuid", "template-uuid", "Block Name", "Type", 1, 5, 30, "Do this", 12, 3, 60, "2025-09-04T12:00:00Z", "admin-uuid")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, reps, series, rest_time_seconds, created_at, created_by\n\t\tFROM public.template_block WHERE id = $1")).
		WithArgs("block-uuid").
		WillReturnRows(row)

	block, err := repo.GetTemplateBlockByID("block-uuid")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if block == nil || block.ID != "block-uuid" {
		t.Errorf("expected block ID 'block-uuid', got '%v'", block)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestGetByTemplateID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	repo := repository.NewTemplateBlockRepository(mockDB)

	rows := sqlmock.NewRows([]string{"id", "template_id", "block_name", "block_type", "block_order", "exercise_count", "estimated_duration_minutes", "instructions", "reps", "series", "rest_time_seconds", "created_at", "created_by"}).
		AddRow("block-uuid-1", "template-uuid", "Block 1", "Type", 1, 5, 30, "Do this", 12, 3, 60, "2025-09-04T12:00:00Z", "admin-uuid").
		AddRow("block-uuid-2", "template-uuid", "Block 2", "Type", 2, 6, 40, "Do that", 15, 4, 90, "2025-09-04T12:01:00Z", "admin-uuid")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, reps, series, rest_time_seconds, created_at, created_by\n\t\tFROM public.template_block WHERE template_id = $1 ORDER BY block_order")).
		WithArgs("template-uuid").
		WillReturnRows(rows)

	blocks, err := repo.GetTemplateBlocksByTemplateID("template-uuid")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(blocks))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestUpdateTemplateBlock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	repo := repository.NewTemplateBlockRepository(mockDB)

	name := "Updated Name"
	typ := "Updated Type"
	order := 2
	exerciseCount := 10
	duration := 45
	instructions := "Updated instructions"
	reps := 15
	series := 4
	restTimeSeconds := 120

	update := &dto.UpdateTemplateBlockDTO{
		BlockName:                &name,
		BlockType:                &typ,
		BlockOrder:               &order,
		ExerciseCount:            &exerciseCount,
		EstimatedDurationMinutes: &duration,
		Instructions:             &instructions,
		Reps:                     &reps,
		Series:                   &series,
		RestTimeSeconds:          &restTimeSeconds,
	}

	row := sqlmock.NewRows([]string{"id", "template_id", "block_name", "block_type", "block_order", "exercise_count", "estimated_duration_minutes", "instructions", "reps", "series", "rest_time_seconds", "created_at", "created_by"}).
		AddRow("block-uuid", "template-uuid", update.BlockName, update.BlockType, update.BlockOrder, update.ExerciseCount, update.EstimatedDurationMinutes, update.Instructions, update.Reps, update.Series, update.RestTimeSeconds, "2025-09-04T12:02:00Z", "admin-uuid")
	mock.ExpectQuery(regexp.QuoteMeta("UPDATE public.template_block\n\t\tSET block_name = $1, block_type = $2, block_order = $3, exercise_count = $4, estimated_duration_minutes = $5, instructions = $6, reps = $7, series = $8, rest_time_seconds = $9\n\t\tWHERE id = $10\n\t\tRETURNING id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, reps, series, rest_time_seconds, created_at, created_by")).
		WithArgs(update.BlockName, update.BlockType, update.BlockOrder, update.ExerciseCount, update.EstimatedDurationMinutes, update.Instructions, update.Reps, update.Series, update.RestTimeSeconds, "block-uuid").
		WillReturnRows(row)

	updatedBlock, err := repo.UpdateTemplateBlock("block-uuid", update)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updatedBlock == nil || updatedBlock.BlockName != "Updated Name" {
		t.Errorf("expected updated block name 'Updated Name', got '%v'", updatedBlock)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestDeleteTemplateBlock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	repo := repository.NewTemplateBlockRepository(mockDB)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM public.template_block WHERE id = $1")).
		WithArgs("block-uuid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteTemplateBlock("block-uuid")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
