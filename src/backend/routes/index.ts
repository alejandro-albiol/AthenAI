import { Router } from "oak";
import { swaggerRouter } from "./swagger.routes.ts";
import { userRouter } from "./user.routes.ts";
import { authRouter } from "./auth.routes.ts";

export function setupRoutes(): Router {
    const router = new Router();

    // Health check route
    router.get("/health", (ctx) => {
        ctx.response.body = { status: "healthy" };
    });    // API v1 routes
    const apiV1Router = new Router();
    apiV1Router.use("/users", userRouter.routes(), userRouter.allowedMethods());
    apiV1Router.use("/auth", authRouter.routes(), authRouter.allowedMethods());

    // Mount all routes
    router.use("/api/v1", apiV1Router.routes(), apiV1Router.allowedMethods());
    router.use("/api/docs", swaggerRouter.routes(), swaggerRouter.allowedMethods());

    return router;
}
