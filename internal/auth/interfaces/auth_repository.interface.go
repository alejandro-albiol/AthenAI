package interfaces

import "github.com/alejandro-albiol/athenai/internal/auth/dto"

// AuthRepository handles authentication-specific data operations
type AuthRepository interface {
	// Platform admin operations
	GetAdminByUsername(username string) (dto.AdminAuthDTO, error)
	GetAdminByID(adminID string) (dto.AdminAuthDTO, error)
	UpdateAdminLastLogin(adminID string) error

	// Token operations
	StoreRefreshToken(userID, token string, userType dto.UserType, gymDomain *string) error
	ValidateRefreshToken(token string) (dto.RefreshTokenDTO, error)
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID string, userType dto.UserType) error
	CleanupExpiredTokens() error

	// Login history
	LogLogin(userID string, userType dto.UserType, gymDomain *string, success bool, ipAddress string) error
}
