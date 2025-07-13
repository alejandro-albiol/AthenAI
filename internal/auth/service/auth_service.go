package service

import (
	"time"

	"github.com/alejandro-albiol/athenai/internal/auth/dto"
	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository interfaces.AuthRepository
	jwtSecret  string
}

func NewAuthService(repository interfaces.AuthRepository, jwtSecret string) interfaces.AuthServiceInterface {
	return &AuthService{
		repository: repository,
		jwtSecret:  jwtSecret,
	}
}

// LoginAdmin handles platform admin authentication (athenai.com)
func (s *AuthService) LoginAdmin(credentials dto.LoginRequestDTO) (dto.LoginResponseDTO, error) {
	// Get admin from public.admin table
	admin, err := s.repository.GetAdminByUsername(credentials.Username)
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid credentials",
			nil,
		)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(credentials.Password)); err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid credentials",
			nil,
		)
	}

	// Check if admin is active
	if !admin.IsActive {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Account is inactive",
			nil,
		)
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":   admin.ID,
		"username":  admin.Username,
		"user_type": dto.UserTypePlatformAdmin,
		"is_active": admin.IsActive,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate access token",
			err,
		)
	}

	// Generate refresh token
	refreshClaims := jwt.MapClaims{
		"user_id":   admin.ID,
		"user_type": dto.UserTypePlatformAdmin,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate refresh token",
			err,
		)
	}

	// Store refresh token (no gym domain for platform admin)
	if err := s.repository.StoreRefreshToken(admin.ID, refreshTokenString, dto.UserTypePlatformAdmin, nil); err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to store refresh token",
			err,
		)
	}

	// Update last login
	if err := s.repository.UpdateAdminLastLogin(admin.ID); err != nil {
		// Log but don't fail the login
	}

	return dto.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Unix(claims["exp"].(int64), 0),
		UserInfo: dto.UserInfoDTO{
			ID:       admin.ID,
			Username: admin.Username,
			Email:    admin.Email,
			UserType: dto.UserTypePlatformAdmin,
			IsActive: &admin.IsActive,
		},
	}, nil
}

// LoginTenantUser handles tenant user authentication ({domain}.athenai.com)
func (s *AuthService) LoginTenantUser(gymDomain string, credentials dto.LoginRequestDTO) (dto.LoginResponseDTO, error) {
	// First, validate that the gym domain exists
	gymInfo, err := s.repository.ValidateGymDomain(gymDomain)
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeNotFound,
			"Gym domain not found",
			err,
		)
	}

	if !gymInfo.IsActive {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Gym is inactive",
			nil,
		)
	}

	// Get tenant user from tenant schema
	user, err := s.repository.GetTenantUserByUsername(gymDomain, credentials.Username)
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid credentials",
			nil,
		)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid credentials",
			nil,
		)
	}

	// Check if user is active
	if !user.IsActive {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Account is inactive",
			nil,
		)
	}

	// Validate role
	role := userrole_enum.UserRole(user.Role)
	if !isValidRole(role) {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid user role",
			nil,
		)
	}

	// Check verification limitations for guests and demo users
	if role == userrole_enum.Guest {
		if hasLimitations, reason := hasVerificationLimitations(user.VerificationStatus); hasLimitations {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeUnauthorized,
				reason,
				nil,
			)
		}
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":             user.ID,
		"username":            user.Username,
		"user_type":           dto.UserTypeTenantUser,
		"gym_domain":          gymDomain,
		"role":                user.Role,
		"verification_status": user.VerificationStatus,
		"iat":                 time.Now().Unix(),
		"exp":                 time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate access token",
			err,
		)
	}

	// Generate refresh token
	refreshClaims := jwt.MapClaims{
		"user_id":    user.ID,
		"user_type":  dto.UserTypeTenantUser,
		"gym_domain": gymDomain,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate refresh token",
			err,
		)
	}

	// Store refresh token with gym domain
	if err := s.repository.StoreRefreshToken(user.ID, refreshTokenString, dto.UserTypeTenantUser, &gymDomain); err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to store refresh token",
			err,
		)
	}

	// Update last login
	if err := s.repository.UpdateTenantUserLastLogin(gymDomain, user.ID); err != nil {
		// Log but don't fail the login
	}

	return dto.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Unix(claims["exp"].(int64), 0),
		UserInfo: dto.UserInfoDTO{
			ID:                 user.ID,
			Username:           user.Username,
			Email:              user.Email,
			UserType:           dto.UserTypeTenantUser,
			GymDomain:          &gymDomain,
			Role:               &user.Role,
			VerificationStatus: &user.VerificationStatus,
		},
	}, nil
}

