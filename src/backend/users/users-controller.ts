import { CreateUserDto } from "./dtos/create-user.dto.ts";
import { UpdateUserDto } from "./dtos/update-user.dto.ts";
import { ApiResponse } from "../shared/utils/response.interface.ts";
import { User } from "./interfaces/users.interface.ts";
import { UsersService } from "./users-service.ts";
import { BaseError } from "../shared/errors/custom-errors.ts";
import { IUserController } from "./interfaces/users-controller.interface.ts";

export class UsersController implements IUserController {
    constructor(private readonly usersService: UsersService) {}

    async createUser(req: Request): Promise<ApiResponse<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>>> {
        
        const dto = req.body as unknown as CreateUserDto;
        if (!dto.username || !dto.email || !dto.password) {
            return {
                success: false,
                status: 400,
                error: "Username, email, and password are required"
            };
        }

        try {
            const user = await this.usersService.createUser(
                dto.username,
                dto.email,
                dto.password
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