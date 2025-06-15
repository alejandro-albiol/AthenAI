import { assert, assertEquals, assertNotEquals } from "std/assert/mod.ts";
import { beforeAll, afterAll, describe, it, beforeEach } from "testing/bdd.ts";
import { Client } from "postgres";
import { UsersRepository } from "../users-repository.ts";
import { getTestClient, cleanupDatabase } from "../../shared/test/test-utils.ts";

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

const _assertRecentDate = (date: Date, message = "Date should be recent") => {
    const now = new Date().getTime();
    const timestamp = date.getTime();
    const fiveMinutesAgo = now - (5 * 60 * 1000);
    assert(timestamp > fiveMinutesAgo && timestamp <= now, message);
};


describe("Users Repository", () => {
    let client: Client;
    let repository: UsersRepository;

    beforeAll(async () => {
        client = await getTestClient();
        await client.connect();
        repository = new UsersRepository(client);
    });

    afterAll(async () => {
        await cleanupDatabase(client);
        await client.end();
    });

    beforeEach(async () => {
        await cleanupDatabase(client);
    });

    describe("create", () => {
        it("should create a user and return correct fields", async () => {
            const userData = {
                username: "testuser",
                email: "test@example.com",
                password_hash: "hashed_password"
            };

            const user = await repository.create(userData);

            assert(user.id, "User should have an ID");
            assertEquals(user.username, userData.username);
            assertEquals(user.email, userData.email);
            assert(!("password_hash" in user), "Should not return password hash");
            assert(user.created_at instanceof Date, "Should have a created_at date");
        });

        it("should create users with unique IDs", async () => {
            const user1 = await repository.create({
                username: "user1",
                email: "user1@test.com",
                password_hash: "hash1"
            });

            const user2 = await repository.create({
                username: "user2",
                email: "user2@test.com",
                password_hash: "hash2"
            });

            assertNotEquals(user1.id, user2.id);
        });

        it("should handle concurrent user creations", async () => {
            const userData = Array.from({ length: 5 }, (_, i) => ({
                username: `concurrent${i}`,
                email: `concurrent${i}@test.com`,
                password_hash: "hash"
            }));

            const users = await Promise.all(
                userData.map(data => repository.create(data))
            );

            assertEquals(users.length, 5);
            // Verify all users have unique IDs
            const userIds = users.map(u => u.id);
            assertEquals(new Set(userIds).size, 5);
        });

        it("should preserve case sensitivity in username and email", async () => {
            const user1 = await repository.create({
                username: "TestUser",
                email: "Test@Example.com",
                password_hash: "hash"
            });

            assertEquals(user1.username, "TestUser");
            assertEquals(user1.email, "Test@Example.com");

            // Verify that a different case version is still treated as unique
            const user2 = await repository.create({
                username: "testuser",
                email: "test@example.com",
                password_hash: "hash"
            });

            assertNotEquals(user1.id, user2.id);
            assertEquals(user2.username, "testuser");
            assertEquals(user2.email, "test@example.com");
        });
    });

    describe("findById", () => {
        let testUserId: string;

        beforeEach(async () => {
            const user = await repository.create({
                username: "findbyid",
                email: "findbyid@test.com",
                password_hash: "hash"
            });
            testUserId = user.id;
        });

        it("should find user by ID", async () => {
            const user = await repository.findById(testUserId);

            assert(user);
            assertEquals(user.id, testUserId);
            assertEquals(user.username, "findbyid");
            assertEquals(user.email, "findbyid@test.com");
        });

        it("should return null for non-existent ID", async () => {
            const user = await repository.findById("00000000-0000-0000-0000-000000000000");
            assertEquals(user, null);
        });

        it("should not return soft-deleted users", async () => {
            await repository.softDelete(testUserId);
            const user = await repository.findById(testUserId);
            assertEquals(user, null);
        });
    });

    describe("findByEmail", () => {
        beforeEach(async () => {
            await repository.create({
                username: "emailuser",
                email: "Find@Test.com",  // Note the case
                password_hash: "hash"
            });
        });

        it("should find user by email exactly", async () => {
            const user = await repository.findByEmail("Find@Test.com");
            assert(user);
            assertEquals(user.email, "Find@Test.com");
        });

        it("should handle case sensitivity in email search", async () => {
            const user1 = await repository.findByEmail("find@test.com");
            assertEquals(user1, null, "Should not find email with different case");

            const user2 = await repository.findByEmail("Find@Test.com");
            assert(user2, "Should find email with exact case match");
        });

        it("should return null for empty email", async () => {
            const user = await repository.findByEmail("");
            assertEquals(user, null);
        });

        it("should return null for non-existent email", async () => {
            const user = await repository.findByEmail("nonexistent@test.com");
            assertEquals(user, null);
        });

        it("should not find soft-deleted users", async () => {
            const user = await repository.findByEmail("Find@Test.com");
            assert(user);
            await repository.softDelete(user.id);
            
            const deletedUser = await repository.findByEmail("Find@Test.com");
            assertEquals(deletedUser, null);
        });
    });

    describe("findByUsername", () => {
        beforeEach(async () => {
            await repository.create({
                username: "TestUser",  // Note the case
                email: "username@test.com",
                password_hash: "hash"
            });
        });

        it("should find user by username exactly", async () => {
            const user = await repository.findByUsername("TestUser");
            assert(user);
            assertEquals(user.username, "TestUser");
        });

        it("should handle case sensitivity in username search", async () => {
            const user1 = await repository.findByUsername("testuser");
            assertEquals(user1, null, "Should not find username with different case");

            const user2 = await repository.findByUsername("TestUser");
            assert(user2, "Should find username with exact case match");
        });

        it("should return null for empty username", async () => {
            const user = await repository.findByUsername("");
            assertEquals(user, null);
        });

        it("should return null for non-existent username", async () => {
            const user = await repository.findByUsername("nonexistent");
            assertEquals(user, null);
        });

        it("should not find soft-deleted users", async () => {
            const user = await repository.findByUsername("TestUser");
            assert(user);
            await repository.softDelete(user.id);
            
            const deletedUser = await repository.findByUsername("TestUser");
            assertEquals(deletedUser, null);
        });
    });

    describe("update", () => {
        let testUserId: string;

        beforeEach(async () => {
            const user = await repository.create({
                username: "updateuser",
                email: "update@test.com",
                password_hash: "hash"
            });
            testUserId = user.id;
        });

        it("should update single field", async () => {
            const updated = await repository.update(testUserId, {
                username: "newusername"
            });

            assertEquals(updated.username, "newusername");
            assertEquals(updated.email, "update@test.com"); // Should remain unchanged
            assertNotEquals(updated.updated_at, updated.created_at);
        });

        it("should update multiple fields", async () => {
            const updated = await repository.update(testUserId, {
                username: "newname",
                email: "newemail@test.com"
            });

            assertEquals(updated.username, "newname");
            assertEquals(updated.email, "newemail@test.com");
        });

        it("should preserve unchanged fields", async () => {
            const original = await repository.findById(testUserId);
            assert(original);

            const updated = await repository.update(testUserId, {
                username: "changedname"
            });

            assertEquals(updated.email, original.email);
            assertEquals(updated.password_hash, original.password_hash);
        });
    });

    describe("softDelete", () => {
        let testUserId: string;

        beforeEach(async () => {
            const user = await repository.create({
                username: "deleteuser",
                email: "delete@test.com",
                password_hash: "hash"
            });
            testUserId = user.id;
        });

        it("should mark user as deleted", async () => {
            await repository.softDelete(testUserId);

            // Verify user can't be found through normal methods
            const deletedUser = await repository.findById(testUserId);
            assertEquals(deletedUser, null);

            // Verify user still exists in DB but is marked as deleted
            const result = await client.queryObject<{ is_deleted: boolean; deleted_at: Date | null; }>(`
                SELECT is_deleted, deleted_at 
                FROM users 
                WHERE id = $1
            `, [testUserId]);

            assert(result.rows[0]);
            assertEquals(result.rows[0].is_deleted, true);
            assert(result.rows[0].deleted_at instanceof Date);
        });

        it("should not affect other users", async () => {
            const otherUser = await repository.create({
                username: "otheruser",
                email: "other@test.com",
                password_hash: "hash"
            });

            await repository.softDelete(testUserId);

            const survivingUser = await repository.findById(otherUser.id);
            assert(survivingUser);
            assertEquals(survivingUser.username, "otheruser");
        });
    });

    describe("list", () => {
        beforeEach(async () => {
            // Create test users with delay to ensure order
            for (let i = 1; i <= 5; i++) {
                await repository.create({
                    username: `user${i}`,
                    email: `user${i}@test.com`,
                    password_hash: "hash"
                });
                await delay(100); // Add delay between creations
            }
        });

        it("should list all users with default pagination", async () => {
            const users = await repository.list();
            assertEquals(users.length, 5);
        });

        it("should respect limit parameter", async () => {
            const users = await repository.list(2);
            assertEquals(users.length, 2);
        });

        it("should respect offset parameter", async () => {
            const users = await repository.list(undefined, 3);
            assertEquals(users.length, 2); // Should return last 2 users
        });

        it("should handle limit and offset together", async () => {
            const users = await repository.list(2, 2);
            assertEquals(users.length, 2);
        });

        it("should only return non-deleted users", async () => {

            const firstUser = (await repository.list(1))[0];
            await repository.softDelete(firstUser.id);

            const users = await repository.list();
            assertEquals(users.length, 4);
        });

        it("should return users ordered by creation date DESC", async () => {
            const users = await repository.list();
            assertEquals(users.length, 5);
            
            // In DESC order, later created users should come first
            for (let i = 1; i < users.length; i++) {
                const prevDate = new Date(users[i-1].created_at).getTime();
                const currDate = new Date(users[i].created_at).getTime();
                assert(
                    prevDate >= currDate,
                    `User at index ${i-1} should have a later or equal creation date than user at index ${i}`
                );
            }
        });
    });
});