# Authentication Examples

This document provides comprehensive examples of using the Athenai authentication API.

## Example Flow: Platform Admin Authentication

### 1. Login as Platform Admin
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@athenai.com",
    "password": "adminPassword123"
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6",
    "user_info": {
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "admin",
      "email": "admin@athenai.com",
      "user_type": "platform_admin",
      "role": null,
      "gym_id": null
    }
  }
}
```

### 2. Use Access Token for Authenticated Requests
```bash
curl -X GET http://localhost:8080/api/v1/gym \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Example Flow: Tenant User Authentication

### 1. Login as Tenant User
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Gym-ID: 456e7890-e89b-12d3-a456-426614174000" \
  -d '{
    "email": "john@olympusgym.com",
    "password": "userPassword123"
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4",
    "user_info": {
      "user_id": "789e0123-e89b-12d3-a456-426614174000",
      "username": "johndoe",
      "email": "john@olympusgym.com",
      "user_type": "tenant_user",
      "role": "admin",
      "gym_id": "456e7890-e89b-12d3-a456-426614174000"
    }
  }
}
```

### 2. Use Access Token with Gym Context
```bash
curl -X GET http://localhost:8080/api/v1/user \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Gym-ID: 456e7890-e89b-12d3-a456-426614174000"
```

## Token Management Examples

### Validate Token
```bash
curl -X GET http://localhost:8080/api/v1/auth/validate \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Refresh Access Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "refresh_token": "d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"
  }'
```

### Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "refresh_token": "d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"
  }'
```

## Error Handling Examples

### Invalid Credentials
```json
{
  "status": "error",
  "message": "Invalid admin credentials",
  "data": {
    "code": "UNAUTHORIZED",
    "error": "sql: no rows in result set"
  }
}
```

### Gym Not Found
```json
{
  "status": "error",
  "message": "Gym not found",
  "data": {
    "code": "NOT_FOUND",
    "error": "sql: no rows in result set"
  }
}
```

### Inactive Gym
```json
{
  "status": "error",
  "message": "Gym is not active",
  "data": {
    "code": "FORBIDDEN"
  }
}
```

### Invalid Token
```json
{
  "status": "error",
  "message": "Invalid token",
  "data": {
    "code": "UNAUTHORIZED",
    "error": "token is expired"
  }
}
```
