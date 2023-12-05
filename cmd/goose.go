package main

import (
	"log"
	"os"

	_ "github.com/quible-io/quible-api/cmd/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {

	var dir, dsn, command, args = "./migrations", os.Getenv("ENV_DSN"), os.Args[1], os.Args[2:]

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("goose: failed to connect to DB server: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.Run(command, db, dir, args...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
