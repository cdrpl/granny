import express from "express";
import WebSocket from "ws";
import { tokenIsValid, requireAuth } from "./auth.js";
import { knex } from "./knex.js";
import { createApp } from "./app.js";
import { User } from "./types.js";
import { Room } from "./room.js";
import type { IncomingMessage } from "http";

/**
 * Rooms that users have hosted.
 */
const rooms: Record<string, Room> = {};

/**
 * Clients links a user ID to a WebSocket connection.
 */
const clients: Record<number, WebSocket> = {};

/**
 * HTTP router.
 */
const router = express.Router();

// Health check route
router.get("/", (req, res) => {
  res.sendStatus(200);
});

// Rooms index
router.get("/rooms", (req, res) => {
  const data = [];

  for (const key in rooms) {
    const room = rooms[key];
    data.push(room);
  }

  res.json({ data });
});

/**
 * Create a new room and return the room ID.
 */
router.post("/rooms", requireAuth, async (req, res) => {
  const { userId } = req;
  const { roomName } = req.body;

  // Grab user data
  const user: User = await knex("users").select("id", "name").where("id", userId).first();

  // Create the room
  const room = await Room.createRoom(roomName, user);
  rooms[room.id] = room;

  res.json({ data: room });
});

router.post("/rooms/join", requireAuth, async (req, res) => {
  const { userId } = req;
  const { roomId } = req.body;
  const room = rooms[roomId];

  // Grab user data
  const user: User = await knex("users").select("id", "name").where("id", userId).first();

  // Room must exist
  if (room !== undefined) {
    if (room.isFull()) {
      res.json({ error: { code: 1, message: "Room is full" } });
    } else {
      room.addUser(user);

      // Create joined message and broadcast to all users in the room
      const msg = { channel: "join", id: user.id, name: user.name };
      room.broadcast(msg, clients);

      res.json({ data: room });
    }
  } else {
    res.json({ error: { code: 2, message: "Room doesn't exist" } });
  }
});

//router.post("/rooms/leave", (req, res) => {});
//router.post("/rooms/ready", (req, res) => {});
//router.post("/rooms/start", (req, res) => {});

/**
 * HTTP server used for authentication.
 */
const port = process.env.LOBBY_PORT || 3001;
const app = createApp(router);
const server = app.listen(port, () => {
  console.log(`lobby listening on port ${port}`);
});

/**
 * WebSocket server.
 */
const wss = new WebSocket.Server({ noServer: true });

/**
 * WebSocket Authentication.
 */
server.on("upgrade", async (req: IncomingMessage, socket, head: Buffer) => {
  if (req.url !== "/ws") {
    socket.write("HTTP/1.1 404 Not Found\r\n\r\n");
    socket.destroy();
  }

  try {
    const id = req.headers.id || ""; // The user ID
    const token = req.headers.token || "";

    if (!Array.isArray(id) && !Array.isArray(token)) {
      const isAuthenticated = await tokenIsValid(id, token);

      if (isAuthenticated) {
        wss.handleUpgrade(req, socket, head, async (ws) => {
          clients[id] = ws;
          wss.emit("connection", ws);
        });

        return;
      }
    }

    socket.write("HTTP/1.1 401 Unauthorized\r\n\r\n");
    socket.destroy();
  } catch (e) {
    console.error(e);
  }
});

/**
 * On new client connected.
 */
wss.on("connection", (ws: WebSocket) => {
  // Heartbeat functionality
  ws["isAlive"] = true;
  ws.on("pong", heartbeat);

  // Message handler
  ws.on("message", (msg) => {
    console.log(`Received message ${msg}`);
  });
});

/**
 * Broadcast message to all connected clients.
 */
/*function broadcastAll(msg: any) {
  wss.clients.forEach((client) => {
    if (client.readyState === WebSocket.OPEN) {
      client.send(msg);
    }
  });
}*/

/**
 * Broadcast message to every client except the sender.
 */
/*function broadcastExclude(ws: WebSocket, msg: any) {
  wss.clients.forEach((client) => {
    if (client !== ws && client.readyState === WebSocket.OPEN) {
      client.send(msg);
    }
  });
}*/

/**
 * Heartbeat used to detect dead connections.
 */
function heartbeat() {
  this["isAlive"] = true;
}

/**
 * Interval between heartbeats.
 */
const interval = setInterval(() => {
  wss.clients.forEach((ws) => {
    if (ws["isAlive"] === false) return ws.terminate();

    ws["isAlive"] = false;
    ws.ping();
  });
}, 30000);

/**
 * Clears interval when WebSocket server is closed.
 */
wss.on("close", () => {
  clearInterval(interval);
});
