package repository_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/dto"
	"github.com/alejandro-albiol/athenai/internal/muscular_group/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateMuscularGroup(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	mg := &dto.CreateMuscularGroupDTO{Name: "Chest", Description: "Upper body", BodyPart: "Torso"}
	mock.ExpectQuery("INSERT INTO public.muscular_group").
		WithArgs(mg.Name, mg.Description, mg.BodyPart).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	id, err := repo.CreateMuscularGroup(mg)
	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, "1", *id)
}

func TestGetAllMuscularGroups(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "body_part", "is_active"}).
		AddRow("1", "Chest", "Upper body", "Torso", true).
		AddRow("2", "Back", "Upper back", "Torso", true)
	mock.ExpectQuery("SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE is_active = TRUE").WillReturnRows(rows)
	groups, err := repo.GetAllMuscularGroups()
	assert.NoError(t, err)
	assert.Len(t, groups, 2)
	assert.Equal(t, "Chest", groups[0].Name)
	assert.Equal(t, "Back", groups[1].Name)
}

func TestGetMuscularGroupByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	row := sqlmock.NewRows([]string{"id", "name", "description", "body_part", "is_active"}).AddRow("1", "Chest", "Upper body", "Torso", true)
	mock.ExpectQuery("SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE id = \\$1 AND is_active = TRUE").WillReturnRows(row)
	mg, err := repo.GetMuscularGroupByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "Chest", mg.Name)
}

func TestGetMuscularGroupByName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	row := sqlmock.NewRows([]string{"id", "name", "description", "body_part", "is_active"}).AddRow("1", "Chest", "Upper body", "Torso", true)
	mock.ExpectQuery("SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE name = \\$1 AND is_active = TRUE").WillReturnRows(row)
	mg, err := repo.GetMuscularGroupByName("Chest")
	assert.NoError(t, err)
	assert.Equal(t, "Chest", mg.Name)
}

func TestUpdateMuscularGroup(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	name := "Chest"
	description := "Updated"
	bodyPart := "Torso"
	update := &dto.UpdateMuscularGroupDTO{Name: &name, Description: &description, BodyPart: &bodyPart}
	row := sqlmock.NewRows([]string{"id", "name", "description", "body_part", "is_active"}).AddRow("1", "Chest", "Updated", "Torso", true)
	mock.ExpectQuery(`UPDATE public\.muscular_group SET name = COALESCE\(\$2, name\), description = COALESCE\(\$3, description\), body_part = COALESCE\(\$4, body_part\) WHERE id = \$1 RETURNING id, name, description, body_part, is_active`).WillReturnRows(row)
	mg, err := repo.UpdateMuscularGroup("1", update)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", mg.Description)
}

func TestDeleteMuscularGroup(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	mock.ExpectExec("UPDATE public.muscular_group SET is_active = FALSE WHERE id = \\$1").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.DeleteMuscularGroup("1")
	assert.NoError(t, err)
}

func TestGetMuscularGroupByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	mock.ExpectQuery("SELECT id, name, description, body_part, is_active FROM public.muscular_group WHERE id = $1 AND is_active = TRUE").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "body_part", "is_active"}))
	mg, err := repo.GetMuscularGroupByID("999")
	assert.Error(t, err)
	assert.Nil(t, mg)
}

func TestCreateMuscularGroup_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := repository.NewMuscularGroupRepository(db)
	mg := &dto.CreateMuscularGroupDTO{Name: "Chest", Description: "Upper body", BodyPart: "Torso"}
	mock.ExpectQuery("INSERT INTO public.muscular_group").WillReturnError(errors.New("insert error"))
	id, err := repo.CreateMuscularGroup(mg)
	assert.Error(t, err)
	assert.Empty(t, id)
}
