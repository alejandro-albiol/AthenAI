import { assert, assertEquals, assertRejects } from "std/assert/mod.ts";
import { beforeAll, afterAll, describe, it } from "testing/bdd.ts";
import { Client } from "postgres";
import { UsersService } from "../users-service.ts";
import { UsersRepository } from "../users-repository.ts";
import { AuthService } from "../../auth/auth-service.ts";
import { getTestClient, cleanupDatabase } from "../../shared/test/test-utils.ts";
import { ConflictError, ValidationError } from "../../shared/errors/custom-errors.ts";
import { UserErrorCode } from "../../shared/enums/error-codes.enum.ts";

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
            const result = await service.createUser(
                "validuser", 
                "valid@test.com", 
                "ValidPass123!"
            );
            
            assertEquals(result.username, "validuser");
            assertEquals(result.email, "valid@test.com");
            assert(result.id);
        });

        it("should reject invalid username format", async () => {
            await assertRejects(
                () => service.createUser("a", "test@test.com", "ValidPass123!"),
                ValidationError,
                UserErrorCode.INVALID_USERNAME_FORMAT
            );
        });

        it("should reject invalid email format", async () => {
            await assertRejects(
                () => service.createUser("testuser2", "invalid-email", "ValidPass123!"),
                ValidationError,
                UserErrorCode.INVALID_EMAIL_FORMAT
            );
        });

        it("should reject invalid password", async () => {
            await assertRejects(
                () => service.createUser("testuser3", "test@test.com", "weak"),
                ValidationError,
                UserErrorCode.INVALID_PASSWORD
            );
        });

        it("should reject duplicate email", async () => {
            // First create a user
            await service.createUser("uniqueuser", "duplicate@test.com", "ValidPass123!");
            
            // Try to create another user with same email
            await assertRejects(
                () => service.createUser("anotheruser", "duplicate@test.com", "ValidPass123!"),
                ConflictError,
                UserErrorCode.EMAIL_ALREADY_EXISTS
            );
        });

        it("should reject duplicate username", async () => {
            // First create a user
            await service.createUser("duplicateuser", "user1@test.com", "ValidPass123!");
            
            // Try to create another user with same username
            await assertRejects(
                () => service.createUser("duplicateuser", "user2@test.com", "ValidPass123!"),
                ConflictError,
                UserErrorCode.USERNAME_ALREADY_EXISTS
            );
        });
    });

    describe("User Updates", () => {
        let testUserId: string;

        beforeAll(async () => {
            const user = await service.createUser(
                "updatetestuser",
                "update@test.com",
                "UpdatePass123!"
            );
            testUserId = user.id;
        });

        it("should update username", async () => {
            const updated = await service.updateUser(testUserId, {
                username: "newusername"
            });
            assertEquals(updated.username, "newusername");
        });

        it("should reject duplicate username update", async () => {
            // Create another user first
            await service.createUser("existinguser", "existing@test.com", "ValidPass123!");
            
            // Try to update first user's username to the existing one
            await assertRejects(
                () => service.updateUser(testUserId, { username: "existinguser" }),
                ConflictError,
                UserErrorCode.USERNAME_ALREADY_EXISTS
            );
        });
    });
});