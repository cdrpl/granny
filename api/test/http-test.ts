import supertest from "supertest";
import { assert } from "chai";
import { describe, it } from "mocha";
import { knex } from "../src/knex.js";
import { redisGet } from "../src/redis.js";
import { createApp } from "../src/app.js";
import { router } from "../src/router.js";

const app = createApp(router);

/**
 * Run the migrations down then up.
 */
async function migrate() {
  await knex.migrate.rollback();
  await knex.migrate.latest();
}

await migrate();

/**
 * Health check tests.
 */
describe("GET /", () => {
  it("returns 200", () => {
    return supertest(app).get("/").expect("Content-Length", "2").expect(200);
  });
});

/**
 * Sign up tests.
 */
describe("POST /sign-up", () => {
  it("adds the user to the database", async () => {
    const form = { name: "test", email: "test@test.com", pass: "password" };

    return supertest(app)
      .post("/sign-up")
      .type("form")
      .send(form)
      .expect(200)
      .then(async () => {
        // User must be in the database
        const user = await knex("users").where({ name: form.name }).first();
        assert.isDefined(user, "user could not be found");

        // Values must match
        assert.strictEqual(form.name, user.name);
        assert.strictEqual(form.email, user.email);

        // Password should be hashed
        assert.notStrictEqual(form.pass, user.pass);
      });
  });
});

/**
 * Sign in tests.
 */
describe("POST /sign-in", () => {
  it("returns 200 with valid credentials", () => {
    const form = { email: "test@test.com", pass: "password" };

    return supertest(app).post("/sign-in").type("form").send(form).expect(200);
  });

  it("returns 401 with invalid credentials", () => {
    const form = { email: "test@test.com", pass: "invalid-pass" };

    return supertest(app).post("/sign-in").type("form").send(form).expect(401);
  });

  it("returns an auth token and user id", () => {
    const form = { email: "test@test.com", pass: "password" };

    return supertest(app)
      .post("/sign-in")
      .type("form")
      .send(form)
      .expect(200)
      .then((res) => {
        assert.hasAllKeys(res.body.data, ["id", "name", "token"]);
      });
  });

  it("adds the auth token to redis", () => {
    const form = { email: "test@test.com", pass: "password" };

    return supertest(app)
      .post("/sign-in")
      .type("form")
      .send(form)
      .expect(200)
      .then(async (res) => {
        const { id, token } = res.body.data;

        const result = await redisGet(id);
        assert.strictEqual(token, result);
      });
  });
});
