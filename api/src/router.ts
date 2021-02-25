import express from "express";
import * as controller from "./controller.js";
import { validate, signUpRules, signInRules } from "./validator.js";

const router = express.Router();

// Health check route
router.get("/", (req, res) => {
  res.sendStatus(200);
});

router.post("/sign-up", signUpRules, validate, controller.signUp);

router.post("/sign-in", signInRules, validate, controller.signIn);

export { router };
