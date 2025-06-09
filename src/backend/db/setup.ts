import { DataBaseConnection } from "./db-connection.ts";
import { DataBaseGenerator } from "./db-generator.ts";
import { createClient } from "./config.ts";
import { load } from "dotenv";

async function loadEnv() {
    try {
        await load({
            export: true,
            envPath: ".env",
            defaultsPath: null
        });
        console.log("Environment variables loaded from .env");
    } catch (error) {
        console.error("Error loading .env file:", error);
        throw error;
    }
}

export async function setupDatabase() {
    await loadEnv();
    
    const client = createClient();
    const dbConnection = new DataBaseConnection(client);
    const tablesGenerator = new DataBaseGenerator(client);
    
    console.log("PG_PASSWORD is set:", !!Deno.env.get("PG_PASSWORD"));
    await dbConnection.connect();
    await tablesGenerator.createDatabase();
    await tablesGenerator.createEnums();
    await tablesGenerator.createTables();
    await dbConnection.close();
}

if (import.meta.main) {
    await setupDatabase();
}