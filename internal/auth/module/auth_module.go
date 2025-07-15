package module

import (
	"database/sql"
	"net/http"
	"os"

	authhandler "github.com/alejandro-albiol/athenai/internal/auth/handler"
	"github.com/alejandro-albiol/athenai/internal/auth/interfaces"
	authrepository "github.com/alejandro-albiol/athenai/internal/auth/repository"
	authrouter "github.com/alejandro-albiol/athenai/internal/auth/router"
	authservice "github.com/alejandro-albiol/athenai/internal/auth/service"
	gymrepository "github.com/alejandro-albiol/athenai/internal/gym/repository"
)

// AuthModule holds the auth service and router
type AuthModule struct {
	Service interfaces.AuthServiceInterface
	Router  http.Handler
}

// NewAuthModule creates a new auth module with all dependencies wired and returns both service and router
func NewAuthModule(db *sql.DB) *AuthModule {
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

	// Create router with all endpoints wired
	router := authrouter.NewAuthRouter(handler)

	return &AuthModule{
		Service: service,
		Router:  router,
	}
}
