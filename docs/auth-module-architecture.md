# Auth Module Architecture

## Overview

This document explains the authentication module architecture for the Athenai multi-tenant gym management platform. The auth module provides a simplified, single-endpoint authentication system that supports both platform administrators and tenant users through header-based routing.

## Architecture Decision

### Simplified Authentication Approach

The auth module implements a **single login endpoint** strategy that automatically routes authentication based on the presence of the `X-Gym-ID` header:

- **Without `X-Gym-ID` header**: Authenticates against `public.admin` table (platform administrators)
- **With `X-Gym-ID` header**: Looks up gym by ID, then authenticates against `{gym_domain}.users` table (tenant users)

## Module Structure

```
internal/auth/
├── dto/
│   ├── login.dto.go        # LoginRequestDTO, LoginResponseDTO, UserInfoDTO
│   ├── token.dto.go        # TokenValidationResponseDTO, ClaimsDTO, RefreshTokenRequestDTO
│   └── repository.dto.go   # AdminAuthDTO, TenantUserAuthDTO, RefreshTokenDTO
├── handler/
│   └── auth_handler.go     # HTTP layer following standard response patterns
├── interfaces/
│   ├── auth_handler.interface.go     # AuthHandler contract
│   ├── auth_repository.interface.go  # AuthRepositoryInterface contract
│   └── auth_service.interface.go     # AuthServiceInterface contract
├── module/
│   └── auth_module.go      # Self-contained dependency injection
├── repository/
│   └── auth_repository.go  # Pure authentication database operations
├── router/
│   └── auth_router.go      # RESTful endpoint definitions
└── service/
    └── auth_service.go     # Business logic orchestration
```

## Key Design Principles

### 1. **Clean Separation of Concerns**

**Auth Repository Responsibilities:**
- Platform admin authentication (`public.admin` table)
- Tenant user authentication (`{gym_domain}.users` table)
- Refresh token management
- **Pure authentication operations only**

**Service Layer Orchestration:**
- Header-based routing logic
- Gym lookup via gym repository
- JWT token generation and validation
- Business logic coordination

**Handler Layer:**
- Standard response patterns (`response.WriteAPIError`, `response.WriteAPISuccess`)
- Proper error code enums (`errorcode_enum.CodeBadRequest`, etc.)
- Request/response transformation

### 2. **Self-Contained Module Architecture**

The auth module follows the established pattern of being completely self-sufficient:

```go
// api.go - Clean and consistent
r.Mount("/auth", authmodule.NewAuthModule(db))
r.Mount("/user", usermodule.NewUserModule(db))
r.Mount("/gym", gymmodule.NewGymModule(db))
```

**Module Benefits:**
- **Encapsulation**: Creates its own gym repository internally
- **Environment-aware**: Reads JWT secret from env with fallback
- **Consistent interface**: Only requires `*sql.DB` parameter
- **Standard pattern**: Returns `http.Handler` like other modules

### 3. **Repository Delegation Pattern**

Instead of duplicating gym operations, the auth service uses the existing gym repository:

```go
// Auth service delegates gym operations
gym, err := s.gymRepo.GetGymByID(gymID)
if err != nil {
    return nil, apierror.New(errorcode_enum.CodeNotFound, "Gym not found", err)
}

// Then authenticates against tenant schema
user, err := s.authRepo.AuthenticateTenantUser(gym.Domain, email, password)
```

This maintains **single responsibility** and avoids code duplication.

## Authentication Flow

### 1. **Platform Admin Login**

```http
POST /auth/login
Content-Type: application/json

{
  "email": "admin@athenai.com",
  "password": "password"
}
```

**Process:**
1. No `X-Gym-ID` header detected
2. Authenticate against `public.admin` table
3. Generate JWT with `user_type: "platform_admin"`
4. Return tokens and user info

### 2. **Tenant User Login**

```http
POST /auth/login
Content-Type: application/json
X-Gym-ID: 123e4567-e89b-12d3-a456-426614174000

{
  "email": "user@olympusgym.com", 
  "password": "password"
}
```

**Process:**
1. `X-Gym-ID` header detected
2. Lookup gym by ID via gym repository
3. Verify gym is active
4. Authenticate against `{gym_domain}.users` table
5. Generate JWT with `user_type: "tenant_user"` and gym context
6. Return tokens and user info

### 3. **Token Structure**

