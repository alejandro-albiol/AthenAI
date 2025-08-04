package repository

import (
	"database/sql"
	"fmt"

	gyminterfaces "github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/lib/pq"
)

type usersRepository struct {
	db      *sql.DB
	gymRepo gyminterfaces.GymRepository
}

func NewUsersRepository(db *sql.DB, gymRepo gyminterfaces.GymRepository) interfaces.UserRepository {
	return &usersRepository{
		db:      db,
		gymRepo: gymRepo,
	}
}

func (r *usersRepository) CreateUser(gymID string, dto dto.UserCreationDTO) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return fmt.Errorf("failed to get gym: %w", err)
	}
	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf(`
        INSERT INTO %s (
            username, email, password_hash, role, gym_id, 
            is_verified, is_active, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, 
            false, true, NOW(), NOW()
        )`, tableName)

	_, err = r.db.Exec(query, dto.Username, dto.Email, dto.Password, dto.Role, gymID)
	return err
}

func (r *usersRepository) GetUserByID(gymID, id string) (dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return dto.UserResponseDTO{}, fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf(`
        SELECT id, username, email, password_hash, role, is_verified, is_active, gym_id, created_at, updated_at 
        FROM %s 
        WHERE id = $1 AND deleted_at IS NULL`, tableName)

	row := r.db.QueryRow(query, id)
	var user dto.UserResponseDTO
	var passwordHash string // We don't return this in the DTO
	err = row.Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role,
		&user.Verified, &user.IsActive, &user.GymID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	return user, nil
}

func (r *usersRepository) GetUserByUsername(gymID, username string) (dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return dto.UserResponseDTO{}, fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("SELECT id, username, email, password_hash, role, is_verified, is_active, gym_id, created_at, updated_at FROM %s WHERE username = $1 AND deleted_at IS NULL", tableName)
	row := r.db.QueryRow(query, username)
	var user dto.UserResponseDTO
	var passwordHash string // We don't return this in the DTO
	err = row.Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.Verified, &user.IsActive, &user.GymID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	return user, nil
}

func (r *usersRepository) GetUserByEmail(gymID, email string) (dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return dto.UserResponseDTO{}, fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("SELECT id, username, email, password_hash, role, is_verified, is_active, gym_id, created_at, updated_at FROM %s WHERE email = $1 AND deleted_at IS NULL", tableName)
	row := r.db.QueryRow(query, email)
	var user dto.UserResponseDTO
	var passwordHash string // We don't return this in the DTO
	err = row.Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.Verified, &user.IsActive, &user.GymID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}
	return user, nil
}

func (r *usersRepository) GetAllUsers(gymID string) ([]dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	users := make([]dto.UserResponseDTO, 0) // Initialize empty slice

	query := fmt.Sprintf(`
        SELECT id, username, email, password_hash, role, is_verified, is_active, gym_id, created_at, updated_at 
        FROM %s 
        WHERE deleted_at IS NULL`, tableName)

	rows, err := r.db.Query(query)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user dto.UserResponseDTO
		var passwordHash string // We don't return this in the DTO
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.Verified, &user.IsActive, &user.GymID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *usersRepository) GetPasswordHashByUsername(gymID, username string) (string, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return "", fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE username = $1 AND deleted_at IS NULL", tableName)
	row := r.db.QueryRow(query, username)
	var passwordHash string
	err = row.Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (r *usersRepository) UpdateUser(gymID string, id string, user dto.UserUpdateDTO) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("UPDATE %s SET username = $1, email = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL", tableName)
	_, err = r.db.Exec(query, user.Username, user.Email, id)
	return err
}

func (r *usersRepository) UpdatePassword(gymID, userID string, newPasswordHash string) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("UPDATE %s SET password_hash = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL", tableName)
	_, err = r.db.Exec(query, newPasswordHash, userID)
	return err
}

func (r *usersRepository) DeleteUser(gymID, id string) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	// Use soft delete
	query := fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", tableName)
	_, err = r.db.Exec(query, id)
	return err
}

func (r *usersRepository) VerifyUser(gymID, userID string) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("UPDATE %s SET is_verified = true, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL", tableName)
	result, err := r.db.Exec(query, userID)
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
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return fmt.Errorf("failed to get gym: %w", err)
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("UPDATE %s SET is_active = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL", tableName)
	result, err := r.db.Exec(query, active, userID)
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
