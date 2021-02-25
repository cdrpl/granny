import { Knex } from "knex";

async function up(knex: Knex) {
  await knex.schema.createTable("users", (table) => {
    table.increments("id");
    table.string("name", 16).notNullable().unique();
    table.string("email", 255).notNullable().unique();
    table.string("pass", 255).notNullable();
    table.timestamp("created_on").notNullable().defaultTo(knex.fn.now());
  });
}

async function down(knex: Knex) {
  await knex.schema.dropTable("users");
}

export { up, down };
