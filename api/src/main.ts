import { createApp } from "./app.js";
import { knex } from "./knex.js";
import { router } from "./router.js";

// Run migrations
await knex.migrate.latest();

// Create app
const app = createApp(router);

// Start listening
app.listen(3000, () => {
  console.log("api listening on port 3000");
});
