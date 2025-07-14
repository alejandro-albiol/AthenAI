package module

import (
	"database/sql"
	"net/http"
	"os"

	authhandler "github.com/alejandro-albiol/athenai/internal/auth/handler"
	authrepository "github.com/alejandro-albiol/athenai/internal/auth/repository"
	authrouter "github.com/alejandro-albiol/athenai/internal/auth/router"
	authservice "github.com/alejandro-albiol/athenai/internal/auth/service"
	gymrepository "github.com/alejandro-albiol/athenai/internal/gym/repository"
)

// NewAuthModule creates a new auth module with all dependencies wired and returns the router
func NewAuthModule(db *sql.DB) http.Handler {
	// Create auth repository
	authRepo := authrepository.NewAuthRepository(db)

	// Create gym repository (needed for gym lookups during login)
	gymRepo := gymrepository.NewGymRepository(db)

	// Get JWT secret from environment or use default
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-in-production" // Default for development
	}

	// Create service with both repositories
	service := authservice.NewAuthService(authRepo, gymRepo, jwtSecret)

	// Create handler
	handler := authhandler.NewAuthHandler(service)

	// Return router with all endpoints wired
	return authrouter.NewAuthRouter(handler)
}
