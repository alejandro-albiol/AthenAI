import { User } from "./users.interface.ts";

export interface IUsersRepository {
    findById(id: string): Promise<User | null>;
    findByEmail(email: string): Promise<User | null>;
    findByUsername(username: string): Promise<User | null>;
    create(user: Omit<User, 'id' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>): Promise<User>;
    update(id: string, user: Partial<User>): Promise<User>;
    softDelete(id: string): Promise<void>;
    list(limit?: number, offset?: number): Promise<User[]>;
}