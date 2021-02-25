import WebSocket from "ws";
import { randToken } from "./rand.js";
import type { User } from "./types.js";

/**
 * Number of bytes in a room ID.
 */
const ROOM_ID_BYTES = 12;

/**
 * Max users allowed in a room.
 */
const ROOM_MAX_USERS = 5;

/**
 * Represents a lobby room.
 */
class Room {
  id: string;
  name: string;
  users: Record<number, RoomUser>;

  constructor(id: string, name: string, host: User) {
    this.id = id;
    this.name = name;
    this.users = {};
    this.addUser(host);
  }

  static async createRoom(name: string, host: User): Promise<Room> {
    const id = await randToken(ROOM_ID_BYTES);
    return new Room(id, name, host);
  }

  /**
   * Add user to the room.
   */
  addUser(user: User) {
    this.users[user.id] = {
      id: user.id,
      name: user.name,
      isReady: false,
    };
  }

  /**
   * Return the number of users in the room.
   */
  numUsers(): number {
    return Object.keys(this.users).length;
  }

  /**
   * Returns true if the room is full.
   */
  isFull(): boolean {
    return this.numUsers() >= ROOM_MAX_USERS;
  }

  broadcast(msg: any, clients: Record<number, WebSocket>) {
    const msgJson = JSON.stringify(msg);

    for (const key in this.users) {
      if (clients[key] !== undefined) {
        clients[key].send(msgJson);
      }
    }
  }
}

/**
 * Represents a user in a room.
 */
interface RoomUser {
  id: number;
  name: string;
  isReady: boolean;
}

export { Room, RoomUser };
