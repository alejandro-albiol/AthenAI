# Contributing to AthenAI

Thank you for your interest in contributing to AthenAI! This guide will help you get started.

## 🚀 Quick Start for Contributors

1. **Fork the repository** and clone your fork
2. **Set up development environment**:
   ```bash
   cp example.env .env
   # Configure your database settings in .env
   go mod tidy
   go run ./cmd/setup-db/main.go
   go run ./cmd/setup-superadmin/main.go
   ```
3. **Start development server**: `air`
4. **Make your changes** following the guidelines below
5. **Test your changes**: `go test ./...`
6. **Submit a pull request**

## 📋 Development Guidelines

### Module Development Pattern

All new modules must follow the standard pattern defined in [`docs/module-pattern.md`](./docs/module-pattern.md):

```
internal/{module_name}/
├── dto/           # Data Transfer Objects
├── handler/       # HTTP layer
├── interfaces/    # Contracts/interfaces
├── module/        # Dependency injection
├── repository/    # Data access layer
├── router/        # Route definitions
└── service/       # Business logic
```

### Code Standards

- **Follow Go conventions**: Use `gofmt`, `golint`, and `govet`
- **Write tests**: Minimum 80% coverage for new code
- **Use interfaces**: Define contracts in `interfaces/` directory
- **DTOs for data transfer**: All API communication through DTOs
- **Error handling**: Use `pkg/apierror` for consistent error responses

### Security Requirements

- **JWT-based authentication**: Follow patterns in `internal/auth/`
- **Tenant isolation**: Ensure all tenant operations use gym context from JWT
- **Input validation**: Validate all inputs using struct tags and business logic
- **No SQL injection**: Use parameterized queries only

### Database Guidelines

- **Multi-tenant aware**: All tenant operations must use `{gym_uuid}` schema
- **Migrations**: Include migration scripts for schema changes
- **Indexes**: Add appropriate indexes for query performance
- **Transactions**: Use transactions for multi-table operations

## 🧪 Testing

### Running Tests

```bash
# All tests
go test ./...

# Specific module
go test ./internal/auth/...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Test Structure

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test repository layer with real database
- **Handler tests**: Test HTTP endpoints with mocked services
- **Use mocks**: Generate mocks for interfaces using `mockery`

### Test Files

- Repository tests: `*_repository_test.go`
- Service tests: `*_service_test.go`
- Handler tests: `*_handler_test.go`
- Integration tests: `tests/` directory

## 📝 Documentation

### Required Documentation

- **Update relevant docs** in `/docs` directory when making architectural changes
- **Update OpenAPI specs** for API changes in `/docs/openapi`
- **Add inline comments** for public functions and complex logic
- **Update README** if adding new features or changing setup process

### Documentation Standards

- Use **clear, concise language**
- Include **code examples** where helpful
- **Keep docs in sync** with code changes
- **Follow existing format** and structure

## 🔧 Environment Setup

### Prerequisites

- **Go 1.24+**
- **PostgreSQL 12+**
- **Air** for live reload (optional): `go install github.com/air-verse/air@latest`

### Configuration

1. **Database setup**: Ensure PostgreSQL is running
2. **Environment variables**: Copy and configure `example.env` to `.env`
3. **Database initialization**: Run setup scripts
4. **Verify setup**: Start server and check http://localhost:8080

## 🚨 Pull Request Process

1. **Create feature branch**: `git checkout -b feature/your-feature-name`
2. **Make changes** following guidelines above
3. **Write/update tests** for your changes
4. **Update documentation** as needed
5. **Test thoroughly**: Run full test suite
6. **Submit PR** with clear description of changes

### PR Requirements

- [ ] All tests pass
- [ ] Code follows project conventions
- [ ] Documentation updated (if applicable)
- [ ] No breaking changes (or clearly documented)
- [ ] Security considerations addressed
- [ ] Performance impact considered

## 📞 Getting Help

- **Documentation**: Check `/docs` directory first
- **Issues**: Search existing issues before creating new ones
- **Architecture questions**: Refer to `docs/backend-architecture.md`
- **Security questions**: Refer to `docs/security-model.md`

## 🎯 Areas for Contribution

- **New modules**: Following the standard pattern
- **AI improvements**: Enhance workout generation algorithms
- **Performance optimization**: Database queries, caching, etc.
- **Testing**: Increase test coverage
- **Documentation**: Improve or expand documentation
- **Frontend enhancements**: UI/UX improvements
- **Security auditing**: Review and improve security measures

Thank you for contributing to AthenAI! 🚀