```go
type ClaimsDTO struct {
    UserID   string  `json:"user_id"`
    Username string  `json:"username"`
    UserType string  `json:"user_type"` // "platform_admin" | "tenant_user"
    GymID    *string `json:"gym_id,omitempty"`    // For tenant users
    Role     *string `json:"role,omitempty"`      // For tenant users
    IsActive bool    `json:"is_active"`
    
    jwt.RegisteredClaims // Standard JWT claims
}
```

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Headers |
|--------|----------|-------------|---------|
| `POST` | `/auth/login` | Single login endpoint | Optional: `X-Gym-ID` |
| `POST` | `/auth/refresh` | Refresh access token | `Authorization: Bearer <token>` |
| `POST` | `/auth/logout` | Revoke refresh token | `Authorization: Bearer <token>` |
| `GET` | `/auth/validate` | Validate JWT token | `Authorization: Bearer <token>` |

### Request/Response Examples

#### Login Response
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6",
    "user_info": {
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "johndoe",
      "email": "john@olympusgym.com",
      "user_type": "tenant_user",
      "role": "admin",
      "gym_id": "456e7890-e89b-12d3-a456-426614174000"
    }
  }
}
```

#### Error Response
```json
{
  "status": "error",
  "message": "Invalid credentials",
  "data": {
    "code": "UNAUTHORIZED",
    "error": "sql: no rows in result set"
  }
}
```

## Database Schema Integration

### Platform Admins
- **Table**: `public.admin`
- **Authentication**: Direct table lookup
- **Scope**: Global platform access

### Tenant Users  
- **Table**: `{gym_domain}.users` (dynamic schema)
- **Authentication**: Gym domain resolution → schema-specific lookup
- **Scope**: Gym-specific access with roles

### Refresh Tokens
- **Table**: `public.refresh_tokens`
- **Purpose**: Logout functionality and token revocation
- **Structure**: Hashed tokens with user context

## Security Features

### 1. **JWT Token Security**
- **HMAC-SHA256** signing
- **24-hour** access token expiration
- **30-day** refresh token expiration
- **Environment-based** secret management

### 2. **Password Security**
- **bcrypt** hashing for refresh tokens
- **Database-level** password validation
- **No plaintext** password storage

### 3. **Multi-tenancy Security**
- **Schema isolation** for tenant data
- **Gym domain validation** before authentication
- **Active gym verification**

## Configuration

### Environment Variables

```bash
# JWT Secret (fallback available for development)
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Database connection
DB_DSN=host=localhost port=5432 user=postgres password=password dbname=athenai sslmode=disable
```

### Development Defaults
- **JWT Secret**: Auto-generated fallback for development
- **Token Expiration**: 24h access, 30d refresh
- **Error Handling**: Comprehensive APIError mapping

## Testing Strategy

### Unit Tests
- **Service Layer**: Mock repositories using testify/mock
- **Handler Layer**: httptest with mock services
- **Repository Layer**: Integration tests with real database

### Integration Tests
- **End-to-end** authentication flows
- **Multi-tenant** authentication scenarios
- **Token lifecycle** testing

## Future Enhancements

### Planned Features
- **Role-based access control** middleware
- **Token blacklisting** for immediate revocation
- **SSO integration** for enterprise clients
- **Rate limiting** for authentication endpoints

### OpenAPI Documentation
- **Security schemes** defined in `/docs/openapi/`
- **Combined authentication** patterns (Bearer + GymID header)
- **Comprehensive examples** for all endpoints

## Comparison with Previous Architecture

### Before: Complex Multi-Middleware
- Multiple authentication middlewares
- Complex subdomain routing
- Duplicated gym lookup logic
- Inconsistent error handling

### After: Simplified Single Endpoint
- ✅ **Single authentication endpoint**
- ✅ **Header-based routing**
- ✅ **Repository delegation pattern**
- ✅ **Standardized responses**
- ✅ **Self-contained modules**

## Integration Points

### With Gym Module
```go
// Auth service uses gym repository for lookups
gym, err := s.gymRepo.GetGymByID(gymID)
```

### With User Module  
```go
// Future: User module will use auth middleware
// r.Use(middleware.RequireAuth)
```

### With API Layer
```go
// Clean integration in api.go
r.Mount("/auth", authmodule.NewAuthModule(db))
```

This architecture provides a **scalable**, **maintainable**, and **secure** foundation for authentication in the multi-tenant AthenAI platform.
