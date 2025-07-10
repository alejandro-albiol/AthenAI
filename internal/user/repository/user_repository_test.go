package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alejandro-albiol/athenai/internal/user/dto"
	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
	"github.com/alejandro-albiol/athenai/internal/user/repository"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock db: %v", err)
	}
	return db, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	userDTO := dto.UserCreationDTO{
		Username: "testuser",
		Email:    "test@test.com",
		Password: "hashedpassword",
		Role:     userrole_enum.User,
	}
	gymID := "gym123"

	// Test successful creation
	mock.ExpectExec("INSERT INTO users").
		WithArgs(userDTO.Username, userDTO.Email, userDTO.Password, userDTO.Role, gymID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateUser(gymID, userDTO)
	assert.NoError(t, err)
}

func TestGetUserByID(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"
	userID := "user123"
	now := time.Now()

	testCases := []struct {
		name      string
		setupMock func()
		wantErr   bool
		want      dto.UserResponseDTO
	}{
		{
			name: "successful fetch",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "verified", "is_active", "created_at", "updated_at"}).
					AddRow(userID, "testuser", "test@test.com", userrole_enum.User, true, true, now, now)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs(userID, gymID).
					WillReturnRows(rows)
			},
			wantErr: false,
			want: dto.UserResponseDTO{
				ID:        userID,
				Username:  "testuser",
				Email:     "test@test.com",
				Role:      userrole_enum.User,
				GymID:     gymID,
				Verified:  true,
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "user not found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs(userID, gymID).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			want:    dto.UserResponseDTO{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			got, err := repo.GetUserByID(gymID, userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestVerifyUser(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"
	userID := "user123"

	mock.ExpectExec("UPDATE users SET verified").
		WithArgs(userID, gymID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.VerifyUser(gymID, userID)
	assert.NoError(t, err)
}

func TestSetUserActive(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"
	userID := "user123"

	testCases := []struct {
		name    string
		active  bool
		wantErr bool
	}{
		{
			name:    "set active true",
			active:  true,
			wantErr: false,
		},
		{
			name:    "set active false",
			active:  false,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectExec("UPDATE users SET is_active").
				WithArgs(tc.active, userID, gymID).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err := repo.SetUserActive(gymID, userID, tc.active)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"
	username := "testuser"

	testCases := []struct {
		name      string
		setupMock func()
		wantErr   bool
		want      dto.UserResponseDTO
	}{
		{
			name: "successful fetch",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "verified", "is_active", "created_at", "updated_at"}).
					AddRow("user123", username, "test@test.com", userrole_enum.User, false, true, time.Time{}, time.Time{})
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username").
					WithArgs(username, gymID).
					WillReturnRows(rows)
			},
			wantErr: false,
			want: dto.UserResponseDTO{
				ID:        "user123",
				Username:  username,
				Email:     "test@test.com",
				Role:      userrole_enum.User,
				GymID:     gymID,
				Verified:  false,
				IsActive:  true,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		{
			name: "user not found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username").
					WithArgs(username, gymID).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			want:    dto.UserResponseDTO{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			got, err := repo.GetUserByUsername(gymID, username)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"

	testCases := []struct {
		name      string
		setupMock func()
		wantErr   bool
		want      []dto.UserResponseDTO
	}{
		{
			name: "successful fetch multiple users",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "verified", "is_active", "created_at", "updated_at"}).
					AddRow("user1", "testuser1", "test1@test.com", userrole_enum.User, true, true, time.Time{}, time.Time{}).
					AddRow("user2", "testuser2", "test2@test.com", userrole_enum.Admin, true, true, time.Time{}, time.Time{})
				mock.ExpectQuery("SELECT (.+) FROM users WHERE gym_id").
					WithArgs(gymID).
					WillReturnRows(rows)
			},
			wantErr: false,
			want: []dto.UserResponseDTO{
				{
					ID:        "user1",
					Username:  "testuser1",
					Email:     "test1@test.com",
					Role:      userrole_enum.User,
					GymID:     gymID,
					Verified:  true,
					IsActive:  true,
					CreatedAt: time.Time{}, // We'll use time.Time{} for consistency in tests
					UpdatedAt: time.Time{},
				},
				{
					ID:        "user2",
					Username:  "testuser2",
					Email:     "test2@test.com",
					Role:      userrole_enum.Admin,
					GymID:     gymID,
					Verified:  true,
					IsActive:  true,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
		},
		{
			name: "no users found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE gym_id").
					WithArgs(gymID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "role", "verified", "is_active", "created_at", "updated_at"}))
			},
			wantErr: false,
			want:    make([]dto.UserResponseDTO, 0), // Initialize empty slice instead of nil
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			got, err := repo.GetAllUsers(gymID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"
	userID := "user123"
	updateDTO := dto.UserUpdateDTO{
		Username: "updateduser",
		Email:    "updated@test.com",
	}

	testCases := []struct {
		name      string
		setupMock func()
		wantErr   bool
	}{
		{
			name: "successful update",
			setupMock: func() {
				mock.ExpectExec("UPDATE users SET").
					WithArgs(updateDTO.Username, updateDTO.Email, userID, gymID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "user not found",
			setupMock: func() {
				mock.ExpectExec("UPDATE users SET").
					WithArgs(updateDTO.Username, updateDTO.Email, userID, gymID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := repo.UpdateUser(gymID, userID, updateDTO)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"
	userID := "user123"
	newPasswordHash := "newhash123"

	testCases := []struct {
		name      string
		setupMock func()
		wantErr   bool
	}{
		{
			name: "successful password update",
			setupMock: func() {
				mock.ExpectExec("UPDATE users SET password_hash").
					WithArgs(newPasswordHash, userID, gymID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "user not found",
			setupMock: func() {
				mock.ExpectExec("UPDATE users SET password_hash").
					WithArgs(newPasswordHash, userID, gymID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := repo.UpdatePassword(gymID, userID, newPasswordHash)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock := setupTest(t)
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	gymID := "gym123"
	userID := "user123"

	testCases := []struct {
		name      string
		setupMock func()
		wantErr   bool
	}{
		{
			name: "successful deletion",
			setupMock: func() {
				mock.ExpectExec("DELETE FROM users").
					WithArgs(userID, gymID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "user not found",
			setupMock: func() {
				mock.ExpectExec("DELETE FROM users").
					WithArgs(userID, gymID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := repo.DeleteUser(gymID, userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
