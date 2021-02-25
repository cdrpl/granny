import type WebSocket from "ws";

/**
 * Represents a user from the database table.
 */
interface User {
  id?: number;
  name?: string;
  email?: string;
  pass?: string;
  created_on?: Date;
}

/**
 * Represents connected users with a WebSocket reference.
 */
interface UserClient {
  id: number;
  name: string;
  ws: WebSocket;
}

export { User, UserClient };
