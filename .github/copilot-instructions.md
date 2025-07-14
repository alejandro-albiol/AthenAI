# AthenAI Copilot Instructions

## Architecture Overview

AthenAI is a **multi-tenant gym management platform** with a Go backend serving both API and static frontend. Key architectural decisions:

### Module-Based Domain Architecture

```
internal/
├── admin/      # Global data (exercises, equipment, muscle groups) - public schema
├── auth/       # Authentication & JWT tokens - handles both platform and tenant users
├── gym/        # Tenant management & multi-tenancy - public schema
├── user/       # Gym-specific users - tenant schemas (gym_domain.users)
└── database/   # Connection utilities & tenant schema creation
```

Each module follows **strict layering**: `repository → service → handler` with interfaces for testability:

- **Repository**: Handles database operations, returns raw SQL errors
- **Service**: Transforms repository data into business logic, maps SQL errors to APIError
- **Handler**: Reads service responses and converts to standardized HTTP responses

### Multi-Tenancy Strategy

- **Public schema**: Gyms, global exercises, equipment (admin module)
- **Tenant schemas**: Each gym gets `gym_domain.*` schema for isolated user data
- **Schema creation**: Automatic tenant schema creation in `GymService.CreateGym()`

### Authentication Architecture

**Header-Based Authentication**: Single login endpoint with automatic routing:

- **No `X-Gym-ID` header** → Platform admin (`public.admin` table)
- **With `X-Gym-ID` header** → Tenant users (`{gym_domain}.users` table)

**User Types**:

- **Platform Admins**: Full platform access, manage global data
- **Tenant Users**: Gym-scoped with roles (admin, user, guest) and verification status (verified, unverified, demo)

**Authentication Module Structure**:

```go
internal/auth/
├── dto/
│   ├── login.dto.go      # LoginRequestDTO, LoginResponseDTO, UserInfoDTO
│   ├── token.dto.go      # TokenValidationResponseDTO, ClaimsDTO, RefreshTokenRequestDTO
│   └── repository.dto.go # AdminAuthDTO, TenantUserAuthDTO, RefreshTokenDTO
├── handler/              # Standard response patterns with error code enums
├── interfaces/           # AuthService, AuthRepository, AuthHandler contracts
├── module/               # Self-contained dependency injection
├── repository/           # Pure authentication operations
├── router/               # RESTful endpoints: /login, /refresh, /logout, /validate
└── service/              # Business logic orchestration with gym repository
```

**Key Patterns**:

- **Single endpoint**: `/auth/login` routes based on `X-Gym-ID` header presence
- **Repository delegation**: Auth service uses gym repository for lookups
- **Self-contained modules**: Each module manages its own dependencies
- **JWT + Refresh tokens**: 24h access tokens, 30d refresh tokens
- **Standard responses**: Uses `response.WriteAPIError()` and `response.WriteAPISuccess()`
  ├── service/ # JWT generation, validation, subdomain routing
  ├── handler/ # HTTP layer with subdomain-aware routing (to be implemented)
  └── repository/ # Database operations for both admin and tenant users (to be implemented)

````

**Key Patterns**:

- Subdomain middleware extracts gym domain automatically
- Single `/auth/login` endpoint routes based on subdomain context
- JWT tokens include user type (platform_admin vs tenant_user) and domain context
- Service validates gym domain exists before tenant authentication

## Critical Development Patterns

### Module Structure (Follow Exactly)

```go
// Each module must have:
internal/module_name/
├── dto/           # Data Transfer Objects (Creation, Update, Response)
├── handler/       # HTTP layer - converts service responses to standardized HTTP responses
├── interfaces/    # Repository, Service, Handler contracts
├── module/        # Dependency injection: repo → service → handler → router
├── repository/    # Data layer - handles database operations, returns raw SQL errors
├── router/        # Chi router setup
└── service/       # Business logic - transforms data, maps SQL errors to APIError
````

### Error Handling Pattern

```go
// Repository: Return raw SQL errors
return "", sql.ErrNoRows

