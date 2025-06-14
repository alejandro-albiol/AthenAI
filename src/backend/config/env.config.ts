import { load } from "std/dotenv/mod.ts";

export class EnvConfig {
    private static instance: EnvConfig;
    private initialized = false;

    private constructor() {}

    static getInstance(): EnvConfig {
        if (!EnvConfig.instance) {
            EnvConfig.instance = new EnvConfig();
        }
        return EnvConfig.instance;
    }

    async init() {
        if (this.initialized) return;

        try {
            await load({
                export: true,
                envPath: ".env",
                defaultsPath: null,
                examplePath: null
            });
            console.log("Environment variables loaded successfully");
            this.initialized = true;
        } catch (error) {
            console.error("Failed to load environment variables:", error);
            throw error;
        }
    }

    get(key: string): string | undefined {
        return Deno.env.get(key);
    }

    getOrThrow(key: string): string {
        const value = this.get(key);
        if (!value) {
            throw new Error(`Environment variable ${key} is not set`);
        }
        return value;
    }

    isDevelopment(): boolean {
        return this.get("NODE_ENV") === "development";
    }

    isProduction(): boolean {
        return this.get("NODE_ENV") === "production";
    }

    isTest(): boolean {
        return this.get("NODE_ENV") === "test";
    }
}

export const envConfig = EnvConfig.getInstance();