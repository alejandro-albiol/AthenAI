import { Client } from "postgres";
import { load } from "dotenv";

export async function getTestClient(): Promise<Client> {
    try {
        await load({ 
            export: true,
            envPath: ".env",
            defaultsPath: null,
            examplePath: null
        });
        
        const client = new Client({
            user: "postgres",
            database: "athenai_test",
            hostname: "localhost",
            port: 5432,
            password: Deno.env.get("PG_PASSWORD")
        });
        return client;
    } catch (error) {
        console.error("Failed to create test client:", error);
        throw error;
    }
}

export async function cleanupDatabase(client: Client) {
    try {
        
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