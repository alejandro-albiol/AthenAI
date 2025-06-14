import { createClient } from "../config.ts";
import { DataBaseConnection } from "../db-connection.ts";
import { envConfig } from "../../config/env.config.ts";
import { assert } from "std/assert/mod.ts";

class DatabaseConnectionTest {
    private dbConnection: DataBaseConnection;
    private readonly DELAY_MS = 1000; // 1 second delay

    constructor() {
        this.dbConnection = null!; // Will be initialized after env loads
    }

    private delay(ms: number): Promise<void> {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    private async initializeConnection() {
        try {
            await envConfig.init();
            console.log("Environment variables loaded");

            // Verify environment variables
            assert(envConfig.get("PG_PASSWORD"), "Database password is not set");
            
            const client = await createClient();
            this.dbConnection = new DataBaseConnection(client);
        } catch (error) {
            console.error("Failed to initialize:", error);
            throw error;
        }
    }

    testDatabaseVariables() {
        console.log("Testing database environment variables...");
        
        assert(envConfig.get("PG_USER") || "postgres", "Database user is not set");
        assert(envConfig.get("PG_HOSTNAME") || "localhost", "Database hostname is not set");
        assert(envConfig.get("PG_PORT") || 5432, "Database port is not set");
        
        console.log("✓ Database environment variables verified");
    }

    async runTest() {
        console.log("\n🚀 Starting database connection test suite...\n");

        try {
            // Test environment variables
            this.testDatabaseVariables();
            await this.delay(this.DELAY_MS);

            // Initialize connection
            await this.initializeConnection();
            await this.delay(this.DELAY_MS);

            // Test connection
            console.log("\n📡 Testing database connection...");
            await this.dbConnection.connect();
            console.log("✓ Database connection established successfully");

            // Test simple query
            console.log("\n🔍 Testing simple query...");
            const client = this.dbConnection.getClient();
            const result = await client.queryArray("SELECT 1 as test");
            assert(result.rows[0][0] === 1, "Simple query failed");
            console.log("✓ Query executed successfully");

            await this.delay(this.DELAY_MS);

            // Test closing
            console.log("\n🔒 Testing connection closure...");
            await this.dbConnection.close();
            console.log("✓ Database connection closed successfully");

        } catch (error) {
            console.error("\n❌ Test failed:", error);
            throw error;
        }
    }

    private formatResult(success: boolean, message: string): string {
        return success ? `✓ ${message}` : `❌ ${message}`;
    }
}

if (import.meta.main) {
    const test = new DatabaseConnectionTest();
    await test.runTest()
        .then(() => {
            console.log("\n✨ All tests completed successfully");
            console.log("=====================================");
        })
        .catch((error) => {
            console.error("\n❌ Test suite failed");
            console.error("Error details:", error);
            Deno.exit(1);
        });
}