package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alejandro-albiol/athenai/api"
	"github.com/alejandro-albiol/athenai/config"
	"github.com/alejandro-albiol/athenai/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT undefined, using 8080")
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
		log.Println("HOST undefined, using localhost")
	}

	// Initialize database connection
	log.Printf("🗄️  Connecting to database...")
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Printf("✅ Database connection established")

	// Setup database tables
	log.Printf("🔧 Setting up database tables...")
	err = database.CreatePublicTables(db)
	if err != nil {
		log.Fatalf("❌ Failed to create database tables: %v", err)
	}
	log.Printf("✅ Database tables ready")

	// Root router
	rootRouter := chi.NewRouter()

	// Add middleware
	rootRouter.Use(middleware.Logger)
	rootRouter.Use(middleware.Recoverer)
	log.Printf("🛡️  Middleware configured (Logger, Recoverer)")

	// Setup Swagger at root level
	api.SetupSwagger(rootRouter)
	log.Printf("📚 Swagger UI available at: http://localhost:%s/swagger-ui/", port)

	// Mount API under /api/v1
	rootRouter.Mount("/api/v1", api.NewAPIModule(db))
	log.Printf("🔌 API mounted at: http://localhost:%s/api/v1", port)

	// Serve static frontend files
	workDir, _ := os.Getwd()
	frontendDir := http.Dir(filepath.Join(workDir, "frontend"))
	FileServer(rootRouter, "/", frontendDir)
	log.Printf("🎨 Frontend served at: http://localhost:%s/", port)

	log.Printf("Server is running on: http://%s:%s", host, port)
	log.Printf("Documentation: http://%s:%s/swagger-ui/", host, port)
	log.Fatal(http.ListenAndServe(":"+port, rootRouter))
}

// FileServer sets up a http.FileServer handler to serve static files from a directory.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
