package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	authdto "github.com/alejandro-albiol/athenai/internal/auth/dto"
	authinterfaces "github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	gyminterfaces "github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
)

type AuthService struct {
	authRepo  authinterfaces.AuthRepositoryInterface
	gymRepo   gyminterfaces.GymRepository
	jwtSecret string
}

func NewAuthService(
	authRepo authinterfaces.AuthRepositoryInterface,
	gymRepo gyminterfaces.GymRepository,
	jwtSecret string,
) authinterfaces.AuthServiceInterface {
	return &AuthService{
		authRepo:  authRepo,
		gymRepo:   gymRepo,
		jwtSecret: jwtSecret,
	}
}

// Login handles the simplified login logic based on X-Gym-ID header
func (s *AuthService) Login(r *http.Request, loginReq authdto.LoginRequestDTO) (*authdto.LoginResponseDTO, *apierror.APIError) {
	// Check for X-Gym-ID header
	gymID := r.Header.Get("X-Gym-ID")

	if gymID != "" {
		// Tenant user login - lookup gym and authenticate in tenant schema
		return s.loginTenantUser(gymID, loginReq)
	} else {
		// Platform admin login - authenticate in public.admin table
		return s.loginPlatformAdmin(loginReq)
	}
}

// loginPlatformAdmin handles platform admin authentication
func (s *AuthService) loginPlatformAdmin(loginReq authdto.LoginRequestDTO) (*authdto.LoginResponseDTO, *apierror.APIError) {
	// Authenticate against public.admin table
	admin, err := s.authRepo.AuthenticatePlatformAdmin(loginReq.Email, loginReq.Password)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid admin credentials",
			err,
		)
	}

	// Generate JWT token for platform admin
	token, err := s.generateJWT(admin.ID, "platform_admin", admin.Username, nil, nil)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate token",
			err,
		)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(admin.ID, "platform_admin", nil)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate refresh token",
			err,
		)
	}

	return &authdto.LoginResponseDTO{
		AccessToken:  token,
		RefreshToken: refreshToken,
		UserInfo: authdto.UserInfoDTO{
			UserID:   admin.ID,
			Username: admin.Username,
			Email:    admin.Email,
			UserType: "platform_admin",
			Role:     nil, // Platform admins don't have roles
			GymID:    nil, // Platform admins are not tied to a specific gym
		},
	}, nil
}

// loginTenantUser handles tenant user authentication
func (s *AuthService) loginTenantUser(gymID string, loginReq authdto.LoginRequestDTO) (*authdto.LoginResponseDTO, *apierror.APIError) {
	// First, lookup the gym to get its domain
	gym, err := s.gymRepo.GetGymByID(gymID)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeNotFound,
			"Gym not found",
			err,
		)
	}

	// Check if gym is active
	if !gym.IsActive {
		return nil, apierror.New(
			errorcode_enum.CodeForbidden,
			"Gym is not active",
			nil,
		)
	}

	// Authenticate against {domain}.users table
	user, err := s.authRepo.AuthenticateTenantUser(gym.ID, loginReq.Email, loginReq.Password)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid user credentials",
			err,
		)
	}

	// Generate JWT token for tenant user
	token, err := s.generateJWT(user.ID, "tenant_user", user.Username, &user.Role, &user.GymID)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate token",
			err,
		)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(user.ID, "tenant_user", &user.GymID)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to generate refresh token",
			err,
		)
	}

	return &authdto.LoginResponseDTO{
		AccessToken:  token,
		RefreshToken: refreshToken,
		UserInfo: authdto.UserInfoDTO{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			UserType: "tenant_user",
			Role:     &user.Role,
			GymID:    &user.GymID,
		},
	}, nil
}

