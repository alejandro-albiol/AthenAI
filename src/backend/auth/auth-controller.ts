import { IAuthController } from "./interfaces/auth-controller.interface.ts";
import { ApiResponse } from "../shared/utils/response.interface.ts";
import { AuthResponse, IAuthService, LoginDto } from "./interfaces/auth-service.interface.ts";
import { UsersService } from "../users/users-service.ts";
import { ValidationError } from "../shared/errors/custom-errors.ts";
import { BaseError } from "../shared/errors/custom-errors.ts";

export class AuthController implements IAuthController {
    constructor(
        private readonly authService: IAuthService,
        private readonly usersService: UsersService,
    ) {}

    async login(request: { body: LoginDto }): Promise<ApiResponse<AuthResponse>> {
        try {
            const { email, password } = request.body;

            // Get user by email
            const user = await this.usersService.getUserByEmail(email);
            if (!user) {
                throw new ValidationError('Invalid credentials');
            }

            // Verify password
            const isValid = await this.authService.comparePasswords(password, user.password_hash);
            if (!isValid) {
                throw new ValidationError('Invalid credentials');
            }

            // Generate tokens
            const tokenPayload = {
                id: user.id,
                email: user.email,
                username: user.username
            };

            const [accessToken, refreshToken] = await Promise.all([
                this.authService.generateToken(tokenPayload),
                this.authService.generateRefreshToken(tokenPayload)
            ]);

            return {
                success: true,
                status: 200,
                data: {
                    accessToken,
                    refreshToken,
                    user: {
                        id: user.id,
                        email: user.email,
                        username: user.username
                    }
                }
            };
        } catch (error) {
            return this.handleError(error);
        }
    }

    async refresh(request: { body: { refreshToken: string } }): Promise<ApiResponse<Omit<AuthResponse, 'user'>>> {
        try {
            const { refreshToken } = request.body;
            const tokens = await this.authService.refreshToken(refreshToken);

            return {
                success: true,
                status: 200,
                data: tokens
            };
        } catch (error) {
            return this.handleError(error);
        }
    }    logout(_request: { body: { refreshToken: string } }): Promise<ApiResponse<void>> {
        // In a real application, you might want to blacklist the token
        // For now, we'll just return a success response as JWT are stateless
        return Promise.resolve({
            success: true,
            status: 200
        });
    }

    async validateToken(token: string): Promise<boolean> {
        try {
            await this.authService.verifyToken(token);
            return true;
        } catch {
            return false;
        }
    }

    private handleError<T>(error: unknown): ApiResponse<T> {
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
            error: error instanceof Error ? error.message : 'Internal Server Error'
        };
    }
}
