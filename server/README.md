# Server

The server handles realtime functionality and network communication is achieved through WebSockets.

### Dependencies

- [Go](https://golang.org/) (1.16.x)
- [Postgres](https://www.postgresql.org/) (13.x)
- [Redis](https://redis.io/) (6.x)

### Environment Variables

The server will attempt to load environment variables from .env and then .env.defaults. Env vars loaded from the .env files will never overwite existing variables.

- `ENV` set this to production when deploying.
- `DB_HOST` this is the IP address of the PostgreSQL database.
- `DB_USER` this is the username used to connect to PostgreSQL.
- `DB_PASS` this is the password used to connect to PostgreSQL.
- `REDIS_HOST` this is the IP address of the Redis server.

### Run with Docker

The server can be run as a Docker container. The container will need access to a Redis and Postgres server, the addresses can be set by using a .env file and passing it to the docker run command. Note that `--network host` is not supported on windows, port 3010 needs to be exposed by replacing `--network host` in the run command with `-p 3010:3010`

1. Build image - `docker build -t server .`
2. Run container - `docker run -itd --env-file .env --restart always --network host --name server server`
