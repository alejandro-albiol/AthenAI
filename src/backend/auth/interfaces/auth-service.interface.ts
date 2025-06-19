export interface IAuthService {
    hashPassword(password: string): Promise<string>;
    comparePasswords(password: string, hashedPassword: string): Promise<boolean>;
    generateToken(payload: TokenPayload): Promise<string>;
    generateRefreshToken(payload: TokenPayload): Promise<string>;
    verifyToken(token: string): Promise<TokenPayload>;
    verifyRefreshToken(token: string): Promise<TokenPayload>;
    refreshToken(refreshToken: string): Promise<{ accessToken: string; refreshToken: string }>;
}

export interface TokenPayload {
    id: string;
    email: string;
    username: string;
}

export interface LoginDto {
    email: string;
    password: string;
}

export interface AuthResponse {
    accessToken: string;
    refreshToken: string;
    user: {
        id: string;
        email: string;
        username: string;
    };
}
