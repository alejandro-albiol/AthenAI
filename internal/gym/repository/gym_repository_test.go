package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/repository"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock db: %v", err)
	}
	return db, mock
}

func TestCreateGym(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	gymDTO := dto.GymCreationDTO{
		Name:        "Test Gym",
		Email:       "test@gym.com",
		Address:     "123 Test St",
		ContactName: "John Doe",
		Phone:       "+1234567890",
		LogoURL:     "https://example.com/logo.png",
	}

	// Test successful creation
	rows := sqlmock.NewRows([]string{"id"}).AddRow("test-gym")
	mock.ExpectQuery("INSERT INTO gym").
		WithArgs(
			gymDTO.Name,
			gymDTO.Email,
			gymDTO.Address,
			gymDTO.Phone,
			sqlmock.AnyArg(),
		).WillReturnRows(rows)

	id, err := repo.CreateGym(gymDTO)
	assert.NoError(t, err)
	assert.Equal(t, "test-gym", id)
}

func TestGetGymByID(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	now := time.Now()

	// Test successful retrieval
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "address", "phone", "is_active", "created_at", "updated_at",
	}).AddRow(
		"gym123", "Test Gym", "test@gym.com", "123 Test St",
		"+1234567890", true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM gym WHERE id").
		WithArgs("gym123").
		WillReturnRows(rows)

	gym, err := repo.GetGymByID("gym123")
	assert.NoError(t, err)
	assert.Equal(t, "gym123", gym.ID)
	assert.Equal(t, "Test Gym", gym.Name)
}

func TestGetAllGyms(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	now := time.Now()

	// Test successful retrieval
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "address", "phone", "is_active", "created_at", "updated_at",
	}).AddRow(
		"gym123", "Test Gym 1", "test1@gym.com", "123 Test St",
		"+1234567890", true, now, now,
	).AddRow(
		"gym456", "Test Gym 2", "test2@gym.com", "456 Test St",
		"+0987654321", true, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM gym").
		WillReturnRows(rows)

	gyms, err := repo.GetAllGyms()
	assert.NoError(t, err)
	assert.Len(t, gyms, 2)
	assert.Equal(t, "Test Gym 1", gyms[0].Name)
	assert.Equal(t, "Test Gym 2", gyms[1].Name)
}

func TestUpdateGym(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	updateDTO := dto.GymUpdateDTO{
		Name:    "Updated Gym",
		Email:   "updated@gym.com",
		Address: "456 Update St",
		Phone:   "+0987654321",
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "address", "phone", "is_active", "created_at", "updated_at",
	}).AddRow(
		"gym123", updateDTO.Name, updateDTO.Email, updateDTO.Address,
		updateDTO.Phone, true, time.Now(), time.Now(),
	)

	mock.ExpectQuery("UPDATE gym").WithArgs(
		updateDTO.Name,
		updateDTO.Email,
		updateDTO.Address,
		updateDTO.Phone,
		sqlmock.AnyArg(), // updated_at
		"gym123",         // id
	).WillReturnRows(rows)

	updatedGym, err := repo.UpdateGym("gym123", updateDTO)
	assert.NoError(t, err)
	assert.Equal(t, "gym123", updatedGym.ID)
}

func TestSetGymActive(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)

	mock.ExpectExec("UPDATE gym SET is_active").
		WithArgs(true, sqlmock.AnyArg(), "gym123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SetGymActive("gym123", true)
	assert.NoError(t, err)
}

func TestDeleteGym(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)

	mock.ExpectExec("UPDATE gym SET deleted_at").
		WithArgs(sqlmock.AnyArg(), "gym123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteGym("gym123")
	assert.NoError(t, err)
}
