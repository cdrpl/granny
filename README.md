# Granny

Granny is an open source online game. Network communication is achieved through HTTP and WebSockets.

### Overview

Granny can be thought of as being split into 3 parts.

- [API](/api) - This is an API that exposes HTTP routes and responds with data in JSON format. Written in TypeScript and ran with Node.js
- [Client](/client) - This is the game client responsible for displaying the game to players as well as receiving input from them. Made with Unity3D.
- [Server](/server) - This is a WebSocket server that handles the realtime backend functionality. Written in Go.

### Languages

- C#
- Go
- TypeScript

### Technologies

- Docker
- NGINX
- Node.js
- PostgreSQL
- Redis
- Unity3D

### Docker Compose

The Granny backend can be easily setup with Docker Compose. Just install Docker and run the following command.

- `docker-compose up -d`
