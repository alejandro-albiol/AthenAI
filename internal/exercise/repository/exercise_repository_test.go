package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/exercise/dto"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *ExerciseRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	repo := NewExerciseRepository(db)
	return db, mock, repo
}

func TestCreateExercise(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ex := dto.ExerciseCreationDTO{
		Name:            "Push Up",
		Synonyms:        []string{"Press Up"},
		MuscularGroups:  []string{"chest", "triceps"},
		Equipment: []string{},
		DifficultyLevel: "beginner",
		ExerciseType:    "strength",
		Instructions:    "Do a push up.",
		VideoURL:        nil,
		ImageURL:        nil,
		CreatedBy:       "admin-uuid",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO public.exercise (`)).
		WithArgs(
			ex.Name,
			pq.Array(ex.Synonyms),
			pq.Array(ex.MuscularGroups),
			pq.Array(ex.Equipment),
			ex.DifficultyLevel,
			ex.ExerciseType,
			ex.Instructions,
			ex.VideoURL,
			ex.ImageURL,
			ex.CreatedBy,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("exercise-uuid"))

	id, err := repo.CreateExercise(ex)
	assert.NoError(t, err)
	assert.Equal(t, "exercise-uuid", id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetExerciseByID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE id = $1 AND is_active = TRUE`)).
		WithArgs("exercise-uuid").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "synonyms", "muscular_groups", "equipment_needed", "difficulty_level", "exercise_type", "instructions", "video_url", "image_url", "created_by", "is_active", "created_at", "updated_at"}).
			AddRow("exercise-uuid", "Push Up", pq.StringArray{"Press Up"}, pq.StringArray{"chest", "triceps"}, pq.StringArray{}, "beginner", "strength", "Do a push up.", nil, nil, "admin-uuid", true, now, now))

	ex, err := repo.GetExerciseByID("exercise-uuid")
	assert.NoError(t, err)
	assert.Equal(t, "Push Up", ex.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllExercises(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE is_active = TRUE`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "synonyms", "muscular_groups", "equipment_needed", "difficulty_level", "exercise_type", "instructions", "video_url", "image_url", "created_by", "is_active", "created_at", "updated_at"}).
			AddRow("exercise-uuid", "Push Up", pq.StringArray{"Press Up"}, pq.StringArray{"chest", "triceps"}, pq.StringArray{}, "beginner", "strength", "Do a push up.", nil, nil, "admin-uuid", true, now, now))

	exs, err := repo.GetAllExercises()
	assert.NoError(t, err)
	assert.Len(t, exs, 1)
	assert.Equal(t, "Push Up", exs[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateExercise(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()
	now := time.Now()
	update := dto.ExerciseUpdateDTO{
		Name:            ptrString("Push Up Updated"),
		Synonyms:        []string{"Press Up"},
		MuscularGroups:  []string{"chest", "triceps"},
		DifficultyLevel: "beginner",
		ExerciseType:    "strength",
		Instructions:    ptrString("Do a push up better."),
		VideoURL:        nil,
		ImageURL:        nil,
		IsActive:        ptrBool(true),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE public.exercise SET`)).
		WithArgs(
			"exercise-uuid",
			update.Name,
			pq.Array(update.Synonyms),
			pq.Array(update.MuscularGroups),
			pq.Array(update.Equipment),
			update.DifficultyLevel,
			update.ExerciseType,
			update.Instructions,
			update.VideoURL,
			update.ImageURL,
			update.IsActive,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "synonyms", "muscular_groups", "equipment_needed", "difficulty_level", "exercise_type", "instructions", "video_url", "image_url", "created_by", "is_active", "created_at", "updated_at"}).
			AddRow("exercise-uuid", "Push Up Updated", pq.StringArray{"Press Up"}, pq.StringArray{"chest", "triceps"}, pq.StringArray{}, "beginner", "strength", "Do a push up better.", nil, nil, "admin-uuid", true, now, now))

	ex, err := repo.UpdateExercise("exercise-uuid", update)
	assert.NoError(t, err)
	assert.Equal(t, "Push Up Updated", ex.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteExercise(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE public.exercise SET is_active = FALSE, updated_at = NOW() WHERE id = $1`)).
		WithArgs("exercise-uuid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteExercise("exercise-uuid")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetExercisesByMuscularGroup(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()
	now := time.Now()
	groups := []string{"chest", "triceps"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE is_active = TRUE AND muscular_groups && $1`)).
		WithArgs(pq.Array(groups)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "synonyms", "muscular_groups", "equipment_needed", "difficulty_level", "exercise_type", "instructions", "video_url", "image_url", "created_by", "is_active", "created_at", "updated_at"}).
			AddRow("exercise-uuid", "Push Up", pq.StringArray{"Press Up"}, pq.StringArray{"chest", "triceps"}, pq.StringArray{}, "beginner", "strength", "Do a push up.", nil, nil, "admin-uuid", true, now, now))

	exs, err := repo.GetExercisesByMuscularGroup(groups)
	assert.NoError(t, err)
	assert.Len(t, exs, 1)
	assert.Equal(t, "Push Up", exs[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetExercisesByEquipment(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()
	now := time.Now()
	equip := []string{"barbell"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, synonyms, muscular_groups, equipment_needed, difficulty_level, exercise_type, instructions, video_url, image_url, created_by, is_active, created_at, updated_at FROM public.exercise WHERE is_active = TRUE AND equipment_needed && $1`)).
		WithArgs(pq.Array(equip)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "synonyms", "muscular_groups", "equipment_needed", "difficulty_level", "exercise_type", "instructions", "video_url", "image_url", "created_by", "is_active", "created_at", "updated_at"}).
			AddRow("exercise-uuid", "Barbell Curl", pq.StringArray{"Curl"}, pq.StringArray{"biceps"}, pq.StringArray{"barbell"}, "beginner", "strength", "Do a curl.", nil, nil, "admin-uuid", true, now, now))

	exs, err := repo.GetExercisesByEquipment(equip)
	assert.NoError(t, err)
	assert.Len(t, exs, 1)
	assert.Equal(t, "Barbell Curl", exs[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }
