# Overview

The undeployed versions of the code are meant for development and debugging purposes. Any branch that is not mapped to one of the deployed environments can be used to create a local development environment (the same can be done with `dev` and `master`, but since they are not intended for direct modification, we exclude them from the scope of this discussion).

# Setup and run instructions

Install Docker desktop software:
- [Install on macOS](https://docs.docker.com/desktop/install/mac-install/)
- [Install on Windows](https://docs.docker.com/desktop/install/windows-install/)
- [Install on Linux](https://docs.docker.com/desktop/install/linux-install/)

Install other CLI dependencies:
- [Stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer) allows for automated creation of methods that satisfy the `fmt.Stringer` interface. Run `go install golang.org/x/tools/cmd/stringer@latest`
- [Swagger](https://pkg.go.dev/github.com/swaggo/swag/v2/cmd/swag) spec generator `go install github.com/swaggo/swag/v2/cmd/swag@latest`

Create `.env` file based on content from `.env.sample` and edit it to define values of the listed variables:

- `ENV_EMAIL_ADDRESS` email address used to access the SMTP service provided by Gmail
- `ENV_EMAIL_PASSWORD` password associated with `ENV_EMAIL_ADDRESS` account
- `ENV_JWT_SECRET` passphrase for JWT signing/verification
- `ENV_RSC_TOKEN` API key for RSC API (sport data retrieval)
- `POSTGRES_USER` DB user to be setup and used for connecting microservices to DB
- `POSTGRES_PASSWORD` password of the DB user (arbitrary good password)
- `POSTGRES_DB` DB name (arbitrary good name)
- `AUTH_PORT` TCP port to run `auth-service`, should not conflict with existing host ports
- `DEMO_PORT` TCP port to run `demo-service`, should not conflict with existing host ports
- `IS_DEVELOPMENT` flag to be set to `"1"` when service is running in development mode

# Extra step for a fresh setup

Setting up the environment from scratch requires DB server to restart a few times before properly settling. That initialization step is recommended to be performed before running microservices relying on DB server. Run this command:
```sh
docker-compose up db
```
to run DB server without running other microservices. 

Wait for the process to settle down by observing that no more messages are printed to the log, and it shows content similar to the following:
```log
...
quible-api-db-1  | PostgreSQL init process complete; ready for start up.
quible-api-db-1  |
quible-api-db-1  | 2023-11-11 03:47:15.989 UTC [1] LOG:  starting PostgreSQL 16.0 (Debian 16.0-1.pgdg120+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
quible-api-db-1  | 2023-11-11 03:47:15.990 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
quible-api-db-1  | 2023-11-11 03:47:15.990 UTC [1] LOG:  listening on IPv6 address "::", port 5432
quible-api-db-1  | 2023-11-11 03:47:16.017 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
quible-api-db-1  | 2023-11-11 03:47:16.157 UTC [75] LOG:  database system was shut down at 2023-11-11 03:47:15 UTC
quible-api-db-1  | 2023-11-11 03:47:16.221 UTC [1] LOG:  database system is ready to accept connections
```
# Start the entire service (all microservices) in Docker

Run `docker-compose up --build` to build and start all microservices. The operation will run in the foreground and can be gracefully terminated by hitting `Ctrl+C`. All services will send their logs to the same console. The **service prefix** can help identify specific service logs. 

Below is the log of an exampled run:
```
...
Attaching to quible-api-auth-1, quible-api-db-1, quible-api-demo-1
quible-api-db-1    |
quible-api-db-1    | PostgreSQL Database directory appears to contain a database; Skipping initialization
quible-api-db-1    |
quible-api-db-1    | 2023-11-11 00:41:43.794 UTC [1] LOG:  starting PostgreSQL 16.0 (Debian 16.0-1.pgdg120+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
quible-api-db-1    | 2023-11-11 00:41:43.794 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
quible-api-db-1    | 2023-11-11 00:41:43.794 UTC [1] LOG:  listening on IPv6 address "::", port 5432
quible-api-db-1    | 2023-11-11 00:41:43.816 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
quible-api-db-1    | 2023-11-11 00:41:43.885 UTC [29] LOG:  database system was shut down at 2023-11-10 23:35:34 UTC
quible-api-db-1    | 2023-11-11 00:41:43.980 UTC [1] LOG:  database system is ready to accept connections
quible-api-demo-1  | 2023/11/11 00:42:12 starting server on port: 8010
quible-api-auth-1  | 2023/11/11 00:42:12 starting server on port: 8001

```
with `quible-api-db-1`, `quible-api-demo-1`, and `quible-api-auth-1` being log prefixes of individual services. 

## Which Dockerfile is used by individual microservices?

There are two ways to configure the composed Docker environment regarding which `Dockerfile` will be picked up by individual microservices.
- Shared `Dockerfile` from the repo root 
- Microservice-specific `Dockerfile` hosted at the root level of the microservice directory. 

The former option (shared config) can be activated by defining the `build` section of the specific microservice in the `docker-compose.yml` as follows:
```yaml
build: 
  context: ./auth-service
  dockerfile: ../Dockerfile
  additional_contexts: 
    lib: ./lib
```
Here, the `dockerfile: ../Dockerfile` references `Dockerfile` from the repo root.

Alternatively, especially when the microservice is not compatible with the shared `Dockerfile`, the build section can be defined as follows:
```yaml
build: 
  context: ./demo-service
  dockerfile: ./Dockerfile.sample
  additional_contexts: 
    lib: ./lib
```
where `dockerfile: ./Dockerfile.sample` references `Dockerfile.sample` file inside microservice directory. 

# Running selected service(s) outside Docker (for speed)

Running all microservices in Docker containers proves their inter-operational compatibility. But there are situations when local development is focused on aspects of one (or only a few) microservice(s), and it can be safely localized without affecting compatibility. In such cases, the overhead of re-building and re-running the docker containers could seriously degrade the development pace. 

The only service that is meant to continue running in Docker is DB. You can start it separately (without modifying any files) by running:
```sh
docker-compose up db
```
while having your previously started Docker setup gracefully stopped by hitting `Ctrl+C`.

Once the DB service is started in a Docker container, it exposes the `TCP:5432` port on the `localhost` interface. If you now start a microservice outside the Docker environment, it will still be able to connect to DB to access data. Under the hood, it will read `.env` file from the root of the repo and will compile the value of `ENV_DSN` environment variable from the variables exposed in that file (more on that can be seen [here](https://gitlab.com/quible/backend/api-server/-/blame/dev/lib/env/setup.go?ref_type=heads#L17))

## Start the service

Some services require an extra step to generate additional files. That step is included in the default `Dockerfile` but has to be explicitly run while working outside the Docker environment
```sh
go generate ./...
```

You might also want to ensure that all the prerequisites are met, e.g., `swagger.yaml` file is present in the microservice's root directory and all modules/dependencies are downloaded. 

Without loss of generality, start the service by running
```sh
go run .
```
