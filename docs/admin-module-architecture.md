# Admin Module Architecture

## Overview

This document explains the architectural decision to create a separate `admin` module for managing public data in the Athenai multitenancy gym application.

## Architecture Decision

### Recommended Structure: Separate Admin Module

```
internal/
├── admin/           # Manages public/global data
│   ├── dto/         # Data Transfer Objects
│   ├── handler/     # HTTP handlers
│   ├── interfaces/  # Repository and service interfaces
│   ├── module/      # Module initialization
│   ├── repository/  # Database operations
│   ├── router/      # Route definitions
│   └── service/     # Business logic
├── gym/            # Manages gym tenants
├── user/           # Manages gym-specific users
└── [future modules] # Other domain-specific modules
```

## Why Admin Module?

### 1. **Clear Separation of Concerns**

**Admin Module Responsibilities:**

- Manage global exercise catalog
- Manage equipment catalog
- Manage muscular groups
- System-wide configurations

**Gym Module Responsibilities:**

- Manage gym tenants
- Handle tenant-specific data
- Tenant isolation

**User Module Responsibilities:**

- Manage users within each gym
- User authentication/authorization
- Role-based access control

### 2. **Security and Access Control**

```go
// Admin routes require admin authentication
/admin/exercises/*     // Only system admins
/admin/equipment/*     // Only system admins
/admin/ai-models/*     // Only system admins

// Gym routes require gym-specific authentication
/gym/*                // Gym owners/admins
/user/*               // Gym users
```

### 3. **Data Ownership**

**Public Data (Admin Module):**

- Exercises available to all gyms
- Equipment catalog
- Muscular groups
- System configurations

**Private Data (Gym/User Modules):**

- Gym-specific user data
- Gym-specific workout plans
- Gym-specific customizations

### 4. **Scalability Benefits**

- **Admin operations**: Less frequent, but critical
- **Gym operations**: High volume, tenant-specific
- **Independent scaling**: Can scale admin vs gym services separately
- **Caching strategies**: Different caching needs for public vs private data

## Implementation Benefits

### 1. **Code Organization**

```go
// Clear module boundaries
adminService.CreateExercise(exercise)     // Global exercise
gymService.CreateGym(gym)                 // New tenant
userService.CreateUser(user)              // Tenant user
```

### 2. **Testing**

```go
// Easier to test in isolation
func TestAdminService_CreateExercise(t *testing.T) { ... }
func TestGymService_CreateGym(t *testing.T) { ... }
func TestUserService_CreateUser(t *testing.T) { ... }
```

### 3. **Deployment**

- Admin features can be deployed independently
- Different teams can work on different modules
- Easier to implement feature flags

### 4. **Database Design**

```sql
-- Public schema (admin module)
public.exercise
public.equipment
public.muscular_group

-- Tenant schemas (gym/user modules)
gym_uuid.users
gym_uuid.workouts
gym_uuid.custom_exercises
```

## Alternative Approaches Considered

### Approach 2: Domain-Specific Modules

```
internal/
├── exercise/        # Both public and private exercises
├── equipment/       # Equipment management
├── muscular_group/  # Muscle groups
├── gym/            # Gym management
└── user/           # User management
```

**Problems with this approach:**

- Mixed responsibilities (public vs private data)
- Harder to implement proper access control
- More complex authentication logic
- Difficult to scale independently
- Unclear ownership boundaries

### Approach 3: Single Module

```
internal/
├── gym/            # Everything in one module
└── user/           # User management
```

**Problems with this approach:**

- Violates single responsibility principle
- Hard to maintain and test
- Security concerns
- Poor scalability

## Migration Strategy

### Phase 1: Create Admin Module

1. Create admin module structure
2. Move public data management to admin module
3. Update API routes

### Phase 2: Update Existing Modules

1. Remove public data management from gym module
2. Update gym module to focus on tenant management
3. Ensure proper separation

### Phase 3: Add New Features

1. Advanced admin features
2. Additional exercise management features
3. Reporting and analytics

## API Structure

### Admin Endpoints

```
POST   /api/v1/admin/exercises
GET    /api/v1/admin/exercises
GET    /api/v1/admin/exercises/{id}
PUT    /api/v1/admin/exercises/{id}
DELETE /api/v1/admin/exercises/{id}

POST   /api/v1/admin/equipment
GET    /api/v1/admin/equipment
PUT    /api/v1/admin/equipment/{id}
DELETE /api/v1/admin/equipment/{id}


```

### Gym Endpoints (Existing)

```
POST   /api/v1/gym
GET    /api/v1/gym/{id}
PUT    /api/v1/gym/{id}
DELETE /api/v1/gym/{id}
```

### User Endpoints (Existing)

```
POST   /api/v1/user
GET    /api/v1/user/{id}
PUT    /api/v1/user/{id}
DELETE /api/v1/user/{id}
```

## Conclusion

The admin module approach provides:

- ✅ Clear separation of concerns
- ✅ Better security and access control
- ✅ Improved scalability
- ✅ Easier maintenance and testing
- ✅ Future-proof architecture

This structure will support the application's growth and make it easier to add new features while maintaining clean code organization and proper security boundaries.
