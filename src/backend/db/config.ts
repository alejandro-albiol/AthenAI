import { Client } from "postgres";

export function createClient() {

    return new Client({
        user: Deno.env.get("PG_USER") || "postgres",
        database: Deno.env.get("PG_DATABASE") || "athenai",
        hostname: Deno.env.get("PG_HOSTNAME") || "localhost",
        port: Number(Deno.env.get("PG_PORT")) || 5432,
        password: Deno.env.get("PG_PASSWORD") || "password",
    });
}