# AthenAI

**Multi-tenant gym management platform with AI-powered workout generation**

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/alejandro-albiol/athenai.git
cd athenai

# Setup environment
cp example.env .env
# Edit .env with your database configuration

# Install dependencies and start development server
go mod tidy
air  # or: go run ./cmd/main.go
```

**Access the application:**

- 🌐 **Frontend**: http://localhost:8080
- 📡 **API**: http://localhost:8080/api/v1
- 📚 **API Docs**: http://localhost:8080/swagger-ui

## 🏗️ What is AthenAI?

AthenAI is a comprehensive gym management platform featuring:

- **🏢 Multi-Tenant Architecture**: Complete data isolation per gym
- **🤖 AI Workout Generation**: Smart workout creation based on user goals
- **👥 User Management**: Role-based access for admins, trainers, and members
- **💪 Exercise Library**: Extensive catalog with custom gym additions
- **📊 Workout Tracking**: Complete workout history and progress monitoring
- **🔐 Secure Authentication**: JWT-based auth with tenant isolation

## � Prerequisites

- **Go**: 1.19 or higher
- **PostgreSQL**: Database for multi-tenant data storage
- **Air** (optional): For live reload development

```bash
# Install Air for live reload
go install github.com/air-verse/air@latest
```

## ⚡ Development Workflow

### Environment Setup

```bash
# Database setup
go run ./cmd/setup-db/main.go

# Create superadmin (platform administrator)
go run ./cmd/setup-superadmin/main.go
```

### Development Commands

```bash
# Development with live reload
air

# Production build
go build -o bin/athenai ./cmd

# Run tests
go test ./...

# Test with coverage
go test -cover ./...
```

### Project Structure

```
├── internal/           # Business modules (auth, gym, user, exercises, etc.)
├── api/               # API setup and routing
├── cmd/               # Application entry points and utilities
├── frontend/          # Static web files (HTML, CSS, JS)
├── docs/              # 📚 Comprehensive documentation
├── pkg/               # Shared utilities and middleware
└── config/            # Configuration management
```

## 🏛️ Architecture Overview

```
┌─────────────────────────────────────────┐
│             AthenAI Platform            │
├─────────────────────────────────────────┤
│  Platform Layer (public schema)        │
│  • System Administration               │
│  • Global Exercise Library             │
│  • Shared Templates & Equipment        │
├─────────────────────────────────────────┤
│  Tenant Layer ({gym_uuid} schemas)     │
│  • Gym-specific Users & Data           │
│  • Custom Exercises & Equipment        │
│  • Workout Instances & Tracking        │
└─────────────────────────────────────────┘
```

**Key Features:**

- 🔐 **Schema-level tenant isolation** for complete data security
- 🎯 **Role-based access control** (Platform Admin, Gym Admin, Trainer, Member)
- 🔄 **Modular architecture** following consistent Go patterns
- 📊 **PostgreSQL multi-tenancy** with shared and isolated data

## 📖 Documentation

**Complete documentation available in [`/docs`](./docs/README.md)**

| Topic                                                     | Description                                         |
| --------------------------------------------------------- | --------------------------------------------------- |
| [📚 Documentation Hub](./docs/README.md)                  | Complete guide and navigation                       |
| [🏗️ Backend Architecture](./docs/backend-architecture.md) | Modules, patterns, and system design                |
| [🗄️ Database Design](./docs/database-design.md)           | Schema, relationships, and multi-tenancy            |
| [🔐 Security Model](./docs/security-model.md)             | Authentication, authorization, and tenant isolation |
| [⚙️ Module Pattern](./docs/module-pattern.md)             | Standard patterns for creating new modules          |
| [🔧 Configuration](./docs/configuration.md)               | Environment setup and deployment                    |
| [📡 API Reference](./docs/openapi/openapi.yaml)           | Complete OpenAPI specification                      |

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Specific module tests
go test ./internal/auth/...
go test ./internal/gym/...

# Integration tests
go test ./tests/...

# Coverage report
go test -cover ./...
```

## 🤝 Contributing

1. **Follow the Module Pattern**: Use the standard structure in [`docs/module-pattern.md`](./docs/module-pattern.md)
2. **Security First**: Review [`docs/security-model.md`](./docs/security-model.md) for security guidelines
3. **Update Documentation**: Keep docs in sync with code changes
4. **Test Coverage**: Include tests for new modules and features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**🚀 Ready to start?** Check out the [complete documentation](./docs/README.md) for detailed guides and architecture information.
