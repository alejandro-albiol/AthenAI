package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_member_workout/dto"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	cleanup := func() { db.Close() }
	return db, mock, cleanup
}

func TestCreateCustomMemberWorkout(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewCustomMemberWorkoutRepository(db)

	mock.ExpectQuery(`INSERT INTO ".*".custom_member_workout`).
		WithArgs("member-id", "member-id", "workout-instance-id", "2025-09-08", nil, nil, "scheduled").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("new-id"))

	input := &dto.CreateCustomMemberWorkoutDTO{
		MemberID:          "member-id",
		WorkoutInstanceID: "workout-instance-id",
		ScheduledDate:     "2025-09-08",
	}
	id, err := repo.Create("gym-id", input)
	assert.NoError(t, err)
	assert.Equal(t, "new-id", *id)
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewCustomMemberWorkoutRepository(db)

	mock.ExpectQuery(`SELECT id, member_id, workout_instance_id, scheduled_date, started_at, completed_at, status, notes, rating, created_at, updated_at`).
		WithArgs("notfound-id").
		WillReturnError(sql.ErrNoRows)

	res, err := repo.GetByID("gym-id", "notfound-id")
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, res)
}

func TestUpdateCustomMemberWorkout_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewCustomMemberWorkoutRepository(db)

	mock.ExpectExec(`UPDATE ".*".custom_member_workout SET`).
		WithArgs(nil, nil, nil, nil, nil, "notfound-id").
		WillReturnResult(sqlmock.NewResult(0, 0))

	input := &dto.UpdateCustomMemberWorkoutDTO{ID: "notfound-id"}
	err := repo.Update("gym-id", input)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDeleteCustomMemberWorkout_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()
	repo := NewCustomMemberWorkoutRepository(db)

	mock.ExpectExec(`DELETE FROM ".*".custom_member_workout`).
		WithArgs("notfound-id").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete("gym-id", "notfound-id")
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
