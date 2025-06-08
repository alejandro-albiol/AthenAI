import { Client } from "postgres";

export class DataBaseConnection {
    private client: Client;

    constructor(client: Client) {
        this.client = client;
    }

    async connect() {
        try {
            await this.client.connect();
            console.log("Database connection established");
        } catch (error) {
            console.error("Database connection failed:", error);
            throw error;
        }
    }

    async close() {
        try {
            await this.client.end();
            console.log("Database connection closed");
        } catch (error) {
            console.error("Failed to close database connection:", error);
            throw error;
        }
    }
}