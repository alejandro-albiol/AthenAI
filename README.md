# AthenAI

A Go-based API application with integrated frontend serving, designed for gymnasiums

## 🏗️ Project Structure

```
athenai/
├── cmd/                    # Application entry points
│   ├── main.go            # Main server with API and frontend serving
│   └── setup-db/          # Database setup utilities
├── internal/              # Private application modules
│   ├── admin/             # Admin management (equipment, exercises, muscle groups)
│   ├── auth/              # Authentication module (JWT, login, tokens)
│   ├── database/          # Database connections and utilities
│   ├── gym/               # Gym management module
│   └── user/              # User management module
├── api/                   # API configuration and setup
│   ├── api.go            # Main API module setup
│   └── swagger.go        # Swagger/OpenAPI documentation setup
├── frontend/              # Static frontend files (served by backend)
│   ├── index.html        # Main frontend entry point
│   ├── css/              # Stylesheets
│   ├── js/               # JavaScript files
│   └── assets/           # Static assets (images, fonts, etc.)
├── docs/                  # Documentation
│   ├── openapi/          # OpenAPI/Swagger specifications
│   └── *.md              # Project documentation
├── pkg/                   # Public packages
│   ├── apierror/         # API error handling
│   ├── middleware/       # Custom middleware
│   ├── response/         # API response utilities
│   └── utils/            # Utility functions
├── config/               # Configuration management
├── scripts/              # Utility scripts
└── .air.toml            # Air live reload configuration
```

## 🚀 Development Setup

### Prerequisites
- Go 1.19 or higher
- PostgreSQL database
- Air (for live reload development)

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd athenai

# Install dependencies
go mod tidy
go mod download

# Install Air for live reload (if not already installed)
go install github.com/air-verse/air@latest
```

### Environment Setup
1. Copy `example.env` to `.env`
2. Configure your database connection and other environment variables

## 🛠️ Development Commands

Since this is a Windows environment without `make`, use these Go commands directly:

### Start Development Server (with live reload)
```bash
air
```

### Run Without Live Reload
```bash
go run ./cmd/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific module tests
go test ./internal/gym/...
go test ./internal/user/...

# Run tests with coverage
go test -cover ./...
```

### Build Application
```bash
go build -o bin/athenai ./cmd
```

### Database Setup
```bash
go run ./cmd/setup-db/main.go
```

## 🌐 API Endpoints

The server runs on `http://localhost:8080` and provides:

- **Frontend**: `http://localhost:8080/` (serves static files from `/frontend`)
- **API**: `http://localhost:8080/api/v1/` (RESTful API endpoints)
- **Documentation**: `http://localhost:8080/swagger-ui/` (Interactive API documentation)

### Available API Routes
- **Authentication**: `/api/v1/auth/*` - Login, logout, token validation, refresh tokens
- **Users**: `/api/v1/user/*` - User management endpoints
- **Gyms**: `/api/v1/gym/*` - Gym management endpoints

## 🔐 Authentication Architecture

### Subdomain-Based Multi-Tenancy
AthenAI uses subdomain-based authentication to automatically route users to the correct authentication context:

- **Platform Admin**: `athenai.com` → Authenticates against `public.admin` table
- **Tenant Users**: `{gym_domain}.athenai.com` → Authenticates against tenant-specific schema

### User Types
1. **Platform Administrators** (`public.admin`)
   - Full platform access across all gyms
   - Manage global data (exercises, equipment, muscle groups)
   - Single user type with active/inactive status

2. **Tenant Users** (per gym schema: `{gym_domain}.users`)
   - **Admin**: Gym owners/managers with full gym access
   - **User**: Gym clients who use the services
   - **Guest**: Demo users with limited trial access

### Authentication Flow
1. User visits subdomain (e.g., `mygym.athenai.com`)
2. Middleware extracts gym domain from subdomain
3. System validates gym domain exists in database
4. Authentication routes to appropriate user table
5. JWT tokens include user type and domain context
6. All subsequent requests use domain-aware authorization

### Module Structure
```
internal/auth/
├── dto/                   # Data Transfer Objects
│   ├── login.dto.go      # Authentication request/response DTOs
│   ├── token.dto.go      # Token validation and refresh DTOs
│   └── repository.dto.go # Database layer DTOs
├── interfaces/           # Service contracts
│   ├── auth_handler.interface.go
│   ├── auth_service.interface.go
│   └── auth_repository.interface.go
├── service/             # Business logic
│   └── auth_service.go  # JWT generation, validation, user authentication
├── handler/             # HTTP layer (to be implemented)
└── repository/          # Data access layer (to be implemented)
```

## 🏛️ Architecture Decisions

### Monolithic Frontend + Backend
We chose to keep the frontend and backend in the same repository because:

1. **Simplicity**: Single deployment, single repository to manage
2. **No Framework**: Since we're not using a frontend framework, there's no complex build process
3. **Shared Configuration**: Environment variables and configuration are shared
4. **Go Static Serving**: Go's `http.FileServer` efficiently serves static files
5. **Development Speed**: Faster iteration without managing multiple repositories

### Module Structure
Each business domain (user, gym, admin) is organized as a module with:
- **DTOs**: Data Transfer Objects for API communication
- **Handlers**: HTTP request/response handling
- **Services**: Business logic implementation
- **Repositories**: Data access layer
- **Interfaces**: Contract definitions for loose coupling

## 📝 API Documentation

The project uses OpenAPI 3.0 specifications with modular YAML files:

- Main spec: `/docs/openapi/openapi.yaml`
- Components: `/docs/openapi/components/`
- Path definitions: `/docs/openapi/paths/`

Access the interactive documentation at: `http://localhost:8080/swagger-ui/`

## 🔧 Configuration

The application uses environment variables for configuration. Key variables:

- `PORT`: Server port (default: 8080)
- Database connection variables (see `example.env`)
- Environment-specific settings

## 🧪 Testing Strategy

- **Unit Tests**: Each module has comprehensive test coverage
- **Integration Tests**: Repository tests with real database connections
- **Handler Tests**: HTTP endpoint testing with mocked services
