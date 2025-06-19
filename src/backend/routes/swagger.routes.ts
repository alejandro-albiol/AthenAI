import { Router, send } from "oak";
import { swaggerDefinition } from "../docs/swagger.ts";

const router = new Router();

router
    .get("/swagger.json", (context) => {
        context.response.body = swaggerDefinition;
    })
    .get("/", async (context) => {
        await send(context, "index.html", {
            root: `${Deno.cwd()}/src/backend/docs/swagger-ui`
        });
    })
    .get("/(.*)", async (context) => {
        const path = context.params[0];
        await send(context, path, {
            root: `${Deno.cwd()}/src/backend/docs/swagger-ui`
        });
    });

export const swaggerRouter = router;
