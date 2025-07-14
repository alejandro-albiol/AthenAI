package api

import (
	"database/sql"
	"os"

	authmodule "github.com/alejandro-albiol/athenai/internal/auth/module"
	authrepository "github.com/alejandro-albiol/athenai/internal/auth/repository"
	authservice "github.com/alejandro-albiol/athenai/internal/auth/service"
	gymmodule "github.com/alejandro-albiol/athenai/internal/gym/module"
	gymrepository "github.com/alejandro-albiol/athenai/internal/gym/repository"
	usermodule "github.com/alejandro-albiol/athenai/internal/user/module"
	userrepository "github.com/alejandro-albiol/athenai/internal/user/repository"
	"github.com/alejandro-albiol/athenai/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func NewAPIModule(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default for development - should be in .env
	}

	// Create shared repositories and services for authentication
	gymRepo := gymrepository.NewGymRepository(db)
	userRepo := userrepository.NewUsersRepository(db)
	authRepo := authrepository.NewAuthRepository(db)
	authService := authservice.NewAuthService(authRepo, gymRepo, userRepo, jwtSecret)

	// PUBLIC ROUTES (no authentication required)
	r.Mount("/auth", authmodule.NewAuthModule(db, gymRepo, userRepo, jwtSecret))

	// PROTECTED ROUTES
	r.Route("/", func(r chi.Router) {
		// Add authentication middleware for all protected routes
		r.Use(middleware.AuthMiddleware(authService))

		// PLATFORM ADMIN ROUTES (JWT required, no gym ID needed)
		r.Route("/gym", func(r chi.Router) {
			r.Use(middleware.RequirePlatformAdmin)
			r.Mount("/", gymmodule.NewGymModule(db))
		})

		// USER ROUTES (Platform admin OR Gym admin with X-Gym-ID)
		r.Route("/user", func(r chi.Router) {
			// Extract gym ID from header first
			r.Use(middleware.OptionalGymID)
			// Custom middleware that allows either:
			// 1. Platform admin (no gym ID required)
			// 2. Gym admin with X-Gym-ID header
			r.Use(middleware.RequireGymAccess)
			r.Mount("/", usermodule.NewUserModule(db))
		})
	})

	return r
}
