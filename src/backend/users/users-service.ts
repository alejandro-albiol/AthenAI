import { IUsersRepository } from "./interfaces/users-repository.interface.ts";
import { User } from "./interfaces/users.interface.ts";
import { UserErrorCode } from "../shared/enums/error-codes.enum.ts";
import { ConflictError, NotFoundError, ValidationError } from "../shared/errors/custom-errors.ts";
import { validateEmail, validatePassword, validateUsername } from "../shared/validators/user.validator.ts";
import { AuthService } from "../auth/auth-service.ts";

export class UsersService {
    constructor(
        private readonly usersRepository: IUsersRepository,
        private readonly authService: AuthService
    ) {}

    async createUser(username: string, email: string, password: string): Promise<Omit<User, 'password_hash' | 'created_at' | 'updated_at' | 'is_deleted' | 'deleted_at'>> {

        if (!validateUsername(username)) {
            throw new ValidationError(UserErrorCode.INVALID_USERNAME_FORMAT);
        }
        if (!validateEmail(email)) {
            throw new ValidationError(UserErrorCode.INVALID_EMAIL_FORMAT);
        }
        if (!validatePassword(password)) {
            throw new ValidationError(UserErrorCode.INVALID_PASSWORD);
        }

        const existingEmail = await this.usersRepository.findByEmail(email);
        if (existingEmail) {
            throw new ConflictError(UserErrorCode.EMAIL_ALREADY_EXISTS);
        }

        const existingUsername = await this.usersRepository.findByUsername(username);
        if (existingUsername) {
            throw new ConflictError(UserErrorCode.USERNAME_ALREADY_EXISTS);
        }

        const passwordHash = await this.authService.hashPassword(password);

        return this.usersRepository.create({
            username,
            email,
            password_hash: passwordHash
        });
    }

    async getUserById(id: string): Promise<User> {
        const user = await this.usersRepository.findById(id);
        if (!user) {
            throw new NotFoundError(UserErrorCode.USER_NOT_FOUND);
        }
        return user;
    }

    async updateUser(id: string, userData: Partial<User>): Promise<User> {
        const user = await this.getUserById(id);

        if (userData.email && userData.email !== user.email) {
            if (!validateEmail(userData.email)) {
                throw new ValidationError(UserErrorCode.INVALID_EMAIL_FORMAT);
            }
            const existingEmail = await this.usersRepository.findByEmail(userData.email);
            if (existingEmail) {
                throw new ConflictError(UserErrorCode.EMAIL_ALREADY_EXISTS);
            }
        }

        if (userData.username && userData.username !== user.username) {
            if (!validateUsername(userData.username)) {
                throw new ValidationError(UserErrorCode.INVALID_USERNAME_FORMAT);
            }
            const existingUsername = await this.usersRepository.findByUsername(userData.username);
            if (existingUsername) {
                throw new ConflictError(UserErrorCode.USERNAME_ALREADY_EXISTS);
            }
        }

        return this.usersRepository.update(id, userData);
    }

    async deleteUser(id: string): Promise<void> {
        const user = await this.getUserById(id);
        if (!user) {
            throw new NotFoundError(UserErrorCode.USER_NOT_FOUND);
        }
        await this.usersRepository.softDelete(id);
    }

}