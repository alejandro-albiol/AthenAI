# Security Model Documentation

## Overview

This document outlines the comprehensive security model for the Athenai multi-tenant gym management platform, focusing on authentication, authorization, and tenant isolation.

## Authentication Architecture

### Single Login Endpoint Strategy

**Endpoint**: `POST /api/v1/auth/login`

**Routing Logic**:

```
Request Headers:
├── No X-Gym-ID header → Platform Admin Authentication
│   └── Authenticate against public.admin table
└── X-Gym-ID header present → Tenant User Authentication
    ├── Lookup gym by ID to get domain
    └── Authenticate against {gym_uuid}.users table
```

**Security Benefits**:

- Single point of authentication control
- Automatic routing based on context
- No complex endpoint management
- Clear separation between admin and tenant flows

## Authorization Model

### JWT-Based Stateless Authorization

**JWT Claims Structure**:

```json
{
  "user_id": "uuid",
  "user_type": "platform_admin|tenant_user",
  "username": "string",
  "role": "admin|user|guest",     // Only for tenant users
  "gym_id": "uuid",               // Only for tenant users
  "is_active": true,
  "exp": timestamp,
  "iat": timestamp
}
```

### Access Control Matrix

| User Type        | Access Level      | Restrictions      | Use Cases                          |
| ---------------- | ----------------- | ----------------- | ---------------------------------- |
| `platform_admin` | **Global Access** | None              | System admin, multi-gym operations |
| `tenant_user`    | **Gym-Scoped**    | Only own gym data | Gym members, trainers, gym admins  |

### Header Usage Policy

**CRITICAL SECURITY RULE**: Headers are ONLY trusted during login

| Endpoint Type             | X-Gym-ID Header           | Authorization Source    |
| ------------------------- | ------------------------- | ----------------------- |
| **Login** (`/auth/login`) | Required for tenant users | Header (one-time trust) |
| **All Other Endpoints**   | **IGNORED**               | JWT claims only         |

## Multi-Tenant Isolation

### Database Schema Isolation

**Public Schema** (`public.*`):

- `admin` - Platform administrators
- `gym` - Gym entities and metadata
- `refresh_tokens` - Cross-tenant token storage
- Global reference data (exercises, equipment)

**Tenant Schemas** (`{gym_uuid}.*`):

- `users` - Gym-specific users
- `user_profiles` - Personal data
- All gym-specific operational data

### Runtime Isolation Enforcement

**Middleware Security**:

```go
// Automatic gym access validation
if !middleware.ValidateGymAccess(r, requestedGymID) {
    return 403 Forbidden
}

// Secure domain lookup when needed
domain, err := middleware.GetGymDomain(r, authService)
```

**Repository Operations**:

```go
// Dynamic table targeting based on JWT gym context
tableName := fmt.Sprintf("%s.users", gymDomain)
query := "SELECT * FROM " + tableName + " WHERE ..."
```

## Security Implementation Patterns

### 1. Handler-Level Authorization

```go
func (h *GymHandler) GetGymByID(w http.ResponseWriter, r *http.Request, id string) {
    // FIRST: Validate access based on JWT claims
    if !middleware.ValidateGymAccess(r, id) {
        response.WriteAPIError(w, apierror.New(
            errorcode_enum.CodeForbidden,
            "Access denied: You can only access your own gym data",
            nil,
        ))
        return
    }

    // THEN: Proceed with business logic
    gym, err := h.service.GetGymByID(id)
    // ...
}
```

### 2. Service-Level Domain Resolution

```go
func (s *UserService) CreateUser(userDTO dto.UserCreationDTO, gymID string) error {
    // Get domain from gym ID (from JWT, not headers)
    gym, err := s.gymRepo.GetGymByID(gymID)
    if err != nil {
        return apierror.New(errorcode_enum.CodeNotFound, "Gym not found", err)
    }

    // Use domain for tenant-specific operations
    return s.userRepo.CreateUserInTenant(gym.Domain, userDTO)
}
```

### 3. Middleware Context Extraction

```go
// Secure context helpers
userID := middleware.GetUserID(r)        // From JWT
userType := middleware.GetUserType(r)    // From JWT
gymID := middleware.GetGymID(r)          // From JWT (NOT headers)

// Authorization helpers
if middleware.IsPlatformAdmin(r) {
    // Global access allowed
}

if middleware.IsGymAdmin(r) {
    // Gym admin operations allowed
}
```

## Attack Prevention

### 1. **Header Tampering Prevention**

- ❌ **Vulnerable**: Using `X-Gym-ID` header for authorization
- ✅ **Secure**: Using gym ID from JWT claims

### 2. **Cross-Tenant Data Access Prevention**

- ❌ **Vulnerable**: Direct database queries without tenant validation
- ✅ **Secure**: `ValidateGymAccess()` before any gym-specific operations

### 3. **Privilege Escalation Prevention**

- ❌ **Vulnerable**: Trusting client-provided role information
- ✅ **Secure**: Role and permissions from JWT claims only

### 4. **Token Replay Prevention**

- ✅ **Secure**: Short-lived access tokens (24h)
- ✅ **Secure**: Refresh token rotation
- ✅ **Secure**: Token revocation on logout

## API Security Examples

### Secure Endpoint Implementation

```http
# Login (X-Gym-ID header required for tenant users)
POST /api/v1/auth/login
X-Gym-ID: 123e4567-e89b-12d3-a456-426614174000
Content-Type: application/json

{
  "username": "johndoe",
  "password": "userPassword123"
}

# All other requests (NO headers needed)
GET /api/v1/gym/123e4567-e89b-12d3-a456-426614174000
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Authorization automatic based on JWT:
# - Platform admin: Can access any gym
# - Tenant user: Can only access gym matching JWT gym_id claim
```

### Security Response Examples

```json
// Platform admin accessing any gym: ✅ 200 OK
{
  "status": "success",
  "message": "Gym found",
  "data": { /* gym data */ }
}

// Tenant user accessing their own gym: ✅ 200 OK
{
  "status": "success",
  "message": "Gym found",
  "data": { /* gym data */ }
}

// Tenant user trying to access other gym: ❌ 403 Forbidden
{
  "status": "error",
  "message": "Access denied: You can only access your own gym data",
  "data": null
}
```

## Best Practices for Developers

### 1. **Always Use JWT Context**

```go
// ✅ CORRECT: Get gym ID from JWT
gymID := middleware.GetGymID(r)

// ❌ WRONG: Get gym ID from header (security risk)
gymID := r.Header.Get("X-Gym-ID")
```

### 2. **Validate Access Before Operations**

```go
// ✅ CORRECT: Validate before proceeding
if !middleware.ValidateGymAccess(r, targetGymID) {
    return 403 Forbidden
}

// ❌ WRONG: Trust client without validation
// gym, err := service.GetGymByID(targetGymID)
```

### 3. **Use Secure Domain Lookup**

```go
// ✅ CORRECT: Get domain from JWT gym ID
domain, err := middleware.GetGymDomain(r, authService)

// ❌ WRONG: Trust domain from headers
// domain := r.Header.Get("X-Gym-Domain")
```

This security model ensures robust multi-tenant isolation while maintaining excellent performance and developer experience.
