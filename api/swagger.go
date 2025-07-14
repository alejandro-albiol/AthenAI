package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupSwagger configures and adds Swagger documentation routes
func SetupSwagger(r chi.Router) {
	// Serve OpenAPI specification (YAML) — no cache
	r.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, filepath.Join("docs", "openapi", "openapi.yaml"))
	})

	// Serve all OpenAPI files (components, paths, etc.) with no-cache headers
	r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		// Add no-cache headers for all swagger files
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		fileServer := http.FileServer(http.Dir("docs/openapi"))
		http.StripPrefix("/swagger", fileServer).ServeHTTP(w, r)
	})

	// Serve Swagger UI, force reload of YAML with timestamp version
	version := fmt.Sprintf("%d", time.Now().Unix())
	swaggerURL := fmt.Sprintf("/swagger/openapi.yaml?v=%s", version)

	r.Get("/swagger", http.RedirectHandler("/swagger-ui/", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/swagger-ui/*", httpSwagger.Handler(
		httpSwagger.URL(swaggerURL),
	))
}
