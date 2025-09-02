package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/custom_equipment/dto"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *CustomEquipmentRepository) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	repo := NewCustomEquipmentRepository(db)
	return db, mock, repo
}

func TestCreateCustomEquipment(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	gymID := "tenant_schema"
	equipment := &dto.CreateCustomEquipmentDTO{
		CreatedBy:   "user123",
		Name:        "Dumbbell",
		Description: "A dumbbell",
		Category:    "weight",
		IsActive:    true,
	}

	query := `INSERT INTO "` + gymID + `".custom_equipment (created_by, name, description, category, is_active) VALUES ($1, $2, $3, $4, $5)`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(equipment.CreatedBy, equipment.Name, equipment.Description, equipment.Category, equipment.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(gymID, equipment)
	assert.NoError(t, err)
}

func TestGetCustomEquipmentByID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	gymID := "tenant_schema"
	expected := &dto.ResponseCustomEquipmentDTO{
		ID:          "eq-1",
		CreatedBy:   "user123",
		Name:        "Dumbbell",
		Description: "A dumbbell",
		Category:    "weight",
		IsActive:    true,
	}

	query := `SELECT id, created_by, name, description, category, is_active FROM "` + gymID + `".custom_equipment WHERE id = $1`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_by", "name", "description", "category", "is_active"}).
			AddRow(expected.ID, expected.CreatedBy, expected.Name, expected.Description, expected.Category, expected.IsActive))

	result, err := repo.GetByID(gymID, expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Category, result.Category)
}

func TestUpdateCustomEquipment(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	gymID := "tenant_schema"
	update := &dto.UpdateCustomEquipmentDTO{
		ID:          "eq-1",
		Name:        ptr("Barbell"),
		Description: ptr("A barbell"),
		Category:    ptr("weight"),
		IsActive:    ptr(true),
	}

	query := `UPDATE "` + gymID + `".custom_equipment SET name = $1, description = $2, category = $3, is_active = $4 WHERE id = $5`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(update.Name, update.Description, update.Category, update.IsActive, update.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(gymID, update)
	assert.NoError(t, err)
}

func TestDeleteCustomEquipment(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	gymID := "tenant_schema"
	id := "eq-1"
	query := `DELETE FROM "` + gymID + `".custom_equipment WHERE id = $1`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(gymID, id)
	assert.NoError(t, err)
}

// Helper for pointer fields
func ptr[T any](v T) *T { return &v }
