.PHONY: dev build test clean deps

# Development mode with live reload
dev:
	air

# Alternative dev without air
dev-simple:
	go run ./cmd/main.go

# Build the application
build:
	go build -o bin/athenai ./cmd

# Run tests
test:
	go test ./...

# Run specific module tests
test-gym:
	go test ./internal/gym/...

test-user:
	go test ./internal/user/...

# Clean build artifacts
clean:
	rm -rf bin/ tmp/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run with specific environment
dev-local:
	ENV=local PORT=8080 go run ./cmd/main.go

# Setup database (if you have migration scripts)
setup-db:
	go run ./cmd/setup-db/main.go
