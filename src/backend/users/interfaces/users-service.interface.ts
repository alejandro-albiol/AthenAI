import { User } from "./users.interface.ts";

export interface IUsersService {
    createUser(username: string, email: string, password: string): Promise<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>>;
    getUserById(id: string): Promise<User | null>;
    getUserByUsername(username: string): Promise<User | null>;
    getUserByEmail(email: string): Promise<User | null>;
    updateUser(id: string, userData: Partial<Omit<User, 'id' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>>): Promise<User>;
    deleteUser(id: string): Promise<void>;
    listUsers(limit?: number, offset?: number): Promise<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>[]>;
}