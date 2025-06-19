export const swaggerDefinition = {
  openapi: "3.0.0",
  info: {
    title: "AthenAI API Documentation",
    version: "1.0.0",
    description: "API documentation for AthenAI project",
    contact: {
      name: "API Support"
    }
  },  servers: [
    {
      url: "http://localhost:3000/api/v1",
      description: "Development server"
    },
    {
      url: "http://localhost:8080/api/v1",
      description: "Production server"
    }
  ],
  tags: [
    {
      name: "Users",
      description: "User management endpoints"
    },
    {
      name: "Auth",
      description: "Authentication endpoints"
    }
  ],  components: {
    securitySchemes: {
      BearerAuth: {
        type: "http",
        scheme: "bearer",
        bearerFormat: "JWT"
      }
    },
    schemas: {
      User: {
        type: "object",
        properties: {
          id: {
            type: "string",
            format: "uuid",
            description: "User's unique identifier"
          },
          username: {
            type: "string",
            description: "User's username"
          },
          email: {
            type: "string",
            format: "email",
            description: "User's email address"
          },
          created_at: {
            type: "string",
            format: "date-time",
            description: "User creation timestamp"
          },
          updated_at: {
            type: "string",
            format: "date-time",
            description: "User last update timestamp"
          },
          is_deleted: {
            type: "boolean",
            description: "Whether the user is deleted"
          },
          deleted_at: {
            type: "string",
            format: "date-time",
            description: "User deletion timestamp",
            nullable: true
          }
        }
      },
      CreateUserDto: {
        type: "object",
        required: ["username", "email", "password"],
        properties: {
          username: {
            type: "string",
            description: "User's username"
          },
          email: {
            type: "string",
            format: "email",
            description: "User's email address"
          },
          password: {
            type: "string",
            format: "password",
            description: "User's password"
          }
        }
      },
      UpdateUserDto: {
        type: "object",
        properties: {
          username: {
            type: "string",
            description: "User's new username"
          },
          email: {
            type: "string",
            format: "email",
            description: "User's new email address"
          },
          password: {
            type: "string",
            format: "password",
            description: "User's new password"
          }
        }
      },
      LoginDto: {
        type: "object",
        required: ["email", "password"],
        properties: {
          email: {
            type: "string",
            format: "email",
            description: "User's email address"
          },
          password: {
            type: "string",
            format: "password",
            description: "User's password"
          }
        }
      },
      AuthResponse: {
        type: "object",
        properties: {
          accessToken: {
            type: "string",
            description: "JWT access token"
          },
          refreshToken: {
            type: "string",
            description: "JWT refresh token"
          },
          user: {
            type: "object",
            properties: {
              id: {
                type: "string",
                format: "uuid"
              },
              username: {
                type: "string"
              },
              email: {
                type: "string",
                format: "email"
              }
            }
          }
        }
      },
      ApiResponse: {
        type: "object",
        properties: {
          success: {
            type: "boolean",
            description: "Whether the request was successful"
          },
          status: {
            type: "integer",
            description: "HTTP status code"
          },
          data: {
            type: "object",
            description: "Response data"
          },
          error: {
            type: "string",
            description: "Error message when success is false"
          }
        }
      }
    }
  },  paths: {
    "/users": {
      get: {
        tags: ["Users"],
        summary: "Get all users",
        parameters: [
          {
            name: "limit",
            in: "query",
            schema: {
              type: "integer",
              minimum: 1
            },
            description: "Number of users to return"
          },
          {
            name: "offset",
            in: "query",
            schema: {
              type: "integer",
              minimum: 0
            },
            description: "Number of users to skip"
          }
        ],
        security: [{ BearerAuth: [] }],
        responses: {
          200: {
            description: "List of users",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          type: "array",
                          items: {
                            $ref: "#/components/schemas/User"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          401: {
            description: "Unauthorized"
          }
        }
      },
      post: {
        tags: ["Users"],
        summary: "Create a new user",
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                $ref: "#/components/schemas/CreateUserDto"
              }
            }
          }
        },
        responses: {
          201: {
            description: "User created successfully",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          $ref: "#/components/schemas/User"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          400: {
            description: "Bad request",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      }
    },
    "/users/{id}": {
      parameters: [
        {
          name: "id",
          in: "path",
          required: true,
          schema: {
            type: "string",
            format: "uuid"
          },
          description: "User ID"
        }
      ],
      get: {
        tags: ["Users"],
        summary: "Get user by ID",
        security: [{ BearerAuth: [] }],
        responses: {
          200: {
            description: "User found",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          $ref: "#/components/schemas/User"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          404: {
            description: "User not found",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      },
      put: {
        tags: ["Users"],
        summary: "Update user",
        security: [{ BearerAuth: [] }],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                $ref: "#/components/schemas/UpdateUserDto"
              }
            }
          }
        },
        responses: {
          200: {
            description: "User updated successfully",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          $ref: "#/components/schemas/User"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          404: {
            description: "User not found",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      },
      delete: {
        tags: ["Users"],
        summary: "Delete user",
        security: [{ BearerAuth: [] }],
        responses: {
          204: {
            description: "User deleted successfully",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          },
          404: {
            description: "User not found",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      }
    },
    "/users/username/{username}": {
      parameters: [
        {
          name: "username",
          in: "path",
          required: true,
          schema: {
            type: "string"
          },
          description: "Username"
        }
      ],
      get: {
        tags: ["Users"],
        summary: "Get user by username",
        security: [{ BearerAuth: [] }],
        responses: {
          200: {
            description: "User found",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          $ref: "#/components/schemas/User"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          404: {
            description: "User not found",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      }
    },
    "/users/email/{email}": {
      parameters: [
        {
          name: "email",
          in: "path",
          required: true,
          schema: {
            type: "string",
            format: "email"
          },
          description: "User email"
        }
      ],
      get: {
        tags: ["Users"],
        summary: "Get user by email",
        security: [{ BearerAuth: [] }],
        responses: {
          200: {
            description: "User found",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          $ref: "#/components/schemas/User"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          404: {
            description: "User not found",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      }
    },
    "/auth/login": {
      post: {
        tags: ["Auth"],
        summary: "User login",
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                $ref: "#/components/schemas/LoginDto"
              }
            }
          }
        },
        responses: {
          200: {
            description: "Login successful",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          $ref: "#/components/schemas/AuthResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          401: {
            description: "Invalid credentials",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      }
    },
    "/auth/refresh": {
      post: {
        tags: ["Auth"],
        summary: "Refresh access token",
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["refreshToken"],
                properties: {
                  refreshToken: {
                    type: "string",
                    description: "Refresh token"
                  }
                }
              }
            }
          }
        },
        responses: {
          200: {
            description: "Tokens refreshed successfully",
            content: {
              "application/json": {
                schema: {
                  allOf: [
                    { $ref: "#/components/schemas/ApiResponse" },
                    {
                      properties: {
                        data: {
                          type: "object",
                          properties: {
                            accessToken: {
                              type: "string"
                            },
                            refreshToken: {
                              type: "string"
                            }
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          401: {
            description: "Invalid refresh token",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      }
    },
    "/auth/logout": {
      post: {
        tags: ["Auth"],
        summary: "Logout user",
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: {
                type: "object",
                required: ["refreshToken"],
                properties: {
                  refreshToken: {
                    type: "string",
                    description: "Refresh token to invalidate"
                  }
                }
              }
            }
          }
        },
        responses: {
          200: {
            description: "Logout successful",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/ApiResponse" }
              }
            }
          }
        }
      }
    }
  }
};
