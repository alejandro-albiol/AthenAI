import { CreateUserDto } from "./dtos/create-user.dto.ts";
import { ApiResponse } from "../shared/utils/response.interface.ts";
import { User } from "./interfaces/users.interface.ts";
import { UsersService } from "./users-service.ts";
import { BaseError } from "../shared/errors/custom-errors.ts";

export class UsersController {
    constructor(private readonly usersService: UsersService) {}

    async createUser(req: Request): Promise<ApiResponse<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>>> {
        try {
            const dto = req.body as unknown as CreateUserDto;

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
            if (error instanceof BaseError) {
                return {
                    success: false,
                    status: error.statusCode,
                    error: error.message
                };
            }

            return {
                success: false,
                status: 500,
                error: error instanceof Error ? error.message : "Internal Server Error"
            };
        }
    }
}

