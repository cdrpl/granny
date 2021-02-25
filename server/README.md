# Server

This is the Idlemon game server. It exposes HTTP routes for user sign up and sign in and includes a WebSocket server for realtime gameplay. The [client](https://github.com/cdrpl/client) is kept in a separate repository.

## Dependencies

- [Go](https://golang.org/) (1.15)
- [Postgres](https://www.postgresql.org/) (13.1)
- [Redis](https://redis.io/) (6.0)

## Documentation

- [API documentation](https://documenter.getpostman.com/view/12308444/T1LLE7wE)

## Environment Variables

The server will attempt to load environment variables from .env and then .env.defaults. Env vars loaded from the .env files will never overwite existing variables. Note that .env files will not be loaded if ENV is set to production, you will need to set the env variables manually when deploying to a production environment.

- `ENV` set this to production when deploying.
- `PORT` this is the port that the API will bind to.
- `DB_HOST` this is the IP address of the PostgreSQL database.
- `DB_USER` this is the username used to connect to PostgreSQL.
- `DB_PASS` this is the password used to connect to PostgreSQL.
- `REDIS_HOST` this is the IP address of the Redis server.

## Dockerfile

You can build and deploy the server manually or use Docker.

- `docker build -t server .`
- `docker run -itd --env-file .env --restart always -p 5000:5000 --name server server`

## Docker Compose

You can also run `docker-compose up` to build the server along with the Redis and Postgres dependencies.

## Authentication process

1. User logs in through the API.
2. The API generates an authentication token and stores it in Redis using the user ID as the key.
3. The user ID and auth token are returned to the client.
4. The client sends an HTTP WebSocket upgrade request to the server along with the user ID and auth token.
5. The server verifies the auth token then fetches the user's data.
6. The server initializes the user then returns the user data to the client.
7. The user is fully authenticated.

## Migrations

Simple migrations are currently supported. Migrations are stored in .sql files in the /migrations directory. The server will run every sql file during startup. Note that since we are still on 0.X.X the migrations will not track changes for existing tables and will only function as schema reconstruction. In order to reflect the latest db schema you will need to drop all tables or manually update existing tables until 1.X.X is reached.

## Data Persistence

Player data is stored in a Postgres database. The server will run an UPDATE query for every online user on set intervals. If the user has logged out, the data is kept around until the next save interval. After the UPDATE query has been run the user data will be removed from memory if they are still offline.
