import { Client } from "postgres";
import { DataBaseGenerator } from "./db-generator.ts";
import { envConfig } from "../config/env.config.ts";

async function setupTestDatabase() {
    try {
        await envConfig.init();

        // Connect to default postgres database first
        const bootstrapClient = new Client({
            user: envConfig.get("PG_USER") || "postgres",
            database: "postgres",
            hostname: envConfig.get("PG_HOSTNAME") || "localhost",
            port: Number(envConfig.get("PG_PORT")) || 5432,
            password: envConfig.getOrThrow("PG_PASSWORD")
        });

        await bootstrapClient.connect();
        console.log("Connected to postgres database");

        // Drop test database if exists and create new one
        await bootstrapClient.queryArray(`DROP DATABASE IF EXISTS athenai_test`);
        await bootstrapClient.queryArray(`CREATE DATABASE athenai_test`);
        await bootstrapClient.end();
        console.log("Test database created");

        // Connect to new test database
        const testClient = new Client({
            user: envConfig.get("PG_USER") || "postgres",
            database: "athenai_test",
            hostname: envConfig.get("PG_HOSTNAME") || "localhost",
            port: Number(envConfig.get("PG_PORT")) || 5432,
            password: envConfig.getOrThrow("PG_PASSWORD")
        });

        await testClient.connect();
        console.log("Connected to test database");

        // Initialize database schema
        const dbGenerator = new DataBaseGenerator(testClient);
        
        // Create enums first
        await dbGenerator.createEnums();
        console.log("Database enums created");

        // Create tables with proper relationships
        await dbGenerator.createTables();
        console.log("Database tables created");

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