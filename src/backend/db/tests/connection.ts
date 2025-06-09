import { createClient } from "../config.ts";
import { DataBaseConnection } from "../db-connection.ts";
import { load } from "dotenv";

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
            // Load environment variables first
            await load({
                export: true,
                envPath: ".env",
                defaultsPath: null
            });
            console.log("Environment variables loaded");

            // Create client after env vars are loaded
            const client = createClient();
            this.dbConnection = new DataBaseConnection(client);
        } catch (error) {
            console.error("Failed to initialize:", error);
            throw error;
        }
    }

    async runTest() {
        console.log("Starting database connection test...");

        try {
            // Initialize connection with proper env vars
            await this.initializeConnection();
            await this.delay(this.DELAY_MS);

            // Test connection
            console.log("Attempting to connect...");
            await this.dbConnection.connect();
            console.log("✓ Database connection established successfully");

            // Wait before next operation
            await this.delay(this.DELAY_MS);

            // Test closing
            console.log("Attempting to close connection...");
            await this.dbConnection.close();
            console.log("✓ Database connection closed successfully");

        } catch (error) {
            console.error("❌ Test failed:", error);
            throw error;
        }
    }
}

if (import.meta.main) {
    const test = new DatabaseConnectionTest();
    await test.runTest()
        .then(() => console.log("✓ All tests completed successfully"))
        .catch(() => {
            console.error("❌ Test suite failed");
            Deno.exit(1);
        });
}