package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/dto"
	"github.com/alejandro-albiol/athenai/internal/custom_workout_instance/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomWorkoutInstance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	createdBy := "user123"

	instance := &dto.CreateCustomWorkoutInstanceDTO{
		Name:           "Test Workout",
		Description:    "Test workout description",
		TemplateSource: "gym",
		GymTemplateID:  stringPtr("template456"),
	}

	mock.ExpectQuery(`INSERT INTO "gym123".custom_workout_instance`).
		WithArgs(
			createdBy,
			instance.Name,
			instance.Description,
			instance.TemplateSource,
			instance.PublicTemplateID,
			instance.GymTemplateID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("instance123"))

	id, err := repo.Create(gymID, createdBy, instance)

	assert.NoError(t, err)
	assert.Equal(t, "instance123", *id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateCustomWorkoutInstance_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	createdBy := "user123"

	instance := &dto.CreateCustomWorkoutInstanceDTO{
		Name:           "Test Workout",
		TemplateSource: "gym",
		GymTemplateID:  stringPtr("template456"),
	}

	mock.ExpectQuery(`INSERT INTO "gym123".custom_workout_instance`).
		WithArgs(
			createdBy,
			instance.Name,
			instance.Description,
			instance.TemplateSource,
			instance.PublicTemplateID,
			instance.GymTemplateID,
		).
		WillReturnError(errors.New("database error"))

	id, err := repo.Create(gymID, createdBy, instance)

	assert.Error(t, err)
	assert.Nil(t, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCustomWorkoutInstanceByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	instanceID := "instance123"

	// Mock basic workout instance query
	now := time.Now()
	basicRows := sqlmock.NewRows([]string{
		"id", "name", "description", "template_source", "public_template_id", "gym_template_id",
		"created_at", "updated_at",
	}).AddRow(
		instanceID, "Test Workout", "Description", "gym", nil, "template456",
		now, now,
	)

	mock.ExpectQuery(`SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at FROM "gym123".custom_workout_instance WHERE id = \$1`).
		WithArgs(instanceID).
		WillReturnRows(basicRows)

	// Mock exercises query
	now2 := time.Now()
	exerciseRows := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg", "duration_seconds",
		"rest_seconds", "notes", "created_at", "updated_at",
	}).AddRow(
		"exercise1", "user123", instanceID, "public", "pub_ex1", nil,
		"main", 1, 3, 8, 12, 50.5, 30,
		60, "Test notes", now2, now2,
	)

	mock.ExpectQuery(`SELECT (.+) FROM "gym123".custom_workout_exercise cwe WHERE cwe.workout_instance_id = \$1`).
		WithArgs(instanceID).
		WillReturnRows(exerciseRows)

	instance, err := repo.GetByID(gymID, instanceID)

	assert.NoError(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, instanceID, instance.ID)
	assert.Equal(t, "Test Workout", instance.Name)
	assert.Equal(t, "gym", instance.TemplateSource)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCustomWorkoutInstanceSummaryByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	instanceID := "instance123"

	// Mock basic workout instance query (same as GetByID)
	now := time.Now()
	basicRows := sqlmock.NewRows([]string{
		"id", "name", "description", "template_source", "public_template_id", "gym_template_id",
		"created_at", "updated_at",
	}).AddRow(
		instanceID, "Test Workout", "Description", "gym", nil, "template456",
		now, now,
	)

	mock.ExpectQuery(`SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at FROM "gym123".custom_workout_instance WHERE id = \$1`).
		WithArgs(instanceID).
		WillReturnRows(basicRows)

	// Mock exercises query (same as GetByID)
	now2 := time.Now()
	exerciseRows := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg", "duration_seconds",
		"rest_seconds", "notes", "created_at", "updated_at",
	}).AddRow(
		"exercise1", "user123", instanceID, "public", "pub_ex1", nil,
		"main", 1, 3, 8, 12, 50.5, 30,
		60, "Test notes", now2, now2,
	)

	mock.ExpectQuery(`SELECT (.+) FROM "gym123".custom_workout_exercise cwe WHERE cwe.workout_instance_id = \$1`).
		WithArgs(instanceID).
		WillReturnRows(exerciseRows)

	summary, err := repo.GetSummaryByID(gymID, instanceID)

	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, instanceID, summary.ID)
	assert.Equal(t, "Test Workout", summary.Name)
	assert.Equal(t, "gym", summary.TemplateSource)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCustomWorkoutInstancesByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	userID := "user123"

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "template_source", "public_template_id", "gym_template_id",
		"created_at", "updated_at",
	}).AddRow(
		"instance1", "Workout 1", "Description 1", "gym", nil, "template1",
		now, now,
	).AddRow(
		"instance2", "Workout 2", "Description 2", "public", "template2", nil,
		now, now,
	)

	mock.ExpectQuery(`SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at FROM "gym123".custom_workout_instance WHERE created_by = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnRows(rows)

	// Mock exercises queries for each instance (the getWorkoutInstanceList method calls getWorkoutExercises for each instance)
	exerciseRows1 := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg", "duration_seconds",
		"rest_seconds", "notes", "created_at", "updated_at",
	}).AddRow(
		"exercise1", "user123", "instance1", "public", "pub_ex1", nil,
		"main", 1, 3, 8, 12, 50.5, 30,
		60, "Test notes", now, now,
	)

	exerciseRows2 := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg", "duration_seconds",
		"rest_seconds", "notes", "created_at", "updated_at",
	}).AddRow(
		"exercise2", "user123", "instance2", "public", "pub_ex2", nil,
		"main", 1, 3, 8, 12, 50.5, 30,
		60, "Test notes", now, now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM "gym123".custom_workout_exercise cwe WHERE cwe.workout_instance_id = \$1`).
		WithArgs("instance1").
		WillReturnRows(exerciseRows1)

	mock.ExpectQuery(`SELECT (.+) FROM "gym123".custom_workout_exercise cwe WHERE cwe.workout_instance_id = \$1`).
		WithArgs("instance2").
		WillReturnRows(exerciseRows2)

	instances, err := repo.GetByUserID(gymID, userID)

	assert.NoError(t, err)
	assert.Len(t, instances, 2)
	assert.Equal(t, "instance1", instances[0].ID)
	assert.Equal(t, "instance2", instances[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastCustomWorkoutInstancesByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	userID := "user123"
	count := 3

	// Mock the main query - note the specific column order
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "template_source", "public_template_id", "gym_template_id", "created_at", "updated_at",
	}).AddRow(
		"instance1", "Workout 1", "Description 1", "gym", nil, "template1", time.Now(), time.Now(),
	)

	expectedSQL := `SELECT id, name, description, template_source, public_template_id, gym_template_id, created_at, updated_at FROM "gym123"\.custom_workout_instance WHERE created_by = \$1 ORDER BY created_at DESC LIMIT \$2`
	mock.ExpectQuery(expectedSQL).
		WithArgs(userID, count).
		WillReturnRows(rows)

	// Mock the exercise query for the returned instance
	exerciseRows := sqlmock.NewRows([]string{
		"id", "created_by", "workout_instance_id", "exercise_source", "public_exercise_id", "gym_exercise_id",
		"block_name", "exercise_order", "sets", "reps_min", "reps_max", "weight_kg", "duration_seconds",
		"rest_seconds", "notes", "created_at", "updated_at",
	}).AddRow(
		"exercise1", "user123", "instance1", "public", "public_ex1", nil,
		"Block A", 1, 3, 8, 12, 50.0, nil, 60, "Test notes", time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT (.+) FROM "gym123"\.custom_workout_exercise cwe WHERE cwe\.workout_instance_id = \$1`).
		WithArgs("instance1").
		WillReturnRows(exerciseRows)

	instances, err := repo.GetLastsByUserID(gymID, userID, count)

	assert.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, "instance1", instances[0].ID)
	assert.Equal(t, "Workout 1", instances[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCustomWorkoutInstance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	instanceID := "instance123"

	updateData := &dto.UpdateCustomWorkoutInstanceDTO{
		Name:        stringPtr("Updated Workout"),
		Description: stringPtr("Updated description"),
	}

	mock.ExpectExec(`UPDATE "gym123"\.custom_workout_instance SET name = COALESCE\(\$1, name\), description = COALESCE\(\$2, description\), template_source = COALESCE\(\$3, template_source\), public_template_id = COALESCE\(\$4, public_template_id\), gym_template_id = COALESCE\(\$5, gym_template_id\), updated_at = NOW\(\) WHERE id = \$6`).
		WithArgs("Updated Workout", "Updated description", (*string)(nil), (*string)(nil), (*string)(nil), instanceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(gymID, instanceID, updateData)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCustomWorkoutInstance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCustomWorkoutInstanceRepository(db)
	gymID := "gym123"
	instanceID := "instance123"

	mock.ExpectExec(`DELETE FROM "gym123"\.custom_workout_instance WHERE id = \$1`).
		WithArgs(instanceID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(gymID, instanceID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
