import crypto from "crypto";
import util from "util";

/**
 * Returns a random hex token with (bytes * 2) number of characters.
 */
async function randToken(bytes: number): Promise<string> {
  const asyncRand = util.promisify(crypto.randomBytes);
  const rand = await asyncRand(bytes);

  return rand.toString("hex");
}

export { randToken };
