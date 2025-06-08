import { Application } from "oak";

const app = new Application();
const isProduction = Deno.env.get("NODE_ENV") === "production";

// Puerto por defecto para producción
const DEFAULT_PORT = 8080;
const port = isProduction 
  ? Number(Deno.env.get("PORT")) || DEFAULT_PORT
  : Number(Deno.env.get("DEV_PORT")) || 3000;

app.use((ctx) => {
  ctx.response.body = { 
    message: "Welcome to AthenAI API",
    mode: isProduction ? "production" : "development"
  };
});

console.log(`Server running in ${isProduction ? "production" : "development"} mode`);
console.log(`Server listening on http://localhost:${port}`);
await app.listen({ port });