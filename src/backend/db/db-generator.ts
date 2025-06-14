import { Client } from "postgres";

export class DataBaseGenerator {
    private client: Client;

    constructor(client: Client) {
        this.client = client;
    }

    async createDatabase() {
        try {
            const result = await this.client.queryArray(`
                SELECT 1 FROM pg_database WHERE datname = 'athenai';
            `);
            if (result.rowCount === 0) {
                await this.client.queryArray(`CREATE DATABASE athenai`);
                console.log("Database created successfully");
            } else {
                console.log("Database already exists");
            }
        } catch (error) {
            console.error("Error creating database:", error);
            throw error;
        }
    }

    async dropDatabase() {
        try {
            await this.client.queryArray(`DROP DATABASE IF EXISTS athenai;`);
            console.log("Database dropped successfully");
        } catch (error) {
            console.error("Error dropping database:", error);
            throw error;
        }
    }

    async createTables() {
        try {
            await this.createUsersTable();
            await this.createExercisesTable();
            await this.createTemplatesTable();
            await this.createTrainingsTable();
            await this.createTrainingBlocksTable();
            await this.createTrainingExercisesTable();
            console.log("Tables created successfully");
        } catch (error) {
            console.error("Error creating tables:", error);
            throw error;
        }
    }

    async createEnums() {
        try {
            await this.client.queryArray(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`);
            console.log("UUID extension enabled successfully");

            await this.client.queryArray(`
                CREATE TYPE muscle_group AS ENUM (
                    'chest', 'back', 'shoulders', 'biceps', 'triceps',
                    'legs', 'calves', 'abs', 'forearms', 'full_body'
                );
            `);

            await this.client.queryArray(`
                CREATE TYPE equipment AS ENUM (
                    'barbell', 'dumbbell', 'machine', 'cable',
                    'bodyweight', 'kettlebell', 'resistance_band', 'other'
                );
            `);

            await this.client.queryArray(`
                CREATE TYPE training_goal AS ENUM (
                    'strength', 'hypertrophy', 'endurance', 
                    'power', 'weight_loss', 'maintenance'
                );
            `);

            await this.client.queryArray(`
                CREATE TYPE intensity_level AS ENUM (
                    'low', 'medium', 'high', 'very_high'
                );
            `);

            console.log("Enums created successfully");
        } catch (error) {
            console.error("Error creating enums:", error);
            throw error;
        }
    }

    private async createUsersTable() {
        const query = `
            CREATE TABLE IF NOT EXISTS users (
                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                username VARCHAR(50) UNIQUE NOT NULL,
                email VARCHAR(100) UNIQUE NOT NULL,
                password_hash VARCHAR(255) NOT NULL,
                created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                is_deleted BOOLEAN DEFAULT FALSE,
                deleted_at TIMESTAMPTZ
            );
        `;
        await this.client.queryArray(query);
    }

    private async createExercisesTable() {
        const query = `
            CREATE TABLE IF NOT EXISTS exercises (
                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                name TEXT NOT NULL,
                synonyms TEXT[],
                description TEXT,
                muscle_group muscle_group,
                equipment equipment,
                created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
            );
        `;
        await this.client.queryArray(query);
    }

    private async createTemplatesTable() {
        const query = `
            CREATE TABLE IF NOT EXISTS templates (
                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                title TEXT NOT NULL,
                description TEXT,
                goal training_goal,
                intensity intensity_level,
                num_blocks INTEGER,
                estimated_time_minutes INTEGER,
                is_public BOOLEAN DEFAULT FALSE,
                created_by UUID REFERENCES users(id),
                created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
            );
        `;
        await this.client.queryArray(query);
    }

    private async createTrainingsTable() {
        const query = `
            CREATE TABLE IF NOT EXISTS trainings (
                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                template_id UUID REFERENCES templates(id),
                title TEXT,
                notes TEXT,
                is_public BOOLEAN DEFAULT FALSE,
                created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
            );
        `;
        await this.client.queryArray(query);
    }

    private async createTrainingBlocksTable() {
        const query = `
            CREATE TABLE IF NOT EXISTS training_blocks (
                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                training_id UUID REFERENCES trainings(id) ON DELETE CASCADE,
                block_number INTEGER NOT NULL,
                description TEXT
            );
        `;
        await this.client.queryArray(query);
    }

    private async createTrainingExercisesTable() {
        const query = `
            CREATE TABLE IF NOT EXISTS training_exercises (
                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                training_block_id UUID REFERENCES training_blocks(id) ON DELETE CASCADE,
                exercise_id UUID REFERENCES exercises(id),
                order_in_block INTEGER NOT NULL,
                sets INTEGER,
                reps TEXT,
                rest_seconds INTEGER,
                tempo TEXT,
                notes TEXT
            );
        `;
        await this.client.queryArray(query);
    }

    async dropEnums() {
        try {
            await this.client.queryArray(`DROP TYPE IF EXISTS muscle_group CASCADE`);
            await this.client.queryArray(`DROP TYPE IF EXISTS equipment CASCADE`);
            await this.client.queryArray(`DROP TYPE IF EXISTS training_goal CASCADE`);
            await this.client.queryArray(`DROP TYPE IF EXISTS intensity_level CASCADE`);
            console.log("Enums dropped successfully");
        } catch (error) {
            console.error("Error dropping enums:", error);
            throw error;
        }
    }

    async dropTables() {
        try {
            await this.dropTrainingExercisesTable();
            await this.dropTrainingBlocksTable();
            await this.dropTrainingsTable();
            await this.dropTemplatesTable();
            await this.dropExercisesTable();
            await this.dropUsersTable();
            console.log("Tables dropped successfully");
        } catch (error) {
            console.error("Error dropping tables:", error);
            throw error;
        }
    }

    private async dropUsersTable() {
        await this.client.queryArray(`DROP TABLE IF EXISTS users CASCADE`);
    }

    private async dropExercisesTable() {
        await this.client.queryArray(`DROP TABLE IF EXISTS exercises CASCADE`);
    }

    private async dropTemplatesTable() {
        await this.client.queryArray(`DROP TABLE IF EXISTS templates CASCADE`);
    }

    private async dropTrainingsTable() {
        await this.client.queryArray(`DROP TABLE IF EXISTS trainings CASCADE`);
    }

    private async dropTrainingBlocksTable() {
        await this.client.queryArray(`DROP TABLE IF EXISTS training_blocks CASCADE`);
    }

    private async dropTrainingExercisesTable() {
        await this.client.queryArray(`DROP TABLE IF EXISTS training_exercises CASCADE`);
    }
}
