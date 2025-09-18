# Development Scripts

This directory contains utility scripts for common development tasks.

## Windows (PowerShell)

### Database Setup

```powershell
.\scripts\setup-db.ps1
```

### Create Superadmin

```powershell
.\scripts\setup-superadmin.ps1
```

### Cleanup Tenant Data

```powershell
.\scripts\cleanup-tenant.ps1 -GymId "gym-uuid-here"
```

## Cross-Platform (Go Commands)

### Database Setup

```bash
go run ./cmd/setup-db/main.go
```

### Create Superadmin

```bash
go run ./cmd/setup-superadmin/main.go
```

### Cleanup Tenant

```bash
go run ./cmd/cleanup-tenant/main.go
```

### Migration Commands

```bash
# Refresh tokens (if needed)
go run ./cmd/migrate-refresh-tokens/main.go
```

## Development Workflow

1. **Initial Setup**:

   ```bash
   cp example.env .env
   # Edit .env with your configuration
   go run ./cmd/setup-db/main.go
   go run ./cmd/setup-superadmin/main.go
   ```

2. **Start Development**:

   ```bash
   air  # Live reload
   # or
   go run ./cmd/main.go  # Direct run
   ```

3. **Testing**:
   ```bash
   go test ./...  # All tests
   go test -v ./internal/auth/...  # Specific module
   go test -cover ./...  # With coverage
   ```

## Notes

- All scripts require proper environment configuration (.env file)
- Database must be accessible with configured credentials
- PowerShell scripts are Windows-specific; use Go commands for cross-platform compatibility
