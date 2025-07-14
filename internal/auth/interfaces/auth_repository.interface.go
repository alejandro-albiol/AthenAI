package interfaces

import "github.com/alejandro-albiol/athenai/internal/auth/dto"

// AuthRepositoryInterface handles only authentication-specific database operations
// Gym operations are handled by GymRepository, User operations by UserRepository
type AuthRepositoryInterface interface {
	// Authentication methods
	AuthenticatePlatformAdmin(email, password string) (*dto.AdminAuthDTO, error)
	AuthenticateTenantUser(domain, email, password string) (*dto.TenantUserAuthDTO, error)

	// Refresh token operations
	StoreRefreshToken(token *dto.RefreshTokenDTO) error
	ValidateRefreshToken(tokenHash string) (*dto.RefreshTokenDTO, error)
	RevokeRefreshToken(tokenHash string) error
	RevokeAllUserTokens(userID, userType string) error
}
