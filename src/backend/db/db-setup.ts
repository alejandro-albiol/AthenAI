import { DataBaseConnection } from "./db-connection.ts";
import { DataBaseGenerator } from "./db-generator.ts";
import { createClient } from "./config.ts";
import { envConfig } from "../config/env.config.ts";


export async function setupDatabase() {
    await envConfig.init();
    
    const client = await createClient();
    const dbConnection = new DataBaseConnection(client);
    const tablesGenerator = new DataBaseGenerator(client);
    
    console.log("PG_PASSWORD is set:", !!Deno.env.get("PG_PASSWORD"));
    await dbConnection.connect();
    await tablesGenerator.createDatabase();
    await tablesGenerator.dropTables();
    await tablesGenerator.dropEnums();
    await tablesGenerator.createEnums();
    await tablesGenerator.createTables();
    await dbConnection.close();
}

if (import.meta.main) {
    await setupDatabase();
}