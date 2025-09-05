# Go Internal Module Development Guide

## Module Structure Pattern

Each internal module should follow this standardized directory structure:

```
internal/{module_name}/
├── dto/
│   ├── create_{entity}.dto.go
│   ├── response_{entity}.dto.go
│   └── update_{entity}.dto.go
├── enum/
│   └── {field_name}.enum.go
├── handler/
│   └── {module}_handler.go
├── interfaces/
│   ├── {module}_handler.interface.go
│   ├── {module}_repository.interface.go
│   └── {module}_service.interface.go
├── repository/
│   ├── {module}_repository.go
│   └── {module}_repository_test.go
├── router/
│   └── {module}_router.go
├── service/
│   └── {module}_service.go
└── module/
    └── {module}_module.go
```

## DTO (Data Transfer Objects)

### Creation DTO

- Use struct tags for validation: `json:"field_name" validate:"required"`
- Use enum types for constrained fields
- Use pointers for optional fields (`*string`, `*bool`)
- Include `CreatedBy` field for audit trails

### Response DTO

- Include all entity fields for complete representation
- Use `time.Time` for timestamp fields
- Convert enum types to strings for JSON serialization
- Include metadata fields (`ID`, `IsActive`, `CreatedAt`, `UpdatedAt`)

### Update DTO

- Use pointers for all updatable fields to distinguish between null and unchanged
- Include validation tags with `omitempty` for optional validations
- Use slice types directly for array fields that replace entirely

## Enums

### Definition Pattern

```go
type EnumName string

const (
    Value1 EnumName = "value1"
    Value2 EnumName = "value2"
)

func (e EnumName) IsValid() bool {
    switch e {
    case Value1, Value2:
        return true
    }
    return false
}
```

### Usage Guidelines

- Use string-based enums for JSON compatibility
- Always implement `IsValid()` method
- Use lowercase string values for consistency
- Validate enum values in service layer

## Interfaces

### Handler Interface

- Define HTTP handler methods with standard signature: `(w http.ResponseWriter, r *http.Request)`
- Group related operations logically
- Use descriptive method names (CreateX, GetXByID, GetXsByFilter, etc.)

### Repository Interface

- Define data access methods with clear return types
- Return entity DTOs, not database models
- Use consistent error handling patterns
- Include filtering and search methods

### Service Interface

- Mirror repository interface but may include additional business logic methods
- Handle cross-entity operations (joins, complex filters)
- Return appropriate DTOs based on operation context

## Repository Layer

### Database Operations

- Use parameterized queries to prevent SQL injection
- Handle `pq.Array()` for PostgreSQL array types
- Implement soft deletes with `is_active` flag
- Use `COALESCE()` for partial updates in UPDATE queries

### Error Handling

- Return descriptive errors that can be handled by service layer
- Don't wrap database-specific errors in repository
- Use consistent query patterns across methods

### Testing

- Use `sqlmock` for database mocking
- Test all CRUD operations
- Verify SQL query structure and parameters
- Include edge cases and error scenarios

## Service Layer

### Business Logic

- Validate input DTOs before repository calls
- Check enum validity using `IsValid()` methods
- Implement duplicate checking where necessary
- Handle cross-service dependencies through injected interfaces

### Error Handling Pattern

```go
// Check for API errors first
var apiErr *apierror.APIError
if errors.As(err, &apiErr) {
    return apiErr
}
// Wrap other errors appropriately
return apierror.New(errorcode_enum.CodeInternal, "Descriptive message", err)
```

### Join Table Management

- Handle relationship creation/updates in service layer
- Use separate services for join table operations
- Remove all links before adding new ones during updates
- Clean up relationships when deleting main entities

## Handler Layer

### Request Processing

- Decode JSON request bodies into appropriate DTOs
- Extract URL parameters using `chi.URLParam()`
- Parse query parameters for filtering operations
- Validate request structure before service calls

### Response Patterns

```go
// Success responses
response.WriteAPISuccess(w, "Message", data)
response.WriteAPICreated(w, "Message", id)

// Error responses
response.WriteAPIError(w, apiError)
```

### Error Handling

- Check for `*apierror.APIError` types first
- Wrap unexpected errors with appropriate error codes
- Provide meaningful error messages to clients
- Log errors appropriately for debugging

## Router Configuration

### Route Organization

- Use RESTful URL patterns
- Group related endpoints logically
- Use HTTP methods appropriately (GET, POST, PUT, DELETE)
- Include search/filter endpoints as needed

### Standard REST Patterns

```
POST   /{resource}              # Create
GET    /{resource}              # List all
GET    /{resource}/search       # Filtered search
GET    /{resource}/{id}         # Get by ID
PUT    /{resource}/{id}         # Update
DELETE /{resource}/{id}         # Delete
```

## Module Assembly

### Dependency Injection

- Initialize dependencies in correct order (repository → service → handler)
- Pass interfaces, not concrete types where possible
- Handle cross-module dependencies through service injection
- Keep module initialization in dedicated module package

### Module Pattern

```go
func New{Module}Module(db *sql.DB, deps ...interface{}) http.Handler {
    // Initialize repository
    repo := repository.New{Module}Repository(db)

    // Initialize external dependencies
    // ...

    // Initialize service with dependencies
    service := service.New{Module}Service(repo, deps...)

    // Initialize handler
    handler := handler.New{Module}Handler(service)

    // Return configured router
    return router.New{Module}Router(handler)
}
```

## General Guidelines

### Naming Conventions

- Use PascalCase for exported types and methods
- Use camelCase for unexported variables and functions
- Use descriptive names that indicate purpose
- Follow Go naming conventions consistently

### Validation Strategy

- Validate at DTO level using struct tags
- Perform business logic validation in service layer
- Return appropriate error codes for different validation failures
- Provide clear error messages for validation failures

### Testing Requirements

- Write unit tests for repository methods using sqlmock, covering all CRUD and query methods (success, not found, and error cases).
- Test service layer business logic with custom mock repositories, ensuring all methods are covered for both positive and negative paths.
- Use httptest and mock services for handler layer tests, checking status codes and response bodies for both success and error scenarios.
- Always provide clear, descriptive assertion messages to make debugging easy.
- Ensure all interface methods are implemented in mocks to avoid compilation errors.
- Use table-driven tests for repetitive patterns when appropriate.
- Run tests after each implementation to verify coverage and correctness.

### Performance Considerations

- Use appropriate database indexes for query patterns
- Implement pagination for list endpoints
- Consider caching strategies for frequently accessed data
- Profile and optimize database queries

This guide ensures consistency across all internal modules while maintaining flexibility for specific business requirements.
