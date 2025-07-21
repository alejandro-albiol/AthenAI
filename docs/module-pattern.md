# AthenAI Module Pattern Guide

## Purpose

This document defines the standard pattern for building new modules in AthenAI. Following this pattern ensures modularity, testability, and maintainability across the codebase.

## Directory Structure

Each module must have the following structure:

```
internal/<module_name>/
├── dto/           # Data Transfer Objects (Creation, Update, Response)
├── handler/       # HTTP layer - converts service responses to standardized HTTP responses
├── interfaces/    # Repository, Service, Handler contracts
├── module/        # Dependency injection: repo → service → handler → router
├── repository/    # Data layer - handles database operations, returns raw SQL errors
├── router/        # Chi router setup
└── service/       # Business logic - transforms data, maps SQL errors to APIError
```

## Wiring Pattern

Each module must wire dependencies in this order:

1. **Repository**: Handles DB operations, returns raw SQL errors.
2. **Service**: Contains business logic, maps SQL errors to APIError.
3. **Handler**: HTTP layer, converts service responses to HTTP responses.
4. **Router**: Defines RESTful endpoints and maps them to handler methods.
5. **Module**: Wires everything together and exposes a `New<ModuleName>Module(db *sql.DB) http.Handler` function.

### Example Module Initializer

```go
func NewExampleModule(db *sql.DB) http.Handler {
	repo := repository.NewExampleRepository(db)
	service := service.NewExampleService(repo)
	handler := handler.NewExampleHandler(service)
	return router.NewExampleRouter(handler)
}
```

## Required Practices

- **Interfaces**: Define interfaces for repository, service, and handler in `interfaces/`.
- **DTOs**: Use DTOs for all data transfer between layers.
- **Error Handling**: Repository returns raw SQL errors; service maps to APIError; handler converts to HTTP response.
- **Standardized Responses**: Use `response.WriteAPISuccess`, `response.WriteAPICreated` and `response.WriteAPIError` for all HTTP responses.
- **Testing**: Mock repositories in service tests; mock services in handler tests; use integration tests for repositories.

## Checklist for New Modules

- [ ] All required folders and files created
- [ ] Interfaces defined for repository, service, handler
- [ ] DTOs defined for all data transfer
- [ ] Repository returns raw SQL errors only
- [ ] Service maps errors to APIError and contains business logic
- [ ] Handler converts APIError to HTTP responses
- [ ] Router maps endpoints to handler methods
- [ ] Module wires everything together as shown above
- [ ] Tests written for service, handler, and repository

## Reference

See `internal/user/` and `internal/gym/` for working examples.
