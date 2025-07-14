package middleware

import (
	"context"
	"net/http"

	"github.com/alejandro-albiol/athenai/pkg/apierror"
	errorcode_enum "github.com/alejandro-albiol/athenai/pkg/apierror/enum"
	"github.com/alejandro-albiol/athenai/pkg/response"
)

type contextKey string

const GymIDKey contextKey = "gymID"

// RequireGymID middleware that enforces X-Gym-ID header presence
func RequireGymID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gymID := r.Header.Get("X-Gym-ID")
		if gymID == "" {
			apiErr := apierror.New(
				errorcode_enum.CodeBadRequest,
				"X-Gym-ID header is required for tenant operations",
				nil,
			)
			response.WriteAPIError(w, apiErr)
			return
		}

		// Store gymID in context
		ctx := context.WithValue(r.Context(), GymIDKey, gymID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalGymID middleware that stores gym ID if present but doesn't enforce it
func OptionalGymID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gymID := r.Header.Get("X-Gym-ID")
		if gymID != "" {
			ctx := context.WithValue(r.Context(), GymIDKey, gymID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// GetGymID helper to get gym ID from request context or header
func GetGymID(r *http.Request) string {
	// First check context (set by middleware)
	if gymID, ok := r.Context().Value(GymIDKey).(string); ok {
		return gymID
	}
	// Fallback to header
	return r.Header.Get("X-Gym-ID")
}
