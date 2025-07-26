# AthenAI Copilot Instructions

## Architecture & Structure

AthenAI is a multi-tenant gym management platform with a Go backend and static frontend. The backend is organized by strict module boundaries under `internal/`, each following a layered pattern:

- **repository → service → handler → router → module**
- Each module (e.g. `gym`, `auth`, `user`, `admin`) has its own `dto/`, `interfaces/`, `repository/`, `service/`, `handler/`, `router/`, and `module/` directories.
- All business logic is in the service layer; handlers only translate service results to HTTP responses.
- All database access is via repositories, which return raw SQL errors only.

### Multi-Tenancy

- Public schema holds gyms and global data; each gym gets its own schema named by UUID (never by domain).
- Tenant schema creation is automatic in `GymService.CreateGym()`.
- Never mix tenant and public data in the same table or schema.

### Authentication

- Single `/auth/login` endpoint, routes by presence of `X-Gym-ID` header.
- Platform admins (no header) vs. tenant users (header, gym UUID).
- JWT tokens include user type and gym UUID context.
- See `internal/auth/` and `docs/auth-module-architecture.md` for details.

## Key Patterns & Conventions

- **Module initialization:** Always use the pattern in `internal/gym/module/gym_module.go` and `internal/auth/module/auth_module.go`.
- **Error handling:**
  - Repository: return raw SQL errors (e.g. `sql.ErrNoRows`).
  - Service: map to APIError using `pkg/apierror`.
  - Handler: use `response.WriteAPIError()` and `response.WriteAPISuccess()`.
- **API responses:** Always use the format in `pkg/response/api_response.go`.
- **Testing:**
  - Service tests: mock repositories with `testify/mock`.
  - Handler tests: call handler methods directly (middleware is not invoked).
  - Repository tests: use a real database.
  - See `internal/gym/handler/gym_handler_test.go` for handler test patterns.
- **Database:**
  - Use pointer fields in DTOs for dynamic updates.
  - Use PostgreSQL arrays for fields like `muscular_groups`.
  - Soft delete for exercises (set `deleted_at`), hard delete for equipment/muscle groups.

## Developer Workflows

- **Build & run:** Use `air` for live reload, or `go run ./cmd/setup-db/main.go` for DB setup.
- **Testing:**
  - `go test ./...` for all tests
  - `go test ./internal/gym/...` for module-specific
  - `go test -v -cover ./...` for verbose with coverage
- **Frontend:** Static files in `frontend/`, see `frontend/js/main.js` for API usage.

## Integration & References

- **API routing:** See `api/api.go` for how modules are mounted.
- **OpenAPI:** Modular YAML in `docs/openapi/`.
- **Key files:**
  - Module pattern: `internal/gym/module/gym_module.go`
  - Error handling: `pkg/apierror/api_error.go`
  - Standard responses: `pkg/response/api_response.go`
  - Auth: `internal/auth/`, `docs/auth-module-architecture.md`
  - Tenancy: `internal/database/tenancy.go`

## Never Do

- Mix tenant and public data in the same table/schema
- Return domain errors from repository layer
- Skip interface definitions for services/repositories
- Hard delete exercises (use soft delete with deleted_at)
- Put business logic in handlers or repositories
