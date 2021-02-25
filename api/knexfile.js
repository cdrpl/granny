import dotenv from "dotenv";

// Load env
dotenv.config();

export default {
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
};
