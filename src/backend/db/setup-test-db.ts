import { Client } from "postgres";
import { load } from "dotenv";
import { DataBaseGenerator } from "./db-generator.ts";

async function setupTestDatabase() {
    try {
        await load({ 
            export: true,
            envPath: ".env",
            defaultsPath: null,
            examplePath: null
        });

        // Connect to default postgres database first
        let testClient = new Client({
            user: "postgres",
            database: "postgres", // Connect to default database
            hostname: "localhost",
            port: 5432,
            password: Deno.env.get("PG_PASSWORD")
        });

        await testClient.connect();

        // Drop test database if exists and create new one
        await testClient.queryArray(`DROP DATABASE IF EXISTS athenai_test`);
        await testClient.queryArray(`CREATE DATABASE athenai_test`);
        await testClient.end();

        // Connect to new test database
        testClient = new Client({
            user: "postgres",
            database: "athenai_test",
            hostname: "localhost",
            port: 5432,
            password: Deno.env.get("PG_PASSWORD")
        });

        await testClient.connect();

        // Initialize database schema
        const dbGenerator = new DataBaseGenerator(testClient);
        await dbGenerator.createEnums();
        await dbGenerator.createTables();

        console.log("Test database setup completed successfully");
        await testClient.end();

    } catch (error) {
        console.error("Error setting up test database:", error);
        throw error;
    }
}

if (import.meta.main) {
    await setupTestDatabase();
}