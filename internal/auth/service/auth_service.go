package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alejandro-albiol/athenai/internal/auth/dto"
	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	gyminterfaces "github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	userrole_enum "github.com/alejandro-albiol/athenai/internal/user/enum"
	userinterfaces "github.com/alejandro-albiol/athenai/internal/user/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository interfaces.AuthRepository
	gymRepository  gyminterfaces.GymRepository
	userRepository userinterfaces.UserRepository
	jwtSecret      string
}

func NewAuthService(
	authRepository interfaces.AuthRepository,
	gymRepository gyminterfaces.GymRepository,
	userRepository userinterfaces.UserRepository,
	jwtSecret string,
) interfaces.AuthServiceInterface {
	return &AuthService{
		authRepository: authRepository,
		gymRepository:  gymRepository,
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

// LoginAdmin handles platform admin authentication (athenai.com)
func (s *AuthService) LoginAdmin(credentials dto.LoginRequestDTO) (dto.LoginResponseDTO, error) {
	// Get admin from public.admin table
	admin, err := s.authRepository.GetAdminByUsername(credentials.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"Invalid credentials",
				nil,
			)
		}
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to validate credentials",
			err,
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
		"exp":       time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenJWT.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate refresh token",
			err,
		)
	}

	// Store refresh token in database
	if err := s.authRepository.StoreRefreshToken(admin.ID, refreshToken, dto.UserTypePlatformAdmin, nil); err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to store refresh token",
			err,
		)
	}

	// Update last login timestamp
	if err := s.authRepository.UpdateAdminLastLogin(admin.ID); err != nil {
		// Log error but don't fail the login
		// TODO: Add proper logging here
	}

	// Create user info for response
	isActive := admin.IsActive
	userInfo := dto.UserInfoDTO{
		ID:       admin.ID,
		Username: admin.Username,
		Email:    admin.Email,
		UserType: dto.UserTypePlatformAdmin,
		IsActive: &isActive,
	}

	return dto.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		UserInfo:     userInfo,
	}, nil
}

// LoginTenantUser handles tenant user authentication (with gym ID)
func (s *AuthService) LoginTenantUser(gymID string, credentials dto.LoginRequestDTO) (dto.LoginResponseDTO, error) {
	// First, validate that the gym ID exists
	gymInfo, err := s.gymRepository.GetGymByID(gymID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeNotFound,
				"Gym not found",
				nil,
			)
		}
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to validate gym domain",
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

	// Get tenant user from user repository
	user, err := s.userRepository.GetUserByUsername(gymInfo.ID, credentials.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"Invalid credentials",
				nil,
			)
		}
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to validate credentials",
			err,
		)
	}

	// Get password hash for verification
	passwordHash, err := s.userRepository.GetPasswordHashByUsername(gymInfo.ID, credentials.Username)
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to validate credentials",
			err,
		)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(credentials.Password)); err != nil {
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
	role := user.Role
	if !isValidRole(role) {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid user role",
			nil,
		)
	}

	// Check verification status for business logic
	// For guest users, you might want to add restrictions
	if role == userrole_enum.Guest && !user.Verified {
		// Add any guest user restrictions here if needed
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_type": dto.UserTypeTenantUser,
		"gym_id":    gymID,
		"role":      string(user.Role),
		"verified":  user.Verified,
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
		"user_id":   user.ID,
		"user_type": dto.UserTypeTenantUser,
		"gym_id":    gymID,
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

	// Store refresh token with gym ID
	if err := s.authRepository.StoreRefreshToken(user.ID, refreshTokenString, dto.UserTypeTenantUser, &gymID); err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to store refresh token",
			err,
		)
	}

	// Update last login - we'll skip this for now as it's not critical
	// This could be added to user repository if needed in the future

	// Create response
	roleStr := string(user.Role)
	verifiedStr := "verified"
	if !user.Verified {
		verifiedStr = "unverified"
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
			GymID:              &gymID,
			Role:               &roleStr,
			VerificationStatus: &verifiedStr,
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
		if gymID, exists := claims["gym_id"]; exists {
			id := gymID.(string)
			claimsDTO.GymID = &id
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
	tokenData, err := s.authRepository.ValidateRefreshToken(refreshReq.RefreshToken)
	if err != nil {
		return dto.LoginResponseDTO{}, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid refresh token",
			err,
		)
	}

	// Check expiration
	if time.Now().After(tokenData.ExpiresAt) {
		s.authRepository.RevokeRefreshToken(refreshReq.RefreshToken) // Clean up expired token
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
		admin, err := s.authRepository.GetAdminByID(tokenData.UserID)
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
		if tokenData.GymID == nil {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeInternal,
				"Missing gym ID in refresh token",
				nil,
			)
		}

		// Get gym info first to validate
		gymInfo, err := s.gymRepository.GetGymByID(*tokenData.GymID)
		if err != nil {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"Gym not found",
				err,
			)
		}

		user, err := s.userRepository.GetUserByID(gymInfo.ID, tokenData.UserID)
		if err != nil {
			return dto.LoginResponseDTO{}, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"User not found",
				err,
			)
		}

		claims = jwt.MapClaims{
			"user_id":   user.ID,
			"username":  user.Username,
			"user_type": dto.UserTypeTenantUser,
			"gym_id":    *tokenData.GymID,
			"role":      string(user.Role),
			"verified":  user.Verified,
			"iat":       time.Now().Unix(),
			"exp":       time.Now().Add(24 * time.Hour).Unix(),
		}

		roleStr := string(user.Role)
		verifiedStr := "verified"
		if !user.Verified {
			verifiedStr = "unverified"
		}

		userInfo = dto.UserInfoDTO{
			ID:                 user.ID,
			Username:           user.Username,
			Email:              user.Email,
			UserType:           dto.UserTypeTenantUser,
			GymID:              tokenData.GymID,
			Role:               &roleStr,
			VerificationStatus: &verifiedStr,
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
	err := s.authRepository.RevokeRefreshToken(logoutReq.RefreshToken)
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
	admin, err := s.authRepository.GetAdminByID(userID)
	if err != nil {
		return false, err
	}

	return admin.IsActive && admin.ID == userID, nil
}

// HasTenantAccess checks if user has access to tenant with required role
func (s *AuthService) HasTenantAccess(userID, gymID string, requiredRole string) (bool, error) {
	// Get gym info first
	gymInfo, err := s.gymRepository.GetGymByID(gymID)
	if err != nil {
		return false, err
	}

	user, err := s.userRepository.GetUserByID(gymInfo.ID, userID)
	if err != nil {
		return false, err
	}

	if user.ID != userID || !user.IsActive {
		return false, nil
	}

	requiredRoleEnum := userrole_enum.UserRole(requiredRole)

	return hasPermission(user.Role, requiredRoleEnum), nil
}

// Helper functions
func isValidRole(role userrole_enum.UserRole) bool {
	return role == userrole_enum.Admin || role == userrole_enum.User || role == userrole_enum.Guest
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
