# Tenant Middleware Architecture

## Overview

The tenant middleware system provides a flexible, modular approach to handle gym domain requirements across different route categories. This ensures that tenant-specific operations require the `X-Gym-ID` header while platform operations remain unaffected.

## Middleware Types

### 1. RequireGymID
- **Purpose**: Enforces `X-Gym-ID` header presence
- **Usage**: Applied automatically to tenant routes
- **Error Response**: Standardized APIError with `BAD_REQUEST` code
- **Context**: Stores gym ID in request context for handlers

### 2. OptionalGymID 
- **Purpose**: Stores gym ID if present, doesn't enforce
- **Usage**: For routes that can work with or without tenant context
- **Behavior**: Continues execution regardless of header presence

### 3. Route Grouping Functions

#### TenantRoutes()
```go
tenantRoutes := middleware.TenantRoutes()
tenantRoutes.Mount("/user", usermodule.NewUserModule(db))
```
- Automatically applies `RequireGymID` middleware
- Use for modules that need tenant context (users, workouts, member-specific data)

#### PlatformRoutes()
```go
platformRoutes := middleware.PlatformRoutes()
platformRoutes.Mount("/auth", authmodule.NewAuthModule(db, gymRepo, userRepo, jwtSecret))
platformRoutes.Mount("/gym", gymmodule.NewGymModule(db))
```
- No gym middleware applied
- Use for platform-level operations (authentication, gym management)

#### MixedRoutes()
```go
mixedRoutes := middleware.MixedRoutes()
mixedRoutes.Mount("/admin", adminmodule.NewAdminModule(db))
```
- Applies `OptionalGymID` middleware
- Use for modules that can work in both contexts

## Implementation Example

### API Module (api/api.go)
```go
func NewAPIModule(db *sql.DB) *chi.Mux {
    r := chi.NewRouter()
    
    // Platform routes - no gym domain required
    platformRoutes := middleware.PlatformRoutes()
    platformRoutes.Mount("/auth", authmodule.NewAuthModule(...))
    platformRoutes.Mount("/gym", gymmodule.NewGymModule(db))
    r.Mount("/", platformRoutes)
    
    // Tenant routes - gym domain required
    tenantRoutes := middleware.TenantRoutes()
    tenantRoutes.Mount("/user", usermodule.NewUserModule(db))
    tenantRoutes.Mount("/workout", workoutmodule.NewWorkoutModule(db))
    r.Mount("/", tenantRoutes)
    
    return r
}
```

### Module Router (no middleware needed)
```go
func NewUsersRouter(handler interfaces.UserHandler) http.Handler {
    r := chi.NewRouter()
    // Middleware applied at API level - no need here
    
    getGymID := func(r *http.Request) string {
        return middleware.GetGymID(r)
    }
    
    r.Post("/", func(w http.ResponseWriter, r *http.Request) {
        handler.RegisterUser(w, r, getGymID(r))
    })
    // ... other routes
}
```

## Route Categories

### Platform Routes (No Gym Middleware)
- `/auth/*` - Authentication endpoints
- `/gym/*` - Gym management (create, update, list gyms)

### Tenant Routes (Require Gym Middleware)
- `/user/*` - User management within a gym
- `/workout/*` - Workout management
- `/exercise/*` - Gym-specific exercises
- `/member/*` - Member operations

### Mixed Routes (Optional Gym Middleware)
- `/admin/*` - Admin operations that might work across gyms
- `/report/*` - Reports that can be gym-specific or global

## Benefits

1. **Modular**: Each module doesn't need to handle middleware internally
2. **Scalable**: Easy to add new modules to appropriate route groups
3. **Clear Separation**: Platform vs tenant operations are clearly distinguished
4. **Consistent**: Standardized error responses using APIError pattern
5. **Flexible**: Easy to change middleware requirements by moving modules between route groups

## Usage Guidelines

1. **New tenant-specific modules**: Add to `tenantRoutes`
2. **New platform modules**: Add to `platformRoutes`
3. **Modules that work in both contexts**: Add to `mixedRoutes`
4. **Don't apply middleware in individual routers**: Let the API level handle it
5. **Use `middleware.GetGymID(r)` in handlers**: Always use this helper for gym ID extraction

## Error Handling

The middleware follows the standard APIError pattern:
```go
apiErr := apierror.New(
    errorcode_enum.CodeBadRequest,
    "X-Gym-ID header is required for tenant operations",
    nil,
)
response.WriteAPIError(w, apiErr)
```

This ensures consistent error responses across all tenant-protected routes.
