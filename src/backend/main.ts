import { Application } from "oak";
import { setupRoutes } from "./routes/index.ts";
import { corsMiddleware } from "./shared/middlewares/cors.middleware.ts";

const app = new Application();
const isProduction = Deno.env.get("NODE_ENV") === "production";

// Default port configuration
const DEFAULT_PORT = 8080;
const port = isProduction 
  ? Number(Deno.env.get("PORT")) || DEFAULT_PORT
  : Number(Deno.env.get("DEV_PORT")) || 3000;

// Add middlewares
app.use(corsMiddleware);

// Setup routes
const router = setupRoutes();

// Use the router
app.use(router.routes());
app.use(router.allowedMethods());

console.log(`Server running in ${isProduction ? "production" : "development"} mode`);
console.log(`Server listening on http://localhost:${port}`);
console.log(`API Documentation available at http://localhost:${port}/api/docs`);
console.log(`API Base URL: http://localhost:${port}/api/v1`);
await app.listen({ port });