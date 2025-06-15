import { assert, assertEquals, assertRejects } from "std/assert/mod.ts";
import { beforeAll, afterAll, describe, it, beforeEach } from "testing/bdd.ts";
import { Client } from "postgres";
import { UsersService } from "../users-service.ts";
import { UsersRepository } from "../users-repository.ts";
import { AuthService } from "../../auth/auth-service.ts";
import { getTestClient, cleanupDatabase } from "../../shared/test/test-utils.ts";
import { ConflictError, NotFoundError, ValidationError } from "../../shared/errors/custom-errors.ts";
import { UserErrorCode } from "../../shared/enums/error-codes.enum.ts";

describe("Users Service", () => {
    let client: Client;
    let repository: UsersRepository;
    let authService: AuthService;
    let service: UsersService;

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

    beforeEach(async () => {
        await cleanupDatabase(client);
    });

    describe("createUser", () => {
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
                () => service.createUser("validuser", "invalid-email", "ValidPass123!"),
                ValidationError,
                UserErrorCode.INVALID_EMAIL_FORMAT
            );
        });

        it("should reject invalid password", async () => {
            await assertRejects(
                () => service.createUser("validuser", "test@test.com", "weak"),
                ValidationError,
                UserErrorCode.INVALID_PASSWORD
            );
        });

        it("should reject duplicate email", async () => {
            await service.createUser("user1", "duplicate@test.com", "ValidPass123!");
            
            await assertRejects(
                () => service.createUser("user2", "duplicate@test.com", "ValidPass123!"),
                ConflictError,
                UserErrorCode.EMAIL_ALREADY_EXISTS
            );
        });

        it("should reject duplicate username", async () => {
            await service.createUser("duplicateuser", "user1@test.com", "ValidPass123!");
            
            await assertRejects(
                () => service.createUser("duplicateuser", "user2@test.com", "ValidPass123!"),
                ConflictError,
                UserErrorCode.USERNAME_ALREADY_EXISTS
            );
        });
    });

    describe("getUserById", () => {
        let testUserId: string;

        beforeEach(async () => {
            const user = await service.createUser(
                "testuser",
                "test@example.com",
                "ValidPass123!"
            );
            testUserId = user.id;
        });

        it("should find user by valid ID", async () => {
            const user = await service.getUserById(testUserId);
            
            assert(user);
            assertEquals(user.username, "testuser");
            assertEquals(user.email, "test@example.com");
        });

        it("should throw NotFoundError for non-existent ID", async () => {
            await assertRejects(
                () => service.getUserById("00000000-0000-0000-0000-000000000000"),
                NotFoundError,
                UserErrorCode.USER_NOT_FOUND
            );
        });
    });

    describe("getUserByUsername", () => {
        beforeEach(async () => {
            await service.createUser(
                "searchuser",
                "search@test.com",
                "ValidPass123!"
            );
        });

        it("should find user by valid username", async () => {
            const user = await service.getUserByUsername("searchuser");
            
            assert(user);
            assertEquals(user.username, "searchuser");
            assertEquals(user.email, "search@test.com");
        });

        it("should throw NotFoundError for non-existent username", async () => {
            await assertRejects(
                () => service.getUserByUsername("nonexistent"),
                NotFoundError,
                UserErrorCode.USER_NOT_FOUND
            );
        });
    });

    describe("getUserByEmail", () => {
        beforeEach(async () => {
            await service.createUser(
                "emailuser",
                "emailsearch@test.com",
                "ValidPass123!"
            );
        });

        it("should find user by valid email", async () => {
            const user = await service.getUserByEmail("emailsearch@test.com");
            
            assert(user);
            assertEquals(user.username, "emailuser");
            assertEquals(user.email, "emailsearch@test.com");
        });

        it("should throw NotFoundError for non-existent email", async () => {
            await assertRejects(
                () => service.getUserByEmail("nonexistent@test.com"),
                NotFoundError,
                UserErrorCode.USER_NOT_FOUND
            );
        });
    });

    describe("updateUser", () => {
        let testUserId: string;

        beforeEach(async () => {
            const user = await service.createUser(
                "updateuser",
                "update@test.com",
                "ValidPass123!"
            );
            testUserId = user.id;
        });

        it("should update username successfully", async () => {
            const updated = await service.updateUser(testUserId, {
                username: "newusername"
            });

            assertEquals(updated.username, "newusername");
            assertEquals(updated.email, "update@test.com"); // Email should remain unchanged
        });

        it("should update email successfully", async () => {
            const updated = await service.updateUser(testUserId, {
                email: "newemail@test.com"
            });

            assertEquals(updated.username, "updateuser"); // Username should remain unchanged
            assertEquals(updated.email, "newemail@test.com");
        });

        it("should reject invalid email format in update", async () => {
            await assertRejects(
                () => service.updateUser(testUserId, { email: "invalid-email" }),
                ValidationError,
                UserErrorCode.INVALID_EMAIL_FORMAT
            );
        });

        it("should reject invalid username format in update", async () => {
            await assertRejects(
                () => service.updateUser(testUserId, { username: "a" }),
                ValidationError,
                UserErrorCode.INVALID_USERNAME_FORMAT
            );
        });

        it("should reject duplicate username in update", async () => {
            // Create another user first
            await service.createUser("existinguser", "existing@test.com", "ValidPass123!");
            
            await assertRejects(
                () => service.updateUser(testUserId, { username: "existinguser" }),
                ConflictError,
                UserErrorCode.USERNAME_ALREADY_EXISTS
            );
        });

        it("should reject duplicate email in update", async () => {
            // Create another user first
            await service.createUser("otheruser", "existing@test.com", "ValidPass123!");
            
            await assertRejects(
                () => service.updateUser(testUserId, { email: "existing@test.com" }),
                ConflictError,
                UserErrorCode.EMAIL_ALREADY_EXISTS
            );
        });
    });

    describe("deleteUser", () => {
        let testUserId: string;

        beforeEach(async () => {
            const user = await service.createUser(
                "deleteuser",
                "delete@test.com",
                "ValidPass123!"
            );
            testUserId = user.id;
        });

        it("should delete user successfully", async () => {
            await service.deleteUser(testUserId);

            // Verify user is deleted by trying to find them
            await assertRejects(
                () => service.getUserById(testUserId),
                NotFoundError,
                UserErrorCode.USER_NOT_FOUND
            );
        });

        it("should throw NotFoundError when deleting non-existent user", async () => {
            await assertRejects(
                () => service.deleteUser("00000000-0000-0000-0000-000000000000"),
                NotFoundError,
                UserErrorCode.USER_NOT_FOUND
            );
        });
    });

    describe("listUsers", () => {
        beforeEach(async () => {
            // Create test users
            for (let i = 1; i <= 5; i++) {
                await service.createUser(
                    `user${i}`,
                    `user${i}@test.com`,
                    "ValidPass123!"
                );
            }
        });

        it("should list all users with default pagination", async () => {
            const users = await service.listUsers();
            
            assertEquals(users.length, 5);
        });

        it("should respect limit parameter", async () => {
            const users = await service.listUsers(2);
            
            assertEquals(users.length, 2);
        });

        it("should respect offset parameter", async () => {
            const users = await service.listUsers(undefined, 3);
            
            assertEquals(users.length, 2); // Should return last 2 users
        });

        it("should handle limit and offset together", async () => {
            const users = await service.listUsers(2, 2);
            
            assertEquals(users.length, 2);
            assert(users[0].username.includes("user"));
            assert(users[1].username.includes("user"));
        });
    });
});