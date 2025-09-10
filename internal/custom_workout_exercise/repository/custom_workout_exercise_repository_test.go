package repository_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_exercise/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomWorkoutExercise(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"

	exercise := &dto.CreateCustomWorkoutExerciseDTO{
		CreatedBy:         "user123",
		WorkoutInstanceID: "workout456",
		ExerciseSource:    "public",
		PublicExerciseID:  stringPtr("exercise789"),
		GymExerciseID:     nil,
		BlockName:         "main",
		ExerciseOrder:     1,
		Sets:              intPtr(3),
		RepsMin:           intPtr(8),
		RepsMax:           intPtr(12),
		WeightKg:          floatPtr(50.5),
		DurationSeconds:   nil,
		RestSeconds:       intPtr(60),
		Notes:             stringPtr("Test exercise"),
	}

	mock.ExpectQuery(`INSERT INTO "gym123".custom_workout_exercise`).
		WithArgs(
			exercise.CreatedBy,
			exercise.WorkoutInstanceID,
			exercise.ExerciseSource,
			exercise.PublicExerciseID,
			exercise.GymExerciseID,
			exercise.BlockName,
			exercise.ExerciseOrder,
			exercise.Sets,
			exercise.RepsMin,
			exercise.RepsMax,
			exercise.WeightKg,
			exercise.DurationSeconds,
			exercise.RestSeconds,
			exercise.Notes,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("exercise123"))

	id, err := repo.Create(gymID, exercise)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "exercise123", *id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateCustomWorkoutExerciseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"

	exercise := &dto.CreateCustomWorkoutExerciseDTO{
		CreatedBy:         "user123",
		WorkoutInstanceID: "workout456",
		ExerciseSource:    "public",
		PublicExerciseID:  stringPtr("exercise789"),
		BlockName:         "main",
		ExerciseOrder:     1,
	}

	mock.ExpectQuery(`INSERT INTO "gym123".custom_workout_exercise`).
		WithArgs(
			exercise.CreatedBy,
			exercise.WorkoutInstanceID,
			exercise.ExerciseSource,
			exercise.PublicExerciseID,
			exercise.GymExerciseID,
			exercise.BlockName,
			exercise.ExerciseOrder,
			exercise.Sets,
			exercise.RepsMin,
			exercise.RepsMax,
			exercise.WeightKg,
			exercise.DurationSeconds,
			exercise.RestSeconds,
			exercise.Notes,
		).
		WillReturnError(errors.New("database error"))

	id, err := repo.Create(gymID, exercise)

	assert.Error(t, err)
	assert.Nil(t, id)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"
	exerciseID := "exercise123"

	rows := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg",
		"duration_seconds", "rest_seconds", "notes", "created_at", "updated_at",
	}).AddRow(
		"exercise123", "user123", "workout456", "public", "exercise789", nil,
		"main", 1, 3, 8, 12, 50.5,
		nil, 60, "Test exercise", "2023-01-01 10:00:00", "2023-01-01 10:00:00",
	)

	mock.ExpectQuery(`SELECT id, created_by, workout_instance_id, exercise_source, public_exercise_id, gym_exercise_id, block_name,`).
		WithArgs(exerciseID).
		WillReturnRows(rows)

	result, err := repo.GetByID(gymID, exerciseID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "exercise123", result.ID)
	assert.Equal(t, "user123", result.CreatedBy)
	assert.Equal(t, "workout456", result.WorkoutInstanceID)
	assert.Equal(t, "public", result.ExerciseSource)
	assert.Equal(t, "exercise789", *result.PublicExerciseID)
	assert.Equal(t, "main", result.BlockName)
	assert.Equal(t, 1, result.ExerciseOrder)
	assert.Equal(t, 3, *result.Sets)
	assert.Equal(t, 8, *result.RepsMin)
	assert.Equal(t, 12, *result.RepsMax)
	assert.Equal(t, 50.5, *result.WeightKg)
	assert.Equal(t, 60, *result.RestSeconds)
	assert.Equal(t, "Test exercise", *result.Notes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListByWorkoutInstanceID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"
	workoutInstanceID := "workout456"

	rows := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg",
		"duration_seconds", "rest_seconds", "notes", "created_at", "updated_at",
	}).
		AddRow("ex1", "user123", "workout456", "public", "exercise789", nil, "warmup", 1, 2, 10, 15, 20.0, nil, 30, "Warmup", "2023-01-01 10:00:00", "2023-01-01 10:00:00").
		AddRow("ex2", "user123", "workout456", "gym", nil, "gym_ex1", "main", 1, 3, 8, 12, 50.5, nil, 60, "Main exercise", "2023-01-01 10:05:00", "2023-01-01 10:05:00")

	mock.ExpectQuery(`SELECT id, created_by, workout_instance_id, exercise_source, public_exercise_id, gym_exercise_id, block_name,`).
		WithArgs(workoutInstanceID).
		WillReturnRows(rows)

	result, err := repo.ListByWorkoutInstanceID(gymID, workoutInstanceID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "ex1", result[0].ID)
	assert.Equal(t, "warmup", result[0].BlockName)
	assert.Equal(t, "public", result[0].ExerciseSource)
	assert.Equal(t, "ex2", result[1].ID)
	assert.Equal(t, "main", result[1].BlockName)
	assert.Equal(t, "gym", result[1].ExerciseSource)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListByMuscularGroupID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"
	muscularGroupID := "muscle456"

	rows := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg",
		"duration_seconds", "rest_seconds", "notes", "created_at", "updated_at",
	}).
		AddRow("ex1", "user123", "workout456", "public", "exercise789", nil, "main", 1, 3, 8, 12, 50.5, nil, 60, "Chest exercise", "2023-01-01 10:00:00", "2023-01-01 10:00:00")

	mock.ExpectQuery(`SELECT cwe.id, cwe.created_by, cwe.workout_instance_id, cwe.exercise_source, cwe.public_exercise_id, cwe.gym_exercise_id,`).
		WithArgs(muscularGroupID).
		WillReturnRows(rows)

	result, err := repo.ListByMuscularGroupID(gymID, muscularGroupID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "ex1", result[0].ID)
	assert.Equal(t, "Chest exercise", *result[0].Notes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListByEquipmentID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"
	equipmentID := "equipment456"

	rows := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg",
		"duration_seconds", "rest_seconds", "notes", "created_at", "updated_at",
	}).
		AddRow("ex1", "user123", "workout456", "public", "exercise789", nil, "main", 1, 3, 8, 12, 50.5, nil, 60, "Barbell exercise", "2023-01-01 10:00:00", "2023-01-01 10:00:00")

	mock.ExpectQuery(`SELECT cwe.id, cwe.created_by, cwe.workout_instance_id, cwe.exercise_source, cwe.public_exercise_id, cwe.gym_exercise_id,`).
		WithArgs(equipmentID).
		WillReturnRows(rows)

	result, err := repo.ListByEquipmentID(gymID, equipmentID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "ex1", result[0].ID)
	assert.Equal(t, "Barbell exercise", *result[0].Notes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCustomWorkoutExercise(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"

	exercise := &dto.UpdateCustomWorkoutExerciseDTO{
		ID:              "exercise123",
		Sets:            intPtr(4),
		RepsMin:         intPtr(6),
		RepsMax:         intPtr(10),
		WeightKg:        floatPtr(60.0),
		DurationSeconds: nil,
		RestSeconds:     intPtr(90),
		Notes:           stringPtr("Updated exercise"),
	}

	mock.ExpectExec(`UPDATE "gym123".custom_workout_exercise SET`).
		WithArgs(
			exercise.Sets,
			exercise.RepsMin,
			exercise.RepsMax,
			exercise.WeightKg,
			exercise.DurationSeconds,
			exercise.RestSeconds,
			exercise.Notes,
			exercise.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(gymID, exercise)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCustomWorkoutExercise(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutExerciseRepository(db)
	gymID := "gym123"
	exerciseID := "exercise123"

	mock.ExpectExec(`DELETE FROM "gym123".custom_workout_exercise WHERE id = \$1`).
		WithArgs(exerciseID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(gymID, exerciseID)

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

func floatPtr(f float64) *float64 {
	return &f
}
