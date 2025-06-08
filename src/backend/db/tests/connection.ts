import { createClient } from "../config.ts";
import { DataBaseConnection } from "../db-connection.ts";

class DatabaseConnectionTest {
    private dbConnection: DataBaseConnection;
    private readonly DELAY_MS = 1000; // 1 second delay

    constructor() {
        const client = createClient();
        this.dbConnection = new DataBaseConnection(client);
    }

    private delay(ms: number): Promise<void> {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    async runTest() {
        console.log("Starting database connection test...");

        try {
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

// Run the test
if (import.meta.main) {
    const test = new DatabaseConnectionTest();
    await test.runTest()
        .then(() => console.log("✓ All tests completed successfully"))
        .catch(() => {
            console.error("❌ Test suite failed");
            Deno.exit(1);
        });
}