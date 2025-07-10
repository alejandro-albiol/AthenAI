package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupSwagger configures and adds Swagger documentation routes
func SetupSwagger(r chi.Router) {
	// Serve OpenAPI specification (YAML)
	r.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, filepath.Join("docs", "openapi", "openapi.yaml"))
	})

	// Serve all OpenAPI files (components, paths, etc.)
	fileServer := http.FileServer(http.Dir("docs/openapi"))
	r.Get("/swagger/*", http.StripPrefix("/swagger", fileServer).ServeHTTP)

	// Serve Swagger UI
	r.Get("/swagger", http.RedirectHandler("/swagger-ui/", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/swagger-ui/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/openapi.yaml"),
	))
}
