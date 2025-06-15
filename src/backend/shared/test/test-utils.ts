import { Client } from "postgres";
import { envConfig } from "../../config/env.config.ts";

export async function getTestClient(): Promise<Client> {
    try {
        await envConfig.init();
        
        const client = new Client({
            user: envConfig.get("PG_USER") || "postgres",
            database: "athenai_test",
            hostname: envConfig.get("PG_HOSTNAME") || "localhost",
            port: Number(envConfig.get("PG_PORT")) || 5432,
            password: envConfig.getOrThrow("PG_PASSWORD")
        });
        return client;
    } catch (error) {
        console.error("Failed to create test client:", error);
        throw error;
    }
}

export async function cleanupDatabase(client: Client) {
    try {
        // Clean all tables
        await client.queryArray(`
            TRUNCATE TABLE 
                users,
                exercises,
                templates,
                trainings,
                training_blocks,
                training_exercises
            CASCADE
        `);
        console.log("Test database cleaned successfully");
    } catch (error) {
        console.error("Failed to cleanup database:", error);
        throw error;
    }
}