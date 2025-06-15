import { ApiResponse } from "../../shared/utils/response.interface.ts";
import { User } from "./users.interface.ts";
import { UpdateUserDto } from "../dtos/update-user.dto.ts";

export interface IUserController {
    createUser(req: Request): Promise<ApiResponse<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>>>;
    findUserById(req: { params: { id?: string } }): Promise<ApiResponse<User | null>>;
    findUserByUsername(req: { params: { username?: string } }): Promise<ApiResponse<User | null>>;
    findUserByEmail(req: { params: { email?: string } }): Promise<ApiResponse<User | null>>;
    updateUser(req: { params: { id?: string }, body: UpdateUserDto }): Promise<ApiResponse<User>>;
    deleteUser(req: { params: { id?: string } }): Promise<ApiResponse<void>>;
    listUsers(req: { query: { limit?: string, offset?: string } }): Promise<ApiResponse<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>[]>>;
}