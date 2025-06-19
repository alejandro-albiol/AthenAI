import { CreateUserDto } from "./dtos/create-user.dto.ts";
import { UpdateUserDto } from "./dtos/update-user.dto.ts";
import { ApiResponse } from "../shared/utils/response.interface.ts";
import { User } from "./interfaces/users.interface.ts";
import { UsersService } from "./users-service.ts";
import { BaseError } from "../shared/errors/custom-errors.ts";
import { IUserController } from "./interfaces/users-controller.interface.ts";
import { validateUsername, validateEmail, validatePassword } from "../shared/validators/user.validator.ts";

export class UsersController implements IUserController {
    constructor(private readonly usersService: UsersService) {}    async createUser(req: Request): Promise<ApiResponse<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>>> {
        try {
            // Parse request body
            let body;
            try {
                body = await req.json();
            } catch (_error) {
                return {
                    success: false,
                    status: 400,
                    error: "Invalid JSON in request body"
                };
            }
            
            // Validate request body structure
            if (!body || typeof body !== 'object') {
                return {
                    success: false,
                    status: 400,
                    error: "Request body must be a JSON object"
                };
            }

            const dto = body as CreateUserDto;
            const errors: string[] = [];

            // Check if required fields exist
            if (!dto.email?.trim()) {
                errors.push("Email is required");
            } else if (!validateEmail(dto.email)) {
                errors.push("Invalid email format");
            }

            if (!dto.username?.trim()) {
                errors.push("Username is required");
            } else if (!validateUsername(dto.username)) {
                errors.push("Username must be 3-20 characters long and can only contain letters, numbers, and underscores");
            }

            if (!dto.password?.trim()) {
                errors.push("Password is required");
            } else if (!validatePassword(dto.password)) {
                errors.push("Password must be at least 8 characters long");
            }

            if (errors.length > 0) {
                return {
                    success: false,
                    status: 400,
                    error: errors.join(", ")
                };
            }

            // Create user with trimmed values
            const user = await this.usersService.createUser(
                dto.username.trim(),
                dto.email.trim(),
                dto.password.trim()
            );

            return {
                success: true,
                status: 201,
                data: user
            };
        } catch (error: unknown) {
            return this.handleError(error);
        }
    }

    async findUserById(req: { params: { id?: string } }): Promise<ApiResponse<User | null>> {
        
        const userId = req.params.id;
        if (!userId) {
            return {
                success: false,
                status: 400,
                error: "User ID is required"
            };
        }

        try {
            const user = await this.usersService.getUserById(userId);

            return {
                success: true,
                status: 200,
                data: user
            };
        } catch (error: unknown) {
            return this.handleError(error);
        }
    }

    async findUserByUsername(req: { params: { username?: string } }): Promise<ApiResponse<User | null>> {
        const username = req.params.username;
        if (!username) {
            return {
                success: false,
                status: 400,
                error: "Username is required"
            };
        }

        try {
            const user = await this.usersService.getUserByUsername(username);
            return {
                success: true,
                status: 200,
                data: user
            };
        } catch (error: unknown) {
            return this.handleError(error);
        }
    }

    async findUserByEmail(req: { params: { email?: string } }): Promise<ApiResponse<User | null>> {
        const email = req.params.email;
        if (!email) {
            return {
                success: false,
                status: 400,
                error: "Email is required"
            };
        }

        try {
            const user = await this.usersService.getUserByEmail(email);
            return {
                success: true,
                status: 200,
                data: user
            };
        } catch (error: unknown) {
            return this.handleError(error);
        }
    }

    async updateUser(req: { params: { id?: string }, body: UpdateUserDto }): Promise<ApiResponse<User>> {
        const userId = req.params.id;
        if (!userId) {
            return {
                success: false,
                status: 400,
                error: "User ID is required"
            };
        }

        try {
            const user = await this.usersService.updateUser(userId, req.body);
            return {
                success: true,
                status: 200,
                data: user
            };
        } catch (error: unknown) {
            return this.handleError(error);
        }
    }

    async deleteUser(req: { params: { id?: string } }): Promise<ApiResponse<void>> {
        const userId = req.params.id;
        if (!userId) {
            return {
                success: false,
                status: 400,
                error: "User ID is required"
            };
        }

        try {
            await this.usersService.deleteUser(userId);
            return {
                success: true,
                status: 204
            };
        } catch (error: unknown) {
            return this.handleError(error);
        }
    }

    async listUsers(req: { query: { limit?: string, offset?: string } }): Promise<ApiResponse<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>[]>> {
        try {
            const limit = req.query.limit ? parseInt(req.query.limit) : undefined;
            const offset = req.query.offset ? parseInt(req.query.offset) : undefined;
            
            const users = await this.usersService.listUsers(limit, offset);
            return {
                success: true,
                status: 200,
                data: users
            };
        } catch (error: unknown) {
            return this.handleError(error);
        }
    }

    private handleError<T>(error: unknown): ApiResponse<T> {
        if (error instanceof BaseError) {
            return {
                success: false,
                status: error.statusCode,
                error: error.message
            } as ApiResponse<T>;
        }

        return {
            success: false,
            status: 500,
            error: error instanceof Error ? error.message : "Internal Server Error"
        } as ApiResponse<T>;
    }
}