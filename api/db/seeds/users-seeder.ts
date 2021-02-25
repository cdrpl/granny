import type { Knex } from "knex";

async function seed(knex: Knex) {
  await knex("users").del();

  await knex("users").insert({
    name: "root",
    email: "root@root.com",
    pass: "c52972808bf407c6876c726a0b3e2267:bc5fd0bcd3e2e421c5ee8dcbaa987bb9afde22407be5d7b98b98de96a9340b3d", // Plaintext == password
  });
}

export { seed };
