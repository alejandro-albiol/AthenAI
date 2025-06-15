import { assert, assertEquals } from "std/assert/mod.ts";
import { beforeAll, afterAll, describe, it, afterEach, beforeEach } from "testing/bdd.ts";
import { Client } from "postgres";
import { UsersController } from "../users-controller.ts";
import { UsersService } from "../users-service.ts";
import { UsersRepository } from "../users-repository.ts";
import { AuthService } from "../../auth/auth-service.ts";
import { getTestClient, cleanupDatabase } from "../../shared/test/test-utils.ts";
import { UserErrorCode } from "../../shared/enums/error-codes.enum.ts";

describe("Users Controller", () => {
    let client: Client;
    let repository: UsersRepository;
    let authService: AuthService;
    let service: UsersService;
    let controller: UsersController;

    beforeAll(async () => {
        // Setup database and services
        client = await getTestClient();
        await client.connect();
        repository = new UsersRepository(client);
        authService = new AuthService();
        service = new UsersService(repository, authService);
        controller = new UsersController(service);

        // Clean start
        await cleanupDatabase(client);
    });

    afterAll(async () => {
        try {
            await cleanupDatabase(client);
            await client.end();
        } catch (error) {
            console.error("Error in afterAll:", error);
        }
    });

    describe("createUser", () => {
        afterEach(async () => {
            // Clean database after each test to avoid conflicts
            await cleanupDatabase(client);
        });

        it("should return 400 if required fields are missing", async () => {
            const req = {
                body: {
                    username: "testuser"
                    // Missing email and password
                }
            } as unknown as Request;

            const response = await controller.createUser(req);

            assertEquals(response.success, false);
            assertEquals(response.status, 400);
            assertEquals(response.error, "Username, email, and password are required");
        });

        it("should create user and return 201 with valid data", async () => {
            const req = {
                body: {
                    username: "validuser",
                    email: "valid@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request;

            const response = await controller.createUser(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 201);
            assert(response.data);
            // Only check non-date fields
            assertEquals(response.data.username, "validuser");
            assertEquals(response.data.email, "valid@test.com");
            assert(response.data.id);
        });

        it("should return 400 for invalid email format", async () => {
            const req = {
                body: {
                    username: "testuser2",
                    email: "invalid-email",
                    password: "ValidPass123!"
                }
            } as unknown as Request;

            const response = await controller.createUser(req);

            assertEquals(response.success, false);
            assertEquals(response.status, 400);
            assertEquals(response.error, UserErrorCode.INVALID_EMAIL_FORMAT);
        });

        it("should return 409 for duplicate email", async () => {
            // First create initial user
            await controller.createUser({
                body: {
                    username: "validuser",
                    email: "valid@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);

            // Try to create user with same email
            const req = {
                body: {
                    username: "anotheruser",
                    email: "valid@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request;

            const response = await controller.createUser(req);

            // Only check the error response fields
            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 409,
                error: UserErrorCode.EMAIL_ALREADY_EXISTS
            });
        });

        it("should return 409 for duplicate username", async () => {
            // First create initial user
            await controller.createUser({
                body: {
                    username: "validuser",
                    email: "first@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);

            // Try to create user with same username
            const req = {
                body: {
                    username: "validuser",
                    email: "second@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request;

            const response = await controller.createUser(req);

            // Only check the error response fields
            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 409,
                error: UserErrorCode.USERNAME_ALREADY_EXISTS
            });
        });
    });

    describe("findUserById", () => {
        let testUserId: string;

        beforeEach(async () => {
            await cleanupDatabase(client);
            // Create a test user for these tests
            const response = await controller.createUser({
                body: {
                    username: "finduser",
                    email: "find@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);
            assert(response.data);
            testUserId = response.data.id;
        });

        it("should return 400 if ID is missing", async () => {
            const req = { params: {} };
            const response = await controller.findUserById(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 400,
                error: "User ID is required"
            });
        });

        it("should return 404 for non-existent user ID", async () => {
            const req = {
                params: { id: "00000000-0000-0000-0000-000000000000" }
            };
            const response = await controller.findUserById(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 404,
                error: UserErrorCode.USER_NOT_FOUND
            });
        });

        it("should return 200 with user data for valid ID", async () => {
            const req = { params: { id: testUserId } };
            const response = await controller.findUserById(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.username, "finduser");
            assertEquals(response.data.email, "find@test.com");
        });
    });

    describe("findUserByUsername", () => {
        beforeEach(async () => {
            await cleanupDatabase(client);
            await controller.createUser({
                body: {
                    username: "searchuser",
                    email: "search@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);
        });

        it("should return 400 if username is missing", async () => {
            const req = { params: {} };
            const response = await controller.findUserByUsername(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 400,
                error: "Username is required"
            });
        });

        it("should return 404 for non-existent username", async () => {
            const req = { params: { username: "nonexistent" } };
            const response = await controller.findUserByUsername(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 404,
                error: UserErrorCode.USER_NOT_FOUND
            });
        });

        it("should return 200 with user data for valid username", async () => {
            const req = { params: { username: "searchuser" } };
            const response = await controller.findUserByUsername(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.username, "searchuser");
            assertEquals(response.data.email, "search@test.com");
        });
    });

    describe("findUserByEmail", () => {
        beforeEach(async () => {
            await cleanupDatabase(client);
            await controller.createUser({
                body: {
                    username: "emailuser",
                    email: "emailsearch@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);
        });

        it("should return 400 if email is missing", async () => {
            const req = { params: {} };
            const response = await controller.findUserByEmail(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 400,
                error: "Email is required"
            });
        });

        it("should return 404 for non-existent email", async () => {
            const req = { params: { email: "nonexistent@test.com" } };
            const response = await controller.findUserByEmail(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 404,
                error: UserErrorCode.USER_NOT_FOUND
            });
        });

        it("should return 200 with user data for valid email", async () => {
            const req = { params: { email: "emailsearch@test.com" } };
            const response = await controller.findUserByEmail(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.username, "emailuser");
            assertEquals(response.data.email, "emailsearch@test.com");
        });
    });

    describe("updateUser", () => {
        let testUserId: string;

        beforeEach(async () => {
            await cleanupDatabase(client);
            const response = await controller.createUser({
                body: {
                    username: "updateuser",
                    email: "update@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);
            assert(response.data);
            testUserId = response.data.id;
        });

        it("should return 400 if ID is missing", async () => {
            const req = {
                params: {},
                body: { username: "newname" }
            };
            const response = await controller.updateUser(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 400,
                error: "User ID is required"
            });
        });

        it("should return 404 for non-existent user ID", async () => {
            const req = {
                params: { id: "00000000-0000-0000-0000-000000000000" },
                body: { username: "newname" }
            };
            const response = await controller.updateUser(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 404,
                error: UserErrorCode.USER_NOT_FOUND
            });
        });

        it("should successfully update user data", async () => {
            const req = {
                params: { id: testUserId },
                body: {
                    username: "updatedname",
                    email: "updated@test.com"
                }
            };
            const response = await controller.updateUser(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.username, "updatedname");
            assertEquals(response.data.email, "updated@test.com");
        });

        it("should return 409 for duplicate username update", async () => {
            // Create another user first
            await controller.createUser({
                body: {
                    username: "existinguser",
                    email: "existing@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);

            // Try to update first user's username to the existing one
            const req = {
                params: { id: testUserId },
                body: { username: "existinguser" }
            };
            const response = await controller.updateUser(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 409,
                error: UserErrorCode.USERNAME_ALREADY_EXISTS
            });
        });
    });

    describe("deleteUser", () => {
        let testUserId: string;

        beforeEach(async () => {
            await cleanupDatabase(client);
            const response = await controller.createUser({
                body: {
                    username: "deleteuser",
                    email: "delete@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request);
            assert(response.data);
            testUserId = response.data.id;
        });

        it("should return 400 if ID is missing", async () => {
            const req = { params: {} };
            const response = await controller.deleteUser(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 400,
                error: "User ID is required"
            });
        });

        it("should return 404 for non-existent user ID", async () => {
            const req = {
                params: { id: "00000000-0000-0000-0000-000000000000" }
            };
            const response = await controller.deleteUser(req);

            assertEquals({
                success: response.success,
                status: response.status,
                error: response.error
            }, {
                success: false,
                status: 404,
                error: UserErrorCode.USER_NOT_FOUND
            });
        });

        it("should successfully delete user", async () => {
            const req = { params: { id: testUserId } };
            const response = await controller.deleteUser(req);

            assertEquals({
                success: response.success,
                status: response.status
            }, {
                success: true,
                status: 204
            });

            // Verify user is deleted by trying to find them
            const findReq = { params: { id: testUserId } };
            const findResponse = await controller.findUserById(findReq);
            assertEquals(findResponse.success, false);
            assertEquals(findResponse.status, 404);
        });
    });

    describe("listUsers", () => {
        beforeEach(async () => {
            await cleanupDatabase(client);
            // Create test users
            for (let i = 1; i <= 5; i++) {
                await controller.createUser({
                    body: {
                        username: `user${i}`,
                        email: `user${i}@test.com`,
                        password: "ValidPass123!"
                    }
                } as unknown as Request);
            }
        });

        it("should list users with default pagination", async () => {
            const req = { query: {} };
            const response = await controller.listUsers(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.length, 5);
        });

        it("should respect limit parameter", async () => {
            const req = { query: { limit: "2" } };
            const response = await controller.listUsers(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.length, 2);
        });

        it("should respect offset parameter", async () => {
            const req = { query: { offset: "2" } };
            const response = await controller.listUsers(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.length, 3); // Should show last 3 users
        });

        it("should handle limit and offset together", async () => {
            const req = { query: { limit: "2", offset: "2" } };
            const response = await controller.listUsers(req);

            assertEquals(response.success, true);
            assertEquals(response.status, 200);
            assert(response.data);
            assertEquals(response.data.length, 2);
        });
    });
});