import { ApiResponse } from "../../shared/utils/response.interface.ts";
import { AuthResponse, LoginDto } from "./auth-service.interface.ts";

export interface IAuthController {
    login(request: { body: LoginDto }): Promise<ApiResponse<AuthResponse>>;
    refresh(request: { body: { refreshToken: string } }): Promise<ApiResponse<Omit<AuthResponse, 'user'>>>;
    logout(request: { body: { refreshToken: string } }): Promise<ApiResponse<void>>;
    validateToken(token: string): Promise<boolean>;
}
