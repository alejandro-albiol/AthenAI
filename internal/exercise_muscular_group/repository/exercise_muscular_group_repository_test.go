package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/exercise_muscular_group/dto"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *ExerciseMuscularGroupRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	repo := &ExerciseMuscularGroupRepository{db: db}
	return db, mock, repo
}

func TestExerciseMuscularGroupRepository_FindByExerciseID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Success case: multiple rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex1").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "muscular_group_id"}).
			AddRow("ex1", "mg1").
			AddRow("ex1", "mg2"))

	links, err := repo.FindByExerciseID("ex1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}

	// Not found case: no rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex2").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "muscular_group_id"}))

	links, err = repo.FindByExerciseID("ex2")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected 0 links, got %d", len(links))
	}

	// DB error case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex3").
		WillReturnError(errors.New("db error"))

	links, err = repo.FindByExerciseID("ex3")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if links != nil {
		t.Errorf("expected nil links, got %+v", links)
	}
}

func TestExerciseMuscularGroupRepository_FindByMuscularGroupID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Success case: multiple rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE muscular_group_id = $1`)).
		WithArgs("mg1").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "muscular_group_id"}).
			AddRow("ex1", "mg1").
			AddRow("ex2", "mg1"))

	links, err := repo.FindByMuscularGroupID("mg1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}

	// Not found case: no rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE muscular_group_id = $1`)).
		WithArgs("mg2").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "muscular_group_id"}))

	links, err = repo.FindByMuscularGroupID("mg2")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected 0 links, got %d", len(links))
	}

	// DB error case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE muscular_group_id = $1`)).
		WithArgs("mg3").
		WillReturnError(errors.New("db error"))

	links, err = repo.FindByMuscularGroupID("mg3")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if links != nil {
		t.Errorf("expected nil links, got %+v", links)
	}
}

func TestExerciseMuscularGroupRepository_CreateLink(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	link := &dto.ExerciseMuscularGroup{ExerciseID: "ex1", MuscularGroupID: "mg1"}

	// Success case
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO exercise_muscular_group (exercise_id, muscular_group_id) VALUES ($1, $2) RETURNING exercise_id`)).
		WithArgs(link.ExerciseID, link.MuscularGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id"}).AddRow("ex1"))

	res, err := repo.CreateLink(link)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if res == nil || *res != "ex1" {
		t.Errorf("expected id 'ex1', got %v", res)
	}

	// Error case
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO exercise_muscular_group (exercise_id, muscular_group_id) VALUES ($1, $2) RETURNING exercise_id`)).
		WithArgs(link.ExerciseID, link.MuscularGroupID).
		WillReturnError(errors.New("insert error"))

	res, err = repo.CreateLink(link)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if res != nil {
		t.Errorf("expected nil id, got %v", res)
	}
}

func TestExerciseMuscularGroupRepository_DeleteLink(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Success case
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteLink("ex1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Error case
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex1").
		WillReturnError(errors.New("delete error"))

	err = repo.DeleteLink("ex1")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestExerciseMuscularGroupRepository_FindByID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Success case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex1").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "muscular_group_id"}).AddRow("ex1", "mg1"))

	link, err := repo.FindByID("ex1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if link == nil || link.ExerciseID != "ex1" || link.MuscularGroupID != "mg1" {
		t.Errorf("unexpected link result: %+v", link)
	}

	// Not found case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex2").
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id", "muscular_group_id"}))

	link, err = repo.FindByID("ex2")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if link != nil {
		t.Errorf("expected nil link, got %+v", link)
	}

	// DB error case
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT exercise_id, muscular_group_id FROM exercise_muscular_group WHERE exercise_id = $1`)).
		WithArgs("ex3").
		WillReturnError(errors.New("db error"))

	link, err = repo.FindByID("ex3")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if link != nil {
		t.Errorf("expected nil link, got %+v", link)
	}
}
