import { validationResult, body } from "express-validator";
import type { Request, Response, NextFunction } from "express";

/**
 * Express-validator middleware.
 */
function validate(req: Request, res: Response, next: NextFunction): void {
  const errors = validationResult(req);

  if (errors.isEmpty()) {
    next();
  } else {
    res.status(400).json({ error: { code: 400, message: errors.array()[0].msg } });
  }
}

/**
 * Validation rules for the sign up route.
 */
const signUpRules = [
  body("name", "name is required").exists().trim(),
  body("name", "name max length is 16").isLength({ max: 16 }),
  body("email", "email is required").exists().trim().normalizeEmail(),
  body("email", "email is too long").isLength({ max: 255 }),
  body("email", "email is not valid").isEmail(),
  body("pass", "password is required").exists(),
  body("pass", "password must have 8 characters").isLength({ min: 8 }),
  body("pass", "password is too long").isLength({ max: 255 }),
];

/**
 * Validation rules for the sign in route.
 */
const signInRules = [
  body("email", "email is required").exists().trim().normalizeEmail(),
  body("email", "email is too long").isLength({ max: 255 }),
  body("email", "email is not valid").isEmail(),
  body("pass", "password is required").exists(),
  body("pass", "password must have 8 characters").isLength({ min: 8 }),
  body("pass", "password is too long").isLength({ max: 255 }),
];

export { validate, signUpRules, signInRules };
