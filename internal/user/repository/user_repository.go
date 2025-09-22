package repository

import (
	"database/sql"
	"errors"
	"fmt"

	gyminterfaces "github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	"github.com/alejandro-albiol/athenai/internal/user/dto"
	"github.com/lib/pq"
)

type UserRepository struct {
	db      *sql.DB
	gymRepo gyminterfaces.GymRepository
}

func NewUsersRepository(db *sql.DB, gymRepo gyminterfaces.GymRepository) *UserRepository {
	return &UserRepository{
		db:      db,
		gymRepo: gymRepo,
	}
}

func (r *UserRepository) CreateUser(gymID string, dto *dto.UserCreationDTO) (*string, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return nil, err
	}
	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	// Set appropriate default values based on user role
	trainingPhase := dto.TrainingPhase
	motivation := dto.Motivation
	specialSituation := dto.SpecialSituation

	// For non-member users, always set valid defaults for member-specific fields
	if dto.Role != "member" {
		// Set defaults for admin, trainer, guest, and other roles
		if trainingPhase == "" {
			trainingPhase = "maintenance"
		}
		if motivation == "" {
			motivation = "self_improvement"
		}
		if specialSituation == "" {
			specialSituation = "none"
		}
	} else {
		// For members, ensure we have valid values (required by business logic)
		if trainingPhase == "" {
			trainingPhase = "maintenance"
		}
		if motivation == "" {
			motivation = "self_improvement"
		}
		if specialSituation == "" {
			specialSituation = "none"
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			username, email, password_hash, role,
			is_verified, is_active, description, training_phase, motivation, special_situation,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			true, true, $5, $6, $7, $8,
			NOW(), NOW()
		) RETURNING id`, tableName)

	var userID string
	err = r.db.QueryRow(query, dto.Username, dto.Email, dto.Password, dto.Role,
		dto.Description, trainingPhase, motivation, specialSituation).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &userID, nil
}

func (r *UserRepository) GetUserByID(gymID, id string) (*dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return nil, err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, role, is_verified, is_active,
			   description, training_phase, motivation, special_situation,
			   created_at, updated_at 
		FROM %s 
		WHERE id = $1 AND is_active = TRUE`, tableName)

	row := r.db.QueryRow(query, id)
	user := &dto.UserResponseDTO{}
	var passwordHash string
	err = row.Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role,
		&user.Verified, &user.IsActive,
		&user.Description, &user.TrainingPhase, &user.Motivation, &user.SpecialSituation,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(gymID, username string) (*dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return nil, err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, role, is_verified, is_active,
			   description, training_phase, motivation, special_situation,
			   created_at, updated_at 
		FROM %s 
		WHERE username = $1 AND is_active = TRUE`, tableName)
	row := r.db.QueryRow(query, username)
	user := &dto.UserResponseDTO{}
	var passwordHash string // We don't return this in the DTO
	err = row.Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.Verified, &user.IsActive,
		&user.Description, &user.TrainingPhase, &user.Motivation, &user.SpecialSituation,
		&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(gymID, email string) (*dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return nil, err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("SELECT id, username, email, password_hash, role, is_verified, is_active, description, training_phase, motivation, special_situation, created_at, updated_at FROM %s WHERE email = $1 AND is_active = TRUE", tableName)
	row := r.db.QueryRow(query, email)
	user := &dto.UserResponseDTO{}
	var passwordHash string // We don't return this in the DTO
	err = row.Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.Verified, &user.IsActive,
		&user.Description, &user.TrainingPhase, &user.Motivation, &user.SpecialSituation,
		&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAllUsers(gymID string) ([]*dto.UserResponseDTO, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		// If gym doesn't exist, return empty users list instead of error
		if errors.Is(err, sql.ErrNoRows) {
			return []*dto.UserResponseDTO{}, nil
		}
		return nil, err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	users := make([]*dto.UserResponseDTO, 0) // Initialize empty slice

	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, role, is_verified, is_active, description, training_phase, motivation, special_situation, created_at, updated_at 
		FROM %s 
		WHERE is_active = TRUE`, tableName)

	rows, err := r.db.Query(query)
	if err != nil {
		// If the table doesn't exist (which can happen for new gyms), return empty list
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &dto.UserResponseDTO{}
		var passwordHash string // We don't return this in the DTO
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role,
			&user.Verified, &user.IsActive,
			&user.Description, &user.TrainingPhase, &user.Motivation, &user.SpecialSituation,
			&user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetPasswordHashByUsername(gymID, username string) (string, error) {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return "", err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE username = $1 AND is_active = TRUE", tableName)
	row := r.db.QueryRow(query, username)
	var passwordHash string
	err = row.Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (r *UserRepository) UpdateUser(gymID string, id string, user *dto.UserUpdateDTO) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	setClause := "SET updated_at = NOW()"
	params := []interface{}{id} // Start with ID for WHERE clause
	paramCount := 1

	if user.Username != "" {
		paramCount++
		setClause += fmt.Sprintf(", username = $%d", paramCount)
		params = append(params, user.Username)
	}
	if user.Email != "" {
		paramCount++
		setClause += fmt.Sprintf(", email = $%d", paramCount)
		params = append(params, user.Email)
	}
	if user.Description != nil {
		paramCount++
		setClause += fmt.Sprintf(", description = $%d", paramCount)
		params = append(params, user.Description)
	}
	if user.TrainingPhase != nil {
		paramCount++
		setClause += fmt.Sprintf(", training_phase = $%d", paramCount)
		params = append(params, user.TrainingPhase)
	}
	if user.Motivation != nil {
		paramCount++
		setClause += fmt.Sprintf(", motivation = $%d", paramCount)
		params = append(params, user.Motivation)
	}
	if user.SpecialSituation != nil {
		paramCount++
		setClause += fmt.Sprintf(", special_situation = $%d", paramCount)
		params = append(params, user.SpecialSituation)
	}

	query := fmt.Sprintf("UPDATE %s %s WHERE id = $1 AND is_active = TRUE", tableName, setClause)
	_, err = r.db.Exec(query, params...)
	return err
}

func (r *UserRepository) UpdatePassword(gymID, userID string, newPasswordHash string) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("UPDATE %s SET password_hash = $1, updated_at = NOW() WHERE id = $2 AND is_active = TRUE", tableName)
	_, err = r.db.Exec(query, newPasswordHash, userID)
	return err
}

func (r *UserRepository) DeleteUser(gymID, id string) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	// Use soft delete
	query := fmt.Sprintf("UPDATE %s SET is_active = FALSE WHERE id = $1 AND is_active = TRUE", tableName)
	_, err = r.db.Exec(query, id)
	return err
}

func (r *UserRepository) VerifyUser(gymID, userID string) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("UPDATE %s SET is_verified = true, updated_at = NOW() WHERE id = $1 AND is_active = TRUE", tableName)
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

func (r *UserRepository) SetUserActive(gymID, userID string, active bool) error {
	// Get gym domain to construct the correct schema table name
	gym, err := r.gymRepo.GetGymByID(gymID)
	if err != nil {
		return err
	}

	// Construct tenant-specific table name
	tableName := pq.QuoteIdentifier(gym.ID) + ".user"

	query := fmt.Sprintf("UPDATE %s SET is_active = $1, updated_at = NOW() WHERE id = $2 AND is_active = TRUE", tableName)
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
