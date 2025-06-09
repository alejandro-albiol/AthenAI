import { hash, compare } from "bcrypt";

export class AuthService {

    async hashPassword(password: string): Promise<string> {
        try {
            return await hash(password);
        } catch (error) {
            console.error('Error hashing password:', error);
            throw new Error('Password hashing failed');
        }
    }

    async verifyPassword(hashedPassword: string, plainPassword: string): Promise<boolean> {
        try {
            return await compare(plainPassword, hashedPassword);
        } catch (error) {
            console.error('Error verifying password:', error);
            return false;
        }
    }
}