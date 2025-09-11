package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomExercise(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewCustomExerciseRepository(db)
	dtoReq := &dto.CustomExerciseCreationDTO{
		Name:            "Push Up",
		Synonyms:        []string{"Pushup"},
		DifficultyLevel: "easy",
		ExerciseType:    "bodyweight",
		Instructions:    "Do a push up",
		VideoURL:        "http://example.com/video",
		ImageURL:        "http://example.com/image",
		MuscularGroups:  []string{"chest", "triceps"},
		CreatedBy:       "user1",
	}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("ex-1"))
	id, err := repo.CreateCustomExercise("tenant1", dtoReq)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "ex-1", *id)
}

func TestGetCustomExerciseByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewCustomExerciseRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(
		sqlmock.NewRows([]string{
			"id", "created_by", "name", "synonyms", "difficulty_level",
			"exercise_type", "instructions", "video_url", "image_url", "is_active",
		}).AddRow(
			"ex-1", "user1", "Push Up", `{Pushup}`, "easy",
			"bodyweight", "Do a push up", "http://example.com/video", "http://example.com/image", true,
		))
	result, err := repo.GetCustomExerciseByID("tenant1", "ex-1")
	assert.NoError(t, err)
	assert.Equal(t, "ex-1", result.ID)
}

func TestListCustomExercises(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewCustomExerciseRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(
		sqlmock.NewRows([]string{
			"id", "created_by", "name", "synonyms", "difficulty_level",
			"exercise_type", "instructions", "video_url", "image_url", "is_active",
		}).AddRow(
			"ex-1", "user1", "Push Up", `{Pushup}`, "easy",
			"bodyweight", "Do a push up", "http://example.com/video", "http://example.com/image", true,
		))
	results, err := repo.ListCustomExercises("tenant1")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "ex-1", results[0].ID)
}

func TestDeleteCustomExercise(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewCustomExerciseRepository(db)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.DeleteCustomExercise("tenant1", "ex-1")
	assert.NoError(t, err)
}
