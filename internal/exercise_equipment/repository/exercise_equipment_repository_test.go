package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/exercise_equipment/dto"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *ExerciseEquipmentRepositoryImpl) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	repo := &ExerciseEquipmentRepositoryImpl{db: db}
	return db, mock, repo
}

func TestExerciseEquipmentRepository_CreateLink(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	link := &dto.ExerciseEquipment{ExerciseID: "ex1", EquipmentID: "eq1"}

	// Success case
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO exercise_equipment (exercise_id, equipment_id) VALUES ($1, $2) RETURNING id`)).
		WithArgs(link.ExerciseID, link.EquipmentID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("id1"))

	res, err := repo.CreateLink(link)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res == nil || *res != "id1" {
		t.Errorf("expected id 'id1', got %v", res)
	}

	// Error case
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO exercise_equipment (exercise_id, equipment_id) VALUES ($1, $2) RETURNING id`)).
		WithArgs(link.ExerciseID, link.EquipmentID).
		WillReturnError(errors.New("insert error"))

	res, err = repo.CreateLink(link)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if res != nil {
		t.Errorf("expected nil id, got %v", res)
	}
}

func TestExerciseEquipmentRepository_DeleteLink(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Success case
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM exercise_equipment WHERE id = $1`)).
		WithArgs("id1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteLink("id1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Error case
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM exercise_equipment WHERE id = $1`)).
		WithArgs("id1").
		WillReturnError(errors.New("delete error"))

	err = repo.DeleteLink("id1")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestExerciseEquipmentRepository_FindByID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Success case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, equipment_id FROM exercise_equipment WHERE id = $1`)).
		WithArgs("id1").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "equipment_id"}).AddRow("ex1", "eq1"))

	link, err := repo.FindByID("id1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if link == nil || link.ExerciseID != "ex1" || link.EquipmentID != "eq1" {
		t.Errorf("unexpected link result: %+v", link)
	}

	// Not found case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, equipment_id FROM exercise_equipment WHERE id = $1`)).
		WithArgs("id2").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "equipment_id"}))

	link, err = repo.FindByID("id2")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if link != nil {
		t.Errorf("expected nil link, got %+v", link)
	}

	// DB error case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, equipment_id FROM exercise_equipment WHERE id = $1`)).
		WithArgs("id3").
		WillReturnError(errors.New("db error"))

	link, err = repo.FindByID("id3")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if link != nil {
		t.Errorf("expected nil link, got %+v", link)
	}
}
