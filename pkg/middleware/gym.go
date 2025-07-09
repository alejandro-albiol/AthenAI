package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/pkg/response"
)

type contextKey string

const GymIDKey contextKey = "gymID"

func RequireGymID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gymID := r.Header.Get("X-Gym-ID")
		if gymID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.APIResponse[any]{
				Status:  "error",
				Message: "X-Gym-ID header is required",
				Data:    nil,
			})
			return
		}

		// Store gymID in context
		ctx := context.WithValue(r.Context(), GymIDKey, gymID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetGymID helper to get gym ID from request context
func GetGymID(r *http.Request) string {
	if gymID, ok := r.Context().Value(GymIDKey).(string); ok {
		return gymID
	}
	// Fallback to header for backward compatibility
	return r.Header.Get("X-Gym-ID")
}