// generateJWT creates a JWT token with the provided claims
func (s *AuthService) generateJWT(userID, userType, username string, role, gymID *string) (string, error) {
	claims := authdto.ClaimsDTO{
		UserID:   userID,
		UserType: userType,
		Username: username,
		Role:     role,
		GymID:    gymID,
		IsActive: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken creates a refresh token and stores it in the database
func (s *AuthService) generateRefreshToken(userID, userType string, gymID *string) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	tokenString := hex.EncodeToString(tokenBytes)

	// Store in database (store the token directly, not hashed)
	refreshTokenDTO := &authdto.RefreshTokenDTO{
		Token:     tokenString,
		UserID:    userID,
		UserType:  userType,
		GymID:     gymID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
		CreatedAt: time.Now(),
	}

	err = s.authRepo.StoreRefreshToken(refreshTokenDTO)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (*authdto.TokenValidationResponseDTO, *apierror.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid token",
			err,
		)
	}

	var claims authdto.ClaimsDTO
	if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Unmarshal custom claims from MapClaims
		if v, ok := tokenClaims["user_id"]; ok {
			claims.UserID = fmt.Sprintf("%v", v)
		}
		if v, ok := tokenClaims["username"]; ok {
			claims.Username = fmt.Sprintf("%v", v)
		}
		if v, ok := tokenClaims["user_type"]; ok {
			claims.UserType = fmt.Sprintf("%v", v)
		}
		if v, ok := tokenClaims["is_active"]; ok {
			switch val := v.(type) {
			case bool:
				claims.IsActive = val
			case string:
				claims.IsActive = val == "true"
			}
		}
		if v, ok := tokenClaims["role"]; ok {
			str := fmt.Sprintf("%v", v)
			claims.Role = &str
		}
		if v, ok := tokenClaims["gym_id"]; ok {
			str := fmt.Sprintf("%v", v)
			claims.GymID = &str
		}
		// RegisteredClaims
		if v, ok := tokenClaims["exp"]; ok {
			switch val := v.(type) {
			case float64:
				claims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(val), 0))
			case int64:
				claims.ExpiresAt = jwt.NewNumericDate(time.Unix(val, 0))
			case json.Number:
				i, _ := val.Int64()
				claims.ExpiresAt = jwt.NewNumericDate(time.Unix(i, 0))
			}
		}
		if v, ok := tokenClaims["iat"]; ok {
			switch val := v.(type) {
			case float64:
				claims.IssuedAt = jwt.NewNumericDate(time.Unix(int64(val), 0))
			case int64:
				claims.IssuedAt = jwt.NewNumericDate(time.Unix(val, 0))
			case json.Number:
				i, _ := val.Int64()
				claims.IssuedAt = jwt.NewNumericDate(time.Unix(i, 0))
			}
		}
		return &authdto.TokenValidationResponseDTO{
			Valid:  true,
			Claims: claims,
		}, nil
	}

	return nil, apierror.New(
		errorcode_enum.CodeUnauthorized,
		"Invalid token claims",
		nil,
	)
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(refreshReq authdto.RefreshTokenRequestDTO) (*authdto.LoginResponseDTO, *apierror.APIError) {
	// Validate the refresh token directly (no hashing needed)
	tokenData, err := s.authRepo.ValidateRefreshToken(refreshReq.RefreshToken)
	if err != nil {
		return nil, apierror.New(
			errorcode_enum.CodeUnauthorized,
			"Invalid refresh token",
			err,
		)
	}

	// Determine user type and get user info
	switch tokenData.UserType {
	case "platform_admin":
		// Get admin info by ID (not by authentication)
		admin, err := s.authRepo.GetPlatformAdminByID(tokenData.UserID)
		if err != nil {
			return nil, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"User not found",
				err,
			)
		}

		// Generate new JWT
		newToken, err := s.generateJWT(admin.ID, "platform_admin", admin.Username, nil, nil)
		if err != nil {
			return nil, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to generate new token",
				err,
			)
		}

		return &authdto.LoginResponseDTO{
			AccessToken:  newToken,
			RefreshToken: refreshReq.RefreshToken, // Keep same refresh token
			UserInfo: authdto.UserInfoDTO{
				UserID:   admin.ID,
				Username: admin.Username,
				Email:    admin.Email,
				UserType: "platform_admin",
				Role:     nil,
				GymID:    nil,
			},
		}, nil

	case "tenant_user":
		// Get gym domain for tenant user
		gym, err := s.gymRepo.GetGymByID(*tokenData.GymID)
		if err != nil {
			return nil, apierror.New(
				errorcode_enum.CodeNotFound,
				"Gym not found",
				err,
			)
		}

		// Get tenant user info by ID
		user, err := s.authRepo.GetTenantUserByID(gym.ID, tokenData.UserID)
		if err != nil {
			return nil, apierror.New(
				errorcode_enum.CodeUnauthorized,
				"User not found",
				err,
			)
		}

		// Generate new JWT
		newToken, err := s.generateJWT(user.ID, "tenant_user", user.Username, &user.Role, &user.GymID)
		if err != nil {
			return nil, apierror.New(
				errorcode_enum.CodeInternal,
				"Failed to generate new token",
				err,
			)
		}

		return &authdto.LoginResponseDTO{
			AccessToken:  newToken,
			RefreshToken: refreshReq.RefreshToken, // Keep same refresh token
			UserInfo: authdto.UserInfoDTO{
				UserID:   user.ID,
				Username: user.Username,
				Email:    user.Email,
				UserType: "tenant_user",
				Role:     &user.Role,
				GymID:    &user.GymID,
			},
		}, nil
	}

	return nil, apierror.New(
		errorcode_enum.CodeBadRequest,
		"Invalid user type in refresh token",
		nil,
	)
}

// Logout revokes a refresh token
func (s *AuthService) Logout(logoutReq authdto.LogoutRequestDTO) *apierror.APIError {
	// Revoke the refresh token directly (no hashing needed)
	err := s.authRepo.RevokeRefreshToken(logoutReq.RefreshToken)
	if err != nil {
		return apierror.New(
			errorcode_enum.CodeInternal,
			"Failed to logout",
			err,
		)
	}

	return nil
}
