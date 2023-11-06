# auth-service 

A backend micro-service providing for user authentication and basic flow for account setup

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Contents

- [Setup local dev environment](#setup-local-dev-environment)
  - [Install dependencies](#install-dependencies)
  - [DBMS server](#dbms-server)
    - [On-fly setup](#on-fly-setup)
    - [Setup for docker-composer](#setup-for-docker-composer)
    - [Setup local DBMS server (macOS)](#setup-local-dbms-server-macos)
- [App server](#app-server)
  - [Run server](#run-server)
  - [Build swagger](#build-swagger)
- [Code review guideline](#code-review-guideline)
  - [Documentation](#documentation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Setup local dev environment

### Install dependencies

To run a local dev environment you need the following components
- [Functional Go (latest) setup](https://go.dev/doc/install)
- Docker desktop software to run local DB server.
  - [Install on macOS](https://docs.docker.com/desktop/install/mac-install/)
  - [Install on Windows](https://docs.docker.com/desktop/install/windows-install/)
  - [Install on Linux](https://docs.docker.com/desktop/install/linux-install/)
- Install [Swag](https://pkg.go.dev/github.com/swaggo/swag/v2) and make sure `$GOPATH/bin` is in your `$PATH`.
  - Run `go install github.com/swaggo/swag/v2/cmd/swag@latest`

### DBMS server

There are several ways of running local DB server:
- Creating a container on-fly using `docker run`
- Create a `docker-compose.yaml` file with service configuration and use `docker compose up` to run it
- Install a local DBMS server from OS package manager (not recommended)

The first two approaches (docker based ones) can be enhanced by allowing DBMS data directory linked to a local host directory, such that DB data will not be wiped after restarting the docker container. For that purpose you can use any existing directory, or create a new one called `dbData` in the root of this project (it is listed in `.gitignore` and its content will not be checked out).

#### On-fly setup

1. Download the latest `postgres` image by running
    ```
    docker pull postgres:latest
    ```
1. Start a detached container:    
    ```
    docker run -itd -e POSTGRES_USER=<username> -e POSTGRES_PASSWORD=<password> -p <local port>:5432 -v <DB data directory>:/var/lib/postgresql/data --name postgresql postgres
    ```

    where
    - `<username>` is the username of a user (to be defined) who is granted to access DB server
    - `<password>` desired password of that user
    - `<local port>` local TCP port exposed to your host system, reasonable default is `5432`
    - `<DB data directory>` am existing directory in your host file system to be used to persist DB data. Reasonable default is `$PWD/dbData` if the command is executed from repo's root and you previously created `dbData` there, as was suggested above.

1. Execute SQL queries from migration script(s) to initialize DB. This step can be done in several ways, for example by using DBMS client like `DBeaver`. Alternatively you can use `docker exec` to utilize `psql` command directly from the container to execute the migration:
    ```
    docker exec -i <container id> psql -d <DB name> -U <username> < <migrations file>
    ```
    where
    - `<container id>` is the ID of the running container. Can be discovered by running `docker ps` and observing the output to capture value of the "CONTAINER ID" column of the row corresponding to the container in question, i.e. PostgreSQL container.
    - `<username>` same as above (in the step describing setting up the DB container)
    - `<migrations file>` a relative/absolute path to the file containing SQL commands.
    - `<DB name>` is the name of DB for which SQL commands are targeted

#### Setup for docker-composer

Create/find a `docker-compose.yml` file in the root of the repo. Example below shows some concrete values

```yml
version: '3'
services:
  db:
    image: postgres
    environment:
      POSTGRES_USER: <username>
      POSTGRES_PASSWORD: <password>
    volumes:
      - <DB data directory>:/var/lib/postgresql/data
    ports:
      - <local port>:5432 
```

Where `<...>` parameters are set by analogy to the on-fly setup. The only exception is with `<DB data directory>` which, without loss of generality, is recommended to be set as `./dbData` (assuming such directory existence in the repo root).

Once `docker-compose.yml` is authored and placed in the repo root, you need to create start the composed service in the daemon mode (hence option `-d`):

```sh
docker-compose up -d
```

The rest of operations, i.e. migration of DB data can be done in the same way as was described above for in the on-fly approach.

#### Setup local DBMS server (macOS)

This approach is *NOT* recommended, it just explains behind the scene, what the manual steps are to setup DB.

- Install: `brew install postgresql`
- Start service: `brew services start postgresql`
  - To stop service later: `brew services stop postgresql`
- Start CLI: `psql postgres`
  - Create user: `CREATE ROLE your_db_username WITH LOGIN PASSWORD 'your_db_password';`
  - Alter user role: `ALTER ROLE newUser CREATEDB;`
  - Create database: `CREATE DATABASE your_db_name;`
    - DB name, username and password are currently hard-coded in `main.go`. It will be addressed later.
  - Choose database: `\c your_db_name`
  - Create table: run `migrations/000001_create_users_table.up.sql`.

## App server

### Run server

To start the app server we first need to install its dependencies. Go to the root of the repo, and run:

```sh
go mod tidy
```

Then run the server:

```sh
go run main.go
```

Or `make` to update swagger docs, then run go server.

All logs (including errors) will be printed to the console.

### Build swagger

To build swagger docs separately, run `swag init --output swagger --outputTypes yaml` under root level.

## Code review guideline

This part is under construction.

### Documentation

All functions should be properly annotated in `controller.go`, and code reviewers should be very cautious, to make sure code changes are reflected in annotations.

Tech tasks or bugs should be raised if inconsistency between documentation and implementation is observed.