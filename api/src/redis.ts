import util from "util";
import redis from "redis";

const host = process.env.REDIS_HOST || "127.0.0.1";

/**
 * The Redis client.
 */
const client = redis.createClient({ host });

/**
 * Async wrapper around Redis get.
 */
function redisGet(key: string): Promise<string> {
  const asyncGet = util.promisify(client.get).bind(client);
  return asyncGet(key);
}

/**
 * Sets a Redis key with an expiration.
 */
async function redisSetEx(key: string, val: string, ttl: number): Promise<void> {
  const asyncSet = util.promisify(client.set).bind(client);
  await asyncSet(key, val, "EX", ttl);
}

export { redisGet, redisSetEx };
