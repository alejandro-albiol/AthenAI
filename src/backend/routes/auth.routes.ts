import { Router } from "oak";
import { AuthController } from "../auth/auth-controller.ts";
import { AuthService } from "../auth/auth-service.ts";
import { UsersService } from "../users/users-service.ts";
import { authGuard } from "../shared/middlewares/auth.middleware.ts";
import { createClient } from "../db/config.ts";
import { UsersRepository } from "../users/users-repository.ts";

const client = await createClient();
const usersRepository = new UsersRepository(client);
const authService = new AuthService();
const usersService = new UsersService(usersRepository, authService);
const authController = new AuthController(authService, usersService);

const router = new Router();

// Auth routes
router
    .post("/login", async (ctx) => {
        try {
            const body = await ctx.request.body({ type: "json" }).value;
            const result = await authController.login({ body });
            ctx.response.status = result.status;
            ctx.response.body = result;
        } catch (_error) {
            ctx.response.status = 400;
            ctx.response.body = {
                success: false,
                error: "Invalid request body"
            };
        }
    })
    .get("/validate", authGuard, (ctx) => {
        // Return the user data from the token
        ctx.response.status = 200;
        ctx.response.body = {
            success: true,
            data: ctx.state.user
        };
    })
    .post("/logout", authGuard, async (ctx) => {
        try {
            const result = await authController.logout(ctx.state.user);
            ctx.response.status = result.status;
            ctx.response.body = result;
        } catch (_error) {
            ctx.response.status = 500;
            ctx.response.body = {
                success: false,
                error: "Internal server error"
            };
        }
    });

export const authRouter = router;
