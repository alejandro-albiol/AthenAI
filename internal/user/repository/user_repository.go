package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/interfaces"
)

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) interfaces.UserRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) CreateUser(gymID string, dto dto.UserCreationDTO) error {
	query := `
        INSERT INTO users (
            username, email, password_hash, role, gym_id, 
            verified, is_active, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, 
            false, true, NOW(), NOW()
        )`
	_, err := r.db.Exec(query, dto.Username, dto.Email, dto.Password, dto.Role, gymID)
	return err
}

func (r *usersRepository) GetUserByID(gymID, id string) (dto.UserResponseDTO, error) {
	query := `
        SELECT id, username, email, role, verified, is_active, created_at, updated_at 
        FROM users 
        WHERE id = $1 AND gym_id = $2 AND is_active = true`
	row := r.db.QueryRow(query, id, gymID)
	var user dto.UserResponseDTO
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Role,
		&user.Verified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	user.GymID = gymID
	return user, nil
}

func (r *usersRepository) GetUserByUsername(gymID, username string) (dto.UserResponseDTO, error) {
	query := "SELECT id, username, email, role, verified, is_active, created_at, updated_at FROM users WHERE username = $1 AND gym_id = $2"
	row := r.db.QueryRow(query, username, gymID)
	var user dto.UserResponseDTO
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Verified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	user.GymID = gymID
	return user, nil
}

func (r *usersRepository) GetUserByEmail(gymID, email string) (dto.UserResponseDTO, error) {
	query := "SELECT id, username, email, role, verified, is_active, created_at, updated_at FROM users WHERE email = $1 AND gym_id = $2"
	row := r.db.QueryRow(query, email, gymID)
	var user dto.UserResponseDTO
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Verified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	user.GymID = gymID
	return user, nil
}

func (r *usersRepository) GetAllUsers(gymID string) ([]dto.UserResponseDTO, error) {
	users := make([]dto.UserResponseDTO, 0) // Initialize empty slice

	rows, err := r.db.Query(`
        SELECT id, username, email, role, verified, is_active, created_at, updated_at 
        FROM users 
        WHERE gym_id = $1`, gymID)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user dto.UserResponseDTO
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Verified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		user.GymID = gymID
		users = append(users, user)
	}
	return users, nil
}

func (r *usersRepository) GetPasswordHashByUsername(gymID, username string) (string, error) {
	query := "SELECT password_hash FROM users WHERE username = $1 AND gym_id = $2"
	row := r.db.QueryRow(query, username, gymID)
	var passwordHash string
	err := row.Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (r *usersRepository) UpdateUser(gymID string, id string, user dto.UserUpdateDTO) error {
	query := "UPDATE users SET username = $1, email = $2 WHERE id = $3 AND gym_id = $4"
	_, err := r.db.Exec(query, user.Username, user.Email, id, gymID)
	return err
}

func (r *usersRepository) UpdatePassword(gymID, userID string, newPasswordHash string) error {
	query := "UPDATE users SET password_hash = $1 WHERE id = $2 AND gym_id = $3"
	_, err := r.db.Exec(query, newPasswordHash, userID, gymID)
	return err
}

func (r *usersRepository) DeleteUser(gymID, id string) error {
	query := "DELETE FROM users WHERE id = $1 AND gym_id = $2"
	_, err := r.db.Exec(query, id, gymID)
	return err
}

func (r *usersRepository) VerifyUser(gymID, userID string) error {
	query := "UPDATE users SET verified = true, updated_at = NOW() WHERE id = $1 AND gym_id = $2"
	result, err := r.db.Exec(query, userID, gymID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *usersRepository) SetUserActive(gymID, userID string, active bool) error {
	query := "UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2 AND gym_id = $3"
	result, err := r.db.Exec(query, active, userID, gymID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
