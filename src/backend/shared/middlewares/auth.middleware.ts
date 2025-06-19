import { Context, Next } from "oak";
import { AuthService } from "../../auth/auth-service.ts";

const authService = new AuthService();

export async function authGuard(ctx: Context, next: Next) {
    try {
        const authorization = ctx.request.headers.get("Authorization");
        
        if (!authorization) {
            ctx.response.status = 401;
            ctx.response.body = { success: false, error: "No authorization token provided" };
            return;
        }

        const [scheme, token] = authorization.split(" ");
        
        if (scheme !== "Bearer" || !token) {
            ctx.response.status = 401;
            ctx.response.body = { success: false, error: "Invalid authorization format" };
            return;
        }

        try {
            const payload = await authService.verifyToken(token);
            // Add user info to context state for use in protected routes
            ctx.state.user = payload;
            await next();
        } catch (_error) {
            ctx.response.status = 401;
            ctx.response.body = { success: false, error: "Invalid or expired token" };
        }
    } catch (_error) {
        ctx.response.status = 500;
        ctx.response.body = { success: false, error: "Internal server error" };
    }
}
