package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
	"github.com/go-chi/chi/v5"
)

const (
	UserIDKey   contextKey = "userID"
	UserTypeKey contextKey = "userType"
	UserRoleKey contextKey = "userRole"
)

// AuthMiddleware handles JWT token validation and user context
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

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePlatformAdmin middleware that only allows platform admins
func RequirePlatformAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, ok := r.Context().Value(UserTypeKey).(string)
		if !ok || userType != "platform_admin" {
			apiErr := apierror.New(
				errorcode_enum.CodeForbidden,
				"Platform admin access required",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireGymAdmin middleware that only allows gym admins for the specified gym
func RequireGymAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, ok := r.Context().Value(UserTypeKey).(string)
		if !ok || userType != "tenant_user" {
			apiErr := apierror.New(
				errorcode_enum.CodeForbidden,
				"Gym access required",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		userRole, ok := r.Context().Value(UserRoleKey).(string)
		if !ok || userRole != "admin" {
			apiErr := apierror.New(
				errorcode_enum.CodeForbidden,
				"Gym admin access required",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		// Ensure the gym ID in context matches the X-Gym-ID header
		gymID := GetGymID(r)
		if gymID == "" {
			apiErr := apierror.New(
				errorcode_enum.CodeBadRequest,
				"X-Gym-ID header required for gym operations",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireGymAccess middleware that allows platform admins OR gym users for tenant operations
// Platform admins can access any tenant, gym users can only access their own tenant
func RequireGymAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, ok := r.Context().Value(UserTypeKey).(string)
		if !ok {
			apiErr := apierror.New(
				errorcode_enum.CodeForbidden,
				"Authentication required",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		// All tenant operations require gym ID (even for platform admins)
		gymID := GetGymID(r)
		if gymID == "" {
			apiErr := apierror.New(
				errorcode_enum.CodeBadRequest,
				"X-Gym-ID header required for tenant operations",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		// Platform admins can access any tenant
		if userType == "platform_admin" {
			next.ServeHTTP(w, r)
			return
		}

		// Tenant users must be accessing their own tenant
		if userType != "tenant_user" {
			apiErr := apierror.New(
				errorcode_enum.CodeForbidden,
				"Tenant access required",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		// TODO: Add validation that tenant user's gym matches the X-Gym-ID header
		// This would require checking the user's gym_id claim against the header

		next.ServeHTTP(w, r)
	})
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

// PlatformAdminRoutes creates a router that requires platform admin authentication
func PlatformAdminRoutes(authService interfaces.AuthServiceInterface) chi.Router {
	r := chi.NewRouter()
	r.Use(AuthMiddleware(authService))
	r.Use(RequirePlatformAdmin)
	return r
}

// GymAdminRoutes creates a router that requires gym admin authentication + X-Gym-ID header
func GymAdminRoutes(authService interfaces.AuthServiceInterface) chi.Router {
	r := chi.NewRouter()
	r.Use(OptionalGymID) // Extract gym ID from header
	r.Use(AuthMiddleware(authService))
	r.Use(RequireGymAdmin)
	return r
}

// GymUserRoutes creates a router that requires any gym user authentication + X-Gym-ID header
func GymUserRoutes(authService interfaces.AuthServiceInterface) chi.Router {
	r := chi.NewRouter()
	r.Use(OptionalGymID) // Extract gym ID from header
	r.Use(AuthMiddleware(authService))
	r.Use(RequireGymAccess)
	return r
}
