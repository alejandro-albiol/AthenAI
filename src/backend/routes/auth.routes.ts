import { Router } from "oak";

const router = new Router();

// Auth routes
router
    .post("/login", (context) => {
        // TODO: Implement login
        context.response.body = { message: "Login endpoint" };
    })
    .post("/refresh", (context) => {
        // TODO: Implement token refresh
        context.response.body = { message: "Refresh token endpoint" };
    })
    .post("/logout", (context) => {
        // TODO: Implement logout
        context.response.body = { message: "Logout endpoint" };
    });

export const authRouter = router;
