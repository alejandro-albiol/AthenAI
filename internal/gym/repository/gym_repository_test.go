package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/gym/dto"
	"github.com/alejandro-albiol/athenai/internal/gym/repository"
	"github.com/lib/pq"
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
		Domain:      "test-gym",
		Email:       "test@gym.com",
		Address:     "123 Test St",
		ContactName: "John Doe",
		Phone:       "+1234567890",
		LogoURL:     "https://example.com/logo.png",
	}

	// Test successful creation
	rows := sqlmock.NewRows([]string{"id"}).AddRow("test-gym")
	mock.ExpectQuery("INSERT INTO gyms").
		WithArgs(
			gymDTO.Name,
			gymDTO.Domain,
			gymDTO.Email,
			gymDTO.Address,
			gymDTO.ContactName,
			gymDTO.Phone,
			gymDTO.LogoURL,
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
		"id", "name", "domain", "email", "address", "contact_name", "phone", "logo_url",
		"description", "business_hours", "social_links", "payment_methods",
		"currency", "timezone_offset", "is_active", "created_at", "updated_at",
	}).AddRow(
		"gym123", "Test Gym", "test-gym", "test@gym.com", "123 Test St",
		"John Doe", "+1234567890", "https://example.com/logo.png",
		sql.NullString{String: "Description", Valid: true},
		pq.Array([]string{}), pq.Array([]string{}), pq.Array([]string{}),
		"USD", "UTC+0", true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM gyms WHERE id").
		WithArgs("gym123").
		WillReturnRows(rows)

	gym, err := repo.GetGymByID("gym123")
	assert.NoError(t, err)
	assert.Equal(t, "gym123", gym.ID)
	assert.Equal(t, "Test Gym", gym.Name)
}

func TestGetGymByDomain(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	now := time.Now()

	// Test successful retrieval
	rows := sqlmock.NewRows([]string{
		"id", "name", "domain", "email", "address", "contact_name", "phone", "logo_url",
		"description", "business_hours", "social_links", "payment_methods", "currency",
		"timezone_offset", "is_active", "created_at", "updated_at",
	}).AddRow(
		"gym123", "Test Gym", "test-gym", "test@gym.com", "123 Test St",
		"John Doe", "+1234567890", "https://example.com/logo.png",
		sql.NullString{String: "Description", Valid: true},
		pq.Array([]string{"9-5"}), pq.Array([]string{"fb.com"}), pq.Array([]string{"cash"}),
		"USD", "UTC+0", true, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM gyms").
		WithArgs("test-gym").
		WillReturnRows(rows)

	gym, err := repo.GetGymByDomain("test-gym")
	assert.NoError(t, err)
	assert.Equal(t, "test-gym", gym.Domain)
}

func TestGetAllGyms(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	now := time.Now()

	// Test successful retrieval
	rows := sqlmock.NewRows([]string{
		"id", "name", "domain", "email", "address", "contact_name", "phone", "logo_url",
		"description", "business_hours", "social_links", "payment_methods", "currency",
		"timezone_offset", "is_active", "created_at", "updated_at",
	}).AddRow(
		"gym123", "Test Gym 1", "test-gym-1", "test1@gym.com", "123 Test St",
		"John Doe", "+1234567890", "https://example.com/logo1.png",
		sql.NullString{String: "Description 1", Valid: true},
		pq.Array([]string{"9-5"}), pq.Array([]string{"fb.com"}), pq.Array([]string{"cash"}),
		"USD", "UTC+0", true, now, now,
	).AddRow(
		"gym456", "Test Gym 2", "test-gym-2", "test2@gym.com", "456 Test St",
		"Jane Doe", "+0987654321", "https://example.com/logo2.png",
		sql.NullString{String: "Description 2", Valid: true},
		pq.Array([]string{"24/7"}), pq.Array([]string{"twitter.com"}), pq.Array([]string{"card"}),
		"EUR", "UTC+1", true, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM gyms").
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
		Name:           "Updated Gym",
		Email:          "updated@gym.com",
		Address:        "456 Update St",
		ContactName:    "Jane Doe",
		Phone:          "+0987654321",
		Description:    "Updated Description",
		BusinessHours:  []string{"9-5"},
		SocialLinks:    []string{"fb.com"},
		PaymentMethods: []string{"cash"},
		Currency:       "USD",
		TimezoneOffset: 0,
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "domain", "email", "address", "contact_name", "phone", "logo_url",
		"description", "business_hours", "social_links", "payment_methods", "currency",
		"timezone_offset", "is_active", "created_at", "updated_at",
	}).AddRow(
		"gym123", updateDTO.Name, "test-gym", updateDTO.Email, updateDTO.Address,
		updateDTO.ContactName, updateDTO.Phone, updateDTO.LogoURL,
		sql.NullString{String: "Description", Valid: true},
		pq.Array([]string{"9-5"}), pq.Array([]string{"fb.com"}), pq.Array([]string{"cash"}),
		"USD", 0, true, time.Now(), time.Now(),
	)

	mock.ExpectQuery("UPDATE gyms").WithArgs(
		updateDTO.Name,
		updateDTO.Email,
		updateDTO.Address,
		updateDTO.ContactName,
		updateDTO.Phone,
		updateDTO.LogoURL,
		sqlmock.AnyArg(),
		pq.Array(updateDTO.BusinessHours),
		pq.Array(updateDTO.SocialLinks),
		pq.Array(updateDTO.PaymentMethods),
		updateDTO.Currency,
		updateDTO.TimezoneOffset,
		sqlmock.AnyArg(),
		"gym123",
	).WillReturnRows(rows)

	updatedGym, err := repo.UpdateGym("gym123", updateDTO)
	assert.NoError(t, err)
	assert.Equal(t, "gym123", updatedGym.ID)
}

func TestSetGymActive(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)

	mock.ExpectExec("UPDATE gyms SET is_active").
		WithArgs(true, sqlmock.AnyArg(), "gym123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SetGymActive("gym123", true)
	assert.NoError(t, err)
}

func TestDeleteGym(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)

	mock.ExpectExec("UPDATE gyms SET deleted_at").
		WithArgs(sqlmock.AnyArg(), "gym123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteGym("gym123")
	assert.NoError(t, err)
}
