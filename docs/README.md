# AthenAI Documentation

Welcome to the comprehensive documentation for AthenAI - A multi-tenant gym management platform with AI-powered workout generation.

## 📚 Documentation Overview

This documentation is organized into logical sections to help developers understand and contribute to the AthenAI platform.

### 🏗️ **Architecture & Design**

- **[Backend Architecture](./backend-architecture.md)** - Complete overview of modules, multi-tenancy, and system design
- **[Database Design](./database-design.md)** - Schema, relationships, tenant isolation, and performance optimization
- **[Security Model](./security-model.md)** - Authentication, authorization, JWT implementation, and access control

### 🔧 **Development Guides**

- **[Module Development Pattern](./module-pattern.md)** - Standard patterns and best practices for creating new modules
- **[Configuration Guide](./configuration.md)** - Environment setup, deployment, and configuration management
- **[API Documentation](./openapi/openapi.yaml)** - Complete OpenAPI specification with examples

### 🚀 **Getting Started**

- **[Local Development Setup](../README.md#development)** - How to set up the development environment
- **[Deployment Guide](./deployment.md)** - Production deployment instructions

## 🏛️ **System Architecture Overview**

AthenAI follows a modular, multi-tenant architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway Layer                        │
│  Authentication • CORS • Rate Limiting • Request Routing   │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   Business Logic Layer                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │   Platform  │ │   Tenant    │ │      AI Workout        │ │
│  │   Modules   │ │   Modules   │ │     Generator          │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                     Data Layer                              │
│  ┌─────────────┐ ┌─────────────────────────────────────────┐ │
│  │   Public    │ │         Tenant Schemas                  │ │
│  │   Schema    │ │    {gym_uuid_1} │ {gym_uuid_2} │ ...   │ │
│  └─────────────┘ └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 🎯 **Key Features**

### **Multi-Tenant Architecture**

- **Tenant Isolation**: Each gym operates in its own database schema
- **Shared Resources**: Common exercises, equipment, and templates in public schema
- **Scalable Design**: Easy to add new gyms without affecting existing ones

### **AI-Powered Workout Generation**

- **Smart Exercise Selection**: Based on user preferences, goals, and available equipment
- **Adaptive Templates**: Dynamic workout generation using AI models
- **Personalization**: Considers user training phase, motivation, and special situations

### **Comprehensive Management**

- **User Management**: Role-based access control for gym staff and members
- **Exercise Library**: Extensive database of exercises with custom additions
- **Workout Tracking**: Complete workout history and progress tracking

## 🔐 **Security Highlights**

- **JWT-Based Authentication**: Stateless, secure token management
- **Multi-Tenant Security**: Automatic tenant isolation and access control
- **Role-Based Authorization**: Fine-grained permissions for different user types
- **Data Protection**: Encrypted passwords, secure API endpoints

## 📖 **Quick Navigation**

| I want to...                                 | Go to...                                          |
| -------------------------------------------- | ------------------------------------------------- |
| Understand the system architecture           | [Backend Architecture](./backend-architecture.md) |
| Learn about multi-tenancy and data isolation | [Database Design](./database-design.md)           |
| Understand authentication and security       | [Security Model](./security-model.md)             |
| Create a new module                          | [Module Development Pattern](./module-pattern.md) |
| Set up development environment               | [Configuration Guide](./configuration.md)         |
| Test and explore the API                     | [OpenAPI Documentation](./openapi/openapi.yaml)   |
| Deploy to production                         | [Configuration Guide](./configuration.md)         |

## 🤝 **Contributing**

When contributing to AthenAI:

1. **Follow the Module Pattern**: Use the standard structure defined in [Module Development Pattern](./module-pattern.md)
2. **Update Documentation**: Keep docs in sync with code changes
3. **Security First**: Follow security guidelines in [Security Model](./security-model.md)
4. **Test Coverage**: Include tests for new modules and features

---

**Need Help?** Check the specific documentation sections above or refer to the code examples in the `/internal` modules.
