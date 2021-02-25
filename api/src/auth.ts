import crypto from "crypto";
import { randToken } from "./rand.js";
import { redisGet, redisSetEx } from "./redis.js";
import type { Request, Response, NextFunction } from "express";

// Auth token
const AUTH_TOKEN_BYTES = 16;
const TOKEN_TTL = 60 * 60; // 1 hour

// Scrypt
const SALT_BYTES = 16;
const SCRYPT_KEY_LEN = 32;

/**
 * Hashes the given text using scrypt.
 */
async function hashText(text: string): Promise<string> {
  const salt = await randToken(SALT_BYTES);

  return new Promise((resolve, reject) => {
    crypto.scrypt(text, salt, SCRYPT_KEY_LEN, (err, derivedKey) => {
      if (err) {
        return reject(err);
      }

      const hash = `${salt}:${derivedKey.toString("hex")}`;
      resolve(hash);
    });
  });
}

/**
 * Compare the text with the hash and return true if matching.
 */
async function compareHash(text: string, hash: string): Promise<boolean> {
  const [salt, key] = hash.split(":");

  return new Promise((resolve, reject) => {
    crypto.scrypt(text, salt, SCRYPT_KEY_LEN, (err, derivedKey) => {
      if (err) {
        return reject(err);
      }

      resolve(key === derivedKey.toString("hex"));
    });
  });
}

/**
 * Generates a random token and stores it in redis using the user id as the key. The token will expire after a set time.
 */
async function createToken(userId: number): Promise<string> {
  const token = await randToken(AUTH_TOKEN_BYTES);
  await redisSetEx(userId.toString(), token, TOKEN_TTL);
  return token;
}

/**
 * Returns true if the given userId and token are valid.
 * For a token to be considered valid, it must be stored in redis using the user ID as the key.
 */
async function tokenIsValid(userId: string, token: string): Promise<boolean> {
  const redisVal = await redisGet(userId);
  return redisVal === token;
}

/**
 * Middleware that requires authentication details in the request header.
 */
async function requireAuth(req: Request, res: Response, next: NextFunction) {
  const { authorization } = req.headers;

  if (authorization === undefined) {
    return res.status(401).json({ error: { code: 401, message: "Unauthorized" } });
  }

  const [id, token] = authorization.split(":");

  if (id !== undefined && token !== undefined) {
    const isAuthorized = await tokenIsValid(id, token);

    if (isAuthorized) {
      req.userId = parseInt(id);
      return next();
    }
  }

  res.status(401).json({ error: { code: 401, message: "Unauthorized" } });
}

export { hashText, compareHash, createToken, tokenIsValid, requireAuth };
