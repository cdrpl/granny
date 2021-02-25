import { knex } from "./knex.js";
import { createToken, hashText, compareHash } from "./auth.js";
import type { Request, Response, NextFunction } from "express";
import type { User } from "./types.js";

/**
 * This route is used to create a new user. It accepts a name, email, and password.
 */
async function signUp(req: Request, res: Response, next: NextFunction) {
  try {
    const { name, email } = req.body;
    let { pass } = req.body;

    // Name must be unique
    let user: User = await knex("users").select("name").where({ name }).first();
    if (user) {
      res.json({ error: { code: 1, message: "Name is not available" } });
      return;
    }

    // Email must be unique
    user = await knex("users").select("email").where({ email }).first();
    if (user) {
      res.json({ error: { code: 2, message: "Email already exists" } });
      return;
    }

    // Hash the user password
    pass = await hashText(pass);

    // Insert user
    await knex("users").insert({ name, email, pass });

    // Response
    res.json({ data: null });
  } catch (e) {
    next(e);
  }
}

/**
 * Returns the user ID and an authentication token if the provided credentials are valid.
 */
async function signIn(req: Request, res: Response, next: NextFunction) {
  try {
    const { email, pass } = req.body;

    // Get the user password
    const user: User = await knex("users").select("id", "name", "pass").where({ email }).first();
    if (!user) {
      res.status(401).json({ error: { status: 401, message: "Invalid credentials" } });
      return;
    }

    // Compare hash
    const validPass = await compareHash(pass, user.pass);
    if (!validPass) {
      res.status(401).json({ error: { status: 401, message: "Invalid credentials" } });
      return;
    }

    // Create auth token
    const token = await createToken(user.id);

    // Response
    const { id, name } = user;
    res.json({ data: { id, name, token } });
  } catch (e) {
    next(e);
  }
}

export { signUp, signIn };
