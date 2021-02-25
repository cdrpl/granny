import { createApp } from "./app.js";
import { knex } from "./knex.js";
import { router } from "./router.js";
import "./lobby.js"; // WebSocket server will be started from importing.

const port = process.env.PORT || 3000;

// Run migrations
await knex.migrate.latest();

// Create app
const app = createApp(router);

// Start listening
app.listen(port, () => {
  console.log(`api listening on port ${port}`);
});
