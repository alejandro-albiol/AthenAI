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

const (
	UserIDKey   contextKey = "userID"
	UserTypeKey contextKey = "userType"
	UserRoleKey contextKey = "userRole"
)

// AuthMiddleware handles JWT token validation and authorization in one place
func AuthMiddleware(authService interfaces.AuthServiceInterface, requiredAccess string) func(http.Handler) http.Handler {
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

			// Handle authorization based on requiredAccess
			userType := string(validationResponse.Claims.UserType)

			switch requiredAccess {
			case "platform_admin":
				// Only platform admins allowed
				if userType != "platform_admin" {
					apiErr := apierror.New(
						errorcode_enum.CodeForbidden,
						"Platform admin access required",
						nil,
					)
					response.WriteAPIError(w, apiErr)
					return
				}

			case "gym_operations":
				// Platform admins (super admin mode) OR gym users with X-Gym-ID
				gymID := r.Header.Get("X-Gym-ID")

				switch userType {
case "platform_admin":
					// Platform admins can access any gym if they provide gym ID
					if gymID == "" {
						apiErr := apierror.New(
							errorcode_enum.CodeBadRequest,
							"X-Gym-ID header required for gym operations",
							nil,
						)
						response.WriteAPIError(w, apiErr)
						return
					}
					ctx = context.WithValue(ctx, GymIDKey, gymID)
					// Platform admin acts as super admin in this gym

				case "tenant_user":
					// Regular gym users need gym ID and must match their gym
					if gymID == "" {
						apiErr := apierror.New(
							errorcode_enum.CodeBadRequest,
							"X-Gym-ID header required for gym operations",
							nil,
						)
						response.WriteAPIError(w, apiErr)
						return
					}
					ctx = context.WithValue(ctx, GymIDKey, gymID)
					// TODO: Validate that tenant user's gym matches the X-Gym-ID header

				default:
					apiErr := apierror.New(
						errorcode_enum.CodeForbidden,
						"Gym access required",
						nil,
					)
					response.WriteAPIError(w, apiErr)
					return
				}

			case "gym_admin":
				// Platform admins (super admin mode) OR gym admins with X-Gym-ID
				gymID := r.Header.Get("X-Gym-ID")
				if gymID == "" {
					apiErr := apierror.New(
						errorcode_enum.CodeBadRequest,
						"X-Gym-ID header required for admin operations",
						nil,
					)
					response.WriteAPIError(w, apiErr)
					return
				}

				ctx = context.WithValue(ctx, GymIDKey, gymID)

				switch userType {
case "platform_admin":
					// Platform admin acts as super admin in any gym

				case "tenant_user":
					// Regular gym user must have admin role
					userRole := ""
					if validationResponse.Claims.Role != nil {
						userRole = *validationResponse.Claims.Role
					}
					if userRole != "admin" {
						apiErr := apierror.New(
							errorcode_enum.CodeForbidden,
							"Gym admin role required",
							nil,
						)
						response.WriteAPIError(w, apiErr)
						return
					}
					// TODO: Validate that tenant user's gym matches the X-Gym-ID header

				default:
					apiErr := apierror.New(
						errorcode_enum.CodeForbidden,
						"Admin access required",
						nil,
					)
					response.WriteAPIError(w, apiErr)
					return
				}
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
