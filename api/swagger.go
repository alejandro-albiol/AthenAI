package api

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupSwagger configures and adds Swagger documentation routes
func SetupSwagger(r chi.Router) {
	// Serve Swagger UI
	r.Get("/swagger", http.RedirectHandler("/swagger/", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("../swagger/doc.json"), // Use relative path
	))

	// Serve OpenAPI specification and its dependencies
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, filepath.Join("docs", "openapi", "openapi.yaml"))
	})

	// Serve component files
	fileServer := http.FileServer(http.Dir("docs/openapi"))
	r.Get("/swagger/components/*", http.StripPrefix("/swagger", fileServer).ServeHTTP)
	r.Get("/swagger/paths/*", http.StripPrefix("/swagger", fileServer).ServeHTTP)
}