// ValidateToken validates JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (dto.TokenValidationResponseDTO, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apierror.New(errorcode_enum.CodeUnauthorized, "Invalid signing method", nil)
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return dto.TokenValidationResponseDTO{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	if !token.Valid {
		return dto.TokenValidationResponseDTO{
			Valid:   false,
			Message: "Token is invalid",
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return dto.TokenValidationResponseDTO{
			Valid:   false,
			Message: "Invalid token claims",
		}, nil
	}

	// Convert claims to DTO
	claimsDTO := &dto.ClaimsDTO{
		UserID:    claims["user_id"].(string),
		Username:  claims["username"].(string),
		UserType:  dto.UserType(claims["user_type"].(string)),
		ExpiresAt: int64(claims["exp"].(float64)),
		IssuedAt:  int64(claims["iat"].(float64)),
	}

	// Add type-specific claims
	if claimsDTO.UserType == dto.UserTypePlatformAdmin {
		if isActive, exists := claims["is_active"]; exists {
			active := isActive.(bool)
			claimsDTO.IsActive = &active
		}
	} else {
		if gymDomain, exists := claims["gym_domain"]; exists {
			domain := gymDomain.(string)
			claimsDTO.GymDomain = &domain
		}
		if role, exists := claims["role"]; exists {
			roleStr := role.(string)
			claimsDTO.Role = &roleStr
		}
		if verificationStatus, exists := claims["verification_status"]; exists {
			status := verificationStatus.(string)
			claimsDTO.VerificationStatus = &status
		}
	}

	return dto.TokenValidationResponseDTO{
		Valid:   true,
		Claims:  claimsDTO,
		Message: "Token is valid",
	}, nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(refreshReq dto.RefreshTokenRequestDTO) (dto.LoginResponseDTO, error) {
	// Validate refresh token
	tokenData, err := s.repository.ValidateRefreshToken(refreshReq.RefreshToken)
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid refresh token",
			err,
		)
	}

	// Check expiration
	if time.Now().After(tokenData.ExpiresAt) {
		s.repository.RevokeRefreshToken(refreshReq.RefreshToken) // Clean up expired token
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Refresh token expired",
			nil,
		)
	}

	// Generate new access token based on user type
	var claims jwt.MapClaims
	var userInfo dto.UserInfoDTO

	if tokenData.UserType == dto.UserTypePlatformAdmin {
		admin, err := s.repository.GetAdminByID(tokenData.UserID)
		if err != nil {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"User not found",
				err,
			)
		}

		claims = jwt.MapClaims{
			"user_id":   admin.ID,
			"username":  admin.Username,
			"user_type": dto.UserTypePlatformAdmin,
			"is_active": admin.IsActive,
			"iat":       time.Now().Unix(),
			"exp":       time.Now().Add(24 * time.Hour).Unix(),
		}

		userInfo = dto.UserInfoDTO{
			ID:       admin.ID,
			Username: admin.Username,
			Email:    admin.Email,
			UserType: dto.UserTypePlatformAdmin,
			IsActive: &admin.IsActive,
		}
	} else {
		// For tenant users
		if tokenData.GymDomain == nil {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeInternal,
				"Missing gym domain in refresh token",
				nil,
			)
		}

		user, err := s.repository.GetTenantUserByID(*tokenData.GymDomain, tokenData.UserID)
		if err != nil {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"User not found",
				err,
			)
		}

		claims = jwt.MapClaims{
			"user_id":             user.ID,
			"username":            user.Username,
			"user_type":           dto.UserTypeTenantUser,
			"gym_domain":          *tokenData.GymDomain,
			"role":                user.Role,
			"verification_status": user.VerificationStatus,
			"iat":                 time.Now().Unix(),
			"exp":                 time.Now().Add(24 * time.Hour).Unix(),
		}

		userInfo = dto.UserInfoDTO{
			ID:                 user.ID,
			Username:           user.Username,
			Email:              user.Email,
			UserType:           dto.UserTypeTenantUser,
			GymDomain:          tokenData.GymDomain,
			Role:               &user.Role,
			VerificationStatus: &user.VerificationStatus,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate access token",
			err,
		)
	}

	return dto.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshReq.RefreshToken, // Keep same refresh token
		ExpiresAt:    time.Unix(claims["exp"].(int64), 0),
		UserInfo:     userInfo,
	}, nil
}

// Logout revokes refresh token
func (s *AuthService) Logout(logoutReq dto.LogoutRequestDTO) (dto.LogoutResponseDTO, error) {
	err := s.repository.RevokeRefreshToken(logoutReq.RefreshToken)
	if err != nil {
		return dto.LogoutResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to logout",
			err,
		)
	}

	return dto.LogoutResponseDTO{
		Success: true,
		Message: "Logout successful",
	}, nil
}

// HasPlatformAccess checks if user has platform admin access
func (s *AuthService) HasPlatformAccess(userID string) (bool, error) {
	admin, err := s.repository.GetAdminByID(userID)
	if err != nil {
		return false, err
	}

	return admin.IsActive && admin.ID == userID, nil
}

// HasTenantAccess checks if user has access to tenant with required role
func (s *AuthService) HasTenantAccess(userID, gymDomain string, requiredRole string) (bool, error) {
	user, err := s.repository.GetTenantUserByID(gymDomain, userID)
	if err != nil {
		return false, err
	}

	if user.ID != userID || !user.IsActive {
		return false, nil
	}

	role := userrole_enum.UserRole(user.Role)
	requiredRoleEnum := userrole_enum.UserRole(requiredRole)

	return hasPermission(role, requiredRoleEnum), nil
}

// Helper functions
func isValidRole(role userrole_enum.UserRole) bool {
	return role == userrole_enum.Admin || role == userrole_enum.User || role == userrole_enum.Guest
}

func hasVerificationLimitations(status string) (bool, string) {
	verificationStatus := userrole_enum.UserVerificationStatus(status)
	if verificationStatus == userrole_enum.Demo {
		// Demo users have time limitations
		return true, "Demo access expired"
	}
	return false, ""
}

func hasPermission(userRole, requiredRole userrole_enum.UserRole) bool {
	// Simple role hierarchy: admin > user > guest
	roleHierarchy := map[userrole_enum.UserRole]int{
		userrole_enum.Admin: 3,
		userrole_enum.User:  2,
		userrole_enum.Guest: 1,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}
