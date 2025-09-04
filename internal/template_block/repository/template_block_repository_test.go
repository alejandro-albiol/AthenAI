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

	block := &dto.CreateTemplateBlockDTO{
		TemplateID:               "template-uuid",
		Name:                     "Block Name",
		Type:                     "Type",
		Order:                    1,
		ExerciseCount:            5,
		EstimatedDurationMinutes: 30,
		Instructions:             "Do this",
	}

	mock.ExpectQuery("INSERT INTO public.template_block").
		WithArgs(block.TemplateID, block.Name, block.Type, block.Order, block.ExerciseCount, block.EstimatedDurationMinutes, block.Instructions).
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

	row := sqlmock.NewRows([]string{"id", "template_id", "name", "type", "order", "exercise_count", "estimated_duration_minutes", "instructions", "created_at"}).
		AddRow("block-uuid", "template-uuid", "Block Name", "Type", 1, 5, 30, "Do this", "2025-09-04T12:00:00Z")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, template_id, name, type, \"order\", exercise_count, estimated_duration_minutes, instructions, created_at FROM public.template_block WHERE id = $1")).
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

	rows := sqlmock.NewRows([]string{"id", "template_id", "name", "type", "order", "exercise_count", "estimated_duration_minutes", "instructions", "created_at"}).
		AddRow("block-uuid-1", "template-uuid", "Block 1", "Type", 1, 5, 30, "Do this", "2025-09-04T12:00:00Z").
		AddRow("block-uuid-2", "template-uuid", "Block 2", "Type", 2, 6, 40, "Do that", "2025-09-04T12:01:00Z")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, template_id, name, type, \"order\", exercise_count, estimated_duration_minutes, instructions, created_at FROM public.template_block WHERE template_id = $1 ORDER BY \"order\"")).
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

	update := &dto.UpdateTemplateBlockDTO{
		Name:  &name,
		Type:  &typ,
		Order: &order,
		ExerciseCount: &exerciseCount,
		EstimatedDurationMinutes: &duration,
		Instructions: &instructions,
	}

	row := sqlmock.NewRows([]string{"id", "template_id", "name", "type", "order", "exercise_count", "estimated_duration_minutes", "instructions", "created_at"}).
		AddRow("block-uuid", "template-uuid", update.Name, update.Type, update.Order, update.ExerciseCount, update.EstimatedDurationMinutes, update.Instructions, "2025-09-04T12:02:00Z")
	mock.ExpectQuery(regexp.QuoteMeta("UPDATE public.template_block SET name = $1, type = $2, \"order\" = $3, exercise_count = $4, estimated_duration_minutes = $5, instructions = $6 WHERE id = $7 RETURNING id, template_id, name, type, \"order\", exercise_count, estimated_duration_minutes, instructions, created_at")).
		WithArgs(update.Name, update.Type, update.Order, update.ExerciseCount, update.EstimatedDurationMinutes, update.Instructions, "block-uuid").
		WillReturnRows(row)

	updatedBlock, err := repo.UpdateTemplateBlock("block-uuid", update)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updatedBlock == nil || updatedBlock.Name != "Updated Name" {
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
