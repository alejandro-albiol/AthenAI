import { assert, assertEquals } from "std/assert/mod.ts";
import { beforeAll, afterAll, describe, it } from "testing/bdd.ts";
import { Client } from "postgres";
import { UsersController } from "../users-controller.ts";
import { UsersService } from "../users-service.ts";
import { UsersRepository } from "../users-repository.ts";
import { AuthService } from "../../auth/auth-service.ts";
import { getTestClient, cleanupDatabase } from "../../shared/test/test-utils.ts";
import { UserErrorCode } from "../../shared/enums/error-codes.enum.ts";

let client: Client;
let repository: UsersRepository;
let authService: AuthService;
let service: UsersService;
let controller: UsersController;

describe("Users Controller", () => {
    beforeAll(async () => {
        client = await getTestClient();
        await client.connect();
        repository = new UsersRepository(client);
        authService = new AuthService();
        service = new UsersService(repository, authService);
        controller = new UsersController(service);
    });

    afterAll(async () => {
        await cleanupDatabase(client);
        await client.end();
    });

    describe("createUser", () => {
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
            // Create second user with same email
            const req = {
                body: {
                    username: "uniqueuser",
                    email: "valid@test.com", // Same as previous test
                    password: "ValidPass123!"
                }
            } as unknown as Request;

            const response = await controller.createUser(req);

            assertEquals(response.success, false);
            assertEquals(response.status, 409);
            assertEquals(response.error, UserErrorCode.EMAIL_ALREADY_EXISTS);
        });

        it("should return 409 for duplicate username", async () => {
            const req = {
                body: {
                    username: "validuser", // Same as first test
                    email: "another@test.com",
                    password: "ValidPass123!"
                }
            } as unknown as Request;

            const response = await controller.createUser(req);

            assertEquals(response.success, false);
            assertEquals(response.status, 409);
            assertEquals(response.error, UserErrorCode.USERNAME_ALREADY_EXISTS);
        });
    });
});
