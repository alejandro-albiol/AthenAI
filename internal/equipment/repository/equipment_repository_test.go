package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/equipment/dto"
)

type mockDB struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func setupMockDB(t *testing.T) *mockDB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	return &mockDB{db: db, mock: mock}
}

func teardownMockDB(m *mockDB) {
	m.db.Close()
}

func TestEquipmentRepository_CreateEquipment(t *testing.T) {
	m := setupMockDB(t)
	defer teardownMockDB(m)
	repo := NewEquipmentRepository(m.db)

	t.Run("success", func(t *testing.T) {
		input := &dto.EquipmentCreationDTO{
			Name:        "Dumbbell",
			Description: "A basic dumbbell",
			Category:    "weights",
		}
		mockID := "123"
		m.mock.ExpectQuery(`INSERT INTO public.equipment`).
			WithArgs(input.Name, input.Description, input.Category).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockID))

		id, err := repo.CreateEquipment(input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if id == nil || *id != mockID {
			t.Errorf("expected id '%s', got %v", mockID, id)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		input := &dto.EquipmentCreationDTO{
			Name:        "Barbell",
			Description: "A basic barbell",
			Category:    "weights",
		}
		dbErr := sql.ErrConnDone
		m.mock.ExpectQuery(`INSERT INTO public.equipment`).
			WithArgs(input.Name, input.Description, input.Category).
			WillReturnError(dbErr)

		id, err := repo.CreateEquipment(input)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if id != nil {
			t.Errorf("expected nil id, got %v", id)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestEquipmentRepository_GetEquipmentByID(t *testing.T) {
	m := setupMockDB(t)
	defer teardownMockDB(m)
	repo := NewEquipmentRepository(m.db)

	t.Run("success", func(t *testing.T) {
		mockID := "id1"
		mockRow := sqlmock.NewRows([]string{"id", "name", "description", "category", "is_active", "created_at", "updated_at"}).
			AddRow(mockID, "Dumbbell", "A basic dumbbell", "weights", true, "2023-01-01", "2023-01-02")
		m.mock.ExpectQuery(`SELECT id, name, description, category, is_active, created_at, updated_at FROM public.equipment WHERE id = \$1`).
			WithArgs(mockID).
			WillReturnRows(mockRow)

		res, err := repo.GetEquipmentByID(mockID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res == nil || res.ID != mockID {
			t.Errorf("expected id '%s', got %v", mockID, res)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockID := "id404"
		m.mock.ExpectQuery(`SELECT id, name, description, category, is_active, created_at, updated_at FROM public.equipment WHERE id = \$1`).
			WithArgs(mockID).
			WillReturnError(sql.ErrNoRows)

		res, err := repo.GetEquipmentByID(mockID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if res != nil {
			t.Errorf("expected nil result, got %v", res)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestEquipmentRepository_GetAllEquipment(t *testing.T) {
	m := setupMockDB(t)
	defer teardownMockDB(m)
	repo := NewEquipmentRepository(m.db)

	t.Run("success", func(t *testing.T) {
		mockRows := sqlmock.NewRows([]string{"id", "name", "description", "category", "is_active", "created_at", "updated_at"}).
			AddRow("id1", "Dumbbell", "A basic dumbbell", "weights", true, "2023-01-01", "2023-01-02").
			AddRow("id2", "Barbell", "A basic barbell", "weights", true, "2023-01-03", "2023-01-04")
		m.mock.ExpectQuery(`SELECT id, name, description, category, is_active, created_at, updated_at FROM public.equipment`).
			WillReturnRows(mockRows)

		res, err := repo.GetAllEquipment()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(res) != 2 {
			t.Errorf("expected 2 results, got %d", len(res))
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		m.mock.ExpectQuery(`SELECT id, name, description, category, is_active, created_at, updated_at FROM public.equipment`).
			WillReturnError(sql.ErrConnDone)

		res, err := repo.GetAllEquipment()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if res != nil {
			t.Errorf("expected nil result, got %v", res)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestEquipmentRepository_UpdateEquipment(t *testing.T) {
	m := setupMockDB(t)
	defer teardownMockDB(m)
	repo := NewEquipmentRepository(m.db)

	t.Run("success", func(t *testing.T) {
		mockID := "id1"
		mockRow := sqlmock.NewRows([]string{"id", "name", "description", "category", "is_active", "created_at", "updated_at"}).
			AddRow(mockID, "Dumbbell", "A basic dumbbell", "weights", true, "2023-01-01", "2023-01-02")
		update := &dto.EquipmentUpdateDTO{
			Name:        nil,
			Description: nil,
			Category:    nil,
			IsActive:    nil,
		}
		m.mock.ExpectQuery(`UPDATE public.equipment SET`).
			WithArgs(mockID, update.Name, update.Description, update.Category, update.IsActive).
			WillReturnRows(mockRow)

		res, err := repo.UpdateEquipment(mockID, update)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res == nil || res.ID != mockID {
			t.Errorf("expected id '%s', got %v", mockID, res)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mockID := "id1"
		update := &dto.EquipmentUpdateDTO{}
		m.mock.ExpectQuery(`UPDATE public.equipment SET`).
			WithArgs(mockID, update.Name, update.Description, update.Category, update.IsActive).
			WillReturnError(sql.ErrConnDone)

		res, err := repo.UpdateEquipment(mockID, update)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if res != nil {
			t.Errorf("expected nil result, got %v", res)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestEquipmentRepository_DeleteEquipment(t *testing.T) {
	m := setupMockDB(t)
	defer teardownMockDB(m)
	repo := NewEquipmentRepository(m.db)

	t.Run("success", func(t *testing.T) {
		mockID := "id1"
		m.mock.ExpectExec(`DELETE FROM public.equipment WHERE id = \$1`).
			WithArgs(mockID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteEquipment(mockID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mockID := "id1"
		m.mock.ExpectExec(`DELETE FROM public.equipment WHERE id = \$1`).
			WithArgs(mockID).
			WillReturnError(sql.ErrConnDone)

		err := repo.DeleteEquipment(mockID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err := m.mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}
