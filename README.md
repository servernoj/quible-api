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

- `ENV_JWT_SECRET` passphrase for JWT signing/verification
- `ENV_RSC_TOKEN` API key for RSC API (sport data retrieval)
- `POSTGRES_USER` DB user to be setup and used for connecting microservices to DB
- `POSTGRES_PASSWORD` password of the DB user (arbitrary good password)
- `POSTGRES_DB` DB name (arbitrary good name)
- `AUTH_PORT` TCP port to run `auth-service`, should not conflict with existing host ports
- `APP_PORT` TCP port to run `app-service`, should not conflict with existing host ports

# Database migrations

Database migrations are defined in the `cmd` module. They are automatically executed up to the highest version when:
- one of the microservices is started in docker environment (docker hosted DB is targeted/used for migrations)
- one of the microservices is deployed to Render (external DB hosted in Render is targeted/used for migrations)

We use `goose` library to perform DB migrations. It supports 2 kinds of migration files:
- pure SQL (for simple migrations)
- Go sources (for complex use cases)

Every migration file satisfies the following conditions
- it is stored in `cmd/migrations`
- it is generated by a migration script `cmd/run.sh migrate create <migration name> <migration type>` and should not be edited unless it is a Go migration
- it has a filename containing `id` (timestamp) and `name` (matching `<migration name>` from the above command). The `id` portion is used to keep the track of applied migrations. The `name` portion is ignored and is used for readability

## Working with migration files

To create/apply/undo/redo migrations use the shell script `cmd/run.sh migrate` and feed it with the commands followed by command's arguments (some usage examples can be found in the official `goose` [documentation](https://github.com/pressly/goose/blob/v3.16.0/README.md#usage)):
- `up` to apply all available migrations 
- `up-to` migrate up to a specific version (timestamp or migration `id`)
- `up-by-one` migrate up a single migration from the current version
- `down` roll back a single migration from the current version  
- `down-to` roll back migrations to a specific version (use `0` as a pseudo version number, corresponding to the initial DB state)
- `redo` roll back the most recently applied migration, then run it again
- `status` print the status of all migrations
- `version` print the current version of the database:

**Important:** the migration utility uses environment variable `ENV_DSN` to locate the DB instance and connect to it. This variable should either be "exported" explicitly or otherwise its value could be implied from `.env` file located in the repo's root. 

## Special Docker service 

A special docker service `migrations` is defined to perform all pending DB migrations every time **when** it is executed explicitly:
```sh
docker-compose up migrations --build
```
or one of the microservices is started in docker environment (see below). 

## Example: add new migration, apply it, and roll it back

Let's say we have a DB with 2 applied migrations. To see those applied migrations we need
1. Run DB service in docker environment: `docker-compose up db` 
1. Run the following command (assuming your current directory matches the repo root)
    ```sh
    cmd/run.sh migrate status
    ```
1. It should print something similar to
    ```sh
    2023/12/04 18:04:54     Applied At                  Migration
    2023/12/04 18:04:54     =======================================
    2023/12/04 18:04:54     Tue Dec  5 00:20:11 2023 -- 20231204150816_init.sql
    2023/12/04 18:04:54     Tue Dec  5 00:20:11 2023 -- 20231204164132_add_image_column_to_users.sql
    ```
    from where we see that 2 migrations have been previously applied and no other migrations are pending

Having initial state of the DB confirmed, we need to keep DB server running in the container and run the following command to create a placeholder for a new SQL migration
```sh
cmd/run.sh migrate create mytest sql
```

This command creates a new file named `xxx_mytest.sql` in `cmd/migrations`. We can re-run the `status` command to see it in `pending status`
```sh
cmd/run.sh migrate status
2023/12/04 18:12:18     Applied At                  Migration
2023/12/04 18:12:18     =======================================
2023/12/04 18:12:18     Tue Dec  5 00:20:11 2023 -- 20231204150816_init.sql
2023/12/04 18:12:18     Tue Dec  5 00:20:11 2023 -- 20231204164132_add_image_column_to_users.sql
2023/12/04 18:12:18     Pending                  -- 20231204181114_mytest.sql
```

The migration file _usually_ needs to be edited to implement the migration logic, but even without editing it can be applied and rolled back as follows:

To apply run:
```sh
cmd/run.sh migrate up
2023/12/04 18:14:34 OK   20231204181114_mytest.sql (28.74ms)
2023/12/04 18:14:34 goose: successfully migrated database to version: 20231204181114
```

To roll back run
```sh
cmd/run.sh migrate down
2023/12/04 18:15:12 OK   20231204181114_mytest.sql (9.43ms)
```

Then to confirm the original state run
```sh
cmd/run.sh migrate status
2023/12/04 18:16:06     Applied At                  Migration
2023/12/04 18:16:06     =======================================
2023/12/04 18:16:06     Tue Dec  5 00:20:11 2023 -- 20231204150816_init.sql
2023/12/04 18:16:06     Tue Dec  5 00:20:11 2023 -- 20231204164132_add_image_column_to_users.sql
2023/12/04 18:16:06     Pending                  -- 20231204181114_mytest.sql
```
which shows that our newly introduced migration is in pending state (since we rolled it back after applying)

# Start the entire service (all microservices) in Docker

Run `docker-compose up --build` to build and start all microservices. The operation will run in the foreground and can be gracefully terminated by hitting `Ctrl+C`. All services will send their logs to the same console. The **service prefix** can help identify specific service logs. 

Below is the log of an exampled run:
```
...
[+] Running 4/4
 ✔ Container quible-api-db-1          Created
 ✔ Container quible-api-migrations-1  Recreated
 ✔ Container quible-api-auth-1        Recreated
 ✔ Container quible-api-app-1         Recreated
Attaching to quible-api-app-1, quible-api-auth-1, quible-api-db-1, quible-api-migrations-1
quible-api-db-1          |
quible-api-db-1          | PostgreSQL Database directory appears to contain a database; Skipping initialization
quible-api-db-1          |
quible-api-db-1          | 2023-12-05 01:25:40.458 UTC [1] LOG:  starting PostgreSQL 16.1 (Debian 16.1-1.pgdg120+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
quible-api-db-1          | 2023-12-05 01:25:40.458 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
quible-api-db-1          | 2023-12-05 01:25:40.459 UTC [1] LOG:  listening on IPv6 address "::", port 5432
quible-api-db-1          | 2023-12-05 01:25:40.464 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
quible-api-db-1          | 2023-12-05 01:25:40.503 UTC [31] LOG:  database system was shut down at 2023-12-05 01:25:28 UTC
quible-api-db-1          | 2023-12-05 01:25:40.557 UTC [1] LOG:  database system is ready to accept connections
quible-api-migrations-1  | 2023/12/05 01:26:51 goose: no migrations to run. current version: 20231204164132
quible-api-migrations-1 exited with code 0
quible-api-auth-1        | 2023/12/05 01:26:51 running in docker...
quible-api-auth-1        | 2023/12/05 01:26:52 starting server on port: 8001
quible-api-app-1         | 2023/12/05 01:26:52 running in docker...
quible-api-app-1         | 2023/12/05 01:26:52 starting server on port: 8002

```
with `quible-api-db-1`, `quible-api-auth-1`, `quible-api-app-1` and `quible-api-migrations-1` being log prefixes of individual services. 

## Which Dockerfile is used by individual microservices?

There are two ways to configure the composed Docker environment regarding which `Dockerfile` will be picked up by individual microservices.
- Shared `Dockerfile` from the repo root 
- Microservice-specific `Dockerfile` hosted at the root level of the microservice directory. 

The former option (shared config) can be activated by defining the `build` section of the specific microservice in the `docker-compose.yml` as follows (see `auth` service for example):
```yaml
build: 
  context: ./auth-service
  dockerfile: ../Dockerfile
  ...
```
Here, the `dockerfile: ../Dockerfile` references `Dockerfile` from the repo root.

Alternatively, especially when the microservice is not compatible with the shared `Dockerfile`, the build section can be defined as follows (see `migrations` service for example):
```yaml
build: 
  context: ./demo-service
  dockerfile: ./Dockerfile
  ...
```
where `dockerfile: ./Dockerfile` references `Dockerfile` file inside microservice directory. 

# Running selected service(s) outside Docker (for speed)

Running all microservices in Docker containers proves their inter-operational compatibility. But there are situations when local development is focused on aspects of one (or only a few) microservice(s), and it can be safely localized without affecting compatibility. In such cases, the overhead of re-building and re-running the docker containers could seriously degrade the development pace. 

The only service that is meant to continue running in Docker is DB. You can start it separately (without modifying any files) by running:
```sh
docker-compose up db
```

Once the DB service is started in a Docker container, it exposes the `TCP:5432` port on the `localhost` interface. If you now start a microservice outside the Docker environment, it will still be able to connect to DB to access data. Under the hood, it will read `.env` file from the root of the repo and will compile the value of `ENV_DSN` environment variable from the variables exposed in that file (see implementation of this logic in `lib/env/setup.go`)

## Start the service

Without loss of generality, start the service by changing into its root directory running the default source as shown below
```sh
cd auth-service
go run .
```
