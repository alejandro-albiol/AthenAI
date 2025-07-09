package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/alejandro-albiol/athenai/pkg/response"
)

// RequireGymID middleware ensures X-Gym-ID header is present
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
		next.ServeHTTP(w, r)
	})
}

// GetGymID helper to get gym ID from request
func GetGymID(r *http.Request) string {
	return r.Header.Get("X-Gym-ID")
}
