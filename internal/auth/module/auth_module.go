package module

import (
	"database/sql"
	"net/http"

	authhandler "github.com/alejandro-albiol/athenai/internal/auth/handler"
	authrepository "github.com/alejandro-albiol/athenai/internal/auth/repository"
	authrouter "github.com/alejandro-albiol/athenai/internal/auth/router"
	authservice "github.com/alejandro-albiol/athenai/internal/auth/service"
	gyminterfaces "github.com/alejandro-albiol/athenai/internal/gym/interfaces"
	userinterfaces "github.com/alejandro-albiol/athenai/internal/user/interfaces"
)

// NewAuthModule creates a new auth module with dependency injection and returns HTTP handler
func NewAuthModule(
	db *sql.DB,
	gymRepository gyminterfaces.GymRepository,
	userRepository userinterfaces.UserRepository,
	jwtSecret string,
) http.Handler {
	authRepository := authrepository.NewAuthRepository(db)
	authService := authservice.NewAuthService(
		authRepository,
		gymRepository,
		userRepository,
		jwtSecret,
	)
	authHandler := authhandler.NewAuthHandler(authService)
	return authrouter.NewAuthRouter(authHandler)
}
