package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_muscular_group/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateLink_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := &CustomExerciseMuscularGroupRepository{db: db}
	mock.ExpectQuery("INSERT INTO tenant1.custom_exercise_muscular_group").
		WithArgs("ex1", "mg1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))
	id, err := repo.CreateLink("tenant1", &dto.CustomExerciseMuscularGroupCreationDTO{CustomExerciseID: "ex1", MuscularGroupID: "mg1"})
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "123", *id)
}

func TestDeleteLink_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := &CustomExerciseMuscularGroupRepository{db: db}
	mock.ExpectExec(`DELETE FROM tenant1.custom_exercise_muscular_group WHERE id = \$1`).
		WithArgs("id1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	err := repo.DeleteLink("tenant1", "id1")
	assert.NoError(t, err)
}

func TestRemoveAllLinksForExercise_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := &CustomExerciseMuscularGroupRepository{db: db}
	mock.ExpectExec(`DELETE FROM tenant1.custom_exercise_muscular_group WHERE custom_exercise_id = \$1`).
		WithArgs("ex1").
		WillReturnResult(sqlmock.NewResult(0, 2))
	err := repo.RemoveAllLinksForExercise("tenant1", "ex1")
	assert.NoError(t, err)
}

func TestFindByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := &CustomExerciseMuscularGroupRepository{db: db}
	mock.ExpectQuery(`SELECT id, custom_exercise_id, muscular_group_id FROM tenant1.custom_exercise_muscular_group WHERE id = \$1`).
		WithArgs("id1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "custom_exercise_id", "muscular_group_id"}).
			AddRow("id1", "ex1", "mg1"))
	link, err := repo.FindByID("tenant1", "id1")
	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, "id1", link.ID)
}

func TestFindByCustomExerciseID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := &CustomExerciseMuscularGroupRepository{db: db}
	mock.ExpectQuery(`SELECT id, custom_exercise_id, muscular_group_id FROM tenant1.custom_exercise_muscular_group WHERE custom_exercise_id = \$1`).
		WithArgs("ex1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "custom_exercise_id", "muscular_group_id"}).
			AddRow("id1", "ex1", "mg1"))
	links, err := repo.FindByCustomExerciseID("tenant1", "ex1")
	assert.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, "id1", links[0].ID)
}

func TestFindByMuscularGroupID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := &CustomExerciseMuscularGroupRepository{db: db}
	mock.ExpectQuery(`SELECT id, custom_exercise_id, muscular_group_id FROM tenant1.custom_exercise_muscular_group WHERE muscular_group_id = \$1`).
		WithArgs("mg1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "custom_exercise_id", "muscular_group_id"}).
			AddRow("id1", "ex1", "mg1"))
	links, err := repo.FindByMuscularGroupID("tenant1", "mg1")
	assert.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, "id1", links[0].ID)
}
