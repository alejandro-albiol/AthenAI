import { Context, Next } from "oak";

export async function corsMiddleware(ctx: Context, next: Next) {
    // Only set essential CORS headers
    ctx.response.headers.set("Access-Control-Allow-Origin", "*");
    
    // Only add other CORS headers for OPTIONS requests
    if (ctx.request.method === "OPTIONS") {
        ctx.response.headers.set(
            "Access-Control-Allow-Methods",
            "GET,POST,PUT,DELETE"
        );
        ctx.response.headers.set(
            "Access-Control-Allow-Headers",
            "Authorization,Content-Type"
        );
    }
    
    await next();
}
