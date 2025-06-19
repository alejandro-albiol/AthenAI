import { hash, compare } from "bcrypt";
import { create, verify, getNumericDate } from "https://deno.land/x/djwt@v2.9.1/mod.ts";
import { IAuthService, TokenPayload } from "./interfaces/auth-service.interface.ts";
import { ValidationError } from "../shared/errors/custom-errors.ts";

export class AuthService implements IAuthService {
    private readonly JWT_SECRET = Deno.env.get("JWT_SECRET") || "your-secret-key";
    private readonly JWT_REFRESH_SECRET = Deno.env.get("JWT_REFRESH_SECRET") || "your-refresh-secret-key";
    private readonly JWT_EXPIRATION = 15 * 60; // 15 minutes in seconds
    private readonly REFRESH_TOKEN_EXPIRATION = 7 * 24 * 60 * 60; // 7 days in seconds

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

    comparePasswords(password: string, hashedPassword: string): Promise<boolean> {
        return this.verifyPassword(hashedPassword, password);
    }

    async generateToken(payload: TokenPayload): Promise<string> {
        const key = await crypto.subtle.importKey(
            "raw",
            new TextEncoder().encode(this.JWT_SECRET),
            { name: "HMAC", hash: "SHA-512" },
            false,
            ["sign", "verify"]
        );

        return await create(
            { alg: "HS512", typ: "JWT" },
            {
                ...payload,
                exp: getNumericDate(this.JWT_EXPIRATION)
            },
            key
        );
    }

    async generateRefreshToken(payload: TokenPayload): Promise<string> {
        const key = await crypto.subtle.importKey(
            "raw",
            new TextEncoder().encode(this.JWT_REFRESH_SECRET),
            { name: "HMAC", hash: "SHA-512" },
            false,
            ["sign", "verify"]
        );

        return await create(
            { alg: "HS512", typ: "JWT" },
            {
                ...payload,
                exp: getNumericDate(this.REFRESH_TOKEN_EXPIRATION)
            },
            key
        );
    }

    async verifyToken(token: string): Promise<TokenPayload> {
        try {
            const key = await crypto.subtle.importKey(
                "raw",
                new TextEncoder().encode(this.JWT_SECRET),
                { name: "HMAC", hash: "SHA-512" },
                false,
                ["sign", "verify"]
            );

            const payload = await verify(token, key);
            return {
                id: payload.id as string,
                email: payload.email as string,
                username: payload.username as string
            };
        } catch (_error) {
            throw new ValidationError('Invalid token');
        }
    }

    async verifyRefreshToken(token: string): Promise<TokenPayload> {
        try {
            const key = await crypto.subtle.importKey(
                "raw",
                new TextEncoder().encode(this.JWT_REFRESH_SECRET),
                { name: "HMAC", hash: "SHA-512" },
                false,
                ["sign", "verify"]
            );

            const payload = await verify(token, key);
            return {
                id: payload.id as string,
                email: payload.email as string,
                username: payload.username as string
            };
        } catch (_error) {
            throw new ValidationError('Invalid refresh token');
        }
    }

    async refreshToken(refreshToken: string): Promise<{ accessToken: string; refreshToken: string }> {
        try {
            const payload = await this.verifyRefreshToken(refreshToken);
            const newAccessToken = await this.generateToken(payload);
            const newRefreshToken = await this.generateRefreshToken(payload);

            return {
                accessToken: newAccessToken,
                refreshToken: newRefreshToken
            };
        } catch (_error) {
            throw new ValidationError('Invalid refresh token');
        }
    }
}