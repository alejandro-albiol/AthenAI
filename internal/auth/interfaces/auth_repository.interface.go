package interfaces

import "github.com/alejandro-albiol/athenai/internal/auth/dto"

// AuthRepository handles authentication data operations
type AuthRepository interface {
	// Domain validation
	ValidateGymDomain(domain string) (dto.GymDomainDTO, error)

	// Platform admin operations
	GetAdminByUsername(username string) (dto.AdminAuthDTO, error)
	GetAdminByID(adminID string) (dto.AdminAuthDTO, error)
	UpdateAdminLastLogin(adminID string) error

	// Tenant user operations
	GetTenantUserByUsername(gymDomain, username string) (dto.TenantUserAuthDTO, error)
	GetTenantUserByID(gymDomain, userID string) (dto.TenantUserAuthDTO, error)
	UpdateTenantUserLastLogin(gymDomain, userID string) error

	// Token operations
	StoreRefreshToken(userID, token string, userType dto.UserType, gymDomain *string) error
	ValidateRefreshToken(token string) (dto.RefreshTokenDTO, error)
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID string, userType dto.UserType) error
	CleanupExpiredTokens() error
}
