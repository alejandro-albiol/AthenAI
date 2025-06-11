import { assert, assertEquals, assertRejects } from "assert/mod.ts";
import { beforeAll, afterAll, describe, it } from "testing/bdd.ts";
import { UsersService } from "../users-service.ts";
import { UsersRepository } from "../users-repository.ts";
import { AuthService } from "../../auth/auth-service.ts";
import { getTestClient, cleanupDatabase } from "../../shared/test/test-utils.ts";
import { ConflictError, ValidationError } from "../../shared/errors/custom-errors.ts";
import { Client } from "postgres";

let client: Client;
let repository: UsersRepository;
let authService: AuthService;
let service: UsersService;

describe("Users Service", () => {
    beforeAll(async () => {
        client = await getTestClient();
        await client.connect();
        repository = new UsersRepository(client);
        authService = new AuthService();
        service = new UsersService(repository, authService);
    });

    afterAll(async () => {
        await cleanupDatabase(client);
        await client.end();
    });

    describe("User Creation", () => {
        it("should create user with valid data", async () => {
            const user = await service.createUser(
                "validuser",
                "valid@example.com",
                "ValidPass123!"
            );
            assertEquals(user.username, "validuser");
            assertEquals(user.email, "valid@example.com");
        });

        it("should reject invalid email", async () => {
            await assertRejects(
                () => service.createUser("user2", "invalid-email", "ValidPass123!"),
                ValidationError
            );
        });

        it("should reject duplicate email", async () => {
            await assertRejects(
                () => service.createUser("user2", "valid@example.com", "ValidPass123!"),
                ConflictError
            );
        });

        it("should reject invalid password", async () => {
            await assertRejects(
                () => service.createUser("user2", "user2@example.com", "weak"),
                ValidationError
            );
        });
    });

    describe("User Updates", () => {
        beforeAll(async () => {
            
            await service.createUser(
                "updatetestuser",
                "update@example.com",
                "UpdatePass123!"
            );
        });

        it("should update username", async () => {
            const user = await repository.findByEmail("update@example.com");
            assert(user);
            const updated = await service.updateUser(user.id, {
                username: "newusername"
            });
            assertEquals(updated.username, "newusername");
        });

        it("should reject duplicate username update", async () => {
            const user = await repository.findByEmail("update@example.com");
            assert(user);
            await assertRejects(
                () => service.updateUser(user.id, { username: "validuser" }),
                ConflictError
            );
        });
    });
});