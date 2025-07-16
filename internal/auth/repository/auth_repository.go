package repository

import (
	"database/sql"
	"errors"

	"github.com/alejandro-albiol/athenai/internal/auth/dto"
	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) interfaces.AuthRepositoryInterface {
	return &AuthRepository{db: db}
}

// AuthenticatePlatformAdmin authenticates against public.admin table
func (r *AuthRepository) AuthenticatePlatformAdmin(username, password string) (*dto.AdminAuthDTO, error) {
	query := `
		SELECT id, username, email, password_hash, is_active 
		FROM public.admin 
		WHERE username = $1 AND is_active IS true
	`

	var admin dto.AdminAuthDTO
	var passwordHash string

	err := r.db.QueryRow(query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&passwordHash,
		&admin.IsActive,
	)

	if err != nil {
		return nil, err
	}

	if !admin.IsActive {
		return nil, errors.New("admin account is not active")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &admin, nil
}

// AuthenticateTenantUser authenticates against {domain}.users table
func (r *AuthRepository) AuthenticateTenantUser(domain, email, password string) (*dto.TenantUserAuthDTO, error) {
	// Construct the table name dynamically
	tableName := pq.QuoteIdentifier(domain) + ".user"

	query := `
		SELECT id, username, email, password_hash, role, is_verified, is_active, gym_id
		FROM ` + tableName + `
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user dto.TenantUserAuthDTO
	var passwordHash string

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.GymID,
	)

	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// StoreRefreshToken stores a refresh token for logout management
func (r *AuthRepository) StoreRefreshToken(token *dto.RefreshTokenDTO) error {
	query := `
		INSERT INTO public.refresh_token (token, user_id, user_type, gym_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := r.db.Exec(query,
		token.Token,
		token.UserID,
		token.UserType,
		token.GymID,
		token.ExpiresAt,
	)

	return err
}

// ValidateRefreshToken validates a refresh token
func (r *AuthRepository) ValidateRefreshToken(tokenHash string) (*dto.RefreshTokenDTO, error) {
	query := `
		SELECT token, user_id, user_type, gym_id, expires_at, created_at
		FROM public.refresh_token 
		WHERE token = $1 AND expires_at > NOW()
	`

	var token dto.RefreshTokenDTO
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.Token,
		&token.UserID,
		&token.UserType,
		&token.GymID,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	return &token, err
}

// RevokeRefreshToken removes a refresh token (for logout)
func (r *AuthRepository) RevokeRefreshToken(tokenHash string) error {
	query := `DELETE FROM public.refresh_token WHERE token = $1`
	_, err := r.db.Exec(query, tokenHash)
	return err
}

// RevokeAllUserTokens removes all refresh tokens for a user (for logout all devices)
func (r *AuthRepository) RevokeAllUserTokens(userID, userType string) error {
	query := `DELETE FROM public.refresh_token WHERE user_id = $1 AND user_type = $2`
	_, err := r.db.Exec(query, userID, userType)
	return err
}

// GetPlatformAdminByID retrieves admin by ID for refresh token validation
func (r *AuthRepository) GetPlatformAdminByID(adminID string) (*dto.AdminAuthDTO, error) {
	query := `
		SELECT id, username, email, is_active 
		FROM public.admin 
		WHERE id = $1 AND is_active IS true
	`

	var admin dto.AdminAuthDTO

	err := r.db.QueryRow(query, adminID).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&admin.IsActive,
	)

	if err != nil {
		return nil, err
	}

	if !admin.IsActive {
		return nil, errors.New("admin account is not active")
	}

	return &admin, nil
}

// GetTenantUserByID retrieves tenant user by ID for refresh token validation
func (r *AuthRepository) GetTenantUserByID(domain, userID string) (*dto.TenantUserAuthDTO, error) {
	// Construct the table name dynamically
	tableName := pq.QuoteIdentifier(domain) + ".user"

	query := `
		SELECT id, username, email, role, is_verified, is_active, gym_id
		FROM ` + tableName + `
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user dto.TenantUserAuthDTO

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.GymID,
	)

	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	return &user, nil
}
