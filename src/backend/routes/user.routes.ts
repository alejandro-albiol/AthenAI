import { Router } from "oak";
import { UsersController } from "../users/users-controller.ts";
import { UsersService } from "../users/users-service.ts";
import { UsersRepository } from "../users/users-repository.ts";
import { AuthService } from "../auth/auth-service.ts";
import { createClient } from "../db/config.ts";

// Initialize dependencies
const client = await createClient();
const authService = new AuthService();
const usersRepository = new UsersRepository(client);
const usersService = new UsersService(usersRepository, authService);
const usersController = new UsersController(usersService);

const router = new Router();

router
    .get("/", async (context) => {
        const limit = context.request.url.searchParams.get("limit");
        const offset = context.request.url.searchParams.get("offset");
        const response = await usersController.listUsers({
            query: {
                limit: limit || undefined,
                offset: offset || undefined
            }
        });
        context.response.status = response.status;
        context.response.body = response;
    })
    .post("/", async (context) => {        const reqBody = await context.request.body().value;
        const response = await usersController.createUser(new Request(context.request.url, {
            method: "POST",
            headers: context.request.headers,
            body: JSON.stringify(reqBody)
        }));
        context.response.status = response.status;
        context.response.body = response;
    })
    .get("/:id", async (context) => {
        const response = await usersController.findUserById({
            params: { id: context.params.id }
        });
        context.response.status = response.status;
        context.response.body = response;
    })
    .get("/username/:username", async (context) => {
        const response = await usersController.findUserByUsername({
            params: { username: context.params.username }
        });
        context.response.status = response.status;
        context.response.body = response;
    })
    .get("/email/:email", async (context) => {
        const response = await usersController.findUserByEmail({
            params: { email: context.params.email }
        });
        context.response.status = response.status;
        context.response.body = response;
    })
    .put("/:id", async (context) => {
        const body = await context.request.body().value;
        const response = await usersController.updateUser({
            params: { id: context.params.id },
            body
        });
        context.response.status = response.status;
        context.response.body = response;
    })
    .delete("/:id", async (context) => {
        const response = await usersController.deleteUser({
            params: { id: context.params.id }
        });
        context.response.status = response.status;
        if (response.status !== 204) {
            context.response.body = response;
        }
    });

export const userRouter = router;
