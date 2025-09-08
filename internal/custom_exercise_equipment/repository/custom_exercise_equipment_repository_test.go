package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_exercise_equipment/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateLink(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewCustomExerciseEquipmentRepository(db)
	link := &dto.CustomExerciseEquipment{CustomExerciseID: "ex-1", EquipmentID: "eq-1"}
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO")).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("link-1"))
	id, err := repo.CreateLink("tenant1", link)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "link-1", *id)
}

func TestDeleteLink(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewCustomExerciseEquipmentRepository(db)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM")).WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.DeleteLink("tenant1", "link-1")
	assert.NoError(t, err)
}

func TestFindByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewCustomExerciseEquipmentRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT custom_exercise_id, equipment_id FROM")).WillReturnRows(sqlmock.NewRows([]string{"custom_exercise_id", "equipment_id"}).AddRow("ex-1", "eq-1"))
	link, err := repo.FindByID("tenant1", "link-1")
	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, "ex-1", link.CustomExerciseID)
}
