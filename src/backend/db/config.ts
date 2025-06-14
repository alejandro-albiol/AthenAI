import { Client } from "postgres";
import { envConfig } from "../config/env.config.ts";

export async function createClient() {
    await envConfig.init();

    return new Client({
        user: envConfig.get("PG_USER") || "postgres",
        database: envConfig.get("PG_DATABASE") || "athenai",
        hostname: envConfig.get("PG_HOSTNAME") || "localhost",
        port: Number(envConfig.get("PG_PORT")) || 5432,
        password: envConfig.getOrThrow("PG_PASSWORD")
    });
}