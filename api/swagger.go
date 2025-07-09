package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupSwagger configures and adds Swagger documentation routes
func SetupSwagger(r chi.Router) {
	// Serve OpenAPI specification
	r.Get("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("docs", "openapi", "openapi.yaml"))
	})

	// Serve Swagger UI
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/swagger.yaml"),
	))
}
