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
	// Load configuration
	cfg := config.Load()

	log.Printf("🚀 Starting AthenAI Server...")
	log.Printf("📊 Environment: %s", cfg.AppEnv)
	log.Printf("🌐 Host: %s", "localhost")
	log.Printf("🔌 Port: %s", cfg.Port)

	if cfg.IsDevelopment() {
		log.Printf("🔧 Development mode: Detailed error logging enabled")
	} else {
		log.Printf("🔒 Production mode: Secure error logging enabled")
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
	log.Printf("📚 Swagger UI available at: http://localhost:%s/swagger-ui/", cfg.Port)

	// Mount API under /api/v1
	rootRouter.Mount("/api/v1", api.NewAPIModule(db))
	log.Printf("🔌 API mounted at: http://localhost:%s/api/v1", cfg.Port)

	// Serve static frontend files
	workDir, _ := os.Getwd()
	frontendDir := http.Dir(filepath.Join(workDir, "frontend"))
	FileServer(rootRouter, "/", frontendDir)
	log.Printf("🎨 Frontend served at: http://localhost:%s/", cfg.Port)

	log.Printf("🎯 Available endpoints:")
	log.Printf("   📖 Frontend:     http://localhost:%s/", cfg.Port)
	log.Printf("   🔌 API:          http://localhost:%s/api/v1", cfg.Port)
	log.Printf("   📚 Swagger:      http://localhost:%s/swagger-ui/", cfg.Port)
	log.Printf("   🔐 Auth:         http://localhost:%s/api/v1/auth", cfg.Port)
	log.Printf("   🏋️  Gym:          http://localhost:%s/api/v1/gym", cfg.Port)
	log.Printf("   👤 Users:        http://localhost:%s/api/v1/user", cfg.Port)

	log.Printf("")
	log.Printf("🌟 AthenAI Server is running!")
	log.Printf("🌐 Server URL: http://localhost:%s", cfg.Port)
	if cfg.IsDevelopment() {
		log.Printf("🔗 Tenant testing: Use X-Gym-ID header with requests")
		log.Printf("📋 Example: curl -H 'X-Gym-ID: your-gym-id' http://localhost:%s/api/v1/user", cfg.Port)
	}
	log.Printf("📡 Press Ctrl+C to stop the server")
	log.Printf("")

	// Bind to 0.0.0.0 to accept requests from any hostname
	serverAddr := "0.0.0.0:" + cfg.Port
	log.Fatal(http.ListenAndServe(serverAddr, rootRouter))
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