// Service: Map to domain errors using pkg/apierror
if errors.Is(err, sql.ErrNoRows) {
    return apierror.New(errorcode_enum.CodeNotFound, "Gym not found", err)
}

// Handler: Convert APIError to HTTP responses
if apiErr, ok := err.(*apierror.APIError); ok {
    response.WriteAPIError(w, apiErr)
}
```

### Standardized API Response Pattern`

```go
// Success responses use WriteAPISuccess
response.WriteAPISuccess(w, "User created successfully", userData)

// Error responses - APIError automatically mapped to HTTP status
response.WriteAPIError(w, apiErr)

// Manual error responses (fallback for non-APIErrors)
w.WriteHeader(http.StatusInternalServerError)
json.NewEncoder(w).Encode(response.APIResponse[any]{
    Status:  "error",
    Message: "Internal server error",
    Data:    nil,
})

// Response format (consistent across all endpoints):
{
    "status": "success|error",
    "message": "Human readable message",
    "data": actual_data_or_error_details
}
```

### Testing Strategy

- **Service tests**: Mock repositories using `testify/mock`
- **Handler tests**: Use `httptest` with mock services
- **Repository tests**: Integration tests with real database
- **Test file naming**: `*_test.go` in separate `*_test` packages

### Database Patterns

```go
// Dynamic updates with pointer fields in DTOs
if updateDTO.Name != nil {
    query += ", name = $X"
    args = append(args, *updateDTO.Name)
}

// PostgreSQL arrays for exercises
pq.Array(exercise.MuscularGroups)  // For insertion
exercise.Synonyms pq.StringArray   // For scanning

// Soft deletes for exercises, hard deletes for equipment/muscle groups
UPDATE table SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
```

## Development Commands

```bash
# Development (with live reload including frontend changes)
air

# Testing
go test ./...                    # All tests
go test ./internal/gym/...       # Module-specific
go test -v -cover ./...          # Verbose with coverage

# Database setup
go run ./cmd/setup-db/main.go    # Creates public tables
```

## API Integration Patterns

### Frontend → Backend Flow

```js
// Frontend API client (frontend/js/main.js)
const response = await fetch("/api/v1/gym", {
  method: "POST",
  body: JSON.stringify(gymData),
});

// Backend routing (api/api.go)
r.Mount("/auth", authmodule.NewAuthModule(db));
r.Mount("/gym", gymmodule.NewGymModule(db));
r.Mount("/user", usermodule.NewUserModule(db));
```

### Module Initialization Pattern

```go
// Every module follows this exact pattern (see gym/module/gym_module.go)
func NewModuleModule(db *sql.DB) http.Handler {
    repo := repository.NewModuleRepository(db)
    service := service.NewModuleService(repo)
    handler := handler.NewModuleHandler(service)
    return router.NewModuleRouter(handler)
}
```

## Key File References

- **Module pattern**: `internal/gym/module/gym_module.go`, `internal/auth/module/auth_module.go`
- **Interface design**: `internal/gym/interfaces/*.go`, `internal/auth/interfaces/*.go`
- **Authentication**: `internal/auth/` (Header-based single-endpoint routing)
- **Auth DTOs**: `internal/auth/dto/` (login, token, repository DTOs)
- **Auth Documentation**: `docs/auth-module-architecture.md` (comprehensive architecture guide)
- **Error handling**: `pkg/apierror/api_error.go`
- **Standard responses**: `pkg/response/api_response.go`
- **Testing examples**: `internal/gym/handler/gym_handler_test.go`
- **Tenancy**: `internal/database/tenancy.go`
- **Frontend integration**: `frontend/js/main.js` (ApiClient class)
- **OpenAPI docs**: `docs/openapi/` (modular YAML structure)

## Never Do

- Mix tenant and public data in same table/schema
- Return domain errors from repository layer
- Skip interface definitions for services/repositories
- Hard delete exercises (use soft delete with deleted_at)
- Put business logic in handlers or repositories
