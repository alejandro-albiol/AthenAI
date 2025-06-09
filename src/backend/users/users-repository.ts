import { Client } from "postgres";
import { User } from "./interfaces/users.interface.ts";
import { IUsersRepository } from "./interfaces/users-repository.interface.ts";

export class UsersRepository implements IUsersRepository {
    constructor(private readonly client: Client) {}

    async findById(id: string): Promise<User | null> {
        try {
            const result = await this.client.queryObject<User>(`
                SELECT * FROM users 
                WHERE id = $1 AND is_deleted = false
            `, [id]);
            
            return result.rows[0] || null;
        } catch (error) {
            console.error('Error in findById:', error);
            throw error;
        }
    }

    async findByUsername(username: string): Promise<User | null> {
        try {
            const result = await this.client.queryObject<User>(`
                SELECT * FROM users 
                WHERE username = $1 AND is_deleted = false
            `, [username]);
            
            return result.rows[0] || null;
        } catch (error) {
            console.error('Error in findByUsername:', error);
            throw error;
        }
    }

    async findByEmail(email: string): Promise<User | null> {
        try {
            const result = await this.client.queryObject<User>(`
                SELECT * FROM users 
                WHERE email = $1 AND is_deleted = false
            `, [email]);
            
            return result.rows[0] || null;
        } catch (error) {
            console.error('Error in findByEmail:', error);
            throw error;
        }
    }

    async create(userData: Omit<User, 'id' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>): Promise<User> {
        try {
            const result = await this.client.queryObject<User>(`
                INSERT INTO users (username, email, password_hash)
                VALUES ($1, $2, $3)
                RETURNING *
            `, [userData.username, userData.email, userData.password_hash]);
            
            return result.rows[0];
        } catch (error) {
            console.error('Error in create:', error);
            throw error;
        }
    }

    async update(id: string, userData: Partial<User>): Promise<User> {
        try {
            const setClause = Object.entries(userData)
                .filter(([key]) => !['id', 'created_at'].includes(key))
                .map(([key], index) => `${key} = $${index + 2}`)
                .join(', ');

            const values = Object.entries(userData)
                .filter(([key]) => !['id', 'created_at'].includes(key))
                .map(([_, value]) => value);

            const result = await this.client.queryObject<User>(`
                UPDATE users 
                SET ${setClause}, updated_at = CURRENT_TIMESTAMP
                WHERE id = $1 AND is_deleted = false
                RETURNING *
            `, [id, ...values]);
            
            return result.rows[0];
        } catch (error) {
            console.error('Error in update:', error);
            throw error;
        }
    }

    async softDelete(id: string): Promise<void> {
        try {
            await this.client.queryObject(`
                UPDATE users 
                SET is_deleted = true, 
                    deleted_at = CURRENT_TIMESTAMP
                WHERE id = $1
            `, [id]);
        } catch (error) {
            console.error('Error in softDelete:', error);
            throw error;
        }
    }

    async list(limit = 10, offset = 0): Promise<User[]> {
        try {
            const result = await this.client.queryObject<User>(`
                SELECT * FROM users 
                WHERE is_deleted = false
                ORDER BY created_at DESC
                LIMIT $1 OFFSET $2
            `, [limit, offset]);
            
            return result.rows;
        } catch (error) {
            console.error('Error in list:', error);
            throw error;
        }
    }
}