import { assert, assertEquals, assertNotEquals } from "https://deno.land/std@0.208.0/assert/mod.ts";
import { UsersRepository } from "../users-repository.ts";
import { getTestClient, cleanupDatabase } from "../../shared/test/test-utils.ts";
import { beforeAll, afterAll, describe, it } from "https://deno.land/std@0.208.0/testing/bdd.ts";
import { Client } from "postgres";

let client: Client;
let repository: UsersRepository;

describe("Users Repository", () => {
    beforeAll(async () => {
        client = await getTestClient();
        await client.connect();
        repository = new UsersRepository(client);
    });

    afterAll(async () => {
        await cleanupDatabase(client);
        await client.end();
    });

    it("should create a user", async () => {
        const userData = {
            username: "testuser",
            email: "test@example.com",
            password_hash: "hashed_password"
        };

        const user = await repository.create(userData);
        assert(user.id);
        assertEquals(user.username, userData.username);
        assertEquals(user.email, userData.email);
    });

    it("should find user by email", async () => {
        const user = await repository.findByEmail("test@example.com");
        assert(user);
        assertEquals(user?.email, "test@example.com");
    });

    it("should find user by username", async () => {
        const user = await repository.findByUsername("testuser");
        assert(user);
        assertEquals(user?.username, "testuser");
    });

    it("should update user", async () => {
        const user = await repository.findByEmail("test@example.com");
        assert(user);

        const updatedUser = await repository.update(user.id, {
            username: "updated_username"
        });

        assertEquals(updatedUser.username, "updated_username");
        assertNotEquals(updatedUser.updated_at, user.updated_at);
    });

    it("should soft delete user", async () => {
        const user = await repository.findByUsername("updated_username");
        assert(user);

        await repository.softDelete(user.id);
        const deletedUser = await repository.findById(user.id);
        assertEquals(deletedUser, null);
    });
});