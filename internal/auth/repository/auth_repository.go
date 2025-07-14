package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/auth/dto"
	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) interfaces.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

// GetAdminByUsername retrieves platform admin by username
func (r *AuthRepository) GetAdminByUsername(username string) (dto.AdminAuthDTO, error) {
	query := `
		SELECT id, username, password_hash, email, is_active, last_login_at, created_at
		FROM public.admin 
		WHERE username = $1 AND is_active = true
	`

	var admin dto.AdminAuthDTO
	err := r.db.QueryRow(query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Email,
		&admin.IsActive,
		&admin.LastLoginAt,
		&admin.CreatedAt,
	)

	if err != nil {
		return dto.AdminAuthDTO{}, err
	}

	return admin, nil
}

// GetAdminByID retrieves platform admin by ID (for refresh token validation)
func (r *AuthRepository) GetAdminByID(adminID string) (dto.AdminAuthDTO, error) {
	query := `
		SELECT id, username, password_hash, email, is_active, last_login_at, created_at
		FROM public.admin 
		WHERE id = $1 AND is_active = true
	`

	var admin dto.AdminAuthDTO
	err := r.db.QueryRow(query, adminID).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Email,
		&admin.IsActive,
		&admin.LastLoginAt,
		&admin.CreatedAt,
	)

	if err != nil {
		return dto.AdminAuthDTO{}, err
	}

	return admin, nil
}

// UpdateAdminLastLogin updates platform admin last login timestamp
func (r *AuthRepository) UpdateAdminLastLogin(adminID string) error {
	query := `
		UPDATE public.admin 
		SET last_login_at = NOW() 
		WHERE id = $1
	`

	_, err := r.db.Exec(query, adminID)
	return err
}

// StoreRefreshToken stores refresh token in auth-specific table
func (r *AuthRepository) StoreRefreshToken(userID, token string, userType dto.UserType, gymID *string) error {
	query := `
		INSERT INTO public.refresh_tokens (user_id, token, user_type, gym_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4, NOW() + INTERVAL '7 days', NOW())
		ON CONFLICT (user_id, user_type, gym_id) 
		DO UPDATE SET 
			token = EXCLUDED.token,
			expires_at = EXCLUDED.expires_at,
			created_at = EXCLUDED.created_at
	`

	_, err := r.db.Exec(query, userID, token, userType, gymID)
	return err
}

// ValidateRefreshToken checks if refresh token exists and is not expired
func (r *AuthRepository) ValidateRefreshToken(token string) (dto.RefreshTokenDTO, error) {
	query := `
		SELECT token, user_id, user_type, gym_id, expires_at, created_at
		FROM public.refresh_tokens 
		WHERE token = $1 AND expires_at > NOW()
	`

	var refreshToken dto.RefreshTokenDTO
	err := r.db.QueryRow(query, token).Scan(
		&refreshToken.Token,
		&refreshToken.UserID,
		&refreshToken.UserType,
		&refreshToken.GymID,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		return dto.RefreshTokenDTO{}, err
	}

	return refreshToken, nil
}

// RevokeRefreshToken removes specific refresh token (for logout)
func (r *AuthRepository) RevokeRefreshToken(token string) error {
	query := `
		DELETE FROM public.refresh_tokens 
		WHERE token = $1
	`

	_, err := r.db.Exec(query, token)
	return err
}

// RevokeAllUserTokens removes all refresh tokens for a user (for security/admin actions)
func (r *AuthRepository) RevokeAllUserTokens(userID string, userType dto.UserType) error {
	query := `
		DELETE FROM public.refresh_tokens 
		WHERE user_id = $1 AND user_type = $2
	`

	_, err := r.db.Exec(query, userID, userType)
	return err
}

// CleanupExpiredTokens removes expired refresh tokens (for maintenance)
func (r *AuthRepository) CleanupExpiredTokens() error {
	query := `
		DELETE FROM public.refresh_tokens 
		WHERE expires_at <= NOW()
	`

	_, err := r.db.Exec(query)
	return err
}

// LogLogin records login attempt in login_history table
func (r *AuthRepository) LogLogin(userID string, userType dto.UserType, gymID *string, success bool, ipAddress string) error {
	query := `
		INSERT INTO public.login_history (user_id, user_type, gym_id, success, ip_address, attempted_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := r.db.Exec(query, userID, userType, gymID, success, ipAddress)
	return err
}
