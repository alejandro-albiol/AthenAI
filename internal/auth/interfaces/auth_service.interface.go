package interfaces

import "github.com/alejandro-albiol/athenai/internal/auth/dto"

// AuthService interface defines authentication business logic
type AuthServiceInterface interface {
	// Platform admin authentication (admin.athenai.com)
	LoginAdmin(credentials dto.LoginRequestDTO) (dto.LoginResponseDTO, error)

	// Tenant user authentication (extracted from subdomain)
	LoginTenantUser(gymDomain string, credentials dto.LoginRequestDTO) (dto.LoginResponseDTO, error)

	// Token operations
	ValidateToken(token string) (dto.TokenValidationResponseDTO, error)
	RefreshToken(refreshReq dto.RefreshTokenRequestDTO) (dto.LoginResponseDTO, error)
	Logout(logoutReq dto.LogoutRequestDTO) (dto.LogoutResponseDTO, error)

	// Authorization checks
	HasPlatformAccess(userID string) (bool, error)
	HasTenantAccess(userID, gymDomain string, requiredRole string) (bool, error)
}
