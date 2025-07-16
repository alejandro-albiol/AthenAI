package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserTypeKey contextKey = "userType"
	UserRoleKey contextKey = "userRole"
	GymIDKey    contextKey = "gymID"
)

// AuthMiddleware handles JWT token validation and puts user context in request
func AuthMiddleware(authService interfaces.AuthServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				apiErr := apierror.New(
					errorcode_enum.CodeUnauthorized,
					"Authorization header required",
					nil,
				)
				response.WriteAPIError(w, apiErr)
				return
			}

			// Extract token from "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				apiErr := apierror.New(
					errorcode_enum.CodeUnauthorized,
					"Invalid authorization header format",
					nil,
				)
				response.WriteAPIError(w, apiErr)
				return
			}

			token := tokenParts[1]
			validationResponse, err := authService.ValidateToken(token)
			if err != nil {
				apiErr := apierror.New(
					errorcode_enum.CodeUnauthorized,
					"Invalid or expired token",
					err,
				)
				response.WriteAPIError(w, apiErr)
				return
			}

			if !validationResponse.Valid {
				apiErr := apierror.New(
					errorcode_enum.CodeUnauthorized,
					"Token validation failed",
					nil,
				)
				response.WriteAPIError(w, apiErr)
				return
			}

			// Store user information in context
			ctx := context.WithValue(r.Context(), UserIDKey, validationResponse.Claims.UserID)
			ctx = context.WithValue(ctx, UserTypeKey, string(validationResponse.Claims.UserType))
			if validationResponse.Claims.Role != nil {
				ctx = context.WithValue(ctx, UserRoleKey, *validationResponse.Claims.Role)
			}

			// Store gym ID from JWT token (not from header for security)
			if validationResponse.Claims.GymID != nil {
				ctx = context.WithValue(ctx, GymIDKey, *validationResponse.Claims.GymID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID helper to get user ID from request context
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserType helper to get user type from request context
func GetUserType(r *http.Request) string {
	if userType, ok := r.Context().Value(UserTypeKey).(string); ok {
		return userType
	}
	return ""
}

// GetUserRole helper to get user role from request context
func GetUserRole(r *http.Request) string {
	if userRole, ok := r.Context().Value(UserRoleKey).(string); ok {
		return userRole
	}
	return ""
}

// GetGymID helper to get gym ID from request context
func GetGymID(r *http.Request) string {
	if gymID, ok := r.Context().Value(GymIDKey).(string); ok {
		return gymID
	}
	return ""
}

// GetGymDomain helper to securely get gym domain from JWT gym ID
func GetGymDomain(r *http.Request, authService interfaces.AuthServiceInterface) (string, error) {
	gymID := GetGymID(r)
	if gymID == "" {
		return "", nil // Platform admin or no gym context
	}
	return authService.GetGymDomain(gymID)
}

// ValidateGymAccess ensures the user has access to the requested gym
func ValidateGymAccess(r *http.Request, requestedGymID string) bool {
	userType := GetUserType(r)

	// Platform admins have access to all gyms
	if userType == "platform_admin" {
		return true
	}

	// Tenant users can only access their own gym
	if userType == "tenant_user" {
		userGymID := GetGymID(r)
		return userGymID == requestedGymID
	}

	return false
}

// IsPlatformAdmin checks if the current user is a platform admin
func IsPlatformAdmin(r *http.Request) bool {
	return GetUserType(r) == "platform_admin"
}

// IsGymAdmin checks if the current user is a gym admin (platform admin OR tenant user with admin role)
func IsGymAdmin(r *http.Request) bool {
	userType := GetUserType(r)
	if userType == "platform_admin" {
		return true
	}
	if userType == "tenant_user" {
		return GetUserRole(r) == "admin"
	}
	return false
}
