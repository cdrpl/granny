# Server

This is the game server.

### Dependencies

- [Go](https://golang.org/) (1.16.x)
- [Postgres](https://www.postgresql.org/) (13.x)
- [Redis](https://redis.io/) (6.x)

### Environment Variables

Env vars will be loaded from the .env file if present. Env vars loaded from the .env file will never overwite existing variables.

- `ENV` set this to production when deploying.
- `DB_HOST` this is the IP address of the PostgreSQL database.
- `DB_USER` this is the username used to connect to PostgreSQL.
- `DB_PASS` this is the password used to connect to PostgreSQL.
- `REDIS_HOST` this is the IP address of the Redis server.

### Run with Docker

The server can be run in a Docker container. The container will need access to a Redis and Postgres server, the addresses can be set by using a .env file and passing it to the docker run command.

1. Build image - `docker build -t server .`
2. Run container - `docker run -itd --env-file .env --restart always --network host --name server server`
