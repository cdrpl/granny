import knex from "knex";

/**
 * Single knex instance used throughout the app.
 */
const knexInstance = knex({
  client: "pg",
  connection: {
    host: process.env.DB_HOST || "127.0.0.1",
    user: process.env.DB_USER || "postgres",
    password: process.env.DB_PASS || "password",
    database: process.env.DB_USER || "postgres",
  },
  migrations: {
    tableName: "migrations",
    directory: "./dist/db/migrations",
  },
  seeds: {
    directory: "./dist/db/seeds",
  },
});

export { knexInstance as knex };
