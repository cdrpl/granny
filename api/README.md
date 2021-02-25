# API

API is a game server and can be thought of as two services called "api" and "lobby".

- api - exposes HTTP routes for user sign up and sign in.
- lobby - is both an HTTP and WebSocket server which allows users to host and join "rooms".

### Documentation

- [API documentation](https://documenter.getpostman.com/view/12308444/T1LLE7wE)

### Dependencies

- [Node.js](https://nodejs.org/en/) (15.x)
- [Postgres](https://www.postgresql.org/) (13.x)
- [Redis](https://redis.io/) (6.x)

### Environment Variables

Make a copy of the .env.example file and name it .env, the .env file will be loaded during startup.

- `NODE_ENV` set this to production when deploying.
- `PORT` this is the port that the API will bind to.
- `LOBBY_PORT` this is the port that the lobby server will bind to.
- `DB_HOST` this is the IP address of the PostgreSQL database.
- `DB_USER` this is the username used to connect to PostgreSQL.
- `DB_PASS` this is the password used to connect to PostgreSQL.
- `REDIS_HOST` this is the IP address of the Redis server.

### Run with Docker

The server can be run as a Docker container. The container will need access to a Redis and Postgres server, the addresses can be set by using a .env file and passing it to the docker run command. Note that `--network host` is not supported on windows, the ports 3000 and 3001 need to be exposed by replacing `--network host` in the run command with `-p 3000:3000 -p 3001:3001`

1. Build Image - `docker build -t api .`
2. Run Container - `docker run -itd --env-file .env --restart always --name api --network host api`

### Docker Compose

Docker Compose is also supported, just run the following command.

- `docker-compose up -d`

### Authentication

When a user logs in with valid credentials a 32 character "auth token" is generated. The token is then stored in Redis using the user's ID as the key with an expiration time set. Any route that requires authentication needs to be passed authentication credentials. The credentials should be included in the Authorization header using the format of "Authorization: userId:token". Note that the Authorization header does not contain any authentication scheme.

### Client

A client has been made with Unity3D. You can make your own client or use the official one.

- [Official Client](https://github.com/cdrpl/client)